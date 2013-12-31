package netcmd

import (
	"esp/shell"
	"fmt"
	"strings"
	"testing"
)

func TestIpLimit(t *testing.T) {
	sh := shell.NewShell()

	white := ""
	black := "127.0.0.1"

	cmd := NewIpLimitCommand()
	cmd.GetWhiteList = func() []string {
		return strings.Split(white, ",")
	}
	cmd.SetWhiteList = func(list []string) {
		white = strings.Join(list, ",")
	}
	cmd.GetBlackList = func() []string {
		return strings.Split(black, ",")
	}
	cmd.SetBlackList = func(list []string) {
		black = strings.Join(list, ",")
	}
	sh.AddCommand(cmd)

	session := shell.NewSession(shell.NewConsoleWriter())
	session.Vars["@WHO"] = "TestCase"

	fmt.Println(">>", "B:", black, "W:", white)
	sh.Process(session, "iplimit")

	sh.Process(session, "iplimit -a 168.0.0.1")
	sh.Process(session, "iplimit -b -a 202.0.0.1")
	fmt.Println(">>", "B:", black, "W:", white)

	sh.Process(session, "iplimit -r 168.0.0.1")
	sh.Process(session, "iplimit -b -r 127.0.0.1")
	fmt.Println(">>", "B:", black, "W:", white)

	session.Close()
}
