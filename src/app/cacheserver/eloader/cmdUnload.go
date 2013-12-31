package eloader

import (
	"esp/shell"
)

const (
	commandNameUnload = "unload"
)

type cmdUnload struct {
	service *LoaderCache
}

func (this *cmdUnload) Name() string {
	return commandNameUnload
}

func (this *cmdUnload) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "loaderName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	lname := fs.Arg(0)

	err := this.service.RemoveLoader(lname)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}
	s.Writeln("unload " + lname + " -> done")
	return true
}
