package gotask

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(5*time.Second, func() {
		os.Exit(-100)
	})
}

func TestTimer(t *testing.T) {
	safeCall()
	gt := NewGoTask()
	defer gt.Close()

	timer := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("i do...")
			case <-gt.C:
				fmt.Println("end!")
				return
			}
		}
	}()
	time.Sleep(1 * time.Second)
}
