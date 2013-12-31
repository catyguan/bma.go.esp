package syncutil

import (
	"sync/atomic"
)

type CloseState struct {
	closed  *Event
	closing uint32
}

func NewCloseState() *CloseState {
	r := new(CloseState)
	r.InitCloseState()
	return r
}

func (this *CloseState) InitCloseState() {
	this.closed = NewManulResetEvent()
}

func (this *CloseState) AskClose() bool {
	r := atomic.CompareAndSwapUint32(&this.closing, 0, 1)
	return r
}

func (this *CloseState) BeginClose() {
	atomic.StoreUint32(&this.closing, 1)
}

func (this *CloseState) IsClosing() bool {
	return atomic.LoadUint32(&this.closing) == 1
}

func (this *CloseState) ResetClosing() {
	atomic.StoreUint32(&this.closing, 0)
}

func (this *CloseState) DoneClose() {
	atomic.StoreUint32(&this.closing, 1)
	this.closed.Done()
}

func (this *CloseState) WaitClosed() bool {
	return this.closed.WaitEvent()
}

func (this *CloseState) CheckClosed() bool {
	return this.closed.CheckEvent()
}

func (this *CloseState) IsClosed() bool {
	return this.closed.CheckEvent()
}

func (this *CloseState) ResetClosed() {
	this.closed.Reset()
}
