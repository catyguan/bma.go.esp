package qpushpull

import (
	"fmt"
	"testing"
	"time"
)

func h4test(req interface{}) {
	fmt.Println("req", "=", req)
	time.Sleep(1000)
}

func TestBase(t *testing.T) {
	q := NewQueuePushPull(3, h4test)
	q.Run()
	for i := 0; i < 20; i++ {
		q.Push(i)
	}
	fmt.Println("main end")
	q.Close()
	q.WaitClose()
}
