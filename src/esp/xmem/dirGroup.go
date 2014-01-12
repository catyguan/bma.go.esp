package xmem

import "esp/shell"

type dirGroup struct {
	shell.ShellDirBase
	service *Service
	name    string
}

func (this *dirGroup) InitDir(s *Service, n string) {
	this.service = s
	this.name = n
	this.DirName = this.service.Name()
	this.Commands = this.MakeCommands
}

func (this *dirGroup) MakeCommands() map[string]shell.ShellProcessor {
	r := make(map[string]shell.ShellProcessor)
	r["delete"] = this.CF(this.commandDelete)
	r["new"] = this.CF(this.commandNew)
	r["edit"] = this.CF(this.commandEdit)
	r["save"] = this.CF(this.commandSave)
	return r
}

func (this *dirGroup) commandNew(s *shell.Session, command string) bool {
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

func (this *dirGroup) commandEdit(s *shell.Session, command string) bool {
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

func (this *dirGroup) commandDelete(s *shell.Session, command string) bool {
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

func (this *dirGroup) commandSave(s *shell.Session, command string) bool {
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
