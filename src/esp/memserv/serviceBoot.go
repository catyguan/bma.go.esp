package memserv

import (
	"boot"
	"fmt"
	"logger"
)

type serviceConfigInfo struct {
	Configs map[string]*MemGoConfig
}

func (this *serviceConfigInfo) Valid() error {
	for k, cfg := range this.Configs {
		err := cfg.Valid()
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
	if len(this.Configs) != len(old.Configs) {
		return boot.CCR_CHANGE
	}
	for k, o := range this.Configs {
		oo, ok := old.Configs[k]
		if ok {
			if !o.Compare(oo) {
				return boot.CCR_CHANGE
			}
		} else {
			return boot.CCR_CHANGE
		}
	}
	return boot.CCR_NONE
}

func (this *MemoryServ) Name() string {
	return this.name
}

func (this *MemoryServ) Prepare() {
}
func (this *MemoryServ) CheckConfig(ctx *boot.BootContext) bool {
	co := ctx.Config
	cfg := new(serviceConfigInfo)
	co.GetBeanConfig(this.name, cfg)
	if err := cfg.Valid(); err != nil {
		logger.Error(tag, "'%s' config error - %s", this.name, err)
		return false
	}
	ccr := boot.NewConfigCheckResult(cfg.Compare(this.config), cfg)
	ctx.CheckFlag = ccr
	return true
}

func (this *MemoryServ) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*serviceConfigInfo)
	this.config = cfg
	return true
}

func (this *MemoryServ) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	return true
}

func (this *MemoryServ) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	return true
}

func (this *MemoryServ) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	return true
}

func (this *MemoryServ) Stop() bool {
	return true
}

func (this *MemoryServ) Close() bool {
	this.CloseAll(true)
	return true
}

func (this *MemoryServ) Cleanup() bool {
	return true
}
