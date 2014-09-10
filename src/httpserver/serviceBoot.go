package httpserver

import (
	"boot"
	"fmt"
	"logger"
	"net"
	"net/http"
)

type HttpServerConfigInfo struct {
	Address   string
	Port      int
	TimeoutMS int
	WhiteIp   string
	BlackIp   string
}

func (this *HttpServerConfigInfo) Valid() error {
	if this.Address == "" {
		if this.Port == 0 {
			this.Port = 80
		}
		if this.Port > 0 {
			this.Address = logger.Sprintf(":%d", this.Port)
		} else {
			return fmt.Errorf("port invalid")
		}
	}
	return nil
}

func (this *HttpServerConfigInfo) Compare(old *HttpServerConfigInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.Address != old.Address {
		return boot.CCR_NEED_START
	}
	if this.TimeoutMS != old.TimeoutMS {
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

func (this *HttpServer) Name() string {
	return this.name
}

func (this *HttpServer) Prepare() {
}
func (this *HttpServer) CheckConfig(ctx *boot.BootContext) bool {
	co := ctx.Config
	cfg := new(HttpServerConfigInfo)
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

func (this *HttpServer) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*HttpServerConfigInfo)
	this.InitConfig(cfg)
	return true
}

func (this *HttpServer) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NEED_START {
		logger.Debug(tag, "%s start listen (%s)", this.name, this.config.Address)
		lis, err := net.Listen("tcp", this.config.Address)
		if err != nil {
			logger.Warn(tag, "%s listen at (%s) fail - %s", this.name, this.config.Address, err)
			return false
		}
		logger.Info(tag, "%s listen at (%s)", this.name, this.config.Address)
		this.listener = lis
	}
	return true
}

func (this *HttpServer) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}

	if this.ownMux {
		logger.Debug(tag, "'%s' build http mux", this.name)
		bl := this.muxBuilders
		mux := http.NewServeMux()
		for _, b := range bl {
			b(mux)
		}
		this.Handler = mux
	}

	addr := this.config.Address
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  this.timeout,
		WriteTimeout: this.timeout,
		Handler:      this,
	}
	go func() {
		defer func() {
			logger.Info(tag, "%s stop (%s)", this.name, addr)
		}()
		s.Serve(this.listener)
	}()
	return true
}

func (this *HttpServer) GraceStop(ctx *boot.BootContext) bool {
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

func (this *HttpServer) Stop() bool {
	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}
	return true
}

func (this *HttpServer) Close() bool {
	return true
}

func (this *HttpServer) Cleanup() bool {
	return true
}
