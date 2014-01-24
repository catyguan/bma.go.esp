package xmemservice

import (
	"boot"
	"esp/sqlite"
	"testing"
)

func TestBinlogSlave(t *testing.T) {

	cfile := "../../../../bin/config/xmem2-config.json"

	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()

	xmemService := NewService("xmemService", sqliteServer)
	boot.QuickDefine(xmemService, "", true)

	f1 := func() {
		cfg := new(MemGroupConfig)
		cfg.NoSave = true
		cfg.BLSlaveConfig = new(BLSlaveConfig)
		cfg.BLSlaveConfig.Address = "127.0.0.1:8080"
		cfg.BLSlaveConfig.TimeoutMS = 1000

		xmemService.UpdateMemGroupConfig("test", cfg)
	}
	if f1 != nil {
	}

	funl := []func(){
		func4tester(xmemService),
		f1,
	}

	boot.TestGo(cfile, 3, funl)
}
