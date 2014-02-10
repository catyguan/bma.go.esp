package espnet

import (
	"boot"
	"fmt"
	"logger"
	"sync"
	"time"
)

type PChannel struct {
	id   uint32
	name string
	cf   ChannelFactory

	lock   sync.RWMutex
	ch     Channel
	lis    MessageListener
	lgroup CloseListenerGroup
	closed bool
}

func NewPChannel(n string, cf ChannelFactory, ch Channel) *PChannel {
	r := new(PChannel)
	r.name = n
	r.id = NextChanneId()
	r.cf = cf
	if ch != nil {
		r.initChannel(ch)
	}
	return r
}

func (this *PChannel) IsBreak() bool {
	return this.channel() == nil
}

func (this *PChannel) IsOpen() bool {
	return this.channel() != nil
}

func (this *PChannel) toChannel() Channel {
	return this
}

func (this *PChannel) channel() Channel {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.ch
}

func (this *PChannel) Id() uint32 {
	return this.id
}

func (this *PChannel) Name() string {
	return this.name
}

func (this *PChannel) String() string {
	return fmt.Sprintf("PChannel[%s]", this.name)
}

func (this *PChannel) AskClose() {
	this.lgroup.OnClose()
}

func (this *PChannel) GetProperty(name string) (interface{}, bool) {
	ch := this.channel()
	if ch != nil {
		return ch.GetProperty(name)
	}
	return nil, false
}

func (this *PChannel) SetProperty(name string, val interface{}) bool {
	ch := this.channel()
	if ch != nil {
		return ch.SetProperty(name, val)
	}
	return false
}

func (this *PChannel) SetMessageListner(rec MessageListener) {
	this.lis = rec
	ch := this.channel()
	if ch != nil {
		ch.SetMessageListner(rec)
	}
}

func (this *PChannel) SendMessage(ev *Message) error {
	ch := this.channel()
	if ch == nil {
		this.reconnect()
	}
	if ch != nil {
		return ch.SendMessage(ev)
	}
	return fmt.Errorf("breaked channel")
}

func (this *PChannel) reconnect() {
	logger.Debug(tag, "%s reconnecting", this)
	this.lock.RLock()
	skip := this.closed || this.ch != nil
	this.lock.RUnlock()
	if skip {
		return
	}
	c, err := this.cf.NewChannel()
	if err != nil {
		logger.Debug(tag, "%s reconnect fail - %s", this, err)
		time.AfterFunc(1*time.Millisecond, this.reconnect)
		return
	}
	logger.Debug(tag, "%s reconnected", this)
	this.initChannel(c)
}

func (this *PChannel) initChannel(c Channel) {
	this.lock.Lock()
	if this.closed || this.ch != nil {
		this.lock.Unlock()
		c.AskClose()
		return
	}
	defer this.lock.Unlock()

	this.ch = c
	c.SetMessageListner(this.lis)
	c.SetCloseListener("pchannel", func() {
		this.lock.Lock()
		if this.ch == c {
			this.ch = nil
		}
		this.lock.Unlock()
		this.reconnect()
	})
}

func (this *PChannel) SetCloseListener(name string, lis func()) error {
	this.lgroup.Set(name, lis)
	return nil
}

func (this *PChannel) Close() bool {
	this.lock.Lock()
	this.closed = true
	this.lock.Unlock()

	if this.ch != nil {
		o := this.ch
		this.ch = nil
		o.SetCloseListener("pchannel", nil)
		o.AskClose()
	}

	if this.cf != nil {
		o := this.cf
		this.cf = nil
		boot.RuntimeStopCloseClean(o, false)
	}
	this.lgroup.OnClose()
	return true
}
