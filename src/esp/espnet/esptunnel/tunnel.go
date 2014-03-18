package esptunnel

import (
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"fmt"
	"logger"
	"sync"
)

const (
	tag = "Tunnel"
)

type pchItem struct {
	ch espchannel.Channel
	cb espchannel.BreakSupport
}

type Tunnel struct {
	id           uint32
	name         string
	CloseOnBreak bool

	lock   sync.RWMutex
	items  []*pchItem
	lis    esnp.MessageListener
	lgroup espchannel.CloseListenerGroup
}

func NewTunnel(n string) *Tunnel {
	r := new(Tunnel)
	r.name = n
	r.id = espchannel.NextChanneId()
	r.items = make([]*pchItem, 0)
	return r
}

func (this *Tunnel) Add(ch espchannel.Channel) int {
	return this.AddAll([]espchannel.Channel{ch})
}

func (this *Tunnel) AddAll(chlist []espchannel.Channel) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	clid := this.cid()
	for _, ch := range chlist {
		item := new(pchItem)
		item.ch = ch
		if cb, ok := ch.(espchannel.BreakSupport); ok {
			item.cb = cb
		}
		this.items = append(this.items, item)
		ch.SetCloseListener(clid, func() {
			go this.onChannelClose(item)
		})
		ch.SetMessageListner(this.myMessageListener)
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "AddChannel(%s)", item.ch.String())
		}
	}
	return len(this.items)
}

func (this *Tunnel) doRemove(item *pchItem) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	for idx, it := range this.items {
		if it == item {
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "onChannelClose(%s)", item.ch.String())
			}
			this.items[idx] = nil
			copy(this.items[idx:], this.items[idx+1:])
			this.items = this.items[0 : len(this.items)-1]
			return len(this.items)
		}
	}
	return -1
}

func (this *Tunnel) onChannelClose(item *pchItem) bool {
	c := this.doRemove(item)
	if c == -1 {
		return false
	}
	if c == 0 && this.CloseOnBreak {
		go this.AskClose()
	}
	return true
}

func (this *Tunnel) myMessageListener(msg *esnp.Message) error {
	if this.lis == nil {
		logger.Warn(tag, "%s not listener", this)
		return nil
	}
	return this.lis(msg)
}

func (this *Tunnel) IsBreak() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, item := range this.items {
		if item.cb != nil {
			if !item.cb.IsBreak() {
				return false
			}
		}
	}
	return true
}

func (this *Tunnel) toChannel() espchannel.Channel {
	return this
}

// Channel
func (this *Tunnel) Id() uint32 {
	return this.id
}

func (this *Tunnel) Name() string {
	return this.name
}

func (this *Tunnel) String() string {
	return fmt.Sprintf("Tunnel[%s,%d]", this.name, len(this.items))
}

func (this *Tunnel) cid() string {
	return fmt.Sprintf("%p", this)
}

func (this *Tunnel) doClose(force bool) {
	this.lgroup.OnClose()
	tmp := func() []*pchItem {
		this.lock.Lock()
		this.lock.Unlock()
		r := make([]*pchItem, 0, len(this.items))
		for _, item := range this.items {
			r = append(r, item)
		}
		this.items = make([]*pchItem, 0)
		return r
	}()

	clid := this.cid()
	for _, item := range tmp {
		ch := item.ch
		if ch != nil {
			ch.SetCloseListener(clid, nil)
			if force {
				ch.ForceClose()
			} else {
				ch.AskClose()
			}
		}
	}
}

func (this *Tunnel) AskClose() {
	this.doClose(false)
}

func (this *Tunnel) ForceClose() {
	this.doClose(true)
}

func (this *Tunnel) GetProperty(name string) (interface{}, bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, item := range this.items {
		ch := item.ch
		if ch != nil {
			v, ok := ch.GetProperty(name)
			if ok {
				return v, true
			}
		}
	}
	return nil, false
}

func (this *Tunnel) SetProperty(name string, val interface{}) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	r := false
	for _, item := range this.items {
		ch := item.ch
		if ch != nil {
			if ch.SetProperty(name, val) {
				r = true
			}
		}
	}
	return r
}

func (this *Tunnel) SetMessageListner(rec esnp.MessageListener) {
	this.lis = rec
}

func (this *Tunnel) doSendMessage(ev *esnp.Message, cb espchannel.ChannelSendCallback) error {
	err, rml := func() (error, []*pchItem) {
		var r []*pchItem
		this.lock.RLock()
		defer this.lock.RUnlock()
		for _, item := range this.items {
			if item.cb == nil || !item.cb.IsBreak() {
				ch := item.ch
				err := ch.SendMessage(ev, cb)
				if err == nil {
					return nil, r
				}
				logger.Debug(tag, "%s send fail - %s", ch, err)
			} else {
				logger.Debug(tag, "%s break, skip", item.ch)
			}
			if r == nil {
				r = []*pchItem{item}
			} else {
				r = append(r, item)
			}
		}
		return fmt.Errorf("breaked channel"), nil
	}()
	if rml != nil {
		for _, item := range rml {
			this.onChannelClose(item)
		}
	}
	return err
}

func (this *Tunnel) PostMessage(ev *esnp.Message) error {
	return this.doSendMessage(ev, nil)
}

func (this *Tunnel) SendMessage(ev *esnp.Message, cb espchannel.ChannelSendCallback) error {
	return this.doSendMessage(ev, cb)
}

func (this *Tunnel) SetCloseListener(name string, lis func()) error {
	this.lgroup.Set(name, lis)
	return nil
}

func (this *Tunnel) Stop() bool {
	this.AskClose()
	return true
}
