package main

import (
	"boot"
	"esp/shell"
	"esp/shell/telnetcmd"
	"esp/sqlite"
	"telnetserver"
)

func main() {
	cfile := "../test/sqlite-config.json"

	shl := shell.NewShell()

	sqlstr := "create table foo (id integer not null primary key, name text);"
	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.AddInit(sqlite.InitTable("", "foo", sqlstr))
	sqliteServer.DefaultBoot()

	shl.AddDir(sqliteServer.NewShellDir())

	tServer := telnetserver.NewTelnetServer("telnetServer", telnetcmd.NewHandler(shl))
	tServer.DefaultBoot(true)
	shl.AddDir(telnetcmd.NewTelnetDir(tServer))

	boot.Go(cfile)
}
