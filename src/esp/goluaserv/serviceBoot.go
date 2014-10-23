package goluaserv

import (
	"boot"
	"fileloader"
	"fmt"
	"golua"
	"logger"
)

type serviceConfigInfo struct {
	GoLua map[string]*goluaConfigInfo
}

func (this *serviceConfigInfo) Valid() error {
	for k, glcfg := range this.GoLua {
		err := glcfg.Valid()
		if err != nil {
			return fmt.Errorf("%s error - %s", k, err)
		}
	}
	return nil
}

func (this *serviceConfigInfo) Compare(old *serviceConfigInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if len(this.GoLua) != len(old.GoLua) {
		return boot.CCR_NEED_START
	}
	r := boot.CCR_NONE
	for k, o := range this.GoLua {
		oo, ok := old.GoLua[k]
		if ok {
			cf := o.Compare(oo)
			if cf == boot.CCR_NEED_START {
				return boot.CCR_NEED_START
			}
		} else {
			return boot.CCR_NEED_START
		}
	}
	return r
}

type goluaConfigInfo struct {
	VM *golua.VMConfig
	FL map[string]interface{}
}

func (this *goluaConfigInfo) Valid() error {
	if this.VM != nil {
		err := this.VM.Valid()
		if err != nil {
			return err
		}
	}
	fac := fileloader.CommonFileLoaderFactory
	if this.FL == nil {
		return fmt.Errorf("empty ScriptSource")
	}
	err1 := fac.Valid(this.FL)
	if err1 != nil {
		return err1
	}
	return nil
}

func (this *goluaConfigInfo) Compare(old *goluaConfigInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	fac := fileloader.CommonFileLoaderFactory
	r2 := fac.Compare(this.FL, old.FL)
	if !r2 {
		return boot.CCR_NEED_START
	}

	if this.VM != nil {
		if old.VM == nil {
			return boot.CCR_CHANGE
		}
		return this.VM.Compare(old.VM)
	} else {
		if old.VM != nil {
			return boot.CCR_CHANGE
		}
	}
	return boot.CCR_NONE
}

func (this *Service) Name() string {
	return this.name
}

func (this *Service) Prepare() {
}
func (this *Service) CheckConfig(ctx *boot.BootContext) bool {
	co := ctx.Config
	cfg := new(serviceConfigInfo)
	if !co.GetBeanConfig(this.name, cfg) {
		logger.Error(tag, "'%s' miss config", this.name)
		return false
	}
	if err := cfg.Valid(); err != nil {
		logger.Error(tag, "'%s' config error - %s", this.name, err)
		return false
	}
	ccr := boot.NewConfigCheckResult(cfg.Compare(this.config), cfg)
	ctx.CheckFlag = ccr
	return true
}

func (this *Service) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*serviceConfigInfo)
	this.config = cfg
	return true
}

func (this *Service) _create(k string, glcfg *goluaConfigInfo) bool {
	fac := fileloader.CommonFileLoaderFactory
	ss, err0 := fac.Create(glcfg.FL)
	if err0 != nil {
		logger.Error(tag, "create ScriptSource['%s', %s] fail %s", k, glcfg.FL, err0)
		return false
	}
	gl := golua.NewGoLua(k, ss, this.vmgInit, glcfg.VM)
	this.gl[k] = gl
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, glcfg := range this.config.GoLua {
		if _, ok := this.gl[k]; ok {
			continue
		}
		if !this._create(k, glcfg) {
			return false
		}
	}
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*serviceConfigInfo)
	for k, _ := range this.config.GoLua {
		if ccr.Type != boot.CCR_NEED_START {
			if cfg.GoLua != nil {
				if _, ok := cfg.GoLua[k]; ok {
					continue
				}
			}
		}
		gl := this.removeGoLua(k)
		if gl != nil {
			gl.Close()
		}
	}
	return true
}

func (this *Service) Stop() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, gl := range this.gl {
		gl.Close()
		delete(this.gl, k)
	}
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
