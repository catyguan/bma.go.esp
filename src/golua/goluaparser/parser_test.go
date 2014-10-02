package goluaparser

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
		s = "a = 1.1"
		// s = "obj:print(1 + 2, true, a.b)"
		// s = "a.b = 1 + 2 - 3"
		// s = "function a(b, c) end"

		p := NewLuaParser1(s)
		node, err := p.Chunk()
		if err != nil {
			fmt.Println("ParseError", err)
		} else {
			node.dump("")
		}

	}
}
