package aclserv

import (
	"acl"
	"boot"
	"logger"
	"strings"
)

type configInfo map[string]map[string]interface{}

func (this configInfo) Valid() error {
	for _, sub := range this {
		err := acl.ValidConfig(sub)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this configInfo) Compare(old configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	for k, cur := range this {
		o, ok := old[k]
		if !ok {
			return boot.CCR_CHANGE
		}
		if !acl.CompareConfig(cur, o) {
			return boot.CCR_CHANGE
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
	tmp := co.GetMapConfig(this.name)
	if tmp == nil {
		logger.Error(tag, "'%s' miss config", this.name)
		return false
	}
	cfg := make(configInfo)
	for k, v := range tmp {
		if mv, ok := v.(map[string]interface{}); ok {
			cfg[k] = mv
		}
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
	cfg := ccr.Config.(configInfo)
	this.config = cfg

	rt := acl.NewRuleTree()

	for path, co := range cfg {
		nlist := strings.Split(path, "/")
		cur := rt
		for _, n := range nlist {
			cur = cur.Node(n)
		}
		rule, err := acl.CreateRule(co)
		if err != nil {
			logger.Error(tag, "create rule '%s' fail - %s", path, err)
			return false
		}
		cur.Append(rule)
	}

	acl.InitRuleTree(rt)

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
