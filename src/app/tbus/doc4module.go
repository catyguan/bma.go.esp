package tbus

import (
	"bytes"
	"errors"
	"esp/shell"
	"fmt"
	"uprop"
)

type doc4Module struct {
	service *TBusService
	name    string
	info    *ThriftServiceInfo
	edit    bool
}

func newDoc4Module(s *TBusService, name string, info *ThriftServiceInfo) *doc4Module {
	r := new(doc4Module)
	r.service = s
	r.name = name
	r.edit = name != ""
	r.info = info
	return r
}

func (this *doc4Module) Title() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("ThriftModule - ")
	if this.name == "" {
		buf.WriteString("*")
	} else {
		buf.WriteString(this.name)
	}
	return buf.String()
}

func (this *doc4Module) commands() map[string]func(s *shell.Session, cmd string) bool {
	r := make(map[string]func(s *shell.Session, cmd string) bool)
	r["set"] = this.commandEdit
	r["addm"] = this.commandAddMethod
	r["delm"] = this.commandDeleteMethod
	return r
}

func (this *doc4Module) ListCommands() []string {
	r := make([]string, 0)
	for k, _ := range this.commands() {
		r = append(r, k)
	}
	return r
}

func (this *doc4Module) HandleCommand(session *shell.Session, cmdline string) (bool, bool) {
	cmd := shell.CommandWord(cmdline)
	f := this.commands()[cmd]
	if f != nil {
		return true, f(session, cmdline)
	}
	return false, false
}

func (this *doc4Module) commandEdit(s *shell.Session, command string) bool {
	name := "set"
	args := "varname varval"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 2, 2) {
		return false
	}
	varn := fs.Arg(0)
	v := fs.Arg(1)
	prop := this.docProp()
	return shell.EditorHelper.DoPropEdit(s, prop, varn, v)
}

func (this *doc4Module) commandAddMethod(s *shell.Session, command string) bool {
	name := "addm"
	args := "methodName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return false
	}

	mname := fs.Arg(0)

	minfo := new(ThriftMethodInfo)
	minfo.Name = mname
	if this.info.Methods == nil {
		this.info.Methods = make(map[string]*ThriftMethodInfo)
	}
	this.info.Methods[mname] = minfo
	s.Writeln("set " + mname + " -> done")
	return true
}

func (this *doc4Module) commandDeleteMethod(s *shell.Session, command string) bool {
	name := "delm"
	args := "methodName"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return false
	}

	mname := fs.Arg(0)

	if this.info.Methods != nil {
		delete(this.info.Methods, mname)
	}
	s.Writeln("delete " + mname + " -> done")
	return true
}

func (this *doc4Module) OnCloseDoc(session *shell.Session) {

}
func (this *doc4Module) docProp() []*uprop.UProperty {
	r := make([]*uprop.UProperty, 0)
	r = append(r, uprop.NewUProperty("name", this.name, false, "module name", func(v string) error {
		if this.edit {
			return errors.New("can't edit name")
		}
		this.name = v
		return nil
	}))
	r = append(r, uprop.NewUProperty("remote", this.info.Remote, false, "remnote name", func(v string) error {
		this.info.Remote = v
		return nil
	}))
	return r
}
func (this *doc4Module) ShowDoc(s *shell.Session) {
	prop := this.docProp()
	shell.EditorHelper.DoPropList(s, "Module", prop)
	idx := 1
	s.Writeln("[ Methods ]")
	if this.info.Methods != nil {
		for _, m := range this.info.Methods {
			s.Writeln(fmt.Sprintf("%d: %s", idx, m.Name))
			idx++
		}
	}
	s.Writeln("[ Methods end ]")
}

func (this *doc4Module) CommitDoc(session *shell.Session) error {
	if this.name == "" {
		return errors.New("module name empty")
	}
	if this.info.Remote == "" {
		return errors.New("remote name empty")
	}
	err := this.service.SetService(this.name, this.info)
	if err != nil {
		return err
	}
	return this.service.save()
}
