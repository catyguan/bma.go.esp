package mcpoint

import (
	"esp/shell"
	"fmt"
)

const (
	commandNameList = "list"
)

type cmdList struct {
	service *MemcachePoint
}

func (this *cmdList) Name() string {
	return commandNameList
}

func (this *cmdList) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := ""
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	rlist := this.service.router
	showList(s, rlist)
	return true
}

func showList(s *shell.Session, rlist []*cacheRouter) {
	s.Writeln("[ router ]")
	for i, cr := range rlist {
		s.Writeln(fmt.Sprintf("%d: %s <- %s", i, cr.group, cr.pattern))
	}
	s.Writeln("[ router end ]")
}
