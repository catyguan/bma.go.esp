package xmemservice

import (
	"bmautil/binlog"
	"bmautil/socket"
	"bmautil/valutil"
	"boot"
	"esp/espnet"
	"esp/xmem/xmemprot"
	"fmt"
	"logger"
	"time"
	"uprop"
)

// Config
type BLSlaveConfig struct {
	Address     string
	SerivceName string
	GroupName   string
	TimeoutMS   int
}

func (this *BLSlaveConfig) Valid() error {
	if this.Address == "" {
		return fmt.Errorf("address empty")
	}
	if this.SerivceName == "" {
		this.SerivceName = "xmem"
	}
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 5 * 1000
	}
	return nil
}

func (this *BLSlaveConfig) GetProperties() []*uprop.UProperty {
	b := new(uprop.UPropertyBuilder)
	b.NewProp("address", "binlog master ip:port").Optional(false).BeValue(this.Address, func(v string) error {
		this.Address = v
		return nil
	})
	b.NewProp("service", "binlog master service name, default xmem").BeValue(this.SerivceName, func(v string) error {
		this.SerivceName = v
		return nil
	})
	b.NewProp("group", "binlog master group name, default same as local name").BeValue(this.GroupName, func(v string) error {
		this.GroupName = v
		return nil
	})
	b.NewProp("timeout", "connect timeout in MS, default 5000").BeValue(this.TimeoutMS, func(v string) error {
		this.TimeoutMS = valutil.ToInt(v, this.TimeoutMS)
		return nil
	})
	return b.AsList()
}

func (this *BLSlaveConfig) GetGroupName(g string) string {
	if this.GroupName != "" {
		return this.GroupName
	}
	return g
}

func (this *BLSlaveConfig) NeedRestart(old *BLSlaveConfig, gname string) bool {
	if this.Address != old.Address {
		return true
	}
	if this.SerivceName != old.SerivceName {
		return true
	}
	if this.GetGroupName(gname) != old.GetGroupName(gname) {
		return true
	}
	return false
}

// BLSlave
type BLSlave struct {
	name    string
	config  *BLSlaveConfig
	service *Service

	dial   *espnet.DialPool
	client *espnet.ChannelClient
}

func (this *BLSlave) Init(name string, cfg *BLSlaveConfig, s *Service) error {
	this.name = name
	this.config = cfg
	this.service = s
	return nil
}

func (this *BLSlave) Run() bool {
	if this.config == nil || this.config.Valid() != nil {
		panic("BLSlave config invalid")
	}
	dcfg := new(espnet.DialPoolConfig)
	dcfg.Dial.Address = this.config.Address
	dcfg.Dial.TimeoutMS = this.config.TimeoutMS
	dcfg.MaxSize = 1
	dcfg.InitSize = 1
	err := dcfg.Valid()
	if err != nil {
		panic("init DialPoolConfig error " + err.Error())
	}
	this.dial = espnet.NewDialPool(this.name+"_blslave", dcfg, this.OnSocketInit)
	if !boot.RuntimeStartRun(this.dial) {
		boot.RuntimeStopCloseClean(this.dial, false)
		this.dial = nil
	}
	return true
}

func (this *BLSlave) OnSocketInit(sock *socket.Socket) error {
	go this.doConnectMaster()
	return nil
}

func (this *BLSlave) doConnectMaster() {
	if this.client != nil {
		this.client.Close()
		this.client = nil
	}
	sock, err := this.dial.GetSocket(time.Duration(this.config.TimeoutMS)*time.Millisecond, true)
	if err != nil {
		logger.Warn(tag, "BLSlave(%s) get socket fail - %s", this.name, err)
		return
	}

	ch := espnet.NewSocketChannel(sock, "espnet")
	cl := new(espnet.ChannelClient)
	err = cl.Connect(ch, true)
	if err != nil {
		logger.Warn(tag, "ChannelClient connect channel fail - %s", err)
		return
	}
	this.client = cl

	cl.SetMessageListner(this.OnMessage)

	var ver binlog.BinlogVer
	err = this.service.executor.DoSync("getSlaveVer", func() error {
		si, err := this.service.doGetGroup(this.name)
		if err != nil {
			return err
		}
		ver = si.group.blver
		return nil
	})
	if err != nil {
		logger.Warn(tag, "BLSlave(%s) get blversion fail and stop - %s", this.name, err)
		return
	}

	msg := espnet.NewMessage()
	msg.SetAddress(espnet.NewAddress(this.config.SerivceName))
	req := new(xmemprot.SHRequestSlaveJoin)
	req.Group = this.config.GetGroupName(this.name)
	req.Version = ver
	req.Write(msg)
	logger.Debug(tag, "BLSlave(%s) send sync request", this.name)
	err = cl.SendMessage(msg)
	if err != nil {
		cl.Close()
	}
}

func (this *BLSlave) OnMessage(msg *espnet.Message) error {
	err := msg.ToError()
	if err != nil {
		logger.Warn(tag, "BLSlave(%s) remote %s error and close - %s", this.name, this.config.Address, err)
		cl := this.client
		if cl != nil {
			cl.Close()
		}
		return nil
	}
	o := new(xmemprot.SHEventBinlog)
	err = o.Read(msg)
	if err != nil {
		logger.Warn(tag, "BLSlave(%s) decode error and close - %s", this.name, err)
		cl := this.client
		if cl != nil {
			cl.Close()
		}
		return nil
	}

	this.service.executor.DoNow("blslave", func() error {
		return this.service.doProcessBinog(this.name, o.Version, o.Data)
	})

	return nil
}

func (this *BLSlave) Stop() bool {
	if this.client != nil {
		this.client.Close()
		this.client = nil
	}
	if this.dial != nil {
		boot.RuntimeStopCloseClean(this.dial, false)
		this.dial = nil
	}
	return true
}

func (this *BLSlave) WaitStop() {

}

// Service
func (this *Service) doStartBinlogSlave(name string, mg *localMemGroup, cfg *MemGroupConfig) error {
	if mg.blslave != nil {
		logger.Debug(tag, "'%s' already start binlog slave, skip", name)
		return nil
	}
	if !cfg.IsEnableBinlogSlave() {
		return logger.Warn(tag, "'%s' binlog slave not enable", name)
	}
	blc := cfg.BLSlaveConfig
	s := new(BLSlave)
	s.Init(name, blc, this)
	if !s.Run() {
		return fmt.Errorf("'%s' binlog slave start fail", name)
	}
	logger.Info(tag, "'%s' start binlog slave - %s,%s,%s", name, blc.Address, blc.SerivceName, blc.GroupName)
	mg.blslave = s
	return nil
}

func (this *Service) doStopBinlogSlave(name string, mg *localMemGroup) error {
	if mg.blslave != nil {
		logger.Info(tag, "'%s' stop binlog slave", name)
		mg.blslave.Stop()
		mg.blslave = nil
	}
	return nil
}
