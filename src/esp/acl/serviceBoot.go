package acl

import (
	"boot"
	"fmt"
	"logger"
)

const (
	defaultReportTimeMS = 5 * 1000
)

type configUserInfo struct {
	Id    string
	Host  []string
	Token string
	Name  string
	Group []string
}

func (this *configUserInfo) GetName(ip string) string {
	if this.Name == "" {
		return this.Id + "@" + ip
	}
	return this.Name
}

func (this *configUserInfo) Valid() error {
	if this.Id == "" {
		return fmt.Errorf("user id empty")
	}
	return nil
}

func (this *configUserInfo) Compare(old *configUserInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	ok := func() bool {
		if this.Id != old.Id {
			return false
		}
		if this.Token != old.Token {
			return false
		}
		if len(this.Host) != len(old.Host) {
			return false
		}
		for i, k := range this.Host {
			if old.Host[i] != k {
				return false
			}
		}
		if this.Name != old.Name {
			return false
		}
		if len(this.Group) != len(old.Group) {
			return false
		}
		for i, k := range this.Group {
			if old.Group[i] != k {
				return false
			}
		}
		return true
	}()
	if ok {
		return boot.CCR_NONE
	}
	return boot.CCR_CHANGE
}

type configPriInfo struct {
	Op  string
	Who []string
}

func (this *configPriInfo) Valid() error {
	if this.Op == "" {
		return fmt.Errorf("pri op empty")
	}
	return nil
}

func (this *configPriInfo) Compare(old *configPriInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	ok := func() bool {
		if this.Op != old.Op {
			return false
		}
		if len(this.Who) != len(old.Who) {
			return false
		}
		for i, k := range this.Who {
			if old.Who[i] != k {
				return false
			}
		}
		return true
	}()
	if ok {
		return boot.CCR_NONE
	}
	return boot.CCR_CHANGE
}

type configInfo struct {
	Users []*configUserInfo
	Ops   []*configPriInfo
}

func (this *configInfo) Valid() error {
	for _, o := range this.Users {
		err := o.Valid()
		if err != nil {
			return err
		}
	}
	for _, o := range this.Ops {
		err := o.Valid()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_CHANGE
	}

	if len(this.Users) != len(old.Users) {
		return boot.CCR_CHANGE
	}
	for i, o := range this.Users {
		if f := o.Compare(old.Users[i]); f != boot.CCR_NONE {
			return f
		}
	}

	if len(this.Ops) != len(old.Ops) {
		return boot.CCR_CHANGE
	}
	for i, o := range this.Ops {
		if f := o.Compare(old.Ops[i]); f != boot.CCR_NONE {
			return f
		}
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
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
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
