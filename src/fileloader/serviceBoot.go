package fileloader

import (
	"boot"
	"fmt"
	"logger"
)

const (
	tag = "fileloader"
)

type Service struct {
	name   string
	config *configInfo
}

func NewService(n string) *Service {
	r := new(Service)
	r.name = n
	return r
}

type configInfo struct {
	MLoader map[string]map[string]interface{}
}

func (this *configInfo) Valid() error {
	for n, mlcfg := range this.MLoader {
		ff, _, err := GetFileLoaderFactoryByType(mlcfg)
		if err != nil {
			return fmt.Errorf("'%s' %s", n, err)
		}
		err2 := ff.Valid(mlcfg)
		if err2 != nil {
			return fmt.Errorf("'%s' %s", n, err2)
		}
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if len(this.MLoader) != len(old.MLoader) {
		return boot.CCR_NEED_START
	}
	r := boot.CCR_NONE
	for k, o := range this.MLoader {
		oo, ok := old.MLoader[k]
		if ok {
			if !commonFileLoaderFactory(0).Compare(o, oo) {
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

func (this *Service) _create(k string, mlcfg map[string]interface{}) bool {
	fl, err := commonFileLoaderFactory(0).Create(mlcfg)
	if err != nil {
		logger.Error(tag, "FileLoader('%s') create fail - %s", k, err)
		return false
	}
	SetModuleFileLoader(k, fl)
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	for k, mlcfg := range this.config.MLoader {
		if this._create(k, mlcfg) {
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
	for k, o := range this.config.MLoader {
		if cfg.MLoader != nil {
			if oo, ok := cfg.MLoader[k]; ok {
				if commonFileLoaderFactory(0).Compare(o, oo) {
					continue
				}
			}
		}
		RemoveModuleFileLoader(k)
	}
	return true
}

func (this *Service) Stop() bool {
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	if this.config != nil {
		for k, _ := range this.config.MLoader {
			RemoveModuleFileLoader(k)
		}
	}
	return true
}
