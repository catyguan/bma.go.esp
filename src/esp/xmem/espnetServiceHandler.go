package xmem

import (
	"bmautil/binlog"
	"bmautil/coder"
	"esp/espnet"
	"fmt"
	"logger"
)

// ServiceHandler
type ServiceHandler struct {
	service *Service
}

func (this *ServiceHandler) Init(s *Service) {
	this.service = s
}

func (this *ServiceHandler) Serv(msg *espnet.Message, rep espnet.ServiceResponser) error {
	v, err := espnet.FrameCoders.XData.Get(msg.ToPackage(), 0, coder.Int8)
	if err != nil {
		return err
	}
	if v == nil {
		return logger.Warn(tag, "unknow SHAction")
	}
	switch v.(int8) {
	case SHA_SLAVE_JOIN:
		req := new(SHRequestSlaveJoin)
		err = req.Read(msg)
		if err != nil {
			return err
		}
		return this.actionSlaveJoin(msg, req, rep)
	}
	return logger.Warn(tag, "unknow Action %d", v)
}

func (this *ServiceHandler) actionSlaveJoin(msg *espnet.Message, req *SHRequestSlaveJoin, rep espnet.ServiceResponser) error {
	ch := rep.GetChannel()
	if ch == nil {
		return fmt.Errorf("ServiceResponser GetChannel nil")
	}
	lis := func(seq binlog.BinlogVer, data []byte, closed bool) {
		if closed {
			ch.AskClose()
		} else {
			ev := new(SHEventBinlog)
			ev.Group = req.Group
			ev.Version = seq
			ev.Data = data
			evm := espnet.NewMessage()
			ev.Write(evm)
			logger.Debug(tag, "'%s' send binlog %d -> %s", req.Group, seq, ch)
			ch.SendMessage(evm)
		}
	}
	rd, err := this.service.SlaveJoin(req.Group, req.Version, lis)
	if err != nil {
		return err
	}
	ch.SetCloseListener("", func() {
		rd.Close()
	})
	return nil
}

// espnet
func (this *Service) CreateHandleRequest() espnet.ServiceHandler {
	sh := new(ServiceHandler)
	sh.Init(this)
	return sh.Serv
}
