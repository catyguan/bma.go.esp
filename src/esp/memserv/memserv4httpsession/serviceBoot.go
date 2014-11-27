package memserv4httpsession

import (
	"boot"
	"logger"
)

type config struct {
	CookieName    string
	SessionPrefix string
	TimeoutMS     int
}

func (this *config) Valid() error {
	if this.CookieName == "" {
		this.CookieName = "__GOLUA_SESSIONID"
	}
	if this.SessionPrefix == "" {
		this.SessionPrefix = "SESSION_"
	}
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 1000 * 30 * 60
	}
	return nil
}

func (this *config) Compare(old *config) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	ok := func() bool {
		if this.CookieName != old.CookieName {
			return false
		}
		if this.SessionPrefix != old.SessionPrefix {
			return false
		}
		if this.TimeoutMS != old.TimeoutMS {
			return false
		}
		return true
	}()
	if ok {
		return boot.CCR_NONE
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
	cfg := new(config)
	co.GetBeanConfig(this.name, cfg)
	if err := cfg.Valid(); err != nil {
		logger.Error(tag, "'%s' config error - %s", this.name, err)
		return false
	}
	ccr := boot.NewConfigCheckResult(cfg.Compare(this.cfg), cfg)
	ctx.CheckFlag = ccr
	return true
}

func (this *Service) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*config)
	this.cfg = cfg
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
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
