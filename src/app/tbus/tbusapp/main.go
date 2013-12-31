package main

import (
	"app/tbus"
	"boot"
	"esp/espnet"
	"esp/shell"
	"esp/shell/telnetcmd"
	"esp/sqlite"
	"telnetserver"
)

func main() {
	// mymysql
	// godrv.Register("set names utf8")
	// tcp:172.19.16.195:3306*db_live2/root/root
	// mysql
	// root:root@tcp(172.19.16.195:3306)/db_live2

	cfile := "config/tbus-config.json"

	shl := shell.NewShell("tbusapp")
	boot.Install("shell", shl)

	// namedsql
	// namedSQL := namedsql.NewNamedSQL("namedSQL")
	// namedSQL.DefaultBoot()
	// shl.AddDir(namedSQL.NewShellDir())

	// sqliteServer
	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()
	shl.AddDir(sqliteServer.NewShellDir())

	// TBusServer
	tbusService := tbus.NewTBusService("tbusService", sqliteServer)
	boot.QuickDefine(tbusService, "", true)
	shl.AddDir(tbusService.NewShellDir())

	// telnetServer
	tServer := telnetserver.NewTelnetServer("telnetServer", telnetcmd.NewHandler(shl))
	tServer.DefaultBoot(true)
	shl.AddDir(telnetcmd.NewShellDir(tServer))

	pointSER := espnet.NewListenPoint("tbusPoint", nil, tbusService.OnSocketAccept)
	boot.QuickDefine(pointSER, "", true)

	boot.Go(cfile)
}
