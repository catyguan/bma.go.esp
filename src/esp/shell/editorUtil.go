package shell

import (
	"fmt"
	"uprop"
)

type EditorSupported interface {
	GetUProperties() []*uprop.UProperty
}

type editorHelper byte

var EditorHelper editorHelper

func (this *editorHelper) DoList(s *Session, title string, o EditorSupported) {
	this.DoPropList(s, title, o.GetUProperties())
}

func (this *editorHelper) DoPropList(s *Session, title string, props []*uprop.UProperty) {
	s.Writeln("[ " + title + " ]")
	for _, p := range props {
		s.Writeln(p.String())
	}
	s.Writeln("[ " + title + " end ]")
}

func (this *editorHelper) DoEdit(s *Session, o EditorSupported, n, v string) bool {
	return this.DoPropEdit(s, o.GetUProperties(), n, v)
}

func (this *editorHelper) DoPropEdit(s *Session, props []*uprop.UProperty, n, v string) bool {
	for _, p := range props {
		if p.Name == n {
			err := p.Commit(v)
			if err != nil {
				s.Writeln(fmt.Sprintf("ERROR: commit %s fail - %s", n, err))
				return false
			}
			return true
		}
	}
	s.Writeln("ERROR: unknow var '" + n + "'")
	return false
}
