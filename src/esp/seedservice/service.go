package seedservice

import (
	"bmautil/socket"
	"bmautil/valutil"
	"boot"
	"config"
	"esp/espnet"
	"esp/espnet/protpack"
	"logger"
	"time"
)

const (
	tag = "seedService"
)

type SeedService struct {
	name string
	node *SeedNode

	// config
	maxPackage int
	cfginfo    *configInfo
}

func NewSeedService(name string) *SeedService {
	this := new(SeedService)
	this.name = name
	return this
}

func (this *SeedService) Name() string {
	return this.name
}

type configInfo struct {
	AcceptTimeoutMS int
	SocketTrace     int
	MaxPackage      string
	Point           espnet.ListenConfig
}

func (this *SeedService) Init() bool {

	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		if err := cfg.Point.Valid(this.name); err != nil {
			logger.Error(tag, "%s Point invalid - %s", this.name, err)
			return false
		}
		if cfg.MaxPackage != "" {
			mp, err := valutil.ToSize(cfg.MaxPackage, 1024, valutil.SizeB)
			if err != nil {
				logger.Error(tag, "%s MaxPackage invalid %s", this.name, err)
				return false
			}
			this.maxPackage = int(mp)
		}
		if this.maxPackage == 0 {
			this.maxPackage = 1012 * 1024
		}
		this.cfginfo = &cfg
	} else {
		logger.Error(tag, "%s config not exists", this.name)
		return false
	}

	return true
}

func (this *SeedService) Start() bool {
	ns := espnet.InstanceOfNet()
	if _, err := ns.CreatePoint(this.name, &this.cfginfo.Point, this.channelAccept, this.socketInit); err != nil {
		logger.Error(tag, "create point fail - %s", err)
		return false
	}

	this.node = NewSeedNode(this.name)
	if !this.node.Run() {
		return false
	}

	return true
}

func (this *SeedService) socketInit(sock *socket.Socket) error {
	if this.cfginfo.AcceptTimeoutMS > 0 {
		sock.Timeout = time.Duration(this.cfginfo.AcceptTimeoutMS) * time.Millisecond
	}
	sock.Trace = this.cfginfo.SocketTrace
	return nil
}

func (this *SeedService) channelAccept(ch espnet.Channel) error {
	h := protpack.PackageEncode(0)
	if true {
		hl := h.NewDecoder(1024)
		if err := h.BindDecode(ch, hl); err != nil {
			return err
		}
	}

	if true {
		if err := h.BindEncode(ch); err != nil {
			return err
		}
	}

	espnet.Connect(ch, this.node)

	return nil
}

func (this *SeedService) Run() bool {
	return true
}

func (this *SeedService) Stop() bool {
	ns := espnet.InstanceOfNet()
	ns.ClosePoint(this.name)

	return true
}

func (this *SeedService) Close() bool {
	if this.node != nil {
		this.node.AskClose()
	}
	return true
}

func (this *SeedService) Cleanup() bool {
	if this.node != nil {
		this.node.WaitStop()
	}
	return true
}

func (this *SeedService) DefaultBoot() {
	boot.Define(boot.INIT, this.name, this.Init)

	boot.Define(boot.START, this.name, this.Start)
	boot.Define(boot.RUN, this.name, this.Run)
	boot.Define(boot.STOP, this.name, this.Stop)
	boot.Define(boot.CLOSE, this.name, this.Close)
	boot.Define(boot.CLEANUP, this.name, this.Cleanup)

	boot.Install(this.name, this)
}
