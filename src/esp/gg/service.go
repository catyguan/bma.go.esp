package gg

import (
	"bmautil/qexec"
	"bmautil/syncutil"
	"boot"
	"config"
	"esp/espnet"
	"logger"
	"net"
	"runtime/debug"
	"strings"
)

const (
	tag = "gg"
)

type GGroupService struct {
	// init
	name string

	// config
	config *configInfo

	// runtime
	executor   *qexec.QueueExecutor
	closeState *syncutil.CloseState
}

type configInfo struct {
	espnet.ListenConfig
	AcceptTimeout int // seconds
}

func NewGGroupService(name string) *GGroupService {
	this := new(GGroupService)
	this.name = name
	this.executor = qexec.NewQueueExecutor(tag, 32, this.requestHandler)
	this.closeState = syncutil.NewCloseState()
	return this
}

func (this *GGroupService) Name() string {
	return this.name
}

func (this *GGroupService) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
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
		if cfg.AcceptTimeout > 0 {
			this.acceptTimeout = cfg.AcceptTimeout
		} else {
			this.acceptTimeout = 5
		}
		return true
	}
	logger.Error(tag, "GetBeanConfig(%s) fail", this.name)
	return false
}

func (this *GGroupService) Start() bool {
	logger.Debug(tag, "start listen (%s %s)", this.network, this.address)
	lis, err := net.Listen(this.network, this.address)
	if err != nil {
		logger.Warn(tag, "listen at (%s %s) fail - %s", this.network, this.address, err)
		return false
	}
	logger.Info(tag, "listen at (%s %s)", this.network, this.address)
	this.listener = lis
	this.group = netutil.NewConnGroup(8)
	return true
}

func (this *GGroupService) run() {
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

func (this *GGroupService) Run() bool {
	go this.run()
	return true
}

func (this *GGroupService) accept(conn net.Conn) {
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
	gch := newGChannel(ch.RemoteAddr().String(), ch)
	gch.passive = true
	this.process(gch)
}

func (this *GGroupService) Stop() bool {
	this.closeState.AskClose()
	if this.listener != nil {
		this.listener.Close()
	}
	return true
}

func (this *GGroupService) Close() bool {
	if !this.group.RemoveAll(func() {
		this.closeState.DoneClose()
	}) {
		this.closeState.DoneClose()
	}
	this.group.Close()
	return true
}

func (this *GGroupService) Cleanup() bool {
	if this.closeState.IsClosing() {
		this.closeState.WaitClosed()
	}
	return true
}

func (this *GGroupService) DefaultBoot(doInit bool) {
	if doInit {
		if this.name == "" {
			panic("GGroupService name not set")
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
