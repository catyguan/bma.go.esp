package mcservercmd

import (
	"esp/shell"
	"esp/shell/netcmd"
	"mcserver"
)

func NewShellDir(server *mcserver.MemcacheServer) shell.ShellDir {
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
