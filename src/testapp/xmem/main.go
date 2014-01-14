package main

import (
	"boot"
	"esp/espnet"
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

	smux := espnet.NewServiceMux(nil, nil)
	smux.AddHandler(espnet.NewAddress("xmem"), xmemService.CreateHandleRequest())

	goService := espnet.NewGoService("SERVICE", smux.Serve)
	boot.QuickDefine(goService, "", true)

	managePoint := espnet.NewListenPoint("managePoint", nil, goService.AcceptESP)
	boot.QuickDefine(managePoint, "", true)

	boot.Go(cfile)
}
