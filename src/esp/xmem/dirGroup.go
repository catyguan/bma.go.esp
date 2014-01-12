package xmem

import (
	"bytes"
	"esp/shell"
	"fmt"
)

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

	str := ""
	err := this.service.executor.DoSync("dump", func() error {
		item, err := this.service.doGetGroup(this.name)
		if err != nil {
			return err
		}
		k := MemKeyFromString(key)
		it, ok := item.group.Get(k)
		if !ok {
			return fmt.Errorf("<%s> not exists", key)
		}
		buf := bytes.NewBuffer([]byte{})
		it.Dump(key, buf, 0, all)
		str = item.group.String() + "\n" + buf.String()
		return nil
	})
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("dump done >>")
		s.Writeln(str)
	}

	return true
}
