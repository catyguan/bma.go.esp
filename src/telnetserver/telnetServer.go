package telnetserver

import (
	"bmautil/netutil"
	"bmautil/syncutil"
	"boot"
	"bufio"
	"config"
	"fmt"
	"logger"
	"net"
	"runtime/debug"
	"strings"
)

const (
	tag = "telnetServer"
)

type TelnetServerHandler func(c *netutil.Channel, msg string) bool

type TelnetServer struct {
	name       string
	network    string
	address    string
	handler    TelnetServerHandler
	listener   net.Listener
	WhiteList  []string
	BlackList  []string
	group      *netutil.ConnGroup
	closeState *syncutil.CloseState
}

type configInfo struct {
	Net     string
	Address string
	Port    int
	WhiteIp string
	BlackIp string
}

func NewTelnetServer(name string, h TelnetServerHandler) *TelnetServer {
	r := new(TelnetServer)
	r.name = name
	r.handler = h
	r.closeState = syncutil.NewCloseState()
	return r
}

func (this *TelnetServer) Name() string {
	return this.name
}

func (this *TelnetServer) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
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
		if cfg.Net != "" {
			this.network = cfg.Net
		} else {
			this.network = "tcp"
		}
		if cfg.WhiteIp != "" {
			this.WhiteList = strings.Split(cfg.WhiteIp, ",")
		}
		if cfg.BlackIp != "" {
			this.BlackList = strings.Split(cfg.BlackIp, ",")
		}
		return true
	}
	logger.Error(tag, "GetBeanConfig(%s) fail", this.name)
	return false
}

func (this *TelnetServer) InitServer(nw string, laddress string, h TelnetServerHandler) bool {
	this.network = nw
	this.address = laddress
	this.handler = h
	return true
}

func (this *TelnetServer) Start() bool {
	logger.Debug(tag, "start listen (%s %s)", this.network, this.address)
	lis, err := net.Listen(this.network, this.address)
	if err != nil {
		logger.Warn(tag, "listen at (%s %s) fail - %s", this.network, this.address, err)
		return false
	}
	logger.Info(tag, "listen at (%s %s)", this.network, this.address)
	this.listener = lis
	this.group = netutil.NewConnGroup(10)
	return true
}

func (this *TelnetServer) run() {
	defer func() {
		logger.Info(tag, "stop (%s %s)", this.network, this.address)
	}()
	for {
		c, err := this.listener.Accept()
		if err == nil {
			addr := c.RemoteAddr().String()
			if ok, msg := netutil.IpAccept(addr, this.WhiteList, this.BlackList, true); !ok {
				logger.Warn(tag, "unaccept(%s) address %s", msg, addr)
				c.Close()
				continue
			}
			this.group.Add(c, nil)
			go this.accept(c) // new connect
		} else {
			if this.closeState.IsClosing() {
				logger.Debug(tag, "closing, exit")
				return
			}
			logger.Debug(tag, "accept fail - ", err)
		}
	}
}

func (this *TelnetServer) Run(backgroup bool) bool {
	if backgroup {
		go this.run()
		return true
	} else {
		this.run()
		return true
	}
}

func (this *TelnetServer) handle(ch *netutil.Channel, msg string) bool {
	defer func() {
		err := recover()
		if err != nil {
			logger.Debug(tag, "handle request fail - %s\n%s", err, debug.Stack())
		}
	}()
	return this.handler(ch, msg)
}

func (this *TelnetServer) accept(conn net.Conn) {
	ch := netutil.NewChannel(conn)
	defer func() {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "connection close - %s", conn.RemoteAddr())
		}
		ch.CloseChannel()
		this.group.Remove(conn, nil)

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
		if !this.handle(ch, strings.TrimSpace(line)) {
			return
		}
	}
}

func (this *TelnetServer) Stop() bool {
	this.closeState.AskClose()
	if this.listener != nil {
		this.listener.Close()
	}
	return true
}

func (this *TelnetServer) Close() bool {
	if !this.group.RemoveAll(func() {
		this.closeState.DoneClose()
	}) {
		this.closeState.DoneClose()
	}
	this.group.Close()
	return true
}

func (this *TelnetServer) Cleanup() bool {
	if this.closeState.IsClosing() {
		this.closeState.WaitClosed()
	}
	return true
}

func (this *TelnetServer) DefaultBoot(doInit bool) {
	if doInit {
		if this.name == "" {
			panic("TelnetServer name not set")
		}
		boot.Define(boot.INIT, this.name, this.Init)
	}
	boot.Define(boot.START, this.name, this.Start)
	boot.Define(boot.RUN, this.name, func() bool {
		this.Run(true)
		return true
	})
	boot.Define(boot.STOP, this.name, this.Stop)
	boot.Define(boot.CLOSE, this.name, this.Close)
	boot.Define(boot.CLEANUP, this.name, this.Cleanup)

	boot.Install(this.name, this)
}
