package syncutil

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	f, fe := NewFuture()
	go func() {
		if fe != nil {
			fe(1234, errors.New("my error"))
		}
	}()
	if !f.Wait(1 * time.Second) {
		t.Error("future wait fail")
		return
	}
	_, v, err := f.Get()
	fmt.Println("future result", v, err)
}
