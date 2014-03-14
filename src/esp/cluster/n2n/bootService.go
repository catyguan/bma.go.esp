package n2n

import (
	"boot"
	"esp/espnet/esnp"
	"fmt"
	"logger"
)

type remoteConfigInfo struct {
	URL string
}

type configInfo struct {
	URL    string
	Remote map[string]*remoteConfigInfo
}

func (this *configInfo) Valid() error {
	if this.URL == "" {
		return fmt.Errorf("URL invalid")
	}
	if true {
		addr, err := esnp.ParseAddress(this.URL)
		if err != nil {
			return err
		}
		if addr.GetHost() == "" {
			return fmt.Errorf("URL host invalid")
		}
	}

	for k, remote := range this.Remote {
		if remote.URL == "" {
			return fmt.Errorf("Remote[%s] invalid", k)
		}
		addr, err := esnp.ParseAddress(this.URL)
		if err != nil {
			return fmt.Errorf("Remote[%s] URL invalid %s", k, err)
		}
		if addr.GetHost() == "" {
			return fmt.Errorf("Remote[%s] URL host invalid")
		}
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.URL != old.URL {
		return boot.CCR_NEED_START
	}
	r := boot.CCR_NONE
	for k, ro := range this.Remote {
		if old.Remote == nil {
			r = boot.CCR_CHANGE
			continue
		}
		oro, ok := old.Remote[k]
		if ok {
			if ro.URL != oro.URL {
				return boot.CCR_NEED_START
			}
		} else {
			r = boot.CCR_CHANGE
		}
	}
	if r == boot.CCR_NONE {
		for k, _ := range old.Remote {
			if this.Remote == nil {
				r = boot.CCR_CHANGE
				break
			}
			_, ok := this.Remote[k]
			if !ok {
				r = boot.CCR_CHANGE
				break
			}
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
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}

	for k, ro := range this.config.Remote {
		this.checkAndConnect(k, ro.URL)
	}
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*configInfo)
	for k, _ := range this.config.Remote {
		if ccr.Type != boot.CCR_NEED_START {
			if cfg.Remote != nil {
				if _, ok := cfg.Remote[k]; ok {
					continue
				}
			}
		}
		this.closeRemote(k)
	}
	return true
}

func (this *Service) Stop() bool {
	this.closeAllRemote()
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
