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
	this.event.Done()
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

func (this *Future) WaitDone() bool {
	return this.event.WaitEvent()
}

func (this *Future) Wait(d time.Duration) bool {
	return this.event.WaitEventTimeout(d)
}

type FutureGroup struct {
	c  int
	fs []*Future
}

func NewFutureGroup() *FutureGroup {
	r := new(FutureGroup)
	r.fs = make([]*Future, 0)
	return r
}

func (this *FutureGroup) Add(f *Future) {
	this.fs = append(this.fs, f)
}

func (this *FutureGroup) WaitAll(d time.Duration) bool {
	st := time.Now()
	d2 := d
	for _, f := range this.fs {
		if d2 <= 0 {
			if !f.IsDone() {
				return false
			}
		} else {
			if !f.Wait(d2) {
				return false
			}
			d2 = d - time.Now().Sub(st)
			// fmt.Println("new wait", d2)
		}
	}
	return true
}

func (this *FutureGroup) GetDone() []*Future {
	r := make([]*Future, 0)
	for _, f := range this.fs {
		if f.IsDone() {
			r = append(r, f)
		}
	}
	return r
}

func (this *FutureGroup) GetNotDone() []*Future {
	r := make([]*Future, 0)
	for _, f := range this.fs {
		if !f.IsDone() {
			r = append(r, f)
		}
	}
	return r
}
