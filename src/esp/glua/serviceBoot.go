package glua

import (
	"boot"
	"fmt"
	"logger"
)

type serviceConfigInfo struct {
	GLua map[string]*ConfigInfo
}

func (this *serviceConfigInfo) Valid() error {
	for k, glcfg := range this.GLua {
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
	r := boot.CCR_NONE
	for k, o := range this.GLua {
		if old.GLua == nil {
			r = boot.CCR_CHANGE
			continue
		}
		oo, ok := old.GLua[k]
		if ok {
			cf := o.Compare(oo)
			if cf == boot.CCR_NEED_START {
				return boot.CCR_NEED_START
			}
		} else {
			r = boot.CCR_CHANGE
		}
	}
	if r == boot.CCR_NONE {
		for k, _ := range old.GLua {
			if this.GLua == nil {
				r = boot.CCR_CHANGE
				break
			}
			_, ok := this.GLua[k]
			if !ok {
				r = boot.CCR_CHANGE
				break
			}
		}
	}

	return r
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

func (this *Service) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, glcfg := range this.config.GLua {
		if _, ok := this.gluas[k]; ok {
			continue
		}
		gl := NewGLua(k, glcfg)
		if this.gluaInit != nil {
			this.gluaInit(gl)
		}
		err := gl.Run()
		if err != nil {
			logger.Error(tag, "start GLua['%s'] fail %s", k, err)
			return false
		}
		this.gluas[k] = gl
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
	for k, _ := range this.config.GLua {
		if ccr.Type != boot.CCR_NEED_START {
			if cfg.GLua != nil {
				if _, ok := cfg.GLua[k]; ok {
					continue
				}
			}
		}
		this.lock.Lock()
		gl := this.gluas[k]
		delete(this.gluas, k)
		this.lock.Unlock()

		if gl != nil {
			gl.Stop()
		}
	}
	return true
}

func (this *Service) Stop() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, gl := range this.gluas {
		gl.Stop()
		delete(this.gluas, k)
	}
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
