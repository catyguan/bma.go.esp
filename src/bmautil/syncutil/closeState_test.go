package syncutil

import (
	"testing"
	"time"
)

func TestCloseState(t *testing.T) {
	var st *CloseState = NewCloseState()

	if st.IsClosing() {
		t.Error("can't closing")
	}
	go func() {
		for {
			if st.IsClosing() {
				t.Log("closing, exit")
				st.DoneClose()
				return
			}
		}
	}()
	time.Sleep(10 * time.Millisecond)
	st.AskClose()
	if !st.IsClosing() {
		t.Error("must closing")
	} else {
		st.WaitClosed()
		t.Log("all done")
	}
}
