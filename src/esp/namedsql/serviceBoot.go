package namedsql

import (
	"boot"
	"fmt"
	"logger"
)

type sqlConfig struct {
	Driver       string
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	DelayOpen    bool
}

func (this *sqlConfig) Valid() error {
	if this.Driver == "" {
		return fmt.Errorf("SQL Driver empty")
	}
	if this.DataSource == "" {
		return fmt.Errorf("SQL DataSource empty")
	}
	return nil
}

func (this *sqlConfig) Compare(old *sqlConfig) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.Driver != old.Driver {
		return boot.CCR_NEED_START
	}
	if this.DataSource != old.DataSource {
		return boot.CCR_NEED_START
	}
	if this.MaxIdleConns != old.MaxIdleConns {
		return boot.CCR_CHANGE
	}
	if this.MaxOpenConns != old.MaxOpenConns {
		return boot.CCR_CHANGE
	}
	if this.DelayOpen != old.DelayOpen {
		return boot.CCR_CHANGE
	}
	return boot.CCR_NONE
}

type configInfo struct {
	SQL map[string]*sqlConfig
}

func (this *configInfo) Valid() error {
	for k, o := range this.SQL {
		err := o.Valid()
		if err != nil {
			return fmt.Errorf("%s %s", k, err)
		}
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if len(this.SQL) != len(this.SQL) {
		return boot.CCR_NEED_START
	}
	r := boot.CCR_NONE
	for k, o := range this.SQL {
		oo, ok := old.SQL[k]
		if ok {
			cf := o.Compare(oo)
			switch cf {
			case boot.CCR_NEED_START:
				return boot.CCR_NEED_START
			case boot.CCR_CHANGE:
				r = cf
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
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for k, dbcfg := range this.config.SQL {
		if dbi, ok := this.dbs[k]; ok {
			dbi.config = dbcfg
			if dbi.db != nil {
				dbi.db.SetMaxIdleConns(dbcfg.MaxIdleConns)
				dbi.db.SetMaxOpenConns(dbcfg.MaxOpenConns)
			}
			continue
		}
		this._create(k, dbcfg)
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
	this.mutex.Lock()
	defer this.mutex.Unlock()
	cfg := ccr.Config.(*configInfo)
	for k, ncfg := range this.config.SQL {
		if cfg.SQL != nil {
			if ocfg, ok := cfg.SQL[k]; ok {
				if ncfg.Compare(ocfg) != boot.CCR_NEED_START {
					continue
				}
			}
		}
		this._remove(k)

	}
	return true
}

func (this *Service) Stop() bool {
	return true
}

func (this *Service) Close() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for k, _ := range this.dbs {
		// fmt.Println("here", 4)
		this._remove(k)
	}
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
