package servproxy

import (
	"bmautil/netutil"
	"logger"
	"net"
	"runtime"
	"sync/atomic"
)

type PortObj struct {
	s        *Service
	name     string
	handler  PortHandler
	cfg      *PortConfigInfo
	listener net.Listener
	closed   uint32
}

func NewPortObj(s *Service, n string, cfg *PortConfigInfo, h PortHandler) *PortObj {
	r := new(PortObj)
	r.s = s
	r.name = n
	r.handler = h
	r.cfg = cfg
	return r
}

func (this *PortObj) Config() *PortConfigInfo {
	return this.cfg
}

func (this *PortObj) IsClose() bool {
	return atomic.LoadUint32(&this.closed) == 1
}

func (this *PortObj) Start() error {
	err := this.cfg.Valid()
	if err != nil {
		return err
	}
	lis, err1 := net.Listen("tcp", this.cfg.Net)
	if err1 != nil {
		return err1
	}
	logger.Info(tag, "%s listen at %s", this.name, this.cfg.Net)
	this.listener = lis

	go func() {
		defer func() {
			logger.Info(tag, "%s stop %s", this.name, this.cfg.Net)
		}()
		for {
			c, err := this.listener.Accept()
			if err == nil {
				addr := c.RemoteAddr().String()
				if ok, msg := netutil.IpAccept(addr, this.cfg.whiteList, this.cfg.blackList, true); !ok {
					logger.Warn(tag, "unaccept(%s) address %s", msg, addr)
					c.Close()
					continue
				}
				go this.accept(c) // new connect
			} else {
				if this.IsClose() {
					logger.Debug(tag, "%s closing, exit", this.name)
					return
				}
				logger.Warn(tag, "accept fail - %s", err)
			}
		}
	}()
	return nil
}

func (this *PortObj) Stop() {
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		if this.listener != nil {
			this.listener.Close()
			this.listener = nil
		}
	}
}

func (this *PortObj) accept(conn net.Conn) {
	defer func() {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "connection close - %s", conn.RemoteAddr())
		}
		err := recover()
		if err != nil {
			trace := make([]byte, 1024)
			runtime.Stack(trace, true)
			logger.Warn(tag, "process panic: %v\n%s", err, trace)
		}
	}()
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "connection accept - %s", conn.RemoteAddr())
	}
	this.handler.Handle(this.s, this, conn)
}
