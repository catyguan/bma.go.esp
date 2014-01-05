package qexec

import (
	"errors"
	"logger"
	"testing"
	"time"
)

func TestQExec(t *testing.T) {

	rhandler := func(req interface{}) (bool, error) {
		logger.Info("handler", "req = %v", req)
		if req.(int) < 0 {
			return true, errors.New("testerror")
		}
		return true, nil
	}
	exec := NewQueueExecutor("test", 10, rhandler)

	exec.Run()

	exec.Do("test", 123, nil)
	ev := make(chan error)
	exec.Do("test", 234, SyncCallback(ev))
	if err := <-ev; err != nil {
		t.Error("error", err)
	}
	exec.Do("test", -1, SyncCallback(ev))
	if err := <-ev; err == nil {
		t.Error("no error")
	}

	time.Sleep(1 * time.Second)
	exec.Stop()
	exec.WaitStop()
}

func TestResize(t *testing.T) {

	var exec *QueueExecutor
	rhandler := func(req interface{}) (bool, error) {
		v := req.(int)
		logger.Info("handler", "req = %v", v)
		if v < 0 {
			if !exec.Resize(-v) {
				panic("BUG")
			}
		}
		return true, nil
	}
	exec = NewQueueExecutor("test", 1, rhandler)

	exec.Run()

	var i int
	for i = 1; i <= 10; i++ {
		exec.DoNow("test", i)
	}
	logger.Debug("T", "here1")
	exec.DoNow("resize", -10)
	for i = 1; i <= 10; i++ {
		exec.DoNow("test", i)
	}
	logger.Debug("T", "here2")

	exec.DoSync("beforeStop", 9999)
	exec.Stop()
	exec.WaitStop()
}

func Test1(t *testing.T) {
	b := make([]byte, 10)
	b2 := b[1:5]
	b3 := b
	logger.Info("TEST", "%p = %p = %p", b, b2, b3)
}
