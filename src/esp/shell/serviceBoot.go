package shell

import (
	"boot"
	"logger"
)

type configInfo struct {
	Paths    []string
	Preloads []string
}

func (this *configInfo) Valid() error {
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	// compare Paths
	same := func() bool {
		if len(this.Paths) != len(old.Paths) {
			return false
		}
		tmp := make(map[string]bool)
		for _, s := range this.Paths {
			tmp[s] = true
		}
		for _, s := range old.Paths {
			if _, ok := tmp[s]; !ok {
				return false
			}
		}
		return true
	}()
	if !same {
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
	/*
		this.lock.Lock()
		defer this.lock.Unlock()
		for k, glcfg := range this.config.GLua {
			if _, ok := this.gluas[k]; ok {
				continue
			}
			if !this._create(k, glcfg) {
				return false
			}
		}
	*/
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
	/*
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
	*/
	return true
}

func (this *Service) Stop() bool {
	/*
		this.lock.Lock()
		defer this.lock.Unlock()

			for k, gl := range this.gluas {
				gl.Stop()
				delete(this.gluas, k)
			}
	*/
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
