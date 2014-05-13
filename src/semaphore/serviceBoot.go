package semaphore

import (
	"boot"
	"fmt"
	"logger"
)

type semInfo struct {
	Name             string
	Limit            int
	ExecuteTimeoutMS int
}

func (this *semInfo) Valid(vn bool) error {
	if vn && this.Name == "" {
		return fmt.Errorf("sem name empty")
	}
	if this.Limit <= 0 {
		return fmt.Errorf("sem limit invalid[%d]", this.Limit)
	}
	if this.ExecuteTimeoutMS <= 0 {
		this.ExecuteTimeoutMS = 30 * 1000
	}
	return nil
}

func (this *semInfo) Compare(old *semInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.Limit != old.Limit {
		return boot.CCR_NEED_START
	}
	if this.ExecuteTimeoutMS != old.ExecuteTimeoutMS {
		return boot.CCR_CHANGE
	}
	return boot.CCR_NONE
}

type configInfo struct {
	DefaultSem *semInfo
	Sems       []*semInfo
}

func (this *configInfo) Find(n string) *semInfo {
	for _, sem := range this.Sems {
		if sem.Name == n {
			return sem
		}
	}
	return nil
}

func (this *configInfo) Valid() error {
	if this.DefaultSem == nil {
		this.DefaultSem = new(semInfo)
		this.DefaultSem.Limit = 32
	}
	if err := this.DefaultSem.Valid(false); err != nil {
		return err
	}
	for _, sem := range this.Sems {
		if err := sem.Valid(true); err != nil {
			return err
		}
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	r := boot.CCR_NONE
	cp1 := this.DefaultSem.Compare(old.DefaultSem)
	if cp1 == boot.CCR_NEED_START {
		return cp1
	}
	if cp1 == boot.CCR_CHANGE {
		r = boot.CCR_CHANGE
	}
	for _, sem := range old.Sems {
		news := this.Find(sem.Name)
		if news == nil {
			return boot.CCR_NEED_START
		}
	}
	for _, sem := range this.Sems {
		olds := old.Find(sem.Name)
		if olds == nil {
			r = boot.CCR_CHANGE
		}
		cp0 := sem.Compare(olds)
		if cp0 == boot.CCR_NEED_START {
			return cp0
		}
		if cp0 == boot.CCR_CHANGE {
			r = boot.CCR_CHANGE
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
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}
	cfg := ccr.Config.(*configInfo)
	this.slock.Lock()
	defer this.slock.Unlock()
	for k, so := range this.sems {
		news := cfg.Find(k)
		if news != nil {
			cr := news.Compare(so.info)
			if cr != boot.CCR_NEED_START {
				continue
			}
		}
		delete(this.sems, k)
		so.Close()
	}
	return true
}

func (this *Service) Stop() bool {
	return true
}

func (this *Service) Close() bool {
	this.slock.Lock()
	defer this.slock.Unlock()
	for k, so := range this.sems {
		delete(this.sems, k)
		so.Close()
	}
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
