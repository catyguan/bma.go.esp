package thrift

import (
	"bmautil/valutil"
	"boot"
	"fmt"
	"logger"
	"strings"
)

type configInfo struct {
	Disable      bool
	MaxFrame     string
	MaxFrameSize uint64
	Address      string
	Port         int
	WhiteIp      string
	BlackIp      string
	WhiteList    []string
	BlackList    []string
}

func (this *configInfo) Valid() error {
	if this.Disable {
		return nil
	}
	if this.Address == "" {
		if this.Port > 0 {
			this.Address = fmt.Sprintf(":%d", this.Port)
		} else {
			return fmt.Errorf("port invalid")
		}
	}
	if this.MaxFrame != "" {
		mf, err := valutil.ToSize(this.MaxFrame, 1024, valutil.SizeB)
		if err != nil {
			return fmt.Errorf("MaxFrame(%s) invalid %s", this.MaxFrame, err)
		}
		this.MaxFrameSize = mf
	}
	if this.WhiteIp != "" {
		this.WhiteList = strings.Split(this.WhiteIp, ",")
	}
	if this.BlackIp != "" {
		this.BlackList = strings.Split(this.BlackIp, ",")
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.Disable != old.Disable {
		return boot.CCR_NEED_START
	}
	if this.MaxFrameSize != old.MaxFrameSize {
		return boot.CCR_NEED_START
	}
	if this.Address != old.Address {
		return boot.CCR_NEED_START
	}
	if this.WhiteIp != old.WhiteIp || this.BlackIp != old.BlackIp {
		return boot.CCR_NEED_START
	}
	return boot.CCR_NONE
}

func (this *ThriftServer) Name() string {
	return this.name
}

func (this *ThriftServer) Prepare() {
}
func (this *ThriftServer) CheckConfig(ctx *boot.BootContext) bool {
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

func (this *ThriftServer) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*configInfo)
	this.config = cfg
	return true
}

func (this *ThriftServer) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	if ccr.Type == boot.CCR_NEED_START {
		if this.config.Disable {
			logger.Info(tag, "thrift server(%s) disabled", this.name)
			return true
		}
		mf := this.config.MaxFrameSize
		transportFactory := NewTFramedTransportFactory(NewTTransportFactory(), func() uint64 {
			return mf
		})
		protocolFactory := NewTBinaryProtocolFactoryDefault()
		//protocolFactory := thrift.NewTCompactProtocolFactory()

		addr := this.config.Address
		serverTransport, err := NewTServerSocket(addr)
		if err != nil {
			logger.Error(tag, "thrift server(%s) start error - %s", this.name, err)
			return false
		}
		wl := this.config.WhiteList
		serverTransport.WhiteList = func() []string {
			return wl
		}
		bl := this.config.BlackList
		serverTransport.BlackList = func() []string {
			return bl
		}

		logger.Info(tag, "thrift server(%s) run - %s", this.name, addr)

		this.server = NewTSimpleServer4(this.processor, serverTransport, transportFactory, protocolFactory)
	}
	return true
}

func (this *ThriftServer) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}
	addr := this.config.Address
	go func() {
		defer func() {
			logger.Info(tag, "thrift server(%s) stop (%s)", this.name, addr)
		}()
		err := this.server.Serve()
		if err != nil {
			logger.Warn(tag, "thrift server(%s) run fail - %s", this.name, err.Error())
		}
	}()
	return true
}

func (this *ThriftServer) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}
	if this.server != nil {
		this.server.Stop()
		this.server = nil
	}
	return true
}

func (this *ThriftServer) Stop() bool {
	if this.server != nil {
		this.server.Stop()
		this.server = nil
	}
	return true
}

func (this *ThriftServer) Close() bool {
	return true
}

func (this *ThriftServer) Cleanup() bool {
	return true
}
