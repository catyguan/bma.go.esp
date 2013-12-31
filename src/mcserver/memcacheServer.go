package mcserver

import (
	"bmautil/netutil"
	"bmautil/syncutil"
	"bmautil/valutil"
	"boot"
	"bufio"
	"config"
	"errors"
	"fmt"
	"logger"
	"net"
	"runtime/debug"
	"strings"
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

type MemcacheServerHandler func(c *netutil.Channel, cmd *MemcacheCommand) (HandleCode, error)

type MemcacheServer struct {
	name       string
	network    string
	address    string
	handler    MemcacheServerHandler
	listener   net.Listener
	WhiteList  []string
	BlackList  []string
	group      *netutil.ConnGroup
	closeState *syncutil.CloseState
	disable    bool
}

type configInfo struct {
	Net     string
	Address string
	Port    int
	WhiteIp string
	BlackIp string
	Disable bool
}

func NewMemcacheServer(name string, h MemcacheServerHandler) *MemcacheServer {
	r := new(MemcacheServer)
	r.name = name
	r.handler = h
	r.closeState = syncutil.NewCloseState()
	return r
}

func (this *MemcacheServer) Name() string {
	return this.name
}

func (this *MemcacheServer) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		if cfg.Disable {
			this.disable = cfg.Disable
		} else {
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
		}
	} else {
		this.disable = true
	}
	if this.disable {
		logger.Info(tag, "%s disabled", this.name)
	}
	return true
}

func (this *MemcacheServer) Start() bool {
	if this.disable {
		return true
	}

	logger.Debug(tag, "start listen (%s %s)", this.network, this.address)
	lis, err := net.Listen(this.network, this.address)
	if err != nil {
		logger.Warn(tag, "listen at (%s %s) fail - %s", this.network, this.address, err)
		return false
	}
	logger.Info(tag, "listen at (%s %s)", this.network, this.address)
	this.listener = lis
	this.group = netutil.NewConnGroup(32)
	return true
}

func (this *MemcacheServer) Run() bool {
	if this.disable {
		return true
	}
	go this.run()
	return true
}

func (this *MemcacheServer) Stop() bool {
	if this.disable {
		return true
	}
	this.closeState.AskClose()
	if this.listener != nil {
		this.listener.Close()
	}
	return true
}

func (this *MemcacheServer) Close() bool {
	if this.disable {
		return true
	}
	if !this.group.RemoveAll(func() {
		this.closeState.DoneClose()
	}) {
		this.closeState.DoneClose()
	}
	this.group.Close()
	return true
}

func (this *MemcacheServer) Cleanup() bool {
	if this.disable {
		return true
	}
	if this.closeState.IsClosing() {
		this.closeState.WaitClosed()
	}
	return true
}

func (this *MemcacheServer) DefaultBoot() {
	boot.Define(boot.INIT, this.name, this.Init)
	boot.Define(boot.START, this.name, this.Start)
	boot.Define(boot.RUN, this.name, this.Run)
	boot.Define(boot.STOP, this.name, this.Stop)
	boot.Define(boot.CLOSE, this.name, this.Close)
	boot.Define(boot.CLEANUP, this.name, this.Cleanup)

	boot.Install(this.name, this)
}

func (this *MemcacheServer) run() {
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

func (this *MemcacheServer) accept(conn net.Conn) {
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
			if !this.closeState.IsClosing() && logger.EnableDebug(tag) {
				logger.Debug(tag, "read fail - (%s %s)", conn.RemoteAddr(), err)
			}
			return
		}
		str := strings.TrimSpace(line)
		logger.Debug(tag, "memcache command << %s", str)
		cmd, err2 := this.decode(ch, in, str)
		if err2 != nil {
			ch.Write([]byte("CLIENT_ERROR " + err2.Error() + "\r\n"))
			continue
		}
		code, err3 := this.handle(ch, cmd)
		if err3 != nil {
			ch.Write([]byte("SERVER_ERROR " + err3.Error() + "\r\n"))
			continue
		}
		if code == UNKNOW_COMMAND {
			ch.Write([]byte("ERROR\r\n"))
			continue
		}
		if code == CLOSE {
			return
		}
	}
}

func (this *MemcacheServer) handle(ch *netutil.Channel, cmd *MemcacheCommand) (hc HandleCode, rerr error) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Warn(tag, "handle request fail - %s\n%s", err, debug.Stack())
			smsg := logger.Sprintf("%s", err)
			rerr = errors.New(smsg)
		}
	}()
	if cmd == nil {
		return UNKNOW_COMMAND, nil
	}
	if cmd.Action == "quit" {
		return CLOSE, nil
	}
	return this.handler(ch, cmd)
}

func (this *MemcacheServer) decode(ch *netutil.Channel, in *bufio.Reader, str string) (*MemcacheCommand, error) {
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
