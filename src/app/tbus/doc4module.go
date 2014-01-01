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

func (this *doc4Module) OnCloseDoc(session *shell.Session) {

}
func (this *doc4Module) GetUProperties() []*uprop.UProperty {
	b := new(uprop.UPropertyBuilder)
	b.NewProp("name", "module name").Optional(false).BeValue(this.name, func(v string) error {
		if this.edit {
			return errors.New("can't edit name")
		}
		this.name = v
		return nil
	})
	b.NewProp("remote", "remote name").Optional(false).BeValue(this.info.Remote, func(v string) error {
		this.info.Remote = v
		return nil
	})
	mlist := b.NewProp("ms", "methods").BeList(this.addMethod, this.removeMethod)
	if this.info.Methods != nil {
		for _, m := range this.info.Methods {
			func(m *ThriftMethodInfo) {
				mlist.Add(m.Name, func(v string) error {
					delete(this.info.Methods, m.Name)
					m.Name = v
					this.info.Methods[v] = m
					return nil
				})
			}(m)
		}
	}
	return b.AsList()
}

func (this *doc4Module) addMethod(ns []string) error {
	for _, mname := range ns {
		minfo := new(ThriftMethodInfo)
		minfo.Name = mname
		if this.info.Methods == nil {
			this.info.Methods = make(map[string]*ThriftMethodInfo)
		}
		this.info.Methods[mname] = minfo
	}
	return nil
}

func (this *doc4Module) removeMethod(ns []string) error {
	if this.info.Methods != nil {
		for _, n := range ns {
			if _, ok := this.info.Methods[n]; ok {
				delete(this.info.Methods, n)
			} else {
				return fmt.Errorf("unknow method name '%s'", n)
			}
		}
	}
	return nil
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
