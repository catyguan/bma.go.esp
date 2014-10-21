package namedsql

import (
	"boot"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func safeCall() {
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func TestServiceBoot(t *testing.T) {
	// safeCall()

	s := NewService("test")
	boot.AddService(s)

	f := func() {
		// s.Write("abc", NewSimpleLog("hello world"))
		// s.Write("abc", NewSimpleLog("hello world"))
		s.Get("test")
		fmt.Printf("i'm here\n")
	}
	boot.TestGo("service_test.json", 1, []func(){f})
}
