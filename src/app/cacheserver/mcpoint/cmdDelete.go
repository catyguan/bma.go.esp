package mcpoint

import (
	"bmautil/valutil"
	"esp/shell"
)

const (
	commandNameDelete = "delete"
)

type cmdDelete struct {
	service *MemcachePoint
}

func (this *cmdDelete) Name() string {
	return commandNameDelete
}

func (this *cmdDelete) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "[pos:int]"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	pos := valutil.ToInt(fs.Arg(0), -1)
	if pos < 0 || pos >= len(this.service.router) {
		s.Writeln("ERROR: position invalid")
		return true
	}

	rlist := this.service.router
	cr := rlist[pos]
	l := len(rlist)
	if l > 1 {
		rlist[pos] = nil
		rlist[pos], rlist[l-1] = rlist[l-1], rlist[pos]
		rlist = rlist[:l-1]
	} else {
		rlist = make([]*cacheRouter, 0)
	}
	this.service.router = rlist

	s.Writeln("delete " + cr.pattern + "," + cr.group + " -> done")
	showList(s, rlist)
	return true
}
