package goluaserv

import (
	"boot"
	"context"
	"fileloader"
	"fmt"
	"golua"
	"logger"
	"time"
)

type serviceConfigInfo struct {
	PoolSize       int
	GetTimeoutMS   int
	StartupRetryMS int
	GoLua          map[string]*goluaConfigInfo
}

func (this *serviceConfigInfo) Valid() error {
	if this.PoolSize <= 0 {
		this.PoolSize = 16
	}
	if this.GetTimeoutMS <= 0 {
		this.GetTimeoutMS = 5000
	}
	if this.StartupRetryMS <= 0 {
		this.StartupRetryMS = 5000
	}
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
	VM       *golua.VMConfig
	FL       map[string]interface{}
	DevMode  int
	FastBoot bool
	Startup  []string
}

func (this *goluaConfigInfo) Valid() error {
	if this.VM != nil {
		err := this.VM.Valid()
		if err != nil {
			return err
		}
	}
	if this.FL == nil {
		return fmt.Errorf("empty ScriptSource")
	}
	err1 := fileloader.DoValid(this.FL)
	if err1 != nil {
		return err1
	}
	return nil
}

func (this *goluaConfigInfo) Compare(old *goluaConfigInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.DevMode != old.DevMode {
		return boot.CCR_NEED_START
	}
	if len(this.Startup) != len(old.Startup) {
		return boot.CCR_NEED_START
	}
	if true {
		tmp := make(map[string]bool)
		for _, k := range this.Startup {
			tmp[k] = true
		}
		for _, k := range old.Startup {
			if _, ok := tmp[k]; !ok {
				return boot.CCR_NEED_START
			}
		}
	}

	if this.FastBoot != old.FastBoot {
		return boot.CCR_NEED_START
	}

	r2 := fileloader.DoCompare(this.FL, old.FL)
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

func (this *Service) startupApp(k string, gli *glInfo, startup []string) {
	go func() {
		this.lock.RLock()
		cgli, ok := this.gli[k]
		this.lock.RUnlock()
		if !ok {
			return
		}
		if cgli != gli {
			return
		}

		err := func() error {
			for _, ss := range startup {
				gl := gli.gl
				ri := golua.NewRequestInfo()
				ri.Script = ss

				ctx := context.Background()
				ctx, _ = context.CreateExecId(ctx)
				ctx = golua.CreateRequest(ctx, ri)
				_, err := gl.Execute(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}()

		this.lock.Lock()
		defer this.lock.Unlock()
		if err != nil {
			logger.Warn(tag, "[%s] startup '%s' fail - %s", k, startup, err)
			gli.startErr = err

			time.AfterFunc(time.Duration(this.config.StartupRetryMS)*time.Millisecond, func() {
				this.startupApp(k, gli, startup)
			})
		} else {
			gli.status = 1
		}
	}()
}

func (this *Service) _create(k string, glcfg *goluaConfigInfo) bool {
	ss, err0 := fileloader.DoCreate(glcfg.FL)
	if err0 != nil {
		logger.Error(tag, "create ScriptSource[%s, %s] fail %s", k, glcfg.FL, err0)
		return false
	}
	gli := new(glInfo)
	gl := golua.NewGoLua(k, this.config.PoolSize, ss, this.glInit, glcfg.VM)
	switch glcfg.DevMode {
	case 0:
		gl.DevMode = boot.DevMode
	case 1:
		gl.DevMode = true
	case -1:
		gl.DevMode = false
	}
	gli.gl = gl
	this.gli[k] = gli

	if len(glcfg.Startup) > 0 {
		this.startupApp(k, gli, glcfg.Startup)
		if glcfg.FastBoot {
			gli.status = 1
		}
	} else {
		gli.status = 1
	}
	smm := new(smmObject)
	smm.s = this
	smm.gl = gl
	gl.ExtSMMApi = smm
	gl.InitSMMApi()

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
		if _, ok := this.gli[k]; ok {
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
	for k, oglcfg := range this.config.GoLua {
		closed := false
		if cfg.GoLua != nil {
			glcfg, ok := cfg.GoLua[k]
			if ok {
				cr := glcfg.Compare(oglcfg)
				if cr != boot.CCR_NONE {
					closed = true
				}
			} else {
				closed = true
			}
		}
		if closed {
			gli := this.removeGoLua(k)
			if gli != nil && gli.gl != nil {
				gli.gl.Close()
				fmt.Printf("close GoLua '%s'\n", k)
			}
		}
	}
	return true
}

func (this *Service) Stop() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, gli := range this.gli {
		if gli.gl != nil {
			gli.gl.Close()
		}
		delete(this.gli, k)
	}
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
