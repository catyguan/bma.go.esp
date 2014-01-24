package espnet

import (
	"fmt"
	"sync"
)

type VChannel struct {
	id             uint32
	name           string
	closeListeners CloseListenerGroup

	peer   MessageListener
	Sender MessageListener

	RemoveChannel func(ch *VChannel)
}

func (this *VChannel) InitVChannel(name string) {
	this.id = NextChanneId()
	this.name = name
}

// Channel Impl
func (this *VChannel) Id() uint32 {
	return this.id
}

func (this *VChannel) Name() string {
	return this.name
}

func (this *VChannel) String() string {
	return fmt.Sprintf("VChanne[%s, %d]", this.name, this.id)
}

func (this *VChannel) AskClose() {
	if this.RemoveChannel != nil {
		this.RemoveChannel(this)
		this.RemoveChannel = nil
	}
	this.closeListeners.OnClose()
}

func (this *VChannel) GetProperty(name string) (interface{}, bool) {
	return nil, false
}

func (this *VChannel) SetProperty(name string, val interface{}) bool {
	return false
}

func (this *VChannel) SetMessageListner(rec MessageListener) {
	this.peer = rec
}

func (this *VChannel) SendMessage(ev *Message) error {
	if this.Sender != nil {
		return this.Sender(ev)
	}
	return nil
}

func (this *VChannel) SetCloseListener(name string, lis func()) error {
	this.closeListeners.Set(name, lis)
	return nil
}

func (this *VChannel) ServiceResponse(msg *Message) error {
	if this.peer != nil {
		return this.peer(msg)
	}
	return nil
}

// VChannelGroup
type VChannelGroup struct {
	lock     sync.Mutex
	channels []*VChannel
}

func (this *VChannelGroup) Len() int {
	if this.channels == nil {
		return 0
	}
	return len(this.channels)
}

func (this *VChannelGroup) Add(ch *VChannel) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.channels == nil {
		this.channels = make([]*VChannel, 0)
	}
	this.channels = append(this.channels, ch)
}

func (this *VChannelGroup) Remove(ch *VChannel) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.channels != nil {
		tmp := make([]*VChannel, 0, len(this.channels))
		for _, lch := range this.channels {
			if ch == lch {
				continue
			}
			tmp = append(tmp, lch)
		}
		this.channels = tmp
	}
	ch.RemoveChannel = nil
	ch.Sender = nil
}

func (this *VChannelGroup) OnClose() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.channels != nil {
		for _, ch := range this.channels {
			ch.closeListeners.OnClose()
		}
		this.channels = nil
	}
}
