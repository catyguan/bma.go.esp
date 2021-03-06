package emnode

import (
	"boot"
	"fmt"
	"logger"
)

type configInfo struct {
	NodeId uint64
}

func (this *configInfo) Valid() error {
	if this.NodeId == 0 {
		return fmt.Errorf("NodeId invalid")
	}
	return nil
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
	ccr := boot.NewConfigCheckResult(boot.CCR_CHANGE, cfg)
	ctx.CheckFlag = ccr
	return true
}

func (this *Service) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*configInfo)
	this.nodeId = NodeId(cfg.NodeId)
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) Stop() bool {
	return true
}

func (this *Service) Close() bool {
	nl := make([]string, 0)
	this.lock.RLock()
	for n, _ := range this.groups {
		nl = append(nl, n)
	}
	this.lock.RUnlock()

	for _, n := range nl {
		this.CloseGroup(n, false)
	}
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
