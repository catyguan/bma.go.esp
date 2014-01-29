package election

import (
	"bmautil/coder"
	"esp/cluster/nodeid"
	"esp/espnet"
	"logger"
)

const (
	OP_VOTE_REQ      = "vr"
	OP_VOTE_RESP     = "vp"
	OP_ANNOUNCE_REQ  = "ar"
	OP_ANNOUNCE_RESP = "ap"
)

// candidateState
func (this *candidateState) WriteState(xd *espnet.MessageXData) error {
	xd.Add(1, this.Id, nodeid.Coder)
	xd.Add(2, this.Epoch, EpochIdCoder)
	xd.Add(3, this.Status, StatusCoder)
	xd.Add(4, this.Leader, nodeid.Coder)
	return nil
}

func (this *candidateState) ReadState(it *espnet.XDataIterator) (bool, error) {
	switch it.Xid() {
	case 1:
		v, err := it.Value(nodeid.Coder)
		if err != nil {
			return true, err
		}
		this.Id = v.(nodeid.NodeId)
	case 2:
		v, err := it.Value(EpochIdCoder)
		if err != nil {
			return true, err
		}
		this.Epoch = v.(EpochId)
	case 3:
		v, err := it.Value(StatusCoder)
		if err != nil {
			return true, err
		}
		this.Status = v.(Status)
	case 4:
		v, err := it.Value(nodeid.Coder)
		if err != nil {
			return true, err
		}
		this.Leader = v.(nodeid.NodeId)
	default:
		return false, nil
	}
	return true, nil
}

// VoteReq
func (this *VoteReq) Write(msg *espnet.Message) error {
	msg.GetAddress().Set(espnet.ADDRESS_OP, OP_VOTE_REQ)
	xd := msg.XDatas()
	this.WriteState(xd)
	xd.Add(5, this.Proposal, nodeid.Coder)
	xd.Add(6, this.Renew, coder.Bool)
	return nil
}

func (this *VoteReq) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		ok, err0 := this.ReadState(it)
		if err0 != nil {
			return err0
		}
		if ok {
			continue
		}
		switch it.Xid() {
		case 5:
			v, err := it.Value(nodeid.Coder)
			if err != nil {
				return err
			}
			this.Leader = v.(nodeid.NodeId)
		case 6:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Renew = v.(bool)
		}
	}
	return nil
}

// ServiceHandler
type ServiceHandler struct {
	service *Service
}

func (this *ServiceHandler) Init(s *Service) {
	this.service = s
}

func (this *ServiceHandler) Serv(msg *espnet.Message, rep espnet.ServiceResponser) error {
	op := msg.GetAddress().Get(espnet.ADDRESS_OP)
	switch op {
	case OP_VOTE_REQ:
		req := new(VoteReq)
		err := req.Read(msg)
		if err != nil {
			return err
		}
		return this.actionVoteReq(msg, req, rep)
	case OP_VOTE_RESP:
	case OP_ANNOUNCE_REQ:
	case OP_ANNOUNCE_RESP:
	}
	return logger.Warn(tag, "unknow op '%s'", op)
}

func (this *ServiceHandler) actionVoteReq(msg *espnet.Message, req *VoteReq, rep espnet.ServiceResponser) error {
	// logger.Debug(tag, "actionList(%v) -> %v", req, resp)
	// rmsg := msg.ReplyMessage()
	// resp.Write(rmsg)
	// rep.SendMessage(rmsg)
	return nil
}
