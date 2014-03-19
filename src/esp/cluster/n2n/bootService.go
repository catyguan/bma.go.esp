package n2n

import (
	"boot"
	"esp/espnet/esnp"
	"fmt"
	"logger"
)

type remoteConfigInfo struct {
	URL  string
	eurl *esnp.URL
}

type configInfo struct {
	URL    string
	eurl   *esnp.URL
	Remote map[string]*remoteConfigInfo
}

func (this *configInfo) Valid() error {
	if this.URL == "" {
		return fmt.Errorf("URL invalid")
	}
	if true {
		v, err := esnp.ParseURL(this.URL)
		if err != nil {
			return fmt.Errorf("URL invalid %s", err)
		}
		host := v.GetHost()
		if host == "" {
			return fmt.Errorf("URL address invalid")
		}
		this.eurl = v
	}

	for k, remote := range this.Remote {
		if remote.URL == "" {
			return fmt.Errorf("Remote[%s] invalid", k)
		}
		v, err := esnp.ParseURL(remote.URL)
		if err != nil {
			return fmt.Errorf("Remote[%s] URL invalid %s", k, err)
		}
		host := v.GetHost()
		if host == "" {
			return fmt.Errorf("Remote[%s] URL address", k)
		}
		remote.eurl = v
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
		this.checkConnector(k, ro.eurl)
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
		this.closeConnector(k)
	}
	return true
}

func (this *Service) Stop() bool {
	this.closeAllConnectors()
	this.closeAllRemote()
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
