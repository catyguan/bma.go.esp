package tbus

import (
	"esp/shell"
	"fmt"
)

type dirModule struct {
	shell.ShellDirBase
	service *TBusService
}

func (this *dirModule) InitDir(s *TBusService) {
	this.service = s
	this.DirName = "module"
	this.Commands = this.MakeCommands
	this.Infos = this.MakeInfos
}

func (this *dirModule) MakeCommands() map[string]shell.ShellProcessor {
	r := make(map[string]shell.ShellProcessor)
	r["delete"] = this.CF(this.commandDelete)
	r["new"] = this.CF(this.commandNew)
	r["edit"] = this.CF(this.commandEdit)
	r["lookup"] = this.CF(this.commandLookup)
	return r
}

func (this *dirModule) MakeInfos() []string {
	this.service.lock.RLock()
	defer this.service.lock.RUnlock()
	r := make([]string, 0)
	for k, _ := range this.service.infos {
		r = append(r, shell.InfoName(k))
	}
	return r
}

func (this *dirModule) commandNew(s *shell.Session, command string) bool {
	name := "name"
	args := ""
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	info := new(ThriftServiceInfo)
	ed := shell.NewEditor()
	ed.Active(s, newDoc4Module(this.service, "", info), false)
	return true
}

func (this *dirModule) commandEdit(s *shell.Session, command string) bool {
	name := "edit"
	args := "moduleName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	mname := fs.Arg(0)

	info := this.service.GetServiceInfo(mname)
	if info == nil {
		s.Writeln("ERROR: module " + mname + " not exists")
		return true
	}

	ed := shell.NewEditor()
	ed.Active(s, newDoc4Module(this.service, mname, info), false)
	return true
}

func (this *dirModule) commandDelete(s *shell.Session, command string) bool {
	name := "delete"
	args := "moduleName"
	fs := shell.NewFlagSet(name)
	sure := ""
	fs.StringVar(&sure, "f", "", "delete confirm word")
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	mname := fs.Arg(0)

	match := shell.CheckConfirmWithAdminWord(s, name, mname, sure, this.service.config.AdminWord)
	if !match {
		word := shell.CreateConfirm(s, name, mname)
		s.Writeln("CONFIRM: " + name + " -f " + word + " " + mname)
		return true
	}

	err := this.service.DeleteService(mname)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}
	s.Writeln("delete " + mname + " -> done")
	return true
}

func (this *dirModule) commandLookup(s *shell.Session, command string) bool {
	name := "lookup"
	args := "thriftInvokeName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	tname := fs.Arg(0)

	module, method := SplitThriftName(tname)
	si, mi, err := this.service.FindServiceAndMethod(module, method)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}
	sin := ""
	min := ""
	if si != nil {
		sin = si.Name
	}
	if mi != nil {
		min = mi.Name
	}
	s.Write(fmt.Sprintf("lookup '%s' => %s, %s", tname, sin, min))
	if this.service.config.DefaultRemote == "" {
		s.Writeln("")
	} else {
		s.Writeln(fmt.Sprintf(", use DEFAULT '%s'", this.service.config.DefaultRemote))
	}
	return true
}
