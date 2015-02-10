package socket

import (
	"bmautil/netutil"
	"bmautil/syncutil"
	"boot"
	"fmt"
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

func (this *ListenConfig) Valid() error {
	if this.Address == "" {
		if this.Port > 0 {
			this.Address = fmt.Sprintf(":%d", this.Port)
		} else {
			return fmt.Errorf("port invalid")
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

func (this *ListenConfig) Compare(old *ListenConfig) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.Net != old.Net || this.Address != old.Address {
		return boot.CCR_NEED_START
	}
	if this.WhiteIp != old.WhiteIp || this.BlackIp != old.BlackIp {
		return boot.CCR_CHANGE
	}
	return boot.CCR_NONE
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

func (this *ListenPoint) Prepare() {
}
func (this *ListenPoint) CheckConfig(ctx *boot.BootContext) bool {
	if this.initConfig != nil {
		return true
	}
	co := ctx.Config
	cfg := new(ListenConfig)
	if !co.GetBeanConfig(this.name, cfg) {
		logger.Error(tag, "'%s' miss config", this.name)
		return false
	}
	if err := cfg.Valid(); err != nil {
		logger.Error(tag, "'%s' config error - %s", this.name, err)
		return false
	}
	ccr := boot.NewConfigCheckResult(cfg.Compare(this.config), cfg)
	ctx.CheckFlag = ccr
	return true
}

func (this *ListenPoint) Init(ctx *boot.BootContext) bool {
	if this.initConfig != nil {
		this.config = this.initConfig
		return true
	}
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	this.config = ccr.Config.(*ListenConfig)
	return true
}

func (this *ListenPoint) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}

	cfg := this.config
	if this.IsClosing() {
		return false
	}

	if ccr.Type == boot.CCR_NEED_START {
		logger.Debug(tag, "%s start listen (%s %s)", this, cfg.Net, cfg.Address)
		lis, err := net.Listen(cfg.Net, cfg.Address)
		if err != nil {
			logger.Warn(tag, "%s listen at (%s %s) fail - %s", this, cfg.Net, cfg.Address, err)
			return false
		}
		this.listener = lis
		logger.Info(tag, "%s listen at (%s %s)", this, cfg.Net, cfg.Address)
	}

	this.blackList = cfg.GetBlackList()
	this.whiteList = cfg.GetWhiteList()

	return true
}

func (this *ListenPoint) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}

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

func (this *ListenPoint) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}

	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}
	return true
}

func (this *ListenPoint) Stop() bool {
	this.AskClose()
	return true
}

func (this *ListenPoint) Close() bool {
	return true
}

func (this *ListenPoint) Cleanup() bool {
	this.WaitClose()
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
