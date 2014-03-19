package n2n

import (
	"bmautil/coder"
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"esp/espnet/espclient"
	"fmt"
	"logger"
	"time"
)

const (
	OP_JOIN = "join"
)

type joinReq struct {
	Id   nodeinfo.NodeId
	Name string
	URL  string
}

func (this *joinReq) Write(msg *esnp.Message) error {
	xd := msg.XDatas()
	xd.Add(1, this.Id, nodeinfo.NodeIdCoder)
	xd.Add(2, this.Name, coder.String)
	xd.Add(3, this.URL, coder.String)
	return nil
}

func (this *joinReq) Read(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(nodeinfo.NodeIdCoder)
			if err != nil {
				return err
			}
			if v != nil {
				this.Id = v.(nodeinfo.NodeId)
			}
		case 2:
			v, err := it.Value(coder.String)
			if err != nil {
				return err
			}
			this.Name = v.(string)
		case 3:
			v, err := it.Value(coder.String)
			if err != nil {
				return err
			}
			this.URL = v.(string)
		}
	}
	return nil
}

func (this *Service) makeJoinReq() *joinReq {
	req := new(joinReq)
	req.Id = this.ninfo.GetId()
	req.Name = this.ninfo.GetNodeName()
	req.URL = this.config.URL
	return req
}

func (this *Service) sendJoinReq(n string, url *esnp.URL, ch espchannel.Channel) {
	go func() {
		logger.Debug(tag, "send joinReq -> (%s : %s)", n, ch)

		req := this.makeJoinReq()

		msg := esnp.NewRequestMessage()
		addr := msg.GetAddress()
		url.BindAddress(addr)
		addr.SetOp(OP_JOIN)
		req.Write(msg)
		cl := espclient.NewChannelClient()
		if err := cl.Connect(ch, false); err != nil {
			logger.Debug(tag, "ChannelClient connect fail - %s", err)
			return
		}
		tm := time.NewTimer(url.GetTimeout(3 * time.Second))
		rmsg, err := cl.Call(msg, tm)
		if err != nil {
			logger.Debug(tag, "%s call fail - %s", ch, err)
			return
		}
		err = this.handleJoin(ch, rmsg, false)
		if err != nil {
			logger.Debug(tag, "%s handle join resp fail - %s", ch, err)
		}
	}()
}

func (this *Service) handleJoin(ch espchannel.Channel, msg *esnp.Message, doReply bool) error {
	req := new(joinReq)
	err := req.Read(msg)
	if err != nil {
		return err
	}
	err2 := this.doJoin(req, ch)
	if err2 != nil {
		return err2
	}
	if doReply {
		logger.Debug(tag, "reply joinReq -> (%s : %s)", req.Name, ch)
		rreq := this.makeJoinReq()
		rmsg := msg.ReplyMessage()
		rreq.Write(rmsg)
		ch.PostMessage(rmsg)
	}
	return nil
}

func (this *Service) Serve(ch espchannel.Channel, msg *esnp.Message) error {
	addr := msg.GetAddress()
	op := addr.GetOp()
	switch op {
	case OP_JOIN:
		err := this.handleJoin(ch, msg, true)
		if err != nil {
			logger.Warn(tag, "%s handle joinReq fail - %s", err)
		}
		return err
	}
	return fmt.Errorf("unknow method '%s'", op)
}
