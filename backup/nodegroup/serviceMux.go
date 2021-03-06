package nodegroup

import (
	"bmautil/coder"
	"esp/cluster/election"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"fmt"
)

const (
	OP_JOIN_PARTNER = "join"
)

type JoinPartnerReq struct {
	State  *election.CandidateState
	Relate bool
}

func (this *JoinPartnerReq) Write(msg *esnp.Message) error {
	msg.GetAddress().SetOp(OP_JOIN_PARTNER)
	xd := msg.XDatas()
	xd.Add(1, this.State, election.CandidateStateCoder)
	xd.Add(2, this.Relate, coder.Bool)
	return nil
}

func (this *JoinPartnerReq) Read(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(election.CandidateStateCoder)
			if err != nil {
				return err
			}
			if v != nil {
				this.State = v.(*election.CandidateState)
			}
		case 2:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Relate = v.(bool)
		}
	}
	return nil
}

func (this *NodeGroup) Serve(msg *esnp.Message, rep espservice.ServiceResponser) error {
	addr := msg.GetAddress()
	op := addr.GetOp()
	switch op {
	case OP_JOIN_PARTNER:
		req := new(JoinPartnerReq)
		err := req.Read(msg)
		if err != nil {
			return err
		}
		return this.executor.DoSync("joinPartner", func() error {
			return this.doJoinPartner(req, rep)
		})
	}
	return fmt.Errorf("unknow method '%s'", op)
}

func (this *NodeGroup) doJoinPartner(req *JoinPartnerReq, rep espservice.ServiceResponser) error {
	// ch := rep.GetChannel()
	if _, ok := this.channels[req.State.Id]; ok {
		return fmt.Errorf("NodeId(%d) exists", req.State.Id)
	}
	if !req.Relate {
		cs := this.candidate.GetState()
		nreq := new(JoinPartnerReq)
		nreq.State = &cs
		nreq.Relate = true

		msg := esnp.NewMessage()
		addr := msg.GetAddress()
		addr.SetService(this.config.ServiceName)
		addr.SetObject(this.name)
		msg.SetKind(esnp.MK_EVENT)
		nreq.Write(msg)
		rep.SendMessage(msg)
	}
	this.channels[req.State.Id] = rep.GetChannel()
	this.candidate.JoinPartner(req.State)
	return nil
}
