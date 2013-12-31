package cacheserver

import (
	"esp/shell"

	"fmt"
)

const (
	commandNameDelete = "delete"
)

type cmdDelete struct {
	service *CacheService
	name    string
}

func (this *cmdDelete) Name() string {
	return commandNameDelete
}

func (this *cmdDelete) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "key"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	key := fs.Arg(0)

	done, err := this.service.Delete(this.name, key)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}

	s.Writeln(fmt.Sprintf("delete %v", done))
	return true
}
