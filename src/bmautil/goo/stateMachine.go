package goo

import (
	"fmt"
	"sync/atomic"
)

type StateInfo struct {
	Id   uint32
	Name string
}

func NewStateInfO(id uint32, name string) *StateInfo {
	r := new(StateInfo)
	r.Id = id
	r.Name = name
	return r
}

func (this *StateInfo) String() string {
	return fmt.Sprintf("%d:%s", this.Id, this.Name)
}

type StateMachine struct {
	state     uint32
	subStates map[uint32]*StateMachine
	states    []*StateInfo
	// helper
	canEnterF   func(o interface{}, state uint32, toState uint32) bool
	afterLeaveF func(o interface{}, state uint32)
	afterEnterF func(o interface{}, state uint32)
}

func (this *StateMachine) String() string {
	var s string
	if this.states != nil {
		for _, si := range this.states {
			if si.Id == this.state {
				s = si.String()
				break
			}
		}
	}
	if s == "" {
		s = fmt.Sprintf("%d", this.state)
	}
	if this.subStates != nil {
		ssm := this.subStates[this.state]
		if ssm != nil {
			s = s + "," + ssm.String()
		}
	}
	return s
}

func (this *StateMachine) InitStateMachine(st uint32, sts []*StateInfo) {
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

func (this *StateMachine) HasSubState(s uint32) bool {
	if this.subStates != nil {
		_, ok := this.subStates[s]
		return ok
	}
	return false
}
func (this *StateMachine) SetSubState(s uint32, sm *StateMachine) {
	if this.subStates == nil {
		this.subStates = make(map[uint32]*StateMachine)
	}
	this.subStates[s] = sm
}
func (this *StateMachine) GetSubState(s uint32) *StateMachine {
	if this.subStates == nil {
		return nil
	}
	return this.subStates[s]
}

func (this *StateMachine) IsState(s uint32) bool {
	return this.GetState() == s
}
func (this *StateMachine) GetState() uint32 {
	return atomic.LoadUint32(&this.state)
}
func (this *StateMachine) doEnter(o interface{}, s uint32, try bool) bool {
	v := this.GetState()
	if v == s {
		return true
	}
	if try && this.canEnterF != nil {
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
	return this.doEnter(o, s, true)
}
func (this *StateMachine) Enter(o interface{}, s uint32) bool {
	return this.doEnter(o, s, false)
}
