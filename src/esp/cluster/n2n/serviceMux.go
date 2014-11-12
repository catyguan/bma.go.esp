package n2n

import (
	"esp/cluster/nodebase"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
)

const (
	SN_N2N  = "espnode.n2n"
	OP_JOIN = "join"
)

type joinReq struct {
	Id   nodebase.NodeId
	Name string
	Host string
}

func (this *joinReq) String() string {
	return fmt.Sprintf("[Id=%d, Name=%s, Host=%s]", this.Id, this.Name, this.Host)
}

func (this *joinReq) Write(msg *esnp.Message) error {
	xd := msg.XDatas()
	xd.Add(1, this.Id, nodebase.NodeIdCoder)
	xd.Add(2, this.Name, esnp.Coders.String)
	xd.Add(3, this.Host, esnp.Coders.String)
	return nil
}

func (this *joinReq) Read(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(nodebase.NodeIdCoder)
			if err != nil {
				return err
			}
			if v != nil {
				this.Id = v.(nodebase.NodeId)
			}
		case 2:
			v, err := it.Value(esnp.Coders.String)
			if err != nil {
				return err
			}
			this.Name = v.(string)
		case 3:
			v, err := it.Value(esnp.Coders.String)
			if err != nil {
				return err
			}
			this.Host = v.(string)
		}
	}
	return nil
}

func (this *Service) makeJoinReq() *joinReq {
	req := new(joinReq)
	req.Id = nodebase.Id
	req.Name = nodebase.Name
	req.Host = this.config.Host
	return req
}

func (this *Service) handleJoin(sock *espsocket.Socket, msg *esnp.Message, doReply bool) error {
	req := new(joinReq)
	err := req.Read(msg)
	if err != nil {
		return err
	}
	err2 := this.doJoin(req, sock)
	if err2 != nil {
		return err2
	}
	if doReply {
		logger.Debug(tag, "reply joinReq -> (%s : %s)", req.Name, sock)
		rreq := this.makeJoinReq()
		rmsg := msg.ReplyMessage()
		rreq.Write(rmsg)
		sock.PostMessage(rmsg)
	}
	return nil
}

func (this *Service) Serve(sock *espsocket.Socket, msg *esnp.Message) error {
	addr := msg.GetAddress()
	op := addr.GetOp()
	switch op {
	case OP_JOIN:
		return this.goo.DoSync(func() error {
			err := this.handleJoin(sock, msg, true)
			if err != nil {
				logger.Warn(tag, "%s handle joinReq fail - %s", err)
			}
			return err
		})
	}
	return fmt.Errorf("unknow method '%s'", op)
}
