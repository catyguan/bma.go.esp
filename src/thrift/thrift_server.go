package thrift

import (
	"bmautil/valutil"
	"boot"
	"config"
	"fmt"
	"logger"
	"strings"
)

type ThriftServer struct {
	// init
	name      string
	processor TProcessor

	// config
	disable bool
	address string

	Maxframe  uint64
	WhiteList []string
	BlackList []string

	// runtime
	server *TSimpleServer
}

type configInfo struct {
	Disable  bool
	MaxFrame string
	Address  string
	Port     int
	WhiteIp  string
	BlackIp  string
}

func NewThriftServer(name string, p TProcessor) *ThriftServer {
	r := new(ThriftServer)
	r.name = name
	r.processor = p
	return r
}

func (this *ThriftServer) Name() string {
	return this.name
}

func (this *ThriftServer) Init() bool {
	cfg := configInfo{}
	co := config.Global
	if co.GetBeanConfig(this.name, &cfg) {
		if cfg.Disable {
			this.disable = cfg.Disable
			logger.Debug(tag, "'%s' disable", this.name)
			return true
		}
		if cfg.Address != "" {
			this.address = cfg.Address
		} else {
			if cfg.Port > 0 {
				this.address = fmt.Sprintf(":%d", cfg.Port)
			} else {
				logger.Error(tag, "config '%s' port invalid", this.name)
				return false
			}
		}
		if cfg.WhiteIp != "" {
			this.WhiteList = strings.Split(cfg.WhiteIp, ",")
		}
		if cfg.BlackIp != "" {
			this.BlackList = strings.Split(cfg.BlackIp, ",")
		}
		if cfg.MaxFrame != "" {
			mf, err := valutil.ToSize(cfg.MaxFrame, 1024, valutil.SizeB)
			if err != nil {
				logger.Error(tag, "config '%s' MaxFrame invalid - %s", this.name, err.Error())
				return false
			}
			this.Maxframe = mf
		}
	} else {
		this.disable = true
	}
	return true
}

func (this *ThriftServer) Start() bool {
	if !this.disable {
		transportFactory := NewTFramedTransportFactory(NewTTransportFactory(), func() uint64 {
			return this.Maxframe
		})
		protocolFactory := NewTBinaryProtocolFactoryDefault()
		//protocolFactory := thrift.NewTCompactProtocolFactory()

		serverTransport, err := NewTServerSocket(this.address)
		if err != nil {
			logger.Error(tag, "start error - %s", err)
			return false
		}
		serverTransport.WhiteList = func() []string {
			return this.WhiteList
		}
		serverTransport.BlackList = func() []string {
			return this.BlackList
		}

		logger.Info(tag, "thrift server run - %s", this.address)

		this.server = NewTSimpleServer4(this.processor, serverTransport, transportFactory, protocolFactory)
	}
	return true
}

func (this *ThriftServer) run() {
	defer func() {
		logger.Info(tag, "stop (%s)", this.address)
	}()
	err := this.server.Serve()
	if err != nil {
		logger.Warn(tag, "run fail - %s", err.Error())
	}
}

func (this *ThriftServer) Run(backgroup bool) bool {
	if this.server != nil {
		if backgroup {
			go this.run()
		} else {
			this.run()

		}
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

func (this *ThriftServer) DefaultBoot(doInit bool) {
	if this.name == "" {
		panic("ThriftServer name not set")
	}
	boot.Define(boot.INIT, this.name, this.Init)

	boot.Define(boot.START, this.name, this.Start)
	boot.Define(boot.RUN, this.name, func() bool {
		this.Run(true)
		return true
	})
	boot.Define(boot.STOP, this.name, this.Stop)
}
