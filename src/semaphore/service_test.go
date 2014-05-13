package semaphore

import (
	"boot"
	"fmt"
	"testing"
	"time"
)

func f4test(s *Service, n string, idx int, wtime, stime time.Duration) func() {
	return func() {
		cf := func() {
			fmt.Println(idx, "executed")
			time.Sleep(stime)
		}
		s.Execute(n, cf, wtime)
	}
}

func TestQueue(t *testing.T) {
	cfgFile := "test.json"

	s := NewService("service")
	boot.AddService(s)

	f := func() {
		for i := 0; i < 10; i++ {
			f := f4test(s, "test2", i+1, 1*time.Second, 100*time.Millisecond)
			go f()
		}
	}
	boot.TestGo(cfgFile, 5, []func(){f})
}

func TestTimeout(t *testing.T) {
	cfgFile := "test.json"

	s := NewService("service")
	boot.AddService(s)

	f := func() {
		f1 := f4test(s, "test1", 1, 1*time.Second, 1100*time.Millisecond)
		go f1()
		f2 := f4test(s, "test1", 2, 1*time.Second, 1100*time.Millisecond)
		go f2()
	}
	boot.TestGo(cfgFile, 5, []func(){f})
}
