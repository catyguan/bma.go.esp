package netserver

import (
	"logger"
	"net"
	"runtime/debug"
)

const (
	tag = "netServer"
)

type ConnHandler func(c net.Conn)

type Service struct {
	name     string
	config   *configInfo
	handler  ConnHandler
	listener net.Listener
}

func NewService(name string, h ConnHandler) *Service {
	r := new(Service)
	r.name = name
	r.handler = h
	return r
}

func (this *Service) Listener() net.Listener {
	return this.listener
}

func (this *Service) accept(conn net.Conn) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Debug(tag, "handle request fail - %s\n%s", err, debug.Stack())
		}
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "connection close - %s", conn.RemoteAddr())
		}
		conn.Close()
	}()
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "connection accept - %s", conn.RemoteAddr())
	}
	this.handler(conn)
}
