package cacheserver

import (
	"esp/shell"
)

const (
	commandNameStartCache = "start"
)

type cmdStartCache struct {
	service *CacheService
}

func (this *cmdStartCache) Name() string {
	return commandNameStartCache
}

func (this *cmdStartCache) Process(s *shell.Session, command string) bool {
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

	err := this.service.StartCache(cname)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}
	s.Writeln("start " + cname + " -> done")
	return true
}
