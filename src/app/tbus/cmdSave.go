package tbus

import (
	"esp/shell"
)

const (
	commandNameSave = "save"
)

type cmdSave struct {
	service *TBusService
}

func (this *cmdSave) Name() string {
	return commandNameSave
}

func (this *cmdSave) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := ""
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	err := this.service.save()
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("save done")
	}

	return true
}
