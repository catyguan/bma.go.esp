package httpserver

import (
	"net"
	"net/http"
	"time"
)

const (
	tag = "httpServer"
)

type HttpServer struct {
	name     string
	config   *HttpServerConfigInfo
	timeout  time.Duration
	listener net.Listener
	Handler  http.Handler
}

func NewHttpServer(name string, h http.Handler) *HttpServer {
	this := new(HttpServer)
	this.name = name
	this.Handler = h
	return this
}

func (this *HttpServer) InitConfig(cfg *HttpServerConfigInfo) {
	this.config = cfg
	if cfg.TimeoutMS > 0 {
		this.timeout = time.Duration(cfg.TimeoutMS) * time.Millisecond
	} else {
		this.timeout = 1 * time.Minute
	}
}
