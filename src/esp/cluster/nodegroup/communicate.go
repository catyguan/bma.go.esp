package nodegroup

import (
	"esp/cluster/election"
	"esp/cluster/nodeid"
	"esp/espnet"
	"fmt"
	"logger"
)

// impl interface
func (this *NodeGroup) asyncPostVote(who nodeid.NodeId, req *election.VoteReq) {
	ch, ok := this.channels[who]
	if !ok {
		err := fmt.Errorf("Node[%d] no channel")
		this.executor.DoNow("err", func() error {
			this.candidate.OnVoteResp(nil, err)
			return nil
		})
		return
	}
	msg := espnet.NewMessage()
	addr := msg.GetAddress()
	addr.Set(espnet.ADDRESS_SERVICE, this.config.ServiceName)
	addr.Set(espnet.ADDRESS_OBJECT, this.name)
	req.Write(msg)
	err := ch.SendMessage(msg)
	if err != nil {
		this.executor.DoNow("err", func() error {
			this.candidate.OnVoteResp(nil, err)
			return nil
		})
		return
	}
	this.doWaitTimeout(who, req.Epoch, true)
}

func (this *NodeGroup) asyncRespVote(who nodeid.NodeId, resp *election.VoteResp) {
	ch, ok := this.channels[who]
	if !ok {
		logger.Warn(tag, "%s respVote fail - Node[%d] no channel", this, who)
		this.candidate.LeavePartner(who)
		return
	}
	msg := espnet.NewMessage()
	addr := msg.GetAddress()
	addr.Set(espnet.ADDRESS_SERVICE, this.config.ServiceName)
	addr.Set(espnet.ADDRESS_OBJECT, this.name)
	resp.Write(msg)
	err := ch.SendMessage(msg)
	if err != nil {
		logger.Warn(tag, "%s respVote fail - %s", this, err)
		this.candidate.LeavePartner(who)
		return
	}
}

func (this *NodeGroup) asyncPostAnnounce(who nodeid.NodeId, req *election.AnnounceReq) {
	ch, ok := this.channels[who]
	if !ok {
		err := fmt.Errorf("Node[%d] no channel")
		this.executor.DoNow("err", func() error {
			this.candidate.OnVoteResp(nil, err)
			return nil
		})
		return
	}
	msg := espnet.NewMessage()
	addr := msg.GetAddress()
	addr.Set(espnet.ADDRESS_SERVICE, this.config.ServiceName)
	addr.Set(espnet.ADDRESS_OBJECT, this.name)
	req.Write(msg)
	err := ch.SendMessage(msg)
	if err != nil {
		this.executor.DoNow("err", func() error {
			this.candidate.OnAnnounceResp(nil, err)
			return nil
		})
	}
	this.doWaitTimeout(who, req.Epoch, false)
}

func (this *NodeGroup) asyncRespAnnounce(who nodeid.NodeId, resp *election.AnnounceResp) {
	ch, ok := this.channels[who]
	if !ok {
		logger.Warn(tag, "%s respAnnounce fail - Node[%d] no channel", this, who)
		this.candidate.LeavePartner(who)
		return
	}

	msg := espnet.NewMessage()
	addr := msg.GetAddress()
	addr.Set(espnet.ADDRESS_SERVICE, this.config.ServiceName)
	addr.Set(espnet.ADDRESS_OBJECT, this.name)
	resp.Write(msg)
	err := ch.SendMessage(msg)
	if err != nil {
		logger.Warn(tag, "%s respAnnounce fail - %s", this, err)
		this.candidate.LeavePartner(who)
	}
}
