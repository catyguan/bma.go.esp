package servicecall

import (
	"boot"
	"fmt"
	"logger"
)

type configInfo struct {
	Services map[string]map[string]interface{}
}

func (this *configInfo) Valid(s *Service) error {
	for n, mlcfg := range this.Services {
		scf, _, err := s.GetServiceCallerFactoryByType(mlcfg)
		if err != nil {
			return fmt.Errorf("'%s' %s", n, err)
		}
		err2 := scf.Valid(mlcfg)
		if err2 != nil {
			return fmt.Errorf("'%s' %s", n, err2)
		}
	}
	return nil
}

func (this *configInfo) Compare(s *Service, old *configInfo) int {
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
			if !s.DoCompare(o, oo) {
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
	if err := cfg.Valid(this); err != nil {
		logger.Error(tag, "'%s' config error - %s", this.name, err)
		return false
	}
	ccr := boot.NewConfigCheckResult(cfg.Compare(this, this.config), cfg)
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
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, mlcfg := range this.config.Services {
		if _, ok := this.services[k]; ok {
			continue
		}
		_, err := this._create(k, mlcfg)
		if err != nil {
			logger.Error(tag, "ServiceCaller('%s') create fail - %s", k, err)
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
	cfg := ccr.Config.(*configInfo)
	for k, o := range this.config.Services {
		if cfg.Services != nil {
			if oo, ok := cfg.Services[k]; ok {
				if this.DoCompare(o, oo) {
					continue
				}
			}
		}
		this.RemoveServiceCall(k)
	}
	return true
}

func (this *Service) Stop() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, s := range this.services {
		delete(this.services, k)
		s.Stop()
	}
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
