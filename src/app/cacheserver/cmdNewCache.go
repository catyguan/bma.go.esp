package cacheserver

import (
	"esp/shell"
)

const (
	commandNameNewCache = "new"
)

type cmdNewCache struct {
	service *CacheService
}

func (this *cmdNewCache) Name() string {
	return commandNameNewCache
}

func (this *cmdNewCache) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "[cacheName] [cacheType]"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 2) {
		return true
	}

	if fs.NArg() < 2 {
		this.showTypes(s)
		return true
	}
	cname := fs.Arg(0)
	ctype := fs.Arg(1)
	_, err := this.service.CreateCache(cname, ctype)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("create " + cname + "," + ctype + " -> done")
	}
	return true
}

func (this *cmdNewCache) showTypes(s *shell.Session) {
	s.Write("cache types : ")
	c := 0
	for k, _ := range factorties {
		if c > 0 {
			s.Write(", ")
		}
		s.Write(k)
		c++
	}
	s.Writeln("")
}
