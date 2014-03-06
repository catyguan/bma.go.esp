package socket

import (
	"bmautil/netutil"
	"bmautil/syncutil"
	"config"
	"logger"
	"net"
	"strings"
)

// Listener
type ListenConfig struct {
	Net     string
	Address string
	Port    int
	WhiteIp string
	BlackIp string
}

func (this *ListenConfig) Valid(name string) error {
	if this.Address == "" {
		if this.Port > 0 {
			this.Address = logger.Sprintf(":%d", this.Port)
		} else {
			return logger.Warn(tag, "ListenConfig '%s' port invalid", name)
		}
	}
	if this.Net == "" {
		this.Net = "tcp"
	}
	return nil
}

func (this *ListenConfig) GetWhiteList() []string {
	if this.WhiteIp != "" {
		return strings.Split(this.WhiteIp, ",")
	}
	return nil
}

func (this *ListenConfig) GetBlackList() []string {
	if this.BlackIp != "" {
		return strings.Split(this.BlackIp, ",")
	}
	return nil
}

// ListenPoint
type ListenPoint struct {
	name       string
	initConfig *ListenConfig
	config     *ListenConfig
	socketInit SocketInit

	listener   net.Listener
	closeState syncutil.CloseState

	// config
	whiteList []string
	blackList []string

	acceptor SocketAcceptor
}

func NewListenPoint(name string, cfg *ListenConfig, sinit SocketInit) *ListenPoint {
	this := new(ListenPoint)
	this.name = name
	this.socketInit = sinit
	this.initConfig = cfg
	this.closeState.InitCloseState()
	return this
}

func (this *ListenPoint) Name() string {
	return this.name
}

func (this *ListenPoint) String() string {
	return "ListenPoint[" + this.name + "]"
}

func (this *ListenPoint) GetListener() net.Listener {
	return this.listener
}

func (this *ListenPoint) Init() bool {
	if this.initConfig != nil {
		this.config = this.initConfig
		return true
	}
	cfg := ListenConfig{}
	if config.GetBeanConfig(this.name, &cfg) {
		this.config = &cfg
	}
	return true
}

func (this *ListenPoint) Start() bool {
	cfg := this.config
	err := cfg.Valid(this.name)
	if err != nil {
		logger.Warn(tag, "%s config invalid %s", this, err)
		return false
	}

	if this.IsClosing() {
		return false
	}

	logger.Debug(tag, "%s start listen (%s %s)", this, cfg.Net, cfg.Address)
	lis, err := net.Listen(cfg.Net, cfg.Address)
	if err != nil {
		logger.Warn(tag, "%s listen at (%s %s) fail - %s", this, cfg.Net, cfg.Address, err)
		return false
	}
	logger.Info(tag, "%s listen at (%s %s)", this, cfg.Net, cfg.Address)

	this.listener = lis
	this.blackList = cfg.GetBlackList()
	this.whiteList = cfg.GetWhiteList()

	return true
}

func (this *ListenPoint) Run() bool {
	lis := this.listener
	go func() {
		defer func() {
			logger.Info(tag, "%s stop (%s %s)", this, this.config.Net, this.config.Address)
			if this.closeState.IsClosing() {
				this.closeState.DoneClose()
			}
		}()
		logger.Info(tag, "%s run (%s %s)", this, this.config.Net, this.config.Address)
		for {
			c, err := lis.Accept()
			if err == nil {
				addr := c.RemoteAddr().String()
				if ok, msg := netutil.IpAccept(addr, this.whiteList, this.blackList, true); !ok {
					logger.Warn(tag, "unaccept(%s) address %s", msg, addr)
					c.Close()
					continue
				}
				sock := NewSocket(c, 32, 0)
				if err := sock.Start(this.socketInit); err != nil {
					logger.Debug(tag, "Socket[%s] start fail", sock)
					return
				}
				if this.acceptor != nil {
					if err := this.acceptor(sock); err != nil {
						logger.Debug(tag, "Socket[%s] acceptor fail - %s", sock, err)
						sock.Close()
						return
					}
				}
			} else {
				return
			}
		}
	}()
	return true
}

func (this *ListenPoint) Stop() bool {
	this.AskClose()
	return true
}

func (this *ListenPoint) Cleanup() bool {
	this.WaitClose()
	return true
}

func (this *ListenPoint) GraceStop() bool {
	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}
	return true
}

func (this *ListenPoint) AskClose() {
	if this.closeState.AskClose() {
		if this.listener != nil {
			this.listener.Close()
			this.listener = nil
		} else {
			this.closeState.DoneClose()
		}
	}
}

func (this *ListenPoint) IsClosing() bool {
	return this.closeState.IsClosing()
}

func (this *ListenPoint) WaitClose() {
	this.closeState.WaitClosed()
}

// socketServer
func (this *ListenPoint) SetAcceptor(sa SocketAcceptor) {
	this.acceptor = sa
}
