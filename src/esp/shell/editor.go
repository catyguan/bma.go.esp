package shell

import "fmt"

type EditorDoc interface {
	Title() string
	ListCommands() []string
	HandleCommand(session *Session, cmdline string) (bool, bool)

	OnCloseDoc(session *Session)
	ShowDoc(session *Session)
	CommitDoc(session *Session) error
}

type Editor struct {
	parent *Editor
	doc    EditorDoc
}

func NewEditor() *Editor {
	this := new(Editor)
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
	this.doc.ShowDoc(s)
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
	this.parent.doc.ShowDoc(s)
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
			session.Write("edit commands: ")
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
			this.doc.ShowDoc(session)
		}
		return true
	}
	h, show := this.doc.HandleCommand(session, command)
	if !h {
		return false
	}
	if show {
		this.doc.ShowDoc(session)
	}
	return true
}
