package httpserver

import (
	"bmautil/netutil"
	"fmt"
	"logger"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	tag = "httpServer"
)

type HttpServer struct {
	name      string
	config    *HttpServerConfigInfo
	timeout   time.Duration
	listener  net.Listener
	Handler   http.Handler
	WhiteList []string
	BlackList []string
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
	if cfg.WhiteIp != "" {
		this.WhiteList = strings.Split(cfg.WhiteIp, ",")
	}
	if cfg.BlackIp != "" {
		this.BlackList = strings.Split(cfg.BlackIp, ",")
	}
}

func (this *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if len(this.WhiteList) > 0 || len(this.BlackList) > 0 {
		if ok, msg := netutil.IpAccept(req.RemoteAddr, this.WhiteList, this.BlackList, true); !ok {
			logger.Info(tag, msg)
			http.Error(w, fmt.Sprintf("IP(%s) FORBIDDEN", req.RemoteAddr), 401)
			return
		}
	}
	this.Handler.ServeHTTP(w, req)
}
