package nodegroup

import (
	// NodeGroup
	"esp/cluster/election"
	"esp/cluster/nodeinfo"
	"esp/espnet/espchannel"
	"logger"
)

const (
	tag = "nodeGroup"
)

// IdleMessage

// NodeGroupConfig
type NodeGroupConfig struct {
	ServiceName  string
	IdleCheckMS  int
	ReqTimeoutMS int
}

func (this *NodeGroupConfig) Valid() error {
	if this.ServiceName == "" {
		this.ServiceName = "cluster"
	}
	if this.IdleCheckMS <= 0 {
		this.IdleCheckMS = 2000
	}
	if this.ReqTimeoutMS <= 0 {
		this.ReqTimeoutMS = 2000
	}
	return nil
}

// NodeGroup
func (this *NodeGroup) Join(ch espchannel.Channel, cs *election.CandidateState) error {
	return this.executor.DoSync("join", func() error {
		return this.doJoin(ch, cs)
	})
}

func (this *NodeGroup) doJoin(ch espchannel.Channel, cs *election.CandidateState) error {
	id := cs.Id
	_, ok := this.channels[id]
	if ok {
		logger.Debug(tag, "%s node(%d) channel exists", id)
		return nil
	}
	this.candidate.JoinPartner(cs)
	this.channels[id] = ch
	ch.SetCloseListener("NodeGroup_"+this.name, func() {
		this.CloseNode(id)
	})
	return nil
}

func (this *NodeGroup) CloseNode(id nodeinfo.NodeId) {
	this.executor.DoNow("onCloseNode", func() error {
		this.doCloseNode(id, true)
		return nil
	})
}

func (this *NodeGroup) doCloseNode(id nodeinfo.NodeId, leave bool) {
	ch, ok := this.channels[id]
	if !ok {
		return
	}
	delete(this.channels, id)
	ch.SetCloseListener("NodeGroup_", nil)
	ch.SetMessageListner(nil)
	ch.AskClose()

	if leave {
		this.candidate.LeavePartner(id)
	}
}
