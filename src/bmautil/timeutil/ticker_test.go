package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func TestClosableTicker(t *testing.T) {
	tk := NewClosableTicker(100 * time.Millisecond)
	go func() {
		for {
			c := <-tk.C
			if c != nil {
				fmt.Println("ticker")
			} else {
				fmt.Println("stop")
			}
		}
	}()
	time.Sleep(1 * time.Second)
	tk.Stop()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("end")
}
