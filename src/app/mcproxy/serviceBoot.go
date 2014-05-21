package main

import (
	"boot"
	"logger"
	"time"
)

const (
	defaultReportTimeMS = 5 * 1000
)

type configInfo struct {
	Version      string
	PoolMax      int
	Trace        int
	Remotes      []string
	ReportTimeMS int
}

func (this *configInfo) Valid() error {
	if this.Version == "" {
		this.Version = "1.0.0"
	}
	if this.PoolMax <= 0 {
		this.PoolMax = 10
	}
	if this.ReportTimeMS <= 0 {
		this.ReportTimeMS = defaultReportTimeMS
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.PoolMax != old.PoolMax {
		return boot.CCR_NEED_START
	}
	return boot.CCR_CHANGE
}

func (this *Service) Name() string {
	return this.name
}

func (this *Service) Prepare() {
	this.reportC <- true
	go func() {
		for {
			doit := <-this.reportC
			if !doit {
				return
			}

			this.report()

			tm := defaultReportTimeMS
			cfg := this.config
			if cfg != nil {
				tm = cfg.ReportTimeMS
			}
			time.AfterFunc(time.Duration(tm)*time.Millisecond, func() {
				defer func() {
					recover()
				}()
				this.reportC <- true
			})
		}
	}()
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
	cfg := this.config
	this.plock.Lock()
	for _, r := range cfg.Remotes {
		if _, ok := this.remotes[r]; ok {
			continue
		}
		this._createRemote(r)
	}
	this.plock.Unlock()
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

	this.plock.Lock()
	for k, rmt := range this.remotes {
		del := true
		if ccr.Type == boot.CCR_CHANGE {
			for _, remote := range cfg.Remotes {
				if remote == k {
					del = false
					break
				}
			}
		}
		if del {
			this._closeRemote(k, rmt)
		}
	}
	this.plock.Unlock()
	return true
}

func (this *Service) Stop() bool {
	this.plock.Lock()
	for k, rmt := range this.remotes {
		this._closeRemote(k, rmt)
	}
	this.plock.Unlock()
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	close(this.reportC)
	return true
}
