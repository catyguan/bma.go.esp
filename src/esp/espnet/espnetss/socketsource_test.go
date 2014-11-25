package espnetss

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func T2estSocketSource(t *testing.T) {
	safeCall()

	ss := NewSocketSource("172.19.16.97:80", "test", 1)
	ss.Add("", "none")
	defer func() {
		ss.Close()
		time.Sleep(100 * time.Millisecond)
	}()
	ss.Start()
	time.Sleep(100 * time.Millisecond)
	sock, err := ss.Open(100)
	if err != nil {
		t.Error(err)
		return
	}
	sock.AskClose()

	time.Sleep(100 * time.Millisecond)
}

func T2estConnFail(t *testing.T) {
	ss := NewSocketSource("127.0.0.1:1080", "test", 1)
	ss.Add("", "none")
	defer func() {
		ss.Close()
		time.Sleep(100 * time.Millisecond)
	}()
	ss.Start()

	time.Sleep(5 * time.Second)
}

func TestMemSS(t *testing.T) {
	safeCall()

	ss := NewSocketSource("172.19.16.97:80", "test", 1)
	ss.Add("", "none")
	defer func() {
		ss.Close()
		time.Sleep(100 * time.Millisecond)
	}()
	ss.Start()
	time.Sleep(100 * time.Millisecond)
	sock, err := ss.Open(100)
	if err != nil {
		t.Error(err)
		return
	}
	sock.AskClose()

	time.Sleep(100 * time.Millisecond)
}
