package servbox

import (
	"boot"
	"logger"
)

type configInfo struct {
	TimeoutMS   int
	MaxConnSize int
	MaxPackage  int
}

func (this *configInfo) Valid() error {
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 5 * 1000
	}
	if this.MaxConnSize <= 0 {
		this.MaxConnSize = 1024
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.MaxConnSize != old.MaxConnSize {
		return boot.CCR_CHANGE
	}
	if this.TimeoutMS != old.TimeoutMS {
		return boot.CCR_CHANGE
	}
	if this.MaxPackage != old.MaxPackage {
		return boot.CCR_CHANGE
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
	return true
}

func (this *Service) Stop() bool {
	return true
}

func (this *Service) Close() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	for n, _ := range this.servs {
		delete(this.servs, n)
	}
	for n, node := range this.nodes {
		delete(this.nodes, n)
		node.Close()
	}
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
