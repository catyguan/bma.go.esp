package goo

import (
	"fmt"
	"sync/atomic"
)

type StateCollection map[uint32]string

func (this StateCollection) ToString(id uint32) string {
	n := ""
	for iv, in := range this {
		if iv == id {
			n = in
			break
		}
	}
	if n == "" {
		return fmt.Sprintf("%d:%d", id, id)
	}
	return fmt.Sprintf("%d:%s", id, n)
}

type StateMachine struct {
	state  uint32
	states StateCollection

	// helper
	canEnterF   func(o interface{}, state uint32, toState uint32) bool
	afterLeaveF func(o interface{}, state uint32)
	afterEnterF func(o interface{}, state uint32)
}

func (this *StateMachine) String() string {
	if this.states == nil {
		return fmt.Sprintf("%d", this.state)
	}
	return this.states.ToString(this.state)
}

func (this *StateMachine) InitStateMachine(st uint32, sts StateCollection) {
	this.state = st
	this.states = sts
}
func (this *StateMachine) SetCanEnterF(f func(o interface{}, state uint32, toState uint32) bool) {
	this.canEnterF = f
}
func (this *StateMachine) SetAfterLeaveF(f func(o interface{}, state uint32)) {
	this.afterLeaveF = f
}
func (this *StateMachine) SetAfterEnterF(f func(o interface{}, state uint32)) {
	this.afterEnterF = f
}

func (this *StateMachine) IsState(s uint32) bool {
	return this.GetState() == s
}
func (this *StateMachine) GetState() uint32 {
	return atomic.LoadUint32(&this.state)
}
func (this *StateMachine) doEnter(o interface{}, e uint32, s uint32, try bool) bool {
	v := this.GetState()
	if v == s {
		return true
	}
	if try && this.canEnterF != nil {
		if e != 0 && e != v {
			return false
		}
		if !this.canEnterF(o, v, s) {
			return false
		}
	}
	if !atomic.CompareAndSwapUint32(&this.state, v, s) {
		return false
	}
	if this.afterLeaveF != nil {
		this.afterLeaveF(o, v)
	}
	if this.afterEnterF != nil {
		this.afterEnterF(o, s)
	}
	return true
}
func (this *StateMachine) TryEnter(o interface{}, s uint32) bool {
	return this.doEnter(o, 0, s, true)
}
func (this *StateMachine) Enter(o interface{}, s uint32) bool {
	return this.doEnter(o, 0, s, false)
}
func (this *StateMachine) CompareAndEnter(o interface{}, expect uint32, s uint32) bool {
	return this.doEnter(o, expect, s, true)
}
