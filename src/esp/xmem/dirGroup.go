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
	r["save"] = this.CF(this.commandSave)
	r["load"] = this.CF(this.commandLoad)
	r["dump"] = this.CF(this.commandDump)
	return r
}

func (this *dirGroup) commandSave(s *shell.Session, command string) bool {
	name := "save"
	args := ""
	fileName := ""
	fs := shell.NewFlagSet(name)
	fs.StringVar(&fileName, "f", "", "file to store snapshot")
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	err := this.service.SaveMemGroup(this.name, fileName)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("save done")
	}

	return true
}

func (this *dirGroup) commandLoad(s *shell.Session, command string) bool {
	name := "load"
	args := ""
	fileName := ""
	fs := shell.NewFlagSet(name)
	fs.StringVar(&fileName, "f", "", "file to load snapshot")
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	err := this.service.LoadMemGroup(this.name, fileName)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("load done")
	}

	return true
}

func (this *dirGroup) commandDump(s *shell.Session, command string) bool {
	name := "dump"
	args := "[key]"
	all := false
	fs := shell.NewFlagSet(name)
	fs.BoolVar(&all, "a", false, "dump all sub items")
	if shell.DoParse(s, command, fs, name, args, 0, 1) {
		return true
	}

	key := ""
	if fs.NArg() > 0 {
		key = fs.Arg(0)
	}

	k := MemKeyFromString(key)
	str, err := this.service.Dump(this.name, k, all)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("dump done >>")
		s.Writeln(str)
	}

	return true
}
