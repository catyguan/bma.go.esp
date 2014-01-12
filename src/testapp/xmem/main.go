package main

import (
	"boot"
	"esp/shell"
	"esp/shell/telnetcmd"
	"esp/sqlite"
	"esp/xmem"
	"telnetserver"
)

func main() {
	cfile := "config/xmem-config.json"

	shl := shell.NewShell("clumem")
	boot.Install("shell", shl)

	// sqliteServer
	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()
	shl.AddDir(sqliteServer.NewShellDir())

	// xmemServer
	xmemService := xmem.NewService("xmemService", sqliteServer)
	boot.QuickDefine(xmemService, "", true)
	shl.AddDir(xmemService.NewShellDir())

	tester := new(Tester)
	tester.xmems = xmemService
	boot.QuickDefine(tester, "", false)

	// telnetServer
	tServer := telnetserver.NewTelnetServer("telnetServer", telnetcmd.NewHandler(shl))
	tServer.DefaultBoot(true)
	shl.AddDir(telnetcmd.NewShellDir(tServer))

	// pointSER := espnet.NewListenPoint("tbusPoint", nil, tbusService.OnSocketAccept)
	// boot.QuickDefine(pointSER, "", true)

	boot.Go(cfile)
}
