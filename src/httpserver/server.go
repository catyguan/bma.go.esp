package httpserver

import (
	"boot"
	"config"
	"logger"
	"net"
	"net/http"
	"time"
)

const (
	tag = "httpServer"
)

type HttpServer struct {
	name     string
	address  string
	timeout  time.Duration
	listener net.Listener
	Handler  http.Handler
}

type HttpServerConfigInfo struct {
	Address   string
	Port      int
	TimeoutMS int
}

func NewHttpServer(name string, h http.Handler) *HttpServer {
	this := new(HttpServer)
	this.name = name
	this.Handler = h
	return this
}

func (this *HttpServer) Name() string {
	return this.name
}

func (this *HttpServer) Init() bool {
	cfg := HttpServerConfigInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		return this.InitConfig(&cfg)
	}
	logger.Error(tag, "GetBeanConfig(%s) fail", this.name)
	return false
}

func (this *HttpServer) InitConfig(cfg *HttpServerConfigInfo) bool {
	if cfg.Address != "" {
		this.address = cfg.Address
	} else {
		if cfg.Port > 0 {
			this.address = logger.Sprintf(":%d", cfg.Port)
		} else {
			logger.Error(tag, "config '%s' port invalid", this.name)
			return false
		}
	}
	if cfg.TimeoutMS > 0 {
		this.timeout = time.Duration(cfg.TimeoutMS) * time.Millisecond
	} else {
		this.timeout = 1 * time.Minute
	}
	return true
}

func (this *HttpServer) Start() bool {
	logger.Debug(tag, "%s start listen (%s)", this.name, this.address)
	lis, err := net.Listen("tcp", this.address)
	if err != nil {
		logger.Warn(tag, "%s listen at (%s) fail - %s", this.name, this.address, err)
		return false
	}
	logger.Info(tag, "%s listen at (%s)", this.name, this.address)
	this.listener = lis
	return true
}

func (this *HttpServer) Run() bool {
	defer func() {
		logger.Info(tag, "%s stop (%s)", this.name, this.address)
	}()
	s := &http.Server{
		Addr:         this.address,
		ReadTimeout:  this.timeout,
		WriteTimeout: this.timeout,
		Handler:      this.Handler,
	}
	go s.Serve(this.listener)
	return true
}

func (this *HttpServer) Stop() bool {
	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}
	return true
}

func (this *HttpServer) DefaultBoot(doInit bool) {
	if doInit {
		boot.Define(boot.INIT, this.name, this.Init)
	}
	boot.Define(boot.START, this.name, this.Start)
	boot.Define(boot.RUN, this.name, this.Run)
	boot.Define(boot.STOP, this.name, this.Stop)
}
