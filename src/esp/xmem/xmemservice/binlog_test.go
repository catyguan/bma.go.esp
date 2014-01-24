package xmemservice

import (
	"bmautil/binlog"
	"boot"
	"esp/sqlite"
	"esp/xmem/xmemprot"
	"fmt"
	"logger"
	"testing"
)

func TestBinlog(t *testing.T) {

	cfile := "../../../bin/config/xmem-config.json"

	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()

	xmemService := NewService("xmemService", sqliteServer)
	boot.QuickDefine(xmemService, "", true)

	f1 := func() {
		cfg := new(MemGroupConfig)
		cfg.NoSave = true
		cfg.BLConfig = new(binlog.BinlogConfig)
		cfg.BLConfig.FileName = "test.blog"

		xmemService.UpdateMemGroupConfig("test", cfg)

		xm, err := xmemService.CreateXMem("test")
		if err != nil {
			logger.Error("test", "CreateXMem error - %s", err)
			return
		}
		fmt.Println("do set")
		xm.Set(xmemprot.MemKey{"a"}, nil, 0)
		xm.Set(xmemprot.MemKey{"a", "b", "c"}, 123, 4)
		xm.Set(xmemprot.MemKey{"a", "b", "d"}, 234, 4)
		xm.Set(xmemprot.MemKey{"a", "e"}, 345, 4)
		xm.Delete(xmemprot.MemKey{"a", "b"})

		fmt.Println("----Dump----")
		str, err := xmemService.Dump("test", xmemprot.MemKey{}, true)
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

func TestBinlogRun(t *testing.T) {

	cfile := "../../../bin/config/xmem-config.json"

	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()

	xmemService := NewService("xmemService", sqliteServer)
	boot.QuickDefine(xmemService, "", true)

	f1 := func() {
		xmemService.RunBinlog("test", "test.blog")

		fmt.Println("----Dump----")
		str, _ := xmemService.Dump("test", xmemprot.MemKey{}, true)
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
