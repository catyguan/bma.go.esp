package netserver

import (
	"bmautil/netutil"
	"boot"
	"fmt"
	"logger"
	"strings"
)

type configInfo struct {
	Net       string
	Address   string
	Port      int
	TimeoutMS int
	WhiteIp   string
	whiteList []string
	BlackIp   string
	blackList []string
	Disable   bool
}

func (this *configInfo) Valid() error {
	if this.Net == "" {
		this.Net = "tcp"
	}
	if this.Address == "" {
		if this.Port > 0 {
			this.Address = logger.Sprintf(":%d", this.Port)
		} else {
			return fmt.Errorf("port invalid")
		}
	}
	if this.WhiteIp != "" {
		this.whiteList = strings.Split(this.WhiteIp, ",")
	}
	if this.BlackIp != "" {
		this.blackList = strings.Split(this.BlackIp, ",")
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.Net != old.Net {
		return boot.CCR_NEED_START
	}
	if this.Address != old.Address {
		return boot.CCR_NEED_START
	}
	if this.TimeoutMS != old.TimeoutMS {
		return boot.CCR_NEED_START
	}
	if this.Disable != old.Disable {
		return boot.CCR_NEED_START
	}
	if this.WhiteIp != old.WhiteIp {
		return boot.CCR_CHANGE
	}
	if this.BlackIp != old.BlackIp {
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
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}
	if this.config.Disable {
		logger.Info(tag, "%s disable, skip start", this.name)
		return true
	}

	logger.Debug(tag, "%s start listen (%s)", this.name, this.config.Address)
	lis, err := netutil.Listen(this.config.Net, this.config.Address)
	if err != nil {
		logger.Warn(tag, "%s listen at (%s, %s) fail - %s", this.name, this.config.Net, this.config.Address, err)
		return false
	}
	logger.Info(tag, "%s listen at (%s, %s)", this.name, this.config.Net, this.config.Address)
	this.listener = lis
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}

	if this.config.Disable {
		logger.Info(tag, "%s disable, skip run", this.name)
		return true
	}

	go func() {
		defer func() {
			logger.Info(tag, "stop (%s, %s)", this.config.Net, this.config.Address)
		}()
		for {
			c, err := this.listener.Accept()
			if err == nil {
				if c == nil {
					return
				}
				addr := c.RemoteAddr().String()
				if ok, msg := netutil.IpAccept(addr, this.config.whiteList, this.config.blackList, true); !ok {
					logger.Warn(tag, "unaccept(%s) address %s", msg, addr)
					c.Close()
					continue
				}
				go this.accept(c) // new connect
			} else {
				logger.Debug(tag, "accept fail and exit - %s", err)
				return
			}
		}
	}()
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}
	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}
	return true
}

func (this *Service) Stop() bool {
	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
