package gom

import (
	"boot"
	"fileloader"
	"fmt"
	"golua"
	"logger"
)

type configInfo struct {
	VM      *golua.VMConfig
	FL      map[string]interface{}
	DevMode int
}

func (this *configInfo) Valid() error {
	if this.VM != nil {
		err := this.VM.Valid()
		if err != nil {
			return err
		}
	}
	if this.FL == nil {
		return fmt.Errorf("empty FileLoader")
	}
	err1 := fileloader.DoValid(this.FL)
	if err1 != nil {
		return err1
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.DevMode != old.DevMode {
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
	cfg := new(configInfo)
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
	cfg := ccr.Config.(*configInfo)
	this.config = cfg
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	if ccr.Type == boot.CCR_NEED_START {
		ss, err0 := fileloader.DoCreate(this.config.FL)
		if err0 != nil {
			logger.Error(tag, "create FileLoader[%s] fail %s", this.config.FL, err0)
			return false
		}
		this.floader = ss
		gl := golua.NewGoLua("gom", 16, ss, func(gl *golua.GoLua) {
			golua.InitCoreLibs(gl)
			InitGoLua(gl)
			if this.gli != nil {
				this.gli(gl)
			}
		}, this.config.VM)
		switch this.config.DevMode {
		case 0:
			gl.DevMode = boot.DevMode
		case 1:
			gl.DevMode = true
		case -1:
			gl.DevMode = false
		}
		this.gl = gl
	}
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}
	if this.gl != nil {
		this.gl.Close()
		this.gl = nil
	}
	return true
}

func (this *Service) Stop() bool {
	if this.gl != nil {
		this.gl.Close()
	}
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
