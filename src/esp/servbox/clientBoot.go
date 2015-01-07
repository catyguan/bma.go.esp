package servbox

import (
	"bmautil/netutil"
	"boot"
	"fmt"
	"io/ioutil"
	"logger"
)

type clientInfo struct {
	NodeName   string
	Net        string
	Address    string
	BoxNet     string
	BoxAddress string
	TimeoutMS  int
	MaxPackage int
	Disable    bool
}

func (this *clientInfo) Valid() error {
	if this.Disable {
		return nil
	}
	if this.NodeName == "" {
		return fmt.Errorf("NodeName empty")
	}
	if this.BoxNet == "" {
		this.BoxNet = "tcp"
	}
	if this.BoxAddress == "" {
		return fmt.Errorf("BoxAddress empty")
	}
	if this.Net == "" {
		this.Net = "tcp"
	}
	if this.Address == "" {
		this.Address = "127.0.0.1:0"
	}
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 5 * 1000
	}
	return nil
}

func (this *clientInfo) Compare(old *clientInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.Disable != old.Disable {
		return boot.CCR_NEED_START
	}
	if this.NodeName != old.NodeName {
		return boot.CCR_NEED_START
	}
	if this.BoxNet != old.BoxNet {
		return boot.CCR_NEED_START
	}
	if this.BoxAddress != old.BoxAddress {
		return boot.CCR_NEED_START
	}
	if this.Net != old.Net {
		return boot.CCR_NEED_START
	}
	if this.Address != old.Address {
		return boot.CCR_NEED_START
	}
	if this.TimeoutMS != old.TimeoutMS {
		return boot.CCR_CHANGE
	}
	if this.MaxPackage != old.MaxPackage {
		return boot.CCR_CHANGE
	}
	return boot.CCR_NONE
}

func (this *Client) Name() string {
	return this.name
}

func (this *Client) Prepare() {
}
func (this *Client) CheckConfig(ctx *boot.BootContext) bool {
	co := ctx.Config
	cfg := new(clientInfo)
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

func (this *Client) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*clientInfo)
	this.config = cfg
	return true
}

func (this *Client) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	if this.config.Disable {
		logger.Info(tag, "%s disable, skip start", this.name)
		return true
	}

	addr := this.config.Address
	if this.config.Net == "unix" {
		tfile, errT := ioutil.TempDir(addr, this.config.NodeName)
		if errT != nil {
			logger.Warn(tag, "%s create tmpfile(%s, %s) fail - %s", addr, this.config.NodeName, errT)
			return false
		}
		addr = tfile
	}
	logger.Debug(tag, "%s start listen (%s, %s)", this.name, this.config.Net, addr)
	lis, err := netutil.Listen(this.config.Net, addr)
	if err != nil {
		logger.Warn(tag, "%s listen at (%s, %s) fail - %s", this.name, this.config.Net, addr, err)
		return false
	}
	logger.Info(tag, "%s listen at (%s)", this.name, lis.Addr())
	this.listener = lis
	return true
}

func (this *Client) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}

	if this.config.Disable {
		logger.Info(tag, "%s disable, skip run", this.name)
		return true
	}

	addr := this.listener.Addr().String()
	go func() {
		defer func() {
			logger.Info(tag, "stop (%s)", addr)
		}()
		for {
			c, err := this.listener.Accept()
			if err == nil {
				if c == nil {
					return
				}
				go this.accept(c) // new connect
			} else {
				if !netutil.IsAcceptClose(err) {
					logger.Debug(tag, "accept fail and exit - %s", err)
				}
				return
			}
		}
	}()

	this.joinBox(this.config, addr)
	return true
}

func (this *Client) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type != boot.CCR_NEED_START {
		return true
	}
	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}
	if this.boxc != nil {
		close(this.boxc)
		this.boxc = nil
	}
	return true
}

func (this *Client) Stop() bool {
	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}
	if this.boxc != nil {
		close(this.boxc)
		this.boxc = nil
	}
	return true
}

func (this *Client) Close() bool {
	return true
}

func (this *Client) Cleanup() bool {
	return true
}
