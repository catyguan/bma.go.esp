package telnetcmd

import (
	"bmautil/netutil"
	"esp/shell"
	"esp/shell/netcmd"
	"telnetserver"
)

func NewHandler(shl *shell.Shell) func(ch *netutil.Channel, msg string) bool {
	return func(ch *netutil.Channel, msg string) bool {
		var session *shell.Session
		if o, ok := ch.Properties["@shell"]; ok {
			session = o.(*shell.Session)
		} else {
			session = shell.NewSession(shell.NewNetWriter(ch))
			session.Vars["@WHO"] = ch.RemoteAddr().String()
			ch.Properties["@shell"] = session
		}
		shl.Process(session, msg)
		return true
	}
}

func NewShellDir(server *telnetserver.TelnetServer) shell.ShellDir {
	r := shell.NewShellDirCommon(server.Name())
	cmd := netcmd.NewIpLimitCommand()
	cmd.GetWhiteList = func() []string {
		return server.WhiteList
	}
	cmd.SetWhiteList = func(list []string) {
		server.WhiteList = list
	}
	cmd.GetBlackList = func() []string {
		return server.BlackList
	}
	cmd.SetBlackList = func(list []string) {
		server.BlackList = list
	}
	r.AddCommand(cmd)
	return r
}
