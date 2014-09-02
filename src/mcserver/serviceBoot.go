package mcserver

import (
	"boot"
	"fmt"
	"logger"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

type configInfo struct {
	Net     string
	Address string
	Port    int
	WhiteIp string
	BlackIp string
	Disable bool
}

func (this *configInfo) Valid() error {
	if this.Disable {
		return nil
	}
	if this.Address == "" {
		if this.Port > 0 {
			this.Address = fmt.Sprintf(":%d", this.Port)
		} else {
			return fmt.Errorf("port[%d] invalid", this.Port)
		}
	}
	if this.Net == "" {
		this.Net = "tcp"
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	r := boot.CCR_NONE
	if this.Disable != old.Disable {
		return boot.CCR_NEED_START
	}
	if this.Address != old.Address || this.Net != old.Net {
		return boot.CCR_NEED_START
	}
	if this.WhiteIp != old.WhiteIp || this.BlackIp != old.BlackIp {
		r = boot.CCR_CHANGE
	}
	return r
}

func (this *MemcacheServer) Name() string {
	return this.name
}

func (this *MemcacheServer) Prepare() {
}
func (this *MemcacheServer) CheckConfig(ctx *boot.BootContext) bool {
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

func (this *MemcacheServer) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*configInfo)
	this.config = cfg
	if cfg.Disable {
		logger.Info(tag, "'%s' disabled", this.name)
		return true
	}
	if cfg.WhiteIp != "" {
		this.whiteList = strings.Split(cfg.WhiteIp, ",")
	}
	if cfg.BlackIp != "" {
		this.blackList = strings.Split(cfg.BlackIp, ",")
	}
	return true
}

func (this *MemcacheServer) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}

	if this.config.Disable {
		return true
	}

	if ccr.Type == boot.CCR_NEED_START {
		logger.Debug(tag, "'%s' start listen (%s %s)", this.name, this.config.Net, this.config.Address)
		lis, err := net.Listen(this.config.Net, this.config.Address)
		if err != nil {
			logger.Warn(tag, "'%s' listen at (%s %s) fail - %s", this.name, this.config.Net, this.config.Address, err)
			return false
		}
		logger.Info(tag, "'%s' listen at (%s %s)", this.name, this.config.Net, this.config.Address)
		this.listener = lis
	}
	return true
}

func (this *MemcacheServer) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	if this.config.Disable {
		return true
	}
	if ccr.Type == boot.CCR_NEED_START {
		atomic.CompareAndSwapUint32(&this.state, 0, 1)
		go this.run(this.listener)
	}
	return true
}

func (this *MemcacheServer) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	if ccr.Type == boot.CCR_NEED_START {
		this.Close()
		for i := 0; i < 1000; i++ {
			if atomic.LoadUint32(&this.state) == 0 {
				break
			}
			// runing
			time.Sleep(1 * time.Millisecond)
		}
	}
	return true
}

func (this *MemcacheServer) Stop() bool {
	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}
	return true
}

func (this *MemcacheServer) Close() bool {
	return true
}

func (this *MemcacheServer) Cleanup() bool {
	return true
}
