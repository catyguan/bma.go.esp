package servproxy

import (
	"boot"
	"fmt"
	"logger"
	"strings"
)

type PortConfigInfo struct {
	Port      int
	Net       string
	Type      string
	GoLua     string
	Script    string
	TimeoutMS int
	WhiteIp   string
	whiteList []string
	BlackIp   string
	blackList []string
}

func (this *PortConfigInfo) Valid() error {
	if this.Net == "" {
		if this.Port > 0 {
			this.Net = fmt.Sprintf(":%d", this.Port)
		} else {
			return fmt.Errorf("Port invalid")
		}
	}
	if this.Type == "" {
		return fmt.Errorf("Type invalid")
	}
	_, err := AssertPortHandler(this.Type)
	if err != nil {
		return err
	}
	if this.GoLua == "" {
		return fmt.Errorf("GoLua invalid")
	}
	if this.Script == "" {
		return fmt.Errorf("Script invalid")
	}
	if this.WhiteIp != "" {
		this.whiteList = strings.Split(this.WhiteIp, ",")
	}
	if this.BlackIp != "" {
		this.blackList = strings.Split(this.BlackIp, ",")
	}
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 30 * 1000
	}
	return nil
}

func (this *PortConfigInfo) Compare(old *PortConfigInfo) bool {
	if this.Net != old.Net {
		return false
	}
	if this.Type != old.Type {
		return false
	}
	if this.GoLua != old.GoLua {
		return false
	}
	if this.Script != old.GoLua {
		return false
	}
	if this.WhiteIp != old.WhiteIp {
		return false
	}
	if this.BlackIp != old.BlackIp {
		return false
	}
	if this.TimeoutMS != old.TimeoutMS {
		return false
	}
	return true
}

type TargetConfigInfo struct {
	Type    string
	Remotes []*RemoteConfigInfo
}

func (this *TargetConfigInfo) Valid() error {
	if this.Type == "" {
		return fmt.Errorf("Type empty")
	}
	h, err := AssertRemoteHandler(this.Type)
	if err != nil {
		return err
	}
	if len(this.Remotes) == 0 {
		return fmt.Errorf("Remote empty")
	}
	for i, r := range this.Remotes {
		err = h.Valid(r)
		if err != nil {
			return fmt.Errorf("Remote(%d) invalid - %s", i, err)
		}
	}
	return nil
}

func (this *TargetConfigInfo) Compare(old *TargetConfigInfo) bool {
	if this.Type != old.Type {
		return false
	}
	h, err := AssertRemoteHandler(this.Type)
	if err != nil {
		return false
	}
	if len(this.Remotes) != len(old.Remotes) {
		for i, r := range this.Remotes {
			or := old.Remotes[i]
			if !r.Compare(or) {
				return false
			}
			if !h.Compare(r, or) {
				return false
			}
		}
	}
	return true
}

type RemoteConfigInfo struct {
	Host        string
	TimeoutMS   int
	Priority    int
	ReadOnly    bool
	FailRetryMS int
	Params      map[string]interface{}
}

func (this *RemoteConfigInfo) Compare(old *RemoteConfigInfo) bool {
	if this.Host != old.Host {
		return false
	}
	if this.TimeoutMS != old.TimeoutMS {
		return false
	}
	if this.Priority != old.Priority {
		return false
	}
	if this.ReadOnly != old.ReadOnly {
		return false
	}
	return true
}

type configInfo struct {
	Ports   map[string]*PortConfigInfo
	Targets map[string]*TargetConfigInfo
}

func (this *configInfo) Valid() error {
	for k, o := range this.Ports {
		err := o.Valid()
		if err != nil {
			return fmt.Errorf("Port(%s) %s", k, err)
		}
	}
	for k, o := range this.Targets {
		err := o.Valid()
		if err != nil {
			return fmt.Errorf("Target(%s) %s", k, err)
		}
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if true {
		if len(this.Ports) != len(old.Ports) {
			return boot.CCR_NEED_START
		}
		for k, o := range this.Ports {
			oo, ok := old.Ports[k]
			if ok {
				if !o.Compare(oo) {
					return boot.CCR_NEED_START
				}
			} else {
				return boot.CCR_NEED_START
			}
		}
	}
	if true {
		if len(this.Targets) != len(old.Targets) {
			return boot.CCR_NEED_START
		}
		for k, o := range this.Targets {
			oo, ok := old.Targets[k]
			if ok {
				if !o.Compare(oo) {
					return boot.CCR_NEED_START
				}
			} else {
				return boot.CCR_NEED_START
			}
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
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, cfg := range this.config.Ports {
		if _, ok := this.ports[k]; ok {
			continue
		}
		err := this._createPort(k, cfg)
		if err != nil {
			logger.Error(tag, "Port(%s) create fail - %s", k, err)
			return false
		}
	}
	for k, cfg := range this.config.Targets {
		if _, ok := this.targets[k]; ok {
			continue
		}
		err := this._createTarget(k, cfg)
		if err != nil {
			logger.Error(tag, "Target(%s) create fail - %s", k, err)
			return false
		}
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
	cfg := ccr.Config.(*configInfo)
	for k, o := range this.config.Ports {
		if cfg.Ports != nil {
			if oo, ok := cfg.Ports[k]; ok {
				if o.Compare(oo) {
					continue
				}
			}
		}
		this.RemovePort(k)
	}
	for k, o := range this.config.Targets {
		if cfg.Targets != nil {
			if oo, ok := cfg.Targets[k]; ok {
				if o.Compare(oo) {
					continue
				}
			}
		}
		this.RemoveTarget(k)
	}
	return true
}

func (this *Service) Stop() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, _ := range this.ports {
		this._removePort(k)
	}
	for k, _ := range this.targets {
		this._removeTarget(k)
	}
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
