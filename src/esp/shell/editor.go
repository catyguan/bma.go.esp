package shell

import (
	"fmt"
	"logger"
	"strings"
	"uprop"
)

type EditorDoc interface {
	Title() string
	ListCommands() []string
	HandleCommand(session *Session, cmdline string) (bool, bool)
	GetUProperties() []*uprop.UProperty

	OnCloseDoc(session *Session)
	CommitDoc(session *Session) error
}

type Editor struct {
	parent   *Editor
	doc      EditorDoc
	selStack []string
}

func NewEditor() *Editor {
	this := new(Editor)
	this.selStack = make([]string, 0)
	return this
}

func (this *Editor) Active(s *Session, doc EditorDoc, stack bool) bool {
	if o := s.Vars["@EDITOR"]; o != nil {
		ed := o.(*Editor)
		if !stack {
			s.Writeln(fmt.Sprintf("ERROR: other editor(%s) actived", ed.doc.Title()))
			return false
		}
		this.parent = o.(*Editor)
	}
	s.Vars["@EDITOR"] = this
	this.doc = doc
	this.Show(s)
	return true
}

func (this *Editor) Close(s *Session) {
	this.doc.OnCloseDoc(s)
	if s.Vars["@EDITOR"] != this {
		s.Writeln("ERROR: close not active editor")
		return
	}
	if this.parent == nil {
		s.Writeln(fmt.Sprintf("editor(%s) closed", this.doc.Title()))
		delete(s.Vars, "@EDITOR")
		return
	}
	s.Vars["@EDITOR"] = this.parent
	this.Show(s)
}

func (this *Editor) Current() ([]*uprop.UProperty, string) {
	props := this.doc.GetUProperties()
	r := ""
	for i, sel := range this.selStack {
		ns := strings.Split(sel, ":")
		_, pv := uprop.Find(props, ns)
		if pv != nil {
			if pv.Expender != nil {
				pr := pv.Expender()
				if pr != nil {
					if r != "" {
						r += "/"
					}
					r += sel
					props = pr
					continue
				}
			}
		}
		logger.Debug(tag, "invalid select stack - %s", sel)
		this.selStack = this.selStack[:i]
		break
	}
	return props, r
}

func (this *Editor) Show(s *Session) {
	props, n := this.Current()
	if n != "" {
		n = " " + n
	}
	EditorHelper.DoPropList(s, this.doc.Title()+n, props)
}

func (this *Editor) Process(session *Session, command string) bool {
	name := "ed"
	cname := CommandWord(command)
	if cname == name {
		fs := NewFlagSet(name)
		args := "[close|commit|list]"
		if DoParse(session, command, fs, name, args, 0, 1) {
			return true
		}
		act := ""
		if fs.NArg() > 0 {
			act = fs.Arg(0)
		}
		switch act {
		case "close":
			this.Close(session)
		case "commit":
			err := this.doc.CommitDoc(session)
			if err == nil {
				session.Writeln("commit done")
				this.Close(session)
			} else {
				session.Writeln("ERROR: commit fail - " + err.Error())
			}
		case "list":
			session.Write("edit commands: set, add, rm, select, unselect")
			fir := true
			for _, k := range this.doc.ListCommands() {
				if fir {
					fir = !fir
				} else {
					session.Write(", ")
				}
				session.Write(k)
			}
			session.Writeln("")
		default:
			this.Show(session)
		}
		return true
	}
	h := true
	show := false
	switch cname {
	case "set":
		show = this.commandSet(session, command)
	case "add":
		show = this.commandAdd(session, command)
	case "rm":
		show = this.commandRemove(session, command)
	case "select":
		show = this.commandSelect(session, command)
	case "unselect":
		show = this.commandUnselect(session, command)
	default:
		h, show = this.doc.HandleCommand(session, command)
	}
	if !h {
		return false
	}
	if show {
		this.Show(session)
	}
	return true
}

func (this *Editor) commandSet(s *Session, command string) bool {
	name := "set"
	args := "propname propval"
	fs := NewFlagSet(name)
	if DoParse(s, command, fs, name, args, 2, 2) {
		return false
	}
	varn := fs.Arg(0)
	v := fs.Arg(1)
	prop, _ := this.Current()
	return EditorHelper.DoPropSet(s, prop, varn, v)
}

func (this *Editor) commandAdd(s *Session, command string) bool {
	name := "add"
	args := "propname [...]"
	fs := NewFlagSet(name)
	if DoParse(s, command, fs, name, args, 1, -1) {
		return false
	}
	varn := fs.Arg(0)
	vs := fs.Args()[1:]
	prop, _ := this.Current()
	return EditorHelper.DoPropAdd(s, prop, varn, vs)
}

func (this *Editor) commandRemove(s *Session, command string) bool {
	name := "remove"
	args := "propname [...]"
	fs := NewFlagSet(name)
	if DoParse(s, command, fs, name, args, 1, -1) {
		return false
	}
	varn := fs.Arg(0)
	vs := fs.Args()[1:]
	prop, _ := this.Current()
	return EditorHelper.DoPropRemove(s, prop, varn, vs)
}

func (this *Editor) commandSelect(s *Session, command string) bool {
	name := "select"
	args := "[propname]"
	fs := NewFlagSet(name)
	if DoParse(s, command, fs, name, args, 0, 1) {
		return false
	}
	varn := fs.Arg(0)
	ns := strings.Split(varn, ":")
	props, _ := this.Current()
	_, pv := uprop.Find(props, ns)
	done := false
	if pv != nil {
		if pv.Expender != nil {
			pr := pv.Expender()
			if pr != nil {
				done = true
			}
		}
	}

	if !done {
		s.Writeln(fmt.Sprintf("ERROR: can't select '%s'", varn))
		return false
	}

	this.selStack = append(this.selStack, varn)

	return true
}

func (this *Editor) commandUnselect(s *Session, command string) bool {
	name := "unselect"
	fs := NewFlagSet(name)
	if DoParse(s, command, fs, name, "", 0, 0) {
		return false
	}

	if len(this.selStack) > 0 {
		this.selStack = this.selStack[:len(this.selStack)-1]
	}
	return true
}
