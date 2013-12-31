package espnet

import (
	"errors"
	"fmt"
)

type PipelineHandler func(in Channel, msg *Message, out Channel) error

// Connector
type Connector struct {
	left               Channel
	right              Channel
	LeftToRightHandler PipelineHandler
	RightToLeftHandler PipelineHandler
	CloseOnBreak       bool
}

func (this *Connector) closerId() string {
	return fmt.Sprintf("PL_%p", this)
}

func (this *Connector) Break() {
	lch := this.left
	rch := this.right
	this.left = nil
	this.right = nil

	if lch != nil {
		lch.SetMessageListner(nil)
	}
	this.doClose(lch)
	if rch != nil {
		rch.SetMessageListner(nil)
	}
	this.doClose(rch)
}

func (this *Connector) Connect(left Channel, right Channel) error {
	if left != nil {
		if this.left != nil {
			return errors.New("left channel exists")
		}
		this.left = left
	}
	if right != nil {
		if this.right != nil {
			return errors.New("right channel exists")
		}
		this.right = right
	}

	cid := this.closerId()
	this.left.SetMessageListner(this.LeftSendMessage)
	this.left.SetCloseListener(cid, this.LeftClose)

	this.right.SetMessageListner(this.RightSendMessage)
	this.right.SetCloseListener(cid, this.RightClose)
	return nil
}

func (this *Connector) IsLeftBreak() bool {
	return this.left == nil
}

func (this *Connector) IsRightBreak() bool {
	return this.right == nil
}

func (this *Connector) LeftSendMessage(msg *Message) error {
	if this.right != nil {
		if this.LeftToRightHandler != nil {
			return this.LeftToRightHandler(this.left, msg, this.right)
		}
		return this.right.SendMessage(msg)
	}
	return errors.New("closed")
}

func (this *Connector) RightSendMessage(msg *Message) error {
	if this.left != nil {
		if this.RightToLeftHandler != nil {
			return this.RightToLeftHandler(this.right, msg, this.left)
		}
		return this.left.SendMessage(msg)
	}
	return errors.New("closed")
}

func (this *Connector) doClose(ch Channel) {
	if ch != nil {
		ch.SetMessageListner(nil)
		ch.SetCloseListener(this.closerId(), nil)
		if this.CloseOnBreak {
			ch.AskClose()
		}
	}
}

func (this *Connector) LeftClose() {
	if this.left != nil {
		this.left.SetMessageListner(nil)
	}
	this.left = nil
	this.doClose(this.right)
}

func (this *Connector) RightClose() {
	if this.right != nil {
		this.right.SetMessageListner(nil)
	}
	this.right = nil
	this.doClose(this.left)
}

// Pipeline
type Pipeline struct {
	name               string
	left               *VChannel
	right              *VChannel
	LeftToRightHandler PipelineHandler
	RightToLeftHandler PipelineHandler
}

func NewPipeline(name string) *Pipeline {
	this := new(Pipeline)

	var ch *VChannel
	ch = new(VChannel)
	ch.InitVChannel(this.name)
	ch.RemoveChannel = this.LeftClose
	ch.Sender = this.LeftSendMessage
	this.left = ch

	ch = new(VChannel)
	ch.InitVChannel(this.name)
	ch.RemoveChannel = this.RightClose
	ch.Sender = this.RightSendMessage
	this.right = ch

	return this
}

func (this *Pipeline) String() string {
	return this.name
}

func (this *Pipeline) GetLeftChannel() Channel {
	return this.left
}

func (this *Pipeline) GetRightChannel() Channel {
	return this.right
}

func (this *Pipeline) LeftSendMessage(msg *Message) error {
	if this.right != nil {
		if this.LeftToRightHandler != nil {
			return this.LeftToRightHandler(this.left, msg, this.right)
		}
		return this.right.SendMessage(msg)
	}
	return errors.New("closed")
}

func (this *Pipeline) RightSendMessage(msg *Message) error {
	if this.left != nil {
		if this.RightToLeftHandler != nil {
			return this.RightToLeftHandler(this.right, msg, this.left)
		}
		return this.left.SendMessage(msg)
	}
	return errors.New("closed")
}

func (this *Pipeline) doClose(ch *VChannel) {
	if ch != nil {
		ch.SetMessageListner(nil)
		ch.closeListeners.OnClose()
	}
}

func (this *Pipeline) LeftClose(ch *VChannel) {
	if this.left != nil {
		this.left.SetMessageListner(nil)
	}
	this.left = nil
	this.doClose(this.right)
}

func (this *Pipeline) RightClose(ch *VChannel) {
	if this.right != nil {
		this.right.SetMessageListner(nil)
	}
	this.right = nil
	this.doClose(this.left)
}
