package servicecallcfg

import (
	"boot"
	"esp/servicecall"
	"fmt"
	"logger"
)

type configInfo struct {
	Services map[string]map[string]interface{}
}

func (this *configInfo) Valid() error {
	for n, mlcfg := range this.Services {
		err := servicecall.DoValid(mlcfg)
		if err != nil {
			return fmt.Errorf("'%s' %s", n, err)
		}
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if len(this.Services) != len(old.Services) {
		return boot.CCR_NEED_START
	}
	r := boot.CCR_NONE
	for k, o := range this.Services {
		oo, ok := old.Services[k]
		if ok {
			if !servicecall.DoCompare(o, oo) {
				return boot.CCR_NEED_START
			}
		} else {
			return boot.CCR_NEED_START
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
	for k, mlcfg := range this.config.Services {
		scid := this.scids[k]
		sok, nscid, err := servicecall.SetServiceCall(k, mlcfg, nil, scid)
		if err != nil {
			logger.Error(tag, "SetServiceCall('%s') fail - %s", k, err)
			return false
		}
		if sok {
			this.scids[k] = nscid
		} else {
			logger.Info(tag, "SetServiceCall('%s') not success", k)
			delete(this.scids, k)
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
	cfg := ccr.Config.(*configInfo)
	for k, o := range this.config.Services {
		if cfg.Services != nil {
			if oo, ok := cfg.Services[k]; ok {
				if servicecall.DoCompare(o, oo) {
					continue
				}
			}
		}
		scid := this.scids[k]
		if scid > 0 {
			if !servicecall.RemoveServiceCall(k, scid) {
				logger.Info(tag, "GraceStop serviceCall(%s) not success", k)
			}
		}
	}
	return true
}

func (this *Service) Stop() bool {
	return true
}

func (this *Service) Close() bool {
	for k, scid := range this.scids {
		servicecall.RemoveServiceCall(k, scid)
	}
	return true
}

func (this *Service) Cleanup() bool {
	servicecall.RemoveAll()
	return true
}
