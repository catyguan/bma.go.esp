package httpmux4goluaserv

import (
	"boot"
	"encoding/json"
	"fmt"
	"logger"
	"strings"
)

type configApp struct {
	Name      string
	Host      string
	Location  string
	IndexName string
	Script    string
	TimeoutMS int
	Skip      []string
}

func (this *configApp) Valid() error {
	if this.Name == "" {
		return fmt.Errorf("empty http golua app name")
	}
	if this.Location == "" {
		this.Location = "/"
	}
	if !strings.HasSuffix(this.Location, "/") {
		this.Location = this.Location + "/"
	}
	if this.IndexName == "" {
		this.IndexName = "index"
	}
	if this.Script == "" {
		this.Script = "main.lua"
	}
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 30 * 1000
	}
	return nil
}

type configInfo struct {
	App                []*configApp
	Skip               []string
	ParseFormMaxMemory int64
}

func (this *configInfo) Valid() error {
	for _, app := range this.App {
		err := app.Valid()
		if err != nil {
			return err
		}
	}
	if this.ParseFormMaxMemory <= 0 {
		this.ParseFormMaxMemory = 1024 * 1024 * 10
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	s1, _ := json.Marshal(this)
	s2, _ := json.Marshal(old)
	if s1 == nil || s2 == nil || string(s1) != string(s2) {
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
