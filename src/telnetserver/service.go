package telnetserver

import (
	"bmautil/syncutil"
	"bufio"
	"logger"
	"net"
	"runtime/debug"
	"strings"
)

const (
	tag = "telnetServer"
)

type ServiceHandler func(c net.Conn, msg string) bool

type Service struct {
	name       string
	config     *configInfo
	handler    ServiceHandler
	listener   net.Listener
	closeState *syncutil.CloseState
}

func NewService(name string, h ServiceHandler) *Service {
	r := new(Service)
	r.name = name
	r.handler = h
	r.closeState = syncutil.NewCloseState()
	return r
}

func (this *Service) handle(c net.Conn, msg string) bool {
	defer func() {
		err := recover()
		if err != nil {
			logger.Debug(tag, "handle request fail - %s\n%s", err, debug.Stack())
		}
	}()
	return this.handler(c, msg)
}

func (this *Service) accept(conn net.Conn) {
	defer func() {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "connection close - %s", conn.RemoteAddr())
		}
		// ch.CloseChannel()
		// this.group.Remove(conn, nil)

		err := recover()
		if err != nil {
			logger.Debug(tag, "process fail - %s", err)
		}
	}()
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "connection accept - %s", conn.RemoteAddr())
	}
	in := bufio.NewReader(conn)
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "read fail - (%s %s)", conn.RemoteAddr(), err)
			}
			return
		}
		if !this.handle(conn, strings.TrimSpace(line)) {
			return
		}
	}
}
