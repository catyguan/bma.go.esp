package espnet

import (
	"bmautil/syncutil"

	// ChannelGroup
	"fmt"
	"sync"
	"sync/atomic"
)

type ChannelGroup struct {
	version    int32
	lock       sync.Mutex
	channels   map[uint32]Channel
	closeState syncutil.CloseState
}

func NewChannelGroup() *ChannelGroup {
	this := new(ChannelGroup)
	this.InitGroup()
	return this
}

func (this *ChannelGroup) InitGroup() {
	this.channels = make(map[uint32]Channel)
	this.closeState.InitCloseState()
}

func (this *ChannelGroup) Add(ch Channel) bool {
	if this.IsClosing() {
		return false
	}

	func() {
		this.lock.Lock()
		defer this.lock.Unlock()
		cid := ch.Id()
		_, ok := this.channels[cid]
		if !ok {
			this.channels[cid] = ch
			atomic.AddInt32(&this.version, 1)
		}
	}()

	name := fmt.Sprintf("CHANNEL_GROUP_%p", this)
	ch.SetCloseListener(name, func() {
		this.lock.Lock()
		defer this.lock.Unlock()
		cid := ch.Id()
		_, ok := this.channels[cid]
		if ok {
			delete(this.channels, cid)
			atomic.AddInt32(&this.version, 1)
			if len(this.channels) == 0 && this.closeState.IsClosing() {
				this.closeState.DoneClose()
			}
		}
	})
	return true
}

func (this *ChannelGroup) Remove(ch Channel) bool {
	if this.IsClosing() {
		return false
	}

	this.lock.Lock()
	defer this.lock.Unlock()
	cid := ch.Id()
	_, ok := this.channels[cid]

	if ok {
		delete(this.channels, cid)
		atomic.AddInt32(&this.version, 1)
		name := fmt.Sprintf("CHANNEL_GROUP_%p", this)
		ch.SetCloseListener(name, nil)
		return true
	}
	return false
}

func (this *ChannelGroup) IsClosing() bool {
	return this.closeState.IsClosing()
}
func (this *ChannelGroup) Close() {
	this.AskClose()
}

func (this *ChannelGroup) AskClose() bool {
	if this.closeState.AskClose() {
		this.lock.Lock()
		defer this.lock.Unlock()
		if len(this.channels) > 0 {
			for _, ch := range this.channels {
				go ch.AskClose()
			}
		} else {
			this.closeState.DoneClose()
		}
		return true
	}
	return false
}

func (this *ChannelGroup) WaitClosed() {
	this.closeState.WaitClosed()
}

func (this *ChannelGroup) Snapshot(mark int64, list []Channel) ([]Channel, int64) {
	ver := atomic.LoadInt32(&this.version)
	if ver == int32(mark) {
		return list, mark
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	r := make([]Channel, 0, len(this.channels))
	for _, ch := range this.channels {
		r = append(r, ch)
	}
	return r, int64(ver)
}

func (this *ChannelGroup) SnapshotMap(mark int64, m map[uint32]Channel) (map[uint32]Channel, int64) {
	ver := atomic.LoadInt32(&this.version)
	if ver == int32(mark) {
		return m, mark
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	r := make(map[uint32]Channel)
	for _, ch := range this.channels {
		r[ch.Id()] = ch
	}
	return r, int64(ver)
}
