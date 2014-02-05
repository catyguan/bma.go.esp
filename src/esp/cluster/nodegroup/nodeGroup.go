package nodegroup

import (
	"bmautil/qexec"
	"esp/cluster/clusterbase"
	"esp/cluster/election"
	"esp/cluster/nodeid"
	"esp/espnet"
	"time"
)

type Action int

const (
	ACTION_JOIN   = Action(1)
	ACTION_LEAVE  = Action(2)
	ACTION_UPDATE = Action(3)
)

// ngSuperior
type ngSuperior struct {
	ng *NodeGroup
}

func (this *ngSuperior) Name() string {
	return this.ng.Name()
}

func (this *ngSuperior) AsyncPostVote(who nodeid.NodeId, vote *election.VoteReq) {
	this.ng.asyncPostVote(who, vote)
}

func (this *ngSuperior) AsyncRespVote(who nodeid.NodeId, resp *election.VoteResp) {
	this.ng.asyncRespVote(who, resp)
}

func (this *ngSuperior) AsyncPostAnnounce(who nodeid.NodeId, ann *election.AnnounceReq) {
	this.ng.asyncPostAnnounce(who, ann)
}

func (this *ngSuperior) AsyncRespAnnounce(who nodeid.NodeId, resp *election.AnnounceResp) {
	this.ng.asyncRespAnnounce(who, resp)
}

func (this *ngSuperior) DoStartLead(old nodeid.NodeId) error {
	return this.ng.doStartLead(old)
}

func (this *ngSuperior) DoStartFollow(lid nodeid.NodeId) error {
	return this.ng.doStartFollow(lid)
}

func (this *ngSuperior) DoStopFollow() error {
	return this.ng.doStopFollow()
}

func (this *ngSuperior) OnCandidateInvalid(id nodeid.NodeId) {
	this.ng.onCandidateInvalid(id)
}

// NodeGroup
type NodeGroup struct {
	name    string
	service *Service
	config  *NodeGroupConfig

	candidate *election.Candidate
	executor  *qexec.QueueExecutor
	channels  map[nodeid.NodeId]espnet.Channel

	role clusterbase.RoleType
}

func newNodeGroup(name string, s *Service, cfg *NodeGroupConfig) *NodeGroup {
	this := new(NodeGroup)
	this.name = name
	this.service = s
	this.config = cfg
	sp := new(ngSuperior)
	sp.ng = this
	this.candidate = election.NewCandidate(s.GetNodeId(), sp)
	this.executor = qexec.NewQueueExecutor(tag, 128, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.channels = make(map[nodeid.NodeId]espnet.Channel)
	this.role = clusterbase.ROLE_NONE
	return this
}

func (this *NodeGroup) Name() string {
	return this.name
}

func (this *NodeGroup) String() string {
	return this.candidate.String()
}

func (this *NodeGroup) Start() bool {
	if !this.executor.Run() {
		return false
	}
	this.executor.DoSync("init", func() error {
		return this.doInit()
	})
	return true
}

func (this *NodeGroup) Stop() bool {
	this.executor.Stop()
	return false
}

func (this *NodeGroup) Cleanup() bool {
	this.executor.WaitStop()
	return true
}

func (this *NodeGroup) WaitStop() {
	this.executor.WaitStop()
}

func (this *NodeGroup) requestHandler(ev interface{}) (bool, error) {
	switch rv := ev.(type) {
	case func() error:
		return true, rv()
	}
	return true, nil
}

func (this *NodeGroup) stopHandler() {
	this.doStopAll()
	for nid, _ := range this.channels {
		this.doCloseNode(nid, false)
	}
}

func (this *NodeGroup) doInit() error {
	this.doStartIdleCheck()
	return nil
}

func (this *NodeGroup) doStartIdleCheck() {
	if this.executor.IsClosing() {
		return
	}
	time.AfterFunc(time.Duration(this.config.IdleCheckMS)*time.Millisecond, func() {
		if this.executor.IsClosing() {
			return
		}
		this.candidate.CheckIdle()
		this.doStartIdleCheck()
	})
}

func (this *NodeGroup) doWaitTimeout(id nodeid.NodeId, epoch election.EpochId, vote bool) {
	if this.executor.IsClosing() {
		return
	}
	time.AfterFunc(time.Duration(this.config.ReqTimeoutMS)*time.Millisecond, func() {
		if this.executor.IsClosing() {
			return
		}
		this.candidate.OnReqTimeout(id, epoch, vote)
	})
}

func (this *NodeGroup) doStopAll() {

}

// interface impl
func (this *NodeGroup) doStartFollow(lid nodeid.NodeId) error {
	return nil
}

func (this *NodeGroup) doStopFollow() error {
	this.doStopAll()
	return nil
}

func (this *NodeGroup) onCandidateInvalid(id nodeid.NodeId) {
	this.doCloseNode(id, false)
}
