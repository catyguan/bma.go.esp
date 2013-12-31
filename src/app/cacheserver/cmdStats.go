package cacheserver

import (
	"esp/shell"
)

const (
	commandNameStats = "stats"
)

type cmdStats struct {
	service *CacheService
	name    string
}

func (this *cmdStats) Name() string {
	return commandNameStats
}

func (this *cmdStats) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := ""
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	info, err := this.service.QueryStats(this.name)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}
	s.Writeln(this.name + " -> ")
	s.Writeln(info)

	return true
}
