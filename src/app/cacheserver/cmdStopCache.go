package cacheserver

import (
	"esp/shell"
)

const (
	commandNameStopCache = "stop"
)

type cmdStopCache struct {
	service *CacheService
}

func (this *cmdStopCache) Name() string {
	return commandNameStopCache
}

func (this *cmdStopCache) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "cacheName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	cname := fs.Arg(0)

	err := this.service.StopCache(cname)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}
	s.Writeln("stop " + cname + " -> done")
	return true
}
