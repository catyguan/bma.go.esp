package tankbat

import (
	"fmt"
	"testing"
	"time"
)

func TestMatrix(t *testing.T) {
	m := NewMatrix(nil, 32)
	m.Run()

	time.Sleep(5 * time.Second)
	m.AskClose()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("END")
}
