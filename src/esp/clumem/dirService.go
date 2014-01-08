package clumem

import "esp/shell"

type dirService struct {
	shell.ShellDirBase
	service *Service
}

func (this *dirService) InitDir(s *Service) {
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
	return r
}

func (this *dirService) MakeDirs() map[string]shell.ShellDir {
	r := make(map[string]shell.ShellDir, 0)
	// clist := this.service.ListCacheName()
	// for _, k := range clist {
	// 	cache, _ := this.service.GetCache(k, false)
	// 	if cache != nil {
	// 		ss, ok := cache.(shell.ShellDirSupported)
	// 		if ok {
	// 			r[k] = ss.CreateShell()
	// 		}
	// 	}
	// }
	return r
}

func (this *dirService) commandNew(s *shell.Session, command string) bool {
	name := "new"
	args := ""
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	cfg := new(MemGroupConfig)
	ed := shell.NewEditor()
	ed.Active(s, newDoc4Cache(this.service, "", cfg), false)
	return true
}

func (this *dirService) commandEdit(s *shell.Session, command string) bool {
	name := "edit"
	args := "cacheName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	// cname := fs.Arg(0)
	// cache, err := this.service.GetCache(cname, false)
	// if err != nil {
	// 	s.Writeln("ERROR: " + err.Error())
	// 	return true
	// }

	// ctype := cache.Type()
	// fac := GetCacheFactory(ctype)
	// if fac == nil {
	// 	s.Writeln(fmt.Sprintf("ERROR : CacheType[%s] not exists", ctype))
	// 	return true
	// }
	// cfg := fac.CreateConfig()
	// ocfg := cache.GetConfig()
	// uprop.Copy(cfg, ocfg)

	// ed := shell.NewEditor()
	// ed.Active(s, newDoc4Cache(this.service, cname, ctype, cfg), false)
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

	vname := fs.Arg(0)

	match := shell.CheckConfirmWithAdminWord(s, name, vname, sure, this.service.config.AdminWord)
	if !match {
		word := shell.CreateConfirm(s, name, vname)
		s.Writeln("CONFIRM: " + name + " -f " + word + " " + vname)
		return true
	}

	// err := this.service.DeleteCache(vname, true)
	// if err != nil {
	// 	s.Writeln("ERROR: " + err.Error())
	// 	return true
	// }
	// s.Writeln("delete " + vname + " -> done")
	// this.service.save()
	return true

}

func (this *dirService) commandSave(s *shell.Session, command string) bool {
	name := "save"
	args := ""
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	err := this.service.Save()
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("save done")
	}

	return true
}
