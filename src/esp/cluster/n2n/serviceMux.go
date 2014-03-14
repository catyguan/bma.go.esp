package n2n

import (
	"bmautil/coder"
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"fmt"
	"logger"
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
	msg.GetAddress().SetOp(OP_JOIN)
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

func (this *Service) serve(msg *esnp.Message, ch espchannel.Channel) {
	err := func() error {
		addr := msg.GetAddress()
		op := addr.GetOp()
		switch op {
		case OP_JOIN:
			req := new(joinReq)
			err := req.Read(msg)
			if err != nil {
				return err
			}
			err2 := this.doJoin(req, ch)
			if err2 != nil {
				return err2
			}
			rmsg := msg.ReplyMessage()
			ch.SendMessage(rmsg)
			return nil
		}
		return fmt.Errorf("unknow method '%s'", op)
	}()
	if err != nil {
		logger.Error(tag, "'%s' serve fail - %s", ch, err)
		rmsg := msg.ReplyMessage()
		rmsg.BeError(err)
		espchannel.CloseAfterSend(rmsg)
	}
}
