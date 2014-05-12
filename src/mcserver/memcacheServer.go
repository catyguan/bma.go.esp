package mcserver

import (
	"bmautil/netutil"
	"fmt"
	"logger"
	"net"
	"runtime/debug"
	"sync/atomic"
)

const (
	tag = "memcacheServer"
)

type HandleCode int

const (
	DONE           = HandleCode(0)
	UNKNOW_COMMAND = HandleCode(1)
	CLOSE          = HandleCode(2)
)

type MemcacheServerHandler interface {
	HandleMemcacheCommand(c net.Conn, cmd *MemcacheCommand) (HandleCode, error)
	OnMemcacheConnOpen(c net.Conn) bool
	OnMemcacheConnClose(c net.Conn)
}

type MemcacheServer struct {
	name    string
	handler MemcacheServerHandler

	config    *configInfo
	whiteList []string
	blackList []string

	listener net.Listener
	state    uint32
}

func NewMemcacheServer(name string, h MemcacheServerHandler) *MemcacheServer {
	r := new(MemcacheServer)
	r.name = name
	r.handler = h
	return r
}

func (this *MemcacheServer) run(lis net.Listener) {
	pnet := this.config.Net
	paddr := this.config.Address
	cg := netutil.NewConnGroup()
	defer func() {
		cg.CloseAll()
		logger.Info(tag, "'%s' stop (%s %s)", this.name, pnet, paddr)
		atomic.CompareAndSwapUint32(&this.state, 1, 0)
	}()
	for {
		c, err := lis.Accept()
		if err == nil {
			addr := c.RemoteAddr().String()
			if ok, msg := netutil.IpAccept(addr, this.whiteList, this.blackList, true); !ok {
				logger.Warn(tag, "unaccept(%s) address %s", msg, addr)
				c.Close()
				continue
			}
			if !this.handler.OnMemcacheConnOpen(c) {
				logger.Warn(tag, "handler unaccept address %s", addr)
				c.Close()
				continue
			}
			cg.Add(c)
			go this.accept(c, cg)
		} else {
			return
		}
	}
}

func (this *MemcacheServer) accept(conn net.Conn, cg *netutil.ConnGroup) {
	defer func() {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "'%s' connection close - %s", this.name, conn.RemoteAddr())
		}
		cg.Remove(conn)
		this.handler.OnMemcacheConnClose(conn)
	}()
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "'%s' connection accept - %s", this.name, conn.RemoteAddr())
	}
	coder := NewMemcacheCoder()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "'%s' read fail - %s - %s", this.name, conn.RemoteAddr(), err)
			}
			return
		}
		coder.Write(buf[:n])
		ok, cmd := coder.DecodeCommand()
		if !ok {
			continue
		}
		code, err3 := this.handle(conn, cmd)
		if err3 != nil {
			conn.Write([]byte("SERVER_ERROR " + err3.Error() + "\r\n"))
			continue
		}
		if code == UNKNOW_COMMAND {
			conn.Write([]byte("ERROR\r\n"))
			continue
		}
		if code == CLOSE {
			conn.Close()
			return
		}
	}
}

func (this *MemcacheServer) handle(c net.Conn, cmd *MemcacheCommand) (hc HandleCode, rerr error) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Warn(tag, "'%s' handle request fail - %s\n%s", this.name, err, debug.Stack())
			rerr = fmt.Errorf("%s", err)
		}
	}()
	if cmd == nil {
		return UNKNOW_COMMAND, nil
	}
	if cmd.Action == "quit" {
		return CLOSE, nil
	}
	return this.handler.HandleMemcacheCommand(c, cmd)
}
