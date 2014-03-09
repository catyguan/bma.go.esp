package boot

import "config"

type BootContext struct {
	IsRestart bool
	Config    config.ConfigObject
	CheckFlag interface{}
}

func (this *BootContext) CheckResult() *ConfigCheckResult {
	return this.CheckFlag.(*ConfigCheckResult)
}

type BootObject interface {
	Prepare()

	CheckConfig(ctx *BootContext) bool

	Init(ctx *BootContext) bool
	Start(ctx *BootContext) bool
	Run(ctx *BootContext) bool
	GraceStop(ctx *BootContext) bool
	Stop() bool
	Close() bool
	Cleanup() bool
}

type SupportName interface {
	Name() string
}

type ConfigCheckResult struct {
	Type   int
	Config interface{}
}

const (
	CCR_NEED_START = 0
	CCR_CHANGE     = 1
	CCR_NONE       = 2
)

func NewConfigCheckResult(t int, v interface{}) *ConfigCheckResult {
	r := new(ConfigCheckResult)
	r.Type = t
	r.Config = v
	return r
}
