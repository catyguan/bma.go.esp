package main

import (
	"boot"
	"fmt"
	"logger"
	"math/rand"
	"time"
)

type configInfo struct {
	EsnpAddress string
	CookiePath  string
	TimeoutMS   int
	ExpiresSec  int
}

func (this *configInfo) Valid() error {
	if this.EsnpAddress == "" {
		return fmt.Errorf("EsnpAddress invalid")
	}
	if this.CookiePath == "" {
		this.CookiePath = "/"
	}
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 5000
	}
	if this.ExpiresSec <= 0 {
		this.ExpiresSec = 5 * 60
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.EsnpAddress != old.EsnpAddress {
		return boot.CCR_NEED_START
	}
	if this.CookiePath != old.CookiePath {
		return boot.CCR_NEED_START
	}
	if this.TimeoutMS != old.TimeoutMS {
		return boot.CCR_CHANGE
	}
	if this.ExpiresSec != old.ExpiresSec {
		return boot.CCR_CHANGE
	}
	return boot.CCR_NONE
}

func (this *Service) Name() string {
	return this.name
}

func (this *Service) Prepare() {
	this.sockets = make(map[string]*SockInfo)
	this.robj = rand.New(rand.NewSource(time.Now().UnixNano()))
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
	this.config = ccr.Config.(*configInfo)
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}
	this.closeAll()
	return true
}

func (this *Service) Stop() bool {
	this.closeAll()
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
