package election

import (
	"esp/cluster/nodeinfo"
	"fmt"
)

type EpochId uint64

type Status uint8

func (O Status) String() string {
	switch O {
	case STATUS_IDLE:
		return "IDLE"
	case STATUS_LOOKING:
		return "LOOKING"
	case STATUS_LEADING:
		return "LEADING"
	case STATUS_FOLLOWING:
		return "FOLLOWING"
	}
	return fmt.Sprintf("UNKNOW-%d", O)
}

func (O Status) IsAnnounced() bool {
	return O == STATUS_LEADING || O == STATUS_FOLLOWING
}

const (
	STATUS_IDLE = iota
	STATUS_LOOKING
	STATUS_LEADING
	STATUS_FOLLOWING
)

type CandidateState struct {
	Id     nodeinfo.NodeId
	Epoch  EpochId
	Status Status
	Leader nodeinfo.NodeId
}

func (this *CandidateState) String() string {
	return fmt.Sprintf("%d,%s,%d:%d", this.Id, this.Status, this.Epoch, this.Leader)
}

func (this *CandidateState) Clone() *CandidateState {
	r := new(CandidateState)
	*r = *this
	return r
}

type VoteReq struct {
	State    *CandidateState
	Proposal nodeinfo.NodeId
	Renew    bool
}

type VoteResp struct {
	Id     nodeinfo.NodeId
	Epoch  EpochId
	Accept bool
	State  *CandidateState
}

func AcceptVote(id nodeinfo.NodeId, req *VoteReq) *VoteResp {
	r := new(VoteResp)
	r.Id = id
	r.Epoch = req.State.Epoch
	r.Accept = true
	return r
}

func RejectVote(req *VoteReq, st *CandidateState) *VoteResp {
	r := new(VoteResp)
	r.Id = st.Id
	r.Epoch = req.State.Epoch
	r.Accept = false
	r.State = st
	return r
}

type AnnounceReq struct {
	State *CandidateState
}

type AnnounceResp struct {
	Id     nodeinfo.NodeId
	Epoch  EpochId
	Accept bool
	State  *CandidateState
}

func AcceptAnnounce(id nodeinfo.NodeId, req *AnnounceReq) *AnnounceResp {
	r := new(AnnounceResp)
	r.Id = id
	r.Epoch = req.State.Epoch
	r.Accept = true
	return r
}

func RejectAnnounce(req *AnnounceReq, st *CandidateState) *AnnounceResp {
	r := new(AnnounceResp)
	r.Id = st.Id
	r.Epoch = req.State.Epoch
	r.Accept = false
	r.State = st
	return r
}

type ISuperior interface {
	Name() string

	OnCandidateInvalid(id nodeinfo.NodeId)

	AsyncPostVote(who nodeinfo.NodeId, req *VoteReq)
	AsyncRespVote(who nodeinfo.NodeId, resp *VoteResp)
	AsyncPostAnnounce(who nodeinfo.NodeId, req *AnnounceReq)
	AsyncRespAnnounce(who nodeinfo.NodeId, resp *AnnounceResp)

	DoStartLead(old nodeinfo.NodeId) error
	DoStartFollow(lid nodeinfo.NodeId) error
	DoStopFollow() error
}
