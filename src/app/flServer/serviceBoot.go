package main

import (
	"boot"
	"fileloader"
	"fmt"
	"logger"
)

type configInfo struct {
	AdminIp string
	Key     string
	FL      map[string]interface{}
}

func (this *configInfo) Valid() error {
	if this.Key == "" {
		return fmt.Errorf("key empty")
	}
	if this.FL == nil {
		return fmt.Errorf("FL empty")
	}
	err := fileloader.CommonFileLoaderFactory.Valid(this.FL)
	if err != nil {
		return err
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if !fileloader.CommonFileLoaderFactory.Compare(this.FL, old.FL) {
		return boot.CCR_NEED_START
	}
	if this.Key != old.Key {
		return boot.CCR_CHANGE
	}
	if this.AdminIp != old.AdminIp {
		return boot.CCR_CHANGE
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
	this.config = ccr.Config.(*configInfo)
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	fl, err := fileloader.CommonFileLoaderFactory.Create(this.config.FL)
	if err != nil {
		logger.Error(tag, "create fileloader fail - %s", err)
	}
	this.fl = fl
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
