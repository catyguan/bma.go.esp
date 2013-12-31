package mcpoint

import (
	"bmautil/valutil"
	"esp/shell"
)

const (
	commandNameMove = "move"
)

type cmdMove struct {
	service *MemcachePoint
}

func (this *cmdMove) Name() string {
	return commandNameMove
}

func (this *cmdMove) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "[up|down] [pos:int]"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 2, 2) {
		return true
	}

	up := false
	if fs.Arg(0) == "up" {
		up = true
	}
	pos := valutil.ToInt(fs.Arg(1), -1)
	if pos < 0 || pos >= len(this.service.router) {
		s.Writeln("ERROR: position invalid")
		return true
	}
	npos := pos + 1
	if up {
		npos = pos - 1
	}

	if npos < 0 || npos >= len(this.service.router) {
		s.Writeln("ERROR: can't move")
		return true
	}

	rlist := this.service.router
	rlist[pos], rlist[npos] = rlist[npos], rlist[pos]

	s.Writeln("move -> done")
	showList(s, rlist)
	return true
}
