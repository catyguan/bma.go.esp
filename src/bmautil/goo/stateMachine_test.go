package goo

import (
	"fmt"
	"testing"
)

const (
	STATE1 = uint32(1)
	STATE2 = uint32(2)
	STATE3 = uint32(3)
	STATE4 = uint32(4)
)

func buildStates() StateCollection {
	return StateCollection{
		STATE1: "ST-1",
		STATE2: "ST-2",
		STATE3: "ST-3",
		STATE4: "ST-4",
	}
}

func canEnter4test(o interface{}, cur uint32, st uint32) bool {
	if st == STATE4 {
		return true
	}
	switch cur {
	case STATE1:
		return st == STATE2
	case STATE2:
		return st == STATE3
	}
	return false
}

func afteLeave4test(o interface{}, st uint32) {
	fmt.Println("leave", st)
}

func afterEnter4test(o interface{}, st uint32) {
	fmt.Println("enter", st)
}

func TestSM(t *testing.T) {
	sm := new(StateMachine)
	sm.InitStateMachine(STATE1, buildStates())
	sm.SetCanEnterF(canEnter4test)
	sm.SetAfterEnterF(afterEnter4test)
	sm.SetAfterLeaveF(afteLeave4test)

	fmt.Printf("%T\n", sm)

	if true {
		ss := sm.String()
		b := sm.TryEnter(nil, STATE2)
		fmt.Println(ss, "enter", "STATE2", b)
	}
	if true {
		ss := sm.String()
		b := sm.TryEnter(nil, STATE1)
		fmt.Println(ss, "enter", "STATE1", b)
	}
	if true {
		ss := sm.String()
		b := sm.TryEnter(nil, STATE4)
		fmt.Println(ss, "enter", "STATE4", b)
	}
	fmt.Println(sm.String())
}
