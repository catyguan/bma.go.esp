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

func TestFutureGroup(t *testing.T) {

	f1, fe1 := NewFuture()
	go func() {
		if fe1 != nil {
			time.Sleep(600 * time.Millisecond)
			fe1(1111, nil)
		}
	}()

	f2, fe2 := NewFuture()
	go func() {
		if fe2 != nil {
			time.Sleep(700 * time.Millisecond)
			fe2(2222, nil)
		}
	}()

	fg := NewFutureGroup()
	fg.Add(f1)
	fg.Add(f2)

	if !fg.WaitAll(1 * time.Second) {
		t.Error("future wait fail")
		return
	}
	_, v1, _ := f1.Get()
	_, v2, _ := f2.Get()
	fmt.Println("future result", v1, v2)

}
