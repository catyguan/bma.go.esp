package cacheserver

import (
	"esp/shell"
)

const (
	commandNameDeleteCache = "delete"
)

type cmdDeleteCache struct {
	service   *CacheService
	adminWord string
}

func (this *cmdDeleteCache) Name() string {
	return commandNameDeleteCache
}

func (this *cmdDeleteCache) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "cacheName"
	fs := shell.NewFlagSet(name)
	sure := ""
	fs.StringVar(&sure, "f", "", "delete confirm word")
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	cname := fs.Arg(0)

	match := false
	if this.adminWord != "" && sure == this.adminWord {
		match = true
	}
	if !match && shell.CheckConfirm(s, name, cname, sure) {
		match = true
	}
	if !match {
		word := shell.CreateConfirm(s, name, cname)
		s.Writeln("CONFIRM: " + name + " -f " + word + " " + cname)
		return true
	}

	err := this.service.DeleteCache(cname, true)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}
	s.Writeln("delete " + cname + " -> done")
	return true
}
