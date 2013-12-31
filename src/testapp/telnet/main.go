package main

import (
	"boot"
	"esp/shell"
	"esp/shell/telnetcmd"
	"telnetserver"
)

func main() {
	cfile := "../test/telnet-config.json"

	shl := shell.NewShell("app")
	boot.QuickDefine(shl, "", true)

	tServer := telnetserver.NewTelnetServer("telnetServer", telnetcmd.NewHandler(shl))
	tServer.DefaultBoot(true)
	shl.AddDir(telnetcmd.NewShellDir(tServer))

	boot.Go(cfile)
}
