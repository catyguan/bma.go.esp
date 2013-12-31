package syncutil

import (
	"testing"
	"time"
)

func TestEvent(t *testing.T) {
	runIt := func(e *Event) {
		time.Sleep(500 * time.Millisecond)
		e.Done()
	}
	e1 := NewAutoEvent()
	go runIt(e1)
	if e1.CheckEvent() {
		t.Error("event can't Set")
		return
	}
	t.Log("waiting", time.Now())
	if e1.WaitEvent() {
		t.Log("Event", time.Now())
	} else {
		t.Error("event must Set")
	}
	if e1.CheckEvent() {
		t.Error("event must Reset auto")
	}

	e2 := NewManulResetEvent()
	go runIt(e2)
	t.Log("waiting", time.Now())
	if e2.WaitEvent() {
		t.Log("Event", time.Now())
	} else {
		t.Error("event must Set")
	}
	if !e2.CheckEvent() {
		t.Error("event must not Reset auto")
	}
	e2.Reset()
	if e2.CheckEvent() {
		t.Error("event must Reset")
	}
}
