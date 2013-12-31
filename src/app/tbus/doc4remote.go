package tbus

import (
	"bmautil/valutil"
	"boot"
	"bytes"
	"errors"
	"esp/espnet/cfprototype"
	"esp/shell"
	"fmt"
	"time"
	"uprop"
)

type doc4Remote struct {
	service *TBusService
	name    string
	kind    string
	info    cfprototype.ChannelFactoryPrototype
	edit    bool
}

func newDoc4Remote(s *TBusService, name, kind string, info cfprototype.ChannelFactoryPrototype) *doc4Remote {
	r := new(doc4Remote)
	r.service = s
	r.name = name
	r.edit = name != ""
	r.kind = kind
	r.info = info
	return r
}

func (this *doc4Remote) Title() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("Remote - ")
	if this.name == "" {
		buf.WriteString("*")
	} else {
		buf.WriteString(this.name)
	}
	buf.WriteString("/")
	buf.WriteString(this.kind)
	return buf.String()
}

func (this *doc4Remote) commands() map[string]func(s *shell.Session, cmd string) bool {
	r := make(map[string]func(s *shell.Session, cmd string) bool)
	r["set"] = this.commandEdit
	r["test"] = this.commandTest
	return r
}

func (this *doc4Remote) ListCommands() []string {
	r := make([]string, 0)
	for k, _ := range this.commands() {
		r = append(r, k)
	}
	return r
}

func (this *doc4Remote) HandleCommand(session *shell.Session, cmdline string) (bool, bool) {
	cmd := shell.CommandWord(cmdline)
	f := this.commands()[cmd]
	if f != nil {
		return true, f(session, cmdline)
	}
	return false, false
}

func (this *doc4Remote) commandEdit(s *shell.Session, command string) bool {
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

func (this *doc4Remote) commandTest(s *shell.Session, command string) bool {
	name := "test"
	args := "[timeoutSec]"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 1) {
		return false
	}

	to := 5
	if fs.NArg() > 0 {
		to = valutil.ToInt(fs.Arg(0), to)
	}
	cf, err := this.info.CreateChannelFactory(this.service, "test", true)
	if err != nil {
		s.Writeln(fmt.Sprintf("ERROR: %s", err))
		return false
	}
	defer func() {
		boot.RuntimeStopCloseClean(cf, false)
	}()

	ch := make(chan error, 1)
	tm := time.NewTimer(time.Duration(to) * time.Second)
	defer tm.Stop()

	go func() {
		c, err := cf.NewChannel()
		if c != nil {
			c.AskClose()
		}
		ch <- err
	}()

	select {
	case err2 := <-ch:
		if err2 != nil {
			s.Writeln(fmt.Sprintf("ERROR: %s", err2))
		} else {
			s.Writeln("test success")
		}
	case <-tm.C:
		s.Writeln("ERROR: test timeout")
	}
	return false
}

func (this *doc4Remote) OnCloseDoc(session *shell.Session) {

}
func (this *doc4Remote) docProp() []*uprop.UProperty {
	p := this.info.GetProperties()
	r := make([]*uprop.UProperty, 2+len(p))
	r[0] = uprop.NewUProperty("name", this.name, false, "module name", func(v string) error {
		if this.edit {
			return errors.New("can't edit name")
		}
		this.name = v
		return nil
	})
	r[1] = uprop.NewUProperty("kind", this.kind, true, "remote kind", func(v string) error {
		return errors.New("can't edit kind")
	})
	copy(r[2:], p)
	return r
}
func (this *doc4Remote) ShowDoc(s *shell.Session) {
	prop := this.docProp()
	shell.EditorHelper.DoPropList(s, "Remote", prop)
}

func (this *doc4Remote) CommitDoc(session *shell.Session) error {
	if this.name == "" {
		return errors.New("remote name empty")
	}
	if err := this.info.Valid(); err != nil {
		return err
	}
	err := this.service.SetRemote(this.name, this.kind, this.info)
	if err != nil {
		return err
	}
	return this.service.save()
}
