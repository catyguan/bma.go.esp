package nodegroup

import (
	"esp/cluster/election"
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
	"fmt"
	"logger"
)

// impl interface
func (this *NodeGroup) asyncPostVote(who nodeinfo.NodeId, req *election.VoteReq) {
	ch, ok := this.channels[who]
	if !ok {
		err := fmt.Errorf("Node[%d] no channel")
		this.executor.DoNow("err", func() error {
			this.candidate.OnVoteResp(nil, err)
			return nil
		})
		return
	}
	msg := esnp.NewMessage()
	addr := msg.GetAddress()
	addr.SetService(this.config.ServiceName)
	addr.SetObject(this.name)
	req.Write(msg)
	msg.SetKind(esnp.MK_EVENT)
	err := ch.SendMessage(msg)
	if err != nil {
		this.executor.DoNow("err", func() error {
			this.candidate.OnVoteResp(nil, err)
			return nil
		})
		return
	}
	this.doWaitTimeout(who, req.State.Epoch, true)
}

func (this *NodeGroup) asyncRespVote(who nodeinfo.NodeId, resp *election.VoteResp) {
	ch, ok := this.channels[who]
	if !ok {
		logger.Warn(tag, "%s respVote fail - Node[%d] no channel", this, who)
		this.candidate.LeavePartner(who)
		return
	}
	msg := esnp.NewMessage()
	addr := msg.GetAddress()
	addr.SetService(this.config.ServiceName)
	addr.SetObject(this.name)
	resp.Write(msg)
	msg.SetKind(esnp.MK_EVENT)
	err := ch.SendMessage(msg)
	if err != nil {
		logger.Warn(tag, "%s respVote fail - %s", this, err)
		this.candidate.LeavePartner(who)
		return
	}
}

func (this *NodeGroup) asyncPostAnnounce(who nodeinfo.NodeId, req *election.AnnounceReq) {
	ch, ok := this.channels[who]
	if !ok {
		err := fmt.Errorf("Node[%d] no channel")
		this.executor.DoNow("err", func() error {
			this.candidate.OnVoteResp(nil, err)
			return nil
		})
		return
	}
	msg := esnp.NewMessage()
	addr := msg.GetAddress()
	addr.SetService(this.config.ServiceName)
	addr.SetObject(this.name)
	req.Write(msg)
	msg.SetKind(esnp.MK_EVENT)
	err := ch.SendMessage(msg)
	if err != nil {
		this.executor.DoNow("err", func() error {
			this.candidate.OnAnnounceResp(nil, err)
			return nil
		})
	}
	this.doWaitTimeout(who, req.Epoch, false)
}

func (this *NodeGroup) asyncRespAnnounce(who nodeinfo.NodeId, resp *election.AnnounceResp) {
	ch, ok := this.channels[who]
	if !ok {
		logger.Warn(tag, "%s respAnnounce fail - Node[%d] no channel", this, who)
		this.candidate.LeavePartner(who)
		return
	}

	msg := esnp.NewMessage()
	addr := msg.GetAddress()
	addr.SetService(this.config.ServiceName)
	addr.SetObject(this.name)
	resp.Write(msg)
	msg.SetKind(esnp.MK_EVENT)
	err := ch.SendMessage(msg)
	if err != nil {
		logger.Warn(tag, "%s respAnnounce fail - %s", this, err)
		this.candidate.LeavePartner(who)
	}
}
