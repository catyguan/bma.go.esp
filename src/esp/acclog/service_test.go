package acclog

import (
	"boot"
	"fmt"
	"testing"
)

func TestServiceBoot(t *testing.T) {
	s := NewService("test")
	boot.Add(s, "", false)

	f := func() {
		s.Write("abc", NewSimpleLog("hello world"))
		s.Write("abc", NewSimpleLog("hello world"))
		fmt.Printf("test done\n")
	}
	boot.TestGo("service_test.json", 3, []func(){f})
}
