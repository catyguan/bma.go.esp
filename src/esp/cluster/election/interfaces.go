package election

import (
	"esp/cluster/nodeid"
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
	Id     nodeid.NodeId
	Epoch  EpochId
	Status Status
	Leader nodeid.NodeId
}

func (this *CandidateState) String() string {
	return fmt.Sprintf("%d,%s,%d:%d", this.Id, this.Status, this.Epoch, this.Leader)
}

type VoteReq struct {
	CandidateState
	Proposal nodeid.NodeId
	Renew    bool
}

type VoteResp struct {
	Id     nodeid.NodeId
	Epoch  EpochId
	Accept bool
	State  *CandidateState
}

func AcceptVote(id nodeid.NodeId, req *VoteReq) *VoteResp {
	r := new(VoteResp)
	r.Id = id
	r.Epoch = req.Epoch
	r.Accept = true
	return r
}

func RejectVote(req *VoteReq, st *CandidateState) *VoteResp {
	r := new(VoteResp)
	r.Id = st.Id
	r.Epoch = req.Epoch
	r.Accept = false
	r.State = st
	return r
}

type AnnounceReq struct {
	CandidateState
}

type AnnounceResp struct {
	Id     nodeid.NodeId
	Epoch  EpochId
	Accept bool
	State  *CandidateState
}

func AcceptAnnounce(id nodeid.NodeId, req *AnnounceReq) *AnnounceResp {
	r := new(AnnounceResp)
	r.Id = id
	r.Epoch = req.Epoch
	r.Accept = true
	return r
}

func RejectAnnounce(req *AnnounceReq, st *CandidateState) *AnnounceResp {
	r := new(AnnounceResp)
	r.Id = st.Id
	r.Epoch = req.Epoch
	r.Accept = false
	r.State = st
	return r
}

type ISuperior interface {
	Name() string

	OnCandidateInvalid(id nodeid.NodeId)

	AsyncPostVote(who nodeid.NodeId, req *VoteReq)
	AsyncRespVote(who nodeid.NodeId, resp *VoteResp)
	AsyncPostAnnounce(who nodeid.NodeId, req *AnnounceReq)
	AsyncRespAnnounce(who nodeid.NodeId, resp *AnnounceResp)

	DoStartLead(old nodeid.NodeId) error
	DoStartFollow(lid nodeid.NodeId) error
	DoStopFollow() error
}
