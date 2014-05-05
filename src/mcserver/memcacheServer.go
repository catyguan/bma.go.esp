package mcserver

import (
	"bmautil/netutil"
	"bmautil/valutil"
	"bufio"
	"fmt"
	"logger"
	"net"
	"runtime/debug"
	"strings"
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

type MemcacheCommand struct {
	Action string
	Params []string
	Data   []byte
}

type MemcacheServerHandler func(c net.Conn, cmd *MemcacheCommand) (HandleCode, error)

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
		logger.Info(tag, "'%s' stop (%s %s)", pnet, paddr)
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
			cg.Add(c)
			go this.accept(c, cg)
		} else {
			return
		}
	}
}

func (this *MemcacheServer) accept(conn net.Conn, cg *netutil.ConnGroup) {
	defer func() {
		cg.Remove(conn)
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "'%s' connection close - %s", this.name, conn.RemoteAddr())
		}
	}()
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "'%s' connection accept - %s", this.name, conn.RemoteAddr())
	}
	in := bufio.NewReader(conn)
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "'%s' read fail - %s - %s", this.name, conn.RemoteAddr(), err)
			}
			return
		}
		str := strings.TrimSpace(line)
		logger.Debug(tag, "memcache command << %s", str)
		cmd, err2 := this.decode(conn, in, str)
		if err2 != nil {
			conn.Write([]byte("CLIENT_ERROR " + err2.Error() + "\r\n"))
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
	return this.handler(c, cmd)
}

func (this *MemcacheServer) decode(c net.Conn, in *bufio.Reader, str string) (*MemcacheCommand, error) {
	cmd := new(MemcacheCommand)
	strlist := strings.Split(str, " ")
	var w string
	if len(strlist) == 0 {
		return nil, nil
	}
	w = strlist[0]
	cmd.Action = w
	cmd.Params = strlist[1:]
	switch w {
	case "set", "add", "replace":
		if len(strlist) != 5 {
			return nil, nil
		}
		sz := valutil.ToInt(strlist[4], 0)

		b := make([]byte, sz)
		var err error
		for i := 0; i < sz; i++ {
			b[i], err = in.ReadByte()
			if err != nil {
				return nil, err
			}
		}
		in.ReadByte()
		in.ReadByte()
		cmd.Data = b
	}
	return cmd, nil
}
