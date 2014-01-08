package main

import (
	"boot"
	"esp/clumem"
	"esp/shell"
	"esp/shell/telnetcmd"
	"esp/sqlite"
	"telnetserver"
)

func main() {
	cfile := "config/clumem-config.json"

	shl := shell.NewShell("clumem")
	boot.Install("shell", shl)

	// sqliteServer
	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()
	shl.AddDir(sqliteServer.NewShellDir())

	// TBusServer
	clumemService := clumem.NewService("clumemService", sqliteServer)
	boot.QuickDefine(clumemService, "", true)
	shl.AddDir(clumemService.NewShellDir())

	// telnetServer
	tServer := telnetserver.NewTelnetServer("telnetServer", telnetcmd.NewHandler(shl))
	tServer.DefaultBoot(true)
	shl.AddDir(telnetcmd.NewShellDir(tServer))

	// pointSER := espnet.NewListenPoint("tbusPoint", nil, tbusService.OnSocketAccept)
	// boot.QuickDefine(pointSER, "", true)

	boot.Go(cfile)
}
