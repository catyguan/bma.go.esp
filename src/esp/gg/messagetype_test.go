package gg

import (
	//	"fmt"
	"testing"
)

func TestMessageTypeBase(t *testing.T) {
	mt := NewMessageType4(1, 2, 3, 4)
	if mt.String() != "1.2.3.4" {
		t.Error("format fail")
	}
	if !mt.Match(NewMessageType4(0, 2, 3, 0)) {
		t.Error("match 1 fail")
	}
	if mt.Match(NewMessageType4(2, 2, 3, 0)) {
		t.Error("match 2 fail")
	}

	mt2 := MessageType(0x04030201)
	if mt2.String() != "4.3.2.1" {
		t.Error("format 2 fail")
	}

	mt3 := mt.Sub3(4, 6)
	if mt3.String() != "1.2.4.6" {
		t.Error("sub3 fail", mt3)
	}
	mt4 := mt.Sub4(6)
	if mt4.String() != "1.2.3.6" {
		t.Error("sub4 fail", mt4)
	}
	mt5 := mt.Sub2(4, 5, 6)
	if mt5.String() != "1.4.5.6" {
		t.Error("sub2 fail", mt5)
	}
}
