package mempipeline

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"sync/atomic"
)

type memPipelineChannel struct {
	name    string
	send    chan *esnp.Message
	receive chan *esnp.Message

	closed        uint32
	listener      esnp.MessageListener
	closeListener func()
}

func (this *memPipelineChannel) String() string {
	return this.name
}

func (this *memPipelineChannel) IsClosing() bool {
	v := atomic.LoadUint32(&this.closed)
	return v != 0
}

func (this *memPipelineChannel) AskClose() {
	this.Shutdown()
}

func (this *memPipelineChannel) Shutdown() {
	defer func() {
		recover()
	}()
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		close(this.send)
		close(this.receive)

		if this.closeListener != nil {
			this.closeListener()
			this.closeListener = nil
		}
		this.listener = nil
	}
}

func (this *memPipelineChannel) GetProperty(name string) (interface{}, bool) {
	if name == espsocket.PROP_SOCKET_REMOTE_ADDR {
		return this.name, true
	}
	return nil, false
}

func (this *memPipelineChannel) SetProperty(name string, val interface{}) bool {
	return false
}

func (this *memPipelineChannel) Bind(rec esnp.MessageListener, closeLis func()) {
	this.listener = rec
	this.closeListener = closeLis
}

func (this *memPipelineChannel) SendMessage(ev *esnp.Message, cb espsocket.SendCallback) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%s", e)
		}
	}()
	this.send <- ev
	if cb != nil {
		cb(nil)
	}
	return nil
}

func (this *memPipelineChannel) run() {
	for {
		msg := <-this.receive
		if msg == nil {
			return
		}
		lis := this.listener
		if lis != nil {
			lis(msg)
		}
	}
}

// MemPipeline
type MemPipeline struct {
	ca  chan *esnp.Message
	cb  chan *esnp.Message
	cha *memPipelineChannel
	chb *memPipelineChannel
}

func NewMemPipeline(n string, sz int) *MemPipeline {
	this := new(MemPipeline)
	this.ca = make(chan *esnp.Message, sz)
	this.cb = make(chan *esnp.Message, sz)

	this.cha = this.initChannle(fmt.Sprintf("%s:a", n), this.ca, this.cb)
	this.chb = this.initChannle(fmt.Sprintf("%s:b", n), this.cb, this.ca)

	return this
}

func (this *MemPipeline) initChannle(n string, s chan *esnp.Message, r chan *esnp.Message) *memPipelineChannel {
	o := new(memPipelineChannel)
	o.name = n
	o.send = s
	o.receive = r
	go o.run()
	return o
}

func (this *MemPipeline) ChannelA() espsocket.Channel {
	return this.cha
}

func (this *MemPipeline) ChannelB() espsocket.Channel {
	return this.chb
}
