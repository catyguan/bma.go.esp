package election

import (
	"bytes"
	"esp/cluster/nodeinfo"
	"fmt"
	"logger"
	"math"
)

const (
	tag = "election"
)

type waitInfo struct {
	id    nodeinfo.NodeId
	epoch EpochId
	vote  bool
}

type Candidate struct {
	super ISuperior
	state CandidateState

	// voting
	partners map[nodeinfo.NodeId]*CandidateState
	votes    map[nodeinfo.NodeId]nodeinfo.NodeId
	waiting  []*waitInfo
}

func NewCandidate(id nodeinfo.NodeId, m ISuperior) *Candidate {
	this := new(Candidate)
	this.super = m
	this.state.Id = id
	this.partners = make(map[nodeinfo.NodeId]*CandidateState)
	this.partners[id] = &this.state
	this.votes = make(map[nodeinfo.NodeId]nodeinfo.NodeId)
	this.waiting = make([]*waitInfo, 0)
	return this
}

func (this *Candidate) String() string {
	return fmt.Sprintf("%s(%s)", this.super.Name(), &this.state)
}

func (this *Candidate) GetState() CandidateState {
	return this.state
}

func (this *Candidate) changeStatus(st Status) {
	if this.state.Status == st {
		return
	}
	old := this.state.Status
	this.state.Status = st
	logger.Info(tag, "%s status %s --> %s", this, old, st)
}

func (this *Candidate) JoinPartner(cs *CandidateState) {
	_, ok := this.partners[cs.Id]
	this.partners[cs.Id] = cs
	if !ok && logger.EnableDebug(tag) {
		ids := make([]nodeinfo.NodeId, 0)
		for nid, _ := range this.partners {
			ids = append(ids, nid)
		}
		logger.Debug(tag, "%s JoinPartner(%d) >> %v", this, cs.Id, ids)
	}
	if cs.Epoch > this.state.Epoch {
		this.keepUp(cs, "Join")
		return
	}
	if cs.Epoch == this.state.Epoch {
		if cs.Status.IsAnnounced() && this.state.Status.IsAnnounced() {
			if cs.Leader != this.state.Leader {
				logger.Debug(tag, "%s diff leader(%d) where (%d) join", this, cs.Leader, cs.Id)
				this.NewLooking(true)
			}
		}
	}
}

func (this *Candidate) UpdatePartnerState(st *CandidateState) {
	if st == nil {
		return
	}

	p, ok := this.partners[st.Id]
	if ok {
		*p = *st
	}
}

func (this *Candidate) LeavePartner(id nodeinfo.NodeId) {
	logger.Debug(tag, "%s leavePartner(%d)", this, id)
	delete(this.partners, id)
	delete(this.votes, id)

	switch this.state.Status {
	case STATUS_LOOKING:
		vid, ok := this.votes[this.state.Id]
		if ok && vid == id {
			// my vote is invalid
			logger.Debug(tag, "%s vote invalid, revote", this)
			this.NewLooking(false)
		} else {
			this.checkVotes()
		}
	case STATUS_FOLLOWING:
		if this.state.Leader == id {
			// leader leave
			this.super.DoStopFollow()
			this.changeStatus(STATUS_IDLE)
			this.state.Leader = nodeinfo.INVALID
			this.NewLooking(false)
		}
	}
}

func (this *Candidate) CheckIdle() {
	if this.state.Status == STATUS_IDLE {
		this.NewLooking(false)
	}
}

func (this *Candidate) NewLooking(renew bool) {
	this.startLooking(this.state.Epoch+1, renew)
}

func (this *Candidate) startLooking(epoch EpochId, renew bool) {
	this.state.Epoch = epoch
	this.changeStatus(STATUS_LOOKING)
	for k, _ := range this.votes {
		delete(this.votes, k)
	}
	for i, _ := range this.waiting {
		this.waiting[i] = nil
	}
	lid := nodeinfo.INVALID
	if !renew {
		for _, c := range this.partners {
			if c.Leader > lid {
				lid = c.Leader
			}
		}
	}
	if lid == 0 {
		for k, _ := range this.partners {
			if k > lid {
				lid = k
			}
		}
	}
	this.doVote(lid, renew)
}

