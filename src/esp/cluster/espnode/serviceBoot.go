package espnode

import (
	"boot"
	"esp/cluster/n2n"
	"esp/cluster/nodebase"
	"fmt"
	"logger"
)

type configInfo struct {
	Id   uint64
	Name string
	N2N  *n2n.ConfigInfo
}

func (this *configInfo) Valid() error {
	if this.Id == 0 {
		return fmt.Errorf("Id invalid")
	}
	if this.Name == "" {
		return fmt.Errorf("Name invalid")
	}
	if this.N2N == nil {
		return fmt.Errorf("N2N invalid")
	}
	err1 := this.N2N.Valid()
	if err1 != nil {
		return err1
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if !this.N2N.Compare(old.N2N) {
		return boot.CCR_NEED_START
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

	if nodebase.Id != 0 {
		if uint64(nodebase.Id) != cfg.Id {
			logger.Error(tag, "'%s' NodeId can't change(%d)", nodebase.Id)
			return false
		}
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
	nodebase.Id = nodebase.NodeId(cfg.Id)
	nodebase.Name = cfg.Name
	this.config = cfg
	this.n2n.InitConfig(cfg.N2N)
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	if !this.n2n.Start() {
		return false
	}
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	if !this.n2n.Run() {
		return false
	}
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*configInfo)
	if !this.n2n.GraceStop(cfg.N2N) {
		return false
	}
	return true
}

func (this *Service) Stop() bool {
	this.n2n.Stop()
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	this.n2n.Cleanup()
	return true
}
