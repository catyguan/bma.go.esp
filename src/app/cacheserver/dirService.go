package cacheserver

import (
	"esp/shell"
	"fmt"
	"uprop"
)

type dirService struct {
	shell.ShellDirBase
	service *CacheService
}

func (this *dirService) InitDir(s *CacheService) {
	this.service = s
	this.DirName = this.service.Name()
	this.Commands = this.MakeCommands
	this.Dirs = this.MakeDirs
}

func (this *dirService) MakeCommands() map[string]shell.ShellProcessor {
	r := make(map[string]shell.ShellProcessor)
	r["delete"] = this.CF(this.commandDelete)
	r["new"] = this.CF(this.commandNew)
	r["edit"] = this.CF(this.commandEdit)
	r["save"] = this.CF(this.commandSave)
	r["start"] = this.CF(this.commandStart)
	r["stop"] = this.CF(this.commandStop)
	return r
}

func (this *dirService) MakeDirs() map[string]shell.ShellDir {
	clist := this.service.ListCacheName()
	r := make(map[string]shell.ShellDir, 0)
	for _, k := range clist {
		cache, _ := this.service.GetCache(k, false)
		if cache != nil {
			ss, ok := cache.(shell.ShellDirSupported)
			if ok {
				r[k] = ss.CreateShell()
			}
		}
	}
	return r
}

func (this *dirService) commandNew(s *shell.Session, command string) bool {
	name := "new"
	args := "[cacheType]"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 1) {
		return true
	}

	if fs.NArg() < 1 {
		s.Writeln("INPUT: new " + args)
		this.showTypes(s)
		return true
	}

	ctype := fs.Arg(0)
	fac := GetCacheFactory(ctype)
	if fac == nil {
		s.Writeln(fmt.Sprintf("ERROR : CacheType[%s] not exists", ctype))
		return true
	}

	cfg := fac.CreateConfig()
	ed := shell.NewEditor()
	ed.Active(s, newDoc4Cache(this.service, "", ctype, cfg), false)
	return true
}

func (this *dirService) showTypes(s *shell.Session) {
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

func (this *dirService) commandEdit(s *shell.Session, command string) bool {
	name := "edit"
	args := "cacheName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	cname := fs.Arg(0)
	cache, err := this.service.GetCache(cname, false)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}

	ctype := cache.Type()
	fac := GetCacheFactory(ctype)
	if fac == nil {
		s.Writeln(fmt.Sprintf("ERROR : CacheType[%s] not exists", ctype))
		return true
	}
	cfg := fac.CreateConfig()
	ocfg := cache.GetConfig()
	uprop.Copy(cfg, ocfg)

	ed := shell.NewEditor()
	ed.Active(s, newDoc4Cache(this.service, cname, ctype, cfg), false)
	return true
}

func (this *dirService) commandDelete(s *shell.Session, command string) bool {
	name := "delete"
	args := "remoteName"
	fs := shell.NewFlagSet(name)
	sure := ""
	fs.StringVar(&sure, "f", "", "delete confirm word")
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	cname := fs.Arg(0)

	match := shell.CheckConfirmWithAdminWord(s, name, cname, sure, this.service.config.AdminWord)
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
	this.service.save()
	return true

}

func (this *dirService) commandSave(s *shell.Session, command string) bool {
	name := "save"
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

func (this *dirService) commandStart(s *shell.Session, command string) bool {
	name := "start"
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

func (this *dirService) commandStop(s *shell.Session, command string) bool {
	name := "stop"
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
