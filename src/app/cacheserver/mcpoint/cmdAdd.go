package mcpoint

import (
	"esp/shell"
	"regexp"
)

const (
	commandNameAdd = "add"
)

type cmdAdd struct {
	service *MemcachePoint
}

func (this *cmdAdd) Name() string {
	return commandNameAdd
}

func (this *cmdAdd) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "[group] [pattern]"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 2, 2) {
		return true
	}

	pat := fs.Arg(1)
	gn := fs.Arg(0)

	cr := new(cacheRouter)
	var err error
	cr.matcher, err = regexp.Compile(pat)
	if err != nil {
		s.Writeln("ERROR: compile '" + pat + "' fail - %s" + err.Error())
		return true
	}
	cr.pattern = pat
	cr.group = gn

	this.service.router = append(this.service.router, cr)

	s.Writeln("add " + gn + "," + pat + " -> done")
	showList(s, this.service.router)
	return true
}
