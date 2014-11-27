package memserv4session

import (
	"esp/memserv"
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

func TestSession(t *testing.T) {
	safeCall()

	s := memserv.NewMemoryServ()
	defer s.CloseAll(true)

	sid := "1234"
	var v interface{}
	var err error

	SetSession(s, sid, "test", 1, 1000)
	MSetSession(s, sid, map[string]interface{}{"abc": true}, 1000)

	v, err = GetSession(s, sid, "test", time.Now())
	fmt.Println("GetSesion", v, err)

	time.Sleep(100 * time.Millisecond)
}