func (this *Candidate) doVote(lid nodeinfo.NodeId, renew bool) {
	logger.Debug(tag, "%s doVote(%d, %v)", this, lid, renew)

	cs := this.state.Clone()

	this.votes[this.state.Id] = lid
	if this.checkVotes() {
		return
	}

	vreq := new(VoteReq)
	vreq.State = cs
	vreq.Proposal = lid
	vreq.Renew = renew

	for k, _ := range this.partners {
		if this.state.Id != k {
			this.newWait(k, cs.Epoch, true)
			this.super.AsyncPostVote(k, vreq)
		}
	}
}

func (this *Candidate) OnVoteReq(req *VoteReq) error {
	logger.Debug(tag, "%s OnVoteReq(%v)", this, req)
	this.UpdatePartnerState(req.State)
	if req.State.Epoch < this.state.Epoch {
		logger.Debug(tag, "%s reject outdate vote %d", this, req.State.Epoch)
		this.super.AsyncRespVote(req.State.Id, RejectVote(req, &this.state))
		return nil
	}
	if req.State.Epoch > this.state.Epoch || this.state.Status == STATUS_IDLE {
		this.super.AsyncRespVote(req.State.Id, AcceptVote(this.state.Id, req))
		this.keepUp(req.State, "voteReq")
		if this.state.Status == STATUS_LOOKING && this.state.Epoch == req.State.Epoch {
			this.putVote(req)
		}
		return nil
	}

	// same epoch
	switch this.state.Status {
	case STATUS_LOOKING:
		this.super.AsyncRespVote(req.State.Id, AcceptVote(this.state.Id, req))
		this.putVote(req)
		return nil
	default:
		if this.state.Leader == req.Proposal {
			this.super.AsyncRespVote(req.State.Id, AcceptVote(this.state.Id, req))
			return nil
		} else {
			logger.Debug(tag, "%s reject finish vote %d", this, req.State.Epoch)
			this.super.AsyncRespVote(req.State.Id, RejectVote(req, &this.state))
			return nil
		}
	}
}

func (this *Candidate) OnVoteResp(resp *VoteResp, err error) {
	logger.Debug(tag, "%s OnVoteResp(%v, %v)", this, resp, err)
	for i, w := range this.waiting {
		if w == nil {
			continue
		}
		if w.vote && w.id == resp.Id && w.epoch == resp.Epoch {
			this.waiting[i] = nil
			break
		}
	}
	if err != nil {
		// handle error
		if this.state.Status == STATUS_LOOKING {
			vid, ok := this.votes[this.state.Id]
			if ok && vid == resp.Id {
				this.LeavePartner(resp.Id)
				this.super.OnCandidateInvalid(resp.Id)
			}
		}
		return
	}
	if resp.Accept {
		return
	}
	// handle reject
	this.UpdatePartnerState(resp.State)
	if resp.State != nil && this.state.Epoch > resp.State.Epoch {
		return
	}
	if this.state.Status == STATUS_LOOKING {
		this.keepUp(resp.State, "voteResp")
	}
}

func (this *Candidate) keepUp(st *CandidateState, why string) {
	logger.Debug(tag, "%s keep-up epoch %d on %s", this, st.Epoch, why)
	switch st.Status {
	case STATUS_IDLE, STATUS_LOOKING:
		this.startLooking(st.Epoch, false)
	case STATUS_FOLLOWING, STATUS_LEADING:
		this.state.Epoch = st.Epoch
		this.Announce(st.Leader)
	}
}

func (this *Candidate) putVote(req *VoteReq) bool {
	sid := req.State.Id
	if _, ok := this.partners[sid]; !ok {
		logger.Warn(tag, "%d not partner", sid)
		return false
	}
	cs := req.State.Clone()
	this.partners[sid] = cs
	this.votes[sid] = req.Proposal
	return this.checkVotes()
}

func (this *Candidate) checkVotes() bool {
	pn := int(math.Ceil(float64(len(this.partners)) / 2))
	cts := make(map[nodeinfo.NodeId]int)
	total := 0
	lid := nodeinfo.INVALID
	for _, vid := range this.votes {
		v, _ := cts[vid]
		v = v + 1
		if v >= pn {
			lid = vid
			break
		}
		cts[vid] = v
		total++
	}
	if lid != nodeinfo.INVALID {
		this.Announce(lid)
		return true
	}
	if logger.EnableDebug(tag) {
		buf := bytes.NewBuffer([]byte{})
		for vid, c := range cts {
			if buf.Len() > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(fmt.Sprintf("%d=%d", vid, c))
		}
		logger.Debug(tag, "%s votes >> %s", this, buf.String())
	}
	if total == len(this.partners) {
		logger.Warn(tag, "%s vote no result, retry", this)
		this.NewLooking(false)
	}
	return false
}

