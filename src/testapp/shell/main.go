package main

import (
	"bmautil/netutil"
	"boot"
	"esp/shell"
	"fmt"
	"telnetserver"
)

func main() {
	cfile := "../test/telnet-config.json"

	shl := shell.NewShell("app")
	boot.QuickDefine(shl, "", true)

	tServer := telnetserver.NewTelnetServer("telnetServer", func(ch *netutil.Channel, msg string) bool {
		fmt.Println("msg =>", msg)
		if msg == "close" {
			return false
		}
		ch.Write([]byte(msg + "\n"))
		return true
	})
	tServer.DefaultBoot(true)

	boot.Go(cfile)
}
