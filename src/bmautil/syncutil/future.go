package syncutil

import (
	"sync/atomic"
	"time"
)

type Future struct {
	event    *Event
	canceled uint32
	result   interface{}
	reserr   error
}

type FutureEnd func(v interface{}, err error)

func NewFuture() (*Future, FutureEnd) {
	r := new(Future)
	r.event = NewManulResetEvent()
	return r, func(v interface{}, err error) {
		r.result = v
		r.reserr = err
		r.event.Done()
	}
}

func (this *Future) Cancel() {
	atomic.StoreUint32(&this.canceled, 1)
}

func (this *Future) Get() (bool, interface{}, error) {
	if this.event.CheckEvent() {
		return true, this.result, this.reserr
	}
	return false, nil, nil
}

func (this *Future) IsCanceled() bool {
	v := atomic.LoadUint32(&this.canceled)
	return v != 0
}

func (this *Future) IsDone() bool {
	return this.event.CheckEvent()
}

func (this *Future) Wait(d time.Duration) bool {
	return this.event.WaitEventTimeout(d)
}
