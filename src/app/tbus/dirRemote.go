package tbus

import (
	"bytes"
	"esp/shell"
	"fmt"
	"uprop"
)

type dirRemote struct {
	shell.ShellDirBase
	service *TBusService
}

func (this *dirRemote) InitDir(s *TBusService) {
	this.service = s
	this.DirName = "remote"
	this.Commands = this.MakeCommands
	this.Infos = this.MakeInfos
}

func (this *dirRemote) MakeCommands() map[string]shell.ShellProcessor {
	r := make(map[string]shell.ShellProcessor)
	r["delete"] = this.CF(this.commandDelete)
	r["new"] = this.CF(this.commandNew)
	r["edit"] = this.CF(this.commandEdit)
	return r
}

func (this *dirRemote) MakeInfos() []string {
	this.service.lock.RLock()
	defer this.service.lock.RUnlock()
	r := make([]string, 0)
	for k, _ := range this.service.remotes {
		r = append(r, shell.InfoName(k))
	}
	return r
}

func (this *dirRemote) commandNew(s *shell.Session, command string) bool {
	name := "new"
	args := "remoteKind"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 1) {
		return true
	}

	if fs.NArg() < 1 {
		this.showRemoteKinds(s)
		return true
	}

	rtype := fs.Arg(0)
	p := CreateChannelFactoryPrototype(rtype)
	if p == nil {
		s.Writeln(fmt.Sprintf("ERROR: unknow remote type '%s'", rtype))
		return true
	}
	ed := shell.NewEditor()
	ed.Active(s, newDoc4Remote(this.service, "", rtype, p), false)
	return true
}

func (this *dirRemote) showRemoteKinds(s *shell.Session) {
	s.Write("remote types : ")
	buf := bytes.NewBuffer(make([]byte, 0))
	c := 0
	for _, k := range ListChannelFactoryPrototype() {
		if c > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(k)
		c++
	}
	s.Writeln(buf.String())
}

func (this *dirRemote) commandEdit(s *shell.Session, command string) bool {
	name := "edit"
	args := "remoteName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	rname := fs.Arg(0)

	kind, info := this.service.GetRemotePrototype(rname)
	if info == nil {
		s.Writeln("ERROR: module " + rname + " not exists")
		return true
	}

	p := CreateChannelFactoryPrototype(kind)
	uprop.Copy(p, info)

	ed := shell.NewEditor()
	ed.Active(s, newDoc4Remote(this.service, rname, kind, p), false)
	return true
}

func (this *dirRemote) commandDelete(s *shell.Session, command string) bool {
	name := "delete"
	args := "remoteName"
	fs := shell.NewFlagSet(name)
	sure := ""
	fs.StringVar(&sure, "f", "", "delete confirm word")
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	rname := fs.Arg(0)

	match := shell.CheckConfirmWithAdminWord(s, name, rname, sure, this.service.config.AdminWord)
	if !match {
		word := shell.CreateConfirm(s, name, rname)
		s.Writeln("CONFIRM: " + name + " -f " + word + " " + rname)
		return true
	}

	err := this.service.DeleteRemote(rname)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}
	s.Writeln("delete " + rname + " -> done")
	return true
}
