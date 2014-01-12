package xmem

import (
	"esp/shell"
	"uprop"
)

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
	r["edit"] = this.CF(this.commandEdit)
	r["savecfg"] = this.CF(this.commandSave)
	r["saveall"] = this.CF(this.commandSaveAll)
	return r
}

func (this *dirService) MakeDirs() map[string]shell.ShellDir {
	r := make(map[string]shell.ShellDir, 0)
	clist, _ := this.service.ListMemGroupName()
	for _, k := range clist {
		g := new(dirGroup)
		g.InitDir(this.service, k)
		r[k] = g
	}
	return r
}

func (this *dirService) commandEdit(s *shell.Session, command string) bool {
	name := "edit"
	args := "memGroupName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	var ocfg *MemGroupConfig
	vname := fs.Arg(0)
	this.service.executor.DoSync("cmdEdit", func() error {
		item, err := this.service.doGetGroup(vname)
		if err != nil {
			return err
		}
		ocfg = item.config
		return nil
	})

	cfg := new(MemGroupConfig)
	uprop.Copy(cfg, ocfg)

	ed := shell.NewEditor()
	ed.Active(s, newDoc4MemGroup(this.service, vname, cfg), false)
	return true
}

func (this *dirService) commandSave(s *shell.Session, command string) bool {
	name := "savecfg"
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

func (this *dirService) commandSaveAll(s *shell.Session, command string) bool {
	name := "saveall"
	args := ""
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	err := this.service.Save()
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("save config done")
	}

	err = this.service.StoreAllMemGroup()
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("save all memgroup done")
	}

	return true
}
