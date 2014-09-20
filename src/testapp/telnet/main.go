package main

import (
	"boot"
	"net"
	"telnetserver"
)

func main() {
	cfile := "../test/telnet-config.json"

	// shl := shell.NewShell("app")
	// boot.QuickDefine(shl, "", true)

	// tServer := telnetserver.NewTelnetServer("telnetServer", telnetcmd.NewHandler(shl))
	// tServer.DefaultBoot(true)
	// shl.AddDir(telnetcmd.NewShellDir(tServer))
	tServer := telnetserver.NewService("telnetServer", func(c net.Conn, msg string) bool {
		c.Write([]byte(msg))
		c.Write([]byte{'\n'})
		return true
	})
	boot.AddService(tServer)

	boot.Go(cfile)
}
