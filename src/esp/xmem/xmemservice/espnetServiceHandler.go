package xmemservice

import (
	"bmautil/binlog"
	"bmautil/coder"
	"esp/espnet"
	"esp/xmem/xmemprot"
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
	case xmemprot.SHA_SLAVE_JOIN:
		req := new(xmemprot.SHRequestSlaveJoin)
		err = req.Read(msg)
		if err != nil {
			return err
		}
		return this.actionSlaveJoin(msg, req, rep)
	case xmemprot.SHA_GET:
		req := new(xmemprot.SHRequestGet)
		err = req.Read(msg)
		if err != nil {
			return err
		}
		return this.actionGet(msg, req, rep)
	case xmemprot.SHA_SET:
		req := new(xmemprot.SHRequestSet)
		err = req.Read(msg)
		if err != nil {
			return err
		}
		return this.actionSet(msg, req, rep)
	case xmemprot.SHA_DELETE:
		req := new(xmemprot.SHRequestDelete)
		err = req.Read(msg)
		if err != nil {
			return err
		}
		return this.actionDelete(msg, req, rep)
	case xmemprot.SHA_LIST:
		req := new(xmemprot.SHRequestList)
		err = req.Read(msg)
		if err != nil {
			return err
		}
		return this.actionList(msg, req, rep)
	}
	return logger.Warn(tag, "unknow Action %d", v)
}

func (this *ServiceHandler) actionList(msg *espnet.Message, req *xmemprot.SHRequestList, rep espnet.ServiceResponser) error {
	xm, err := this.service.CreateXMem(req.Group)
	if err != nil {
		return err
	}
	nlist, done, rerr := xm.List(xmemprot.MemKeyFromString(req.Key))
	if rerr != nil {
		return rerr
	}
	resp := new(xmemprot.SHResponseList)
	resp.Names = nlist
	resp.Miss = !done

	logger.Debug(tag, "actionList(%v) -> %v", req, resp)

	rmsg := msg.ReplyMessage()
	resp.Write(rmsg)

	rep.SendMessage(rmsg)

	return nil
}

func (this *ServiceHandler) actionDelete(msg *espnet.Message, req *xmemprot.SHRequestDelete, rep espnet.ServiceResponser) error {
	xm, err := this.service.CreateXMem(req.Group)
	if err != nil {
		return err
	}
	var rb bool
	var rerr error
	if req.Version.Valid() {
		rb, rerr = xm.CompareAndDelete(xmemprot.MemKeyFromString(req.Key), req.Version)
	} else {
		rb, rerr = xm.Delete(xmemprot.MemKeyFromString(req.Key))
	}
	if rerr != nil {
		return rerr
	}
	resp := new(xmemprot.SHResponseDelete)
	resp.Done = rb

	logger.Debug(tag, "actionDelete(%v) -> %v", req, resp)

	rmsg := msg.ReplyMessage()
	resp.Write(rmsg)

	rep.SendMessage(rmsg)

	return nil
}

func (this *ServiceHandler) actionGet(msg *espnet.Message, req *xmemprot.SHRequestGet, rep espnet.ServiceResponser) error {
	xm, err := this.service.CreateXMem(req.Group)
	if err != nil {
		return err
	}
	val, ver, hit, rerr := xm.Get(xmemprot.MemKeyFromString(req.Key))
	if rerr != nil {
		return rerr
	}
	resp := new(xmemprot.SHResponseGet)
	resp.Miss = !hit
	resp.Value = val
	resp.Version = ver

	logger.Debug(tag, "actionGet(%v) -> %v", req, resp)

	rmsg := msg.ReplyMessage()
	resp.Write(rmsg)

	rep.SendMessage(rmsg)

	return nil
}

func (this *ServiceHandler) actionSet(msg *espnet.Message, req *xmemprot.SHRequestSet, rep espnet.ServiceResponser) error {
	xm, err := this.service.CreateXMem(req.Group)
	if err != nil {
		return err
	}
	key := xmemprot.MemKeyFromString(req.Key)
	ver := xmemprot.VERSION_INVALID
	if req.Absent {
		ver, err = xm.SetIfAbsent(key, req.Value, req.Size)
	} else {
		if req.Version == xmemprot.VERSION_INVALID {
			ver, err = xm.Set(key, req.Value, req.Size)
		} else {
			ver, err = xm.CompareAndSet(key, req.Value, req.Size, req.Version)
		}
	}
	if err != nil {
		return err
	}
	resp := new(xmemprot.SHResponseSet)
	resp.Version = ver

	logger.Debug(tag, "actionSet(%v) -> %v", req, resp)

	rmsg := msg.ReplyMessage()
	resp.Write(rmsg)

	rep.SendMessage(rmsg)

	return nil
}

func (this *ServiceHandler) actionSlaveJoin(msg *espnet.Message, req *xmemprot.SHRequestSlaveJoin, rep espnet.ServiceResponser) error {
	ch := rep.GetChannel()
	if ch == nil {
		return fmt.Errorf("ServiceResponser GetChannel nil")
	}
	lis := func(seq binlog.BinlogVer, data []byte, closed bool) {
		if closed {
			ch.AskClose()
		} else {
			ev := new(xmemprot.SHEventBinlog)
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
