package servproxy

import (
	"boot"
	"esp/goluaserv"
	"fmt"
	"golua"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	gls := goluaserv.NewService("goluaServ", func(gl *golua.GoLua) {
		golua.InitCoreLibs(gl)
	})
	boot.AddService(gls)

	s := NewService("test", gls)
	boot.AddService(s)

	f := func() {
		fmt.Printf("test done\n")
	}
	boot.TestGo("service_test.json", 5, []func(){f})
	time.Sleep(100 * time.Millisecond)
}
