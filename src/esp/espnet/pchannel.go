package espnet

import (
	"fmt"
	"logger"
	"sync"
	"time"
)

type pchItem struct {
	pos  int
	cf   ChannelFactory
	ch   Channel
	cing bool
}

type PChannel struct {
	id   uint32
	name string

	lock    sync.RWMutex
	current *pchItem
	items   []*pchItem
	lis     MessageListener
	lgroup  CloseListenerGroup
	closed  bool
}

func NewPChannel(n string) *PChannel {
	r := new(PChannel)
	r.name = n
	r.id = NextChanneId()
	r.items = make([]*pchItem, 0)
	return r
}

func (this *PChannel) Add(cf ChannelFactory) {
	this.AddAll([]ChannelFactory{cf})
}

func (this *PChannel) AddAll(cflist []ChannelFactory) {
	for _, cf := range cflist {
		item := new(pchItem)
		item.pos = len(this.items)
		item.cf = cf
		this.items = append(this.items, item)
		if this.current == nil {
			this.current = item
		}
	}
}

func (this *PChannel) OnReady() {
	for _, item := range this.items {
		this.reconnect(item)
	}
}

func (this *PChannel) IsBreak() *bool {
	v := this.channel() == nil
	return &v
}

func (this *PChannel) IsOpen() bool {
	return this.channel() != nil
}

func (this *PChannel) toChannel() Channel {
	return this
}

func (this *PChannel) cur() *pchItem {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.current
}

func (this *PChannel) next() {
	l := len(this.items)
	if l == 0 {
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	pos := (this.current.pos + 1) % l
	this.current = this.items[pos]
}

func (this *PChannel) channel() Channel {
	l := len(this.items)
	if l < 1 {
		return nil
	}
	for i := 0; i < l; i++ {
		item := this.cur()
		if item.ch == nil {
			go this.reconnect(item)
			this.next()
			continue
		}
		if item.ch != nil {
			if cb, ok := item.ch.(BreakSupport); ok {
				bv := cb.IsBreak()
				if bv != nil && *bv {
					this.recover(item)
					this.next()
					continue
				}
			}
		}
		return item.ch
	}
	return nil
}

func (this *PChannel) Id() uint32 {
	return this.id
}

func (this *PChannel) Name() string {
	return this.name
}

func (this *PChannel) String() string {
	return fmt.Sprintf("PChannel[%s,%d]", this.name, len(this.items))
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
	l := len(this.items)
	for i := 0; i < l; i++ {
		ch := this.channel()
		if ch == nil {
			break
		}
		err := ch.SendMessage(ev)
		if err == nil {
			return nil
		}
		logger.Debug(tag, "%s send fail - %s", ch, err)
		CloseForce(ch)
	}
	return fmt.Errorf("breaked channel")
}

func (this *PChannel) cing(item *pchItem, v bool) {
	this.lock.Lock()
	item.cing = v
	this.lock.Unlock()
}

func (this *PChannel) recover(item *pchItem) {
	func() {
		this.lock.Lock()
		defer this.lock.Unlock()
		if item.ch != nil {
			item.ch = nil
		}
	}()
	time.AfterFunc(5*time.Millisecond, func() {
		this.reconnect(item)
	})
}

func (this *PChannel) reconnect(item *pchItem) {
	this.lock.RLock()
	skip := this.closed || item.ch != nil || item.cing
	this.lock.RUnlock()
	if skip {
		return
	}

	this.cing(item, true)

	if fb, ok := item.cf.(BreakSupport); ok {
		vb := fb.IsBreak()
		if vb != nil && *vb {
			this.cing(item, false)
			this.recover(item)
			return
		}
	}

	logger.Debug(tag, "%s reconnecting[%d]", this, item.pos)
	c, err := item.cf.NewChannel()

	if err != nil {
		this.cing(item, false)
		logger.Debug(tag, "%s reconnect[%d] fail - %s", this, item.pos, err)
		this.recover(item)
		return
	}
	this.lock.Lock()
	item.cing = false
	if this.closed || item.ch != nil {
		this.lock.Unlock()
		c.AskClose()
		return
	}
	item.ch = c
	c.SetMessageListner(this.myMessageListener)
	c.SetCloseListener("pchannel", func() {
		iscur := false
		this.lock.RLock()
		if this.current == item {
			iscur = true
		}
		this.lock.RUnlock()
		if iscur {
			this.lgroup.OnClose()
		}
		this.recover(item)
	})
	this.lock.Unlock()

	logger.Debug(tag, "%s reconnected", this)
}

func (this *PChannel) myMessageListener(msg *Message) error {
	if this.lis == nil {
		logger.Warn(tag, "%s not listener", this)
		return nil
	}
	return this.lis(msg)
}

func (this *PChannel) SetCloseListener(name string, lis func()) error {
	this.lgroup.Set(name, lis)
	return nil
}

func (this *PChannel) Close() bool {
	this.lock.Lock()
	this.closed = true
	tmp := this.items
	this.items = make([]*pchItem, 0)
	this.lock.Unlock()

	for _, item := range tmp {
		item.cf = nil
		if item.ch != nil {
			o := item.ch
			item.ch = nil
			o.SetCloseListener("pchannel", nil)
			o.AskClose()
		}
	}

	this.lgroup.OnClose()
	return true
}