func (this *Candidate) Announce(lid nodeinfo.NodeId) {
	logger.Debug(tag, "%s Annouce(%d)", this, lid)
	old := this.state.Leader
	this.state.Leader = lid
	var err error
	if this.state.Id == lid {
		this.changeStatus(STATUS_LEADING)
		err = this.super.DoStartLead(old)
	} else {
		this.changeStatus(STATUS_FOLLOWING)
		err = this.super.DoStartFollow(lid)
	}
	if err != nil {
		logger.Error(tag, "%s announce fail %s", this, err)
		this.changeStatus(STATUS_IDLE)
	}
	for k, _ := range this.partners {
		if this.state.Id != k {
			areq := new(AnnounceReq)
			areq.State = this.state.Clone()
			this.newWait(k, areq.State.Epoch, false)
			this.super.AsyncPostAnnounce(k, areq)
		}
	}
}

func (this *Candidate) OnAnnounceReq(req *AnnounceReq) error {
	logger.Debug(tag, "%s OnAnnounceReq(%v)", this, req)
	this.UpdatePartnerState(req.State)
	if req.State.Epoch < this.state.Epoch {
		logger.Debug(tag, "%s reject outdate announce %d", this, req.State.Epoch)
		this.super.AsyncRespAnnounce(req.State.Id, RejectAnnounce(req, &this.state))
		return nil
	}
	if req.State.Epoch > this.state.Epoch || this.state.Status == STATUS_IDLE {
		this.super.AsyncRespAnnounce(req.State.Id, AcceptAnnounce(this.state.Id, req))
		this.keepUp(req.State, "AnnounceReq")
		return nil
	}

	// same epoch
	switch this.state.Status {
	case STATUS_LOOKING:
		this.super.AsyncRespAnnounce(req.State.Id, AcceptAnnounce(this.state.Id, req))
		this.Announce(req.State.Leader)
		return nil
	default:
		if this.state.Leader == req.State.Leader {
			this.super.AsyncRespAnnounce(req.State.Id, AcceptAnnounce(this.state.Id, req))
			return nil
		} else {
			logger.Debug(tag, "%s reject diff announce vote %d:%d", this, req.State.Epoch, req.State.Leader)
			this.super.AsyncRespAnnounce(req.State.Id, RejectAnnounce(req, &this.state))
			this.NewLooking(true)
			return nil
		}
	}
}

func (this *Candidate) OnAnnounceResp(resp *AnnounceResp, err error) {
	logger.Debug(tag, "%s OnAnnounceResp(%v, %v)", this, resp, err)
	for i, w := range this.waiting {
		if w == nil {
			continue
		}
		if !w.vote && w.id == resp.Id && w.epoch == resp.Epoch {
			this.waiting[i] = nil
			break
		}
	}
	if err != nil {
		return
	}
	if resp.Accept {
		return
	}
	// handle reject
	this.UpdatePartnerState(resp.State)
	if resp.State != nil && this.state.Epoch >= resp.State.Epoch {
		return
	}
	this.keepUp(resp.State, "AnnounceResp")
}

func (this *Candidate) newWait(id nodeinfo.NodeId, epoch EpochId, vote bool) {
	for i, w := range this.waiting {
		if w == nil {
			this.waiting[i] = &waitInfo{id, epoch, vote}
			return
		}
	}
	this.waiting = append(this.waiting, &waitInfo{id, epoch, vote})
}

func (this *Candidate) OnReqTimeout(id nodeinfo.NodeId, epoch EpochId, vote bool) {
	for _, w := range this.waiting {
		if w == nil {
			continue
		}
		if w.id == id && w.epoch == epoch && w.vote == vote {
			err := fmt.Errorf("timeout")
			if vote {
				p := new(VoteResp)
				p.Id = id
				p.Epoch = epoch
				this.OnVoteResp(p, err)
			} else {
				p := new(AnnounceResp)
				p.Id = id
				p.Epoch = epoch
				this.OnAnnounceResp(p, err)
			}
			return
		}
	}
}
