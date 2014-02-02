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
	for i := 0; i < 100; i++ {
		q.Push(i)
	}
	fmt.Println("main end1")
	time.Sleep(1 * time.Second)
	fmt.Println("main end2")
	q.Close()
	q.WaitClose()
}
