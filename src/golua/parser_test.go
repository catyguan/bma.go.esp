package golua

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func TestParser1(t *testing.T) {
	if true {
		safeCall()

		s := ""
		s = "obj:print(1 + 2, true, a.b)"
		// s = "abc = 1 + 2"

		p := NewLuaParser1(s)
		node, err := p.Chunk()
		if err != nil {
			fmt.Println("error", err)
		} else {
			node.dump("")
		}

	}
}
