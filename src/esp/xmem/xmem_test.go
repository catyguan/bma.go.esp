package xmem

import (
	"boot"
	"esp/sqlite"
	"fmt"
	"logger"
	"testing"
)

func func4tester(s *Service) func() {
	return func() {
		prof := new(MemGroupProfile)
		prof.Name = "test"
		prof.Coder = SimpleCoder(0)
		err := s.EnableMemGroup(prof)
		if err != nil {
			logger.Warn("test", "error - %s", err)
		}
	}
}

func TestXMem4Service(t *testing.T) {

	cfile := "../../../bin/config/xmem-config.json"

	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()

	xmemService := NewService("xmemService", sqliteServer)
	boot.QuickDefine(xmemService, "", true)

	f1 := func() {
		xm, err := xmemService.CreateXMem("test")
		if err != nil {
			logger.Error("test", "CreateXMem error - %s", err)
			return
		}
		_, _, b, err := xm.Get(MemKey{"a"})
		if !b && err == nil {
			fmt.Println("do init set")
			xm.Set(MemKey{"a"}, nil, 0)
			xm.Set(MemKey{"a", "b", "c"}, 123, 4)
			xm.Set(MemKey{"a", "b", "d"}, 234, 4)
			xm.Set(MemKey{"a", "e"}, 345, 4)
		}

		fmt.Println("----Dump----")
		str, err := xmemService.Dump("test", MemKey{}, true)
		fmt.Print(str)
	}
	if f1 != nil {
	}

	funl := []func(){
		func4tester(xmemService),
		f1,
	}

	boot.TestGo(cfile, 2, funl)
}
