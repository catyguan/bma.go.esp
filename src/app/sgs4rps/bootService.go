package sgs4rps

import (
	"boot"
	"logger"
)

type configInfo struct {
	RobotNum     int
	GameDuration int
}

func (this *configInfo) Valid() error {
	if this.RobotNum <= 0 {
		this.RobotNum = 5
	}
	if this.GameDuration <= 0 {
		this.GameDuration = 10
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	return boot.CCR_CHANGE
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
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	this.goo.DoSync(func() {
		this.startMatrix()
	})
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	return this.Stop()
}

func (this *Service) Stop() bool {
	this.goo.DoSync(func() {
		this.stopMatrix()
	})
	return true
}

func (this *Service) Close() bool {
	this.goo.Stop()
	return true
}

func (this *Service) Cleanup() bool {
	this.goo.StopAndWait()
	return true
}
