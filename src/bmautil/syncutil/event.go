package syncutil

import (
	"sync"
	"time"
)

type Event struct {
	m         sync.Mutex
	cond      *sync.Cond
	autoReset bool
	beSign    bool
}

func NewEvent(initSign, autoReset bool) *Event {
	r := new(Event)
	r.cond = sync.NewCond(&r.m)
	r.autoReset = autoReset
	r.beSign = initSign
	return r
}

func NewAutoEvent() *Event {
	return NewEvent(false, true)
}

func NewManulResetEvent() *Event {
	return NewEvent(false, false)
}

func (this *Event) Done() {
	this.SetEvent(true)
}

func (this *Event) SetEvent(flag bool) {
	this.m.Lock()
	defer this.m.Unlock()
	this.beSign = flag
	this.cond.Broadcast()

}

func (this *Event) Reset() {
	this.SetEvent(false)
}

func (this *Event) doCheckEvent(wait bool) bool {
	this.m.Lock()
	defer this.m.Unlock()
	r := this.beSign
	if !r {
		if wait {
			this.cond.Wait()
		}
		r = this.beSign
	}
	if r && this.autoReset {
		this.beSign = false
	}
	return r
}

func (this *Event) CheckEvent() bool {
	return this.doCheckEvent(false)
}

func (this *Event) WaitEvent() bool {
	return this.doCheckEvent(true)
}

func (this *Event) WaitEventTimeout(timeout time.Duration) bool {
	tm := time.AfterFunc(timeout, func() {
		this.m.Lock()
		defer this.m.Unlock()
		this.cond.Broadcast()
	})
	defer tm.Stop()
	return this.doCheckEvent(true)
}
