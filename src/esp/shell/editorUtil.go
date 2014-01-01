package shell

import (
	"bytes"
	"fmt"
	"strings"
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
		this.DoPropShow(s, "", 0, p, p.Value)
	}
	s.Writeln("[ " + title + " end ]")
}

func (this *editorHelper) DoPropShow(s *Session, prex string, level int, p *uprop.UProperty, val *uprop.UPropertyValue) {
	buf := bytes.NewBuffer(make([]byte, 0))
	for i := 0; i < level; i++ {
		buf.WriteString("\t")
	}
	n := ""
	if p != nil {
		if p.IsOptional {
			buf.WriteString("*")
		}
		if prex == "" {
			n = p.Name
		} else {
			n = prex + p.Name
		}
	} else {
		n = prex
	}

	kind := uprop.UPK_VALUE
	if val != nil {
		kind = val.Kind
	}
	switch kind {
	case uprop.UPK_VALUE:
		var v interface{}
		if val != nil {
			v = val.Value
		}
		buf.WriteString(fmt.Sprintf("%s=%v", n, v))
	case uprop.UPK_FOLD:
		var v interface{}
		if val != nil {
			v = val.Value
		}
		buf.WriteString(fmt.Sprintf("%s {%v}", n, v))
	case uprop.UPK_LIST:
		buf.WriteString(fmt.Sprintf("%s [...]", n))
	case uprop.UPK_MAP:
		buf.WriteString(fmt.Sprintf("%s {...}", n))
	}
	if p != nil && p.Tips != "" {
		buf.WriteString(" (")
		buf.WriteString(p.Tips)
		buf.WriteString(")")
	}
	s.Writeln(buf.String())

	switch kind {
	case uprop.UPK_LIST:
		if val != nil && val.Value != nil {
			list := val.Value.([]*uprop.UPropertyValue)
			for i, child := range list {
				cprex := fmt.Sprintf("%d: ", i+1)
				this.DoPropShow(s, cprex, level+1, nil, child)
			}
		}
	case uprop.UPK_MAP:
		if val != nil && val.Value != nil {
			m := val.Value.([]*uprop.UProperty)
			cprex := p.Name + ": "
			for _, child := range m {
				this.DoPropShow(s, cprex, level+1, child, child.Value)
			}
		}
	}

}

func (this *editorHelper) DoSet(s *Session, o EditorSupported, n, v string) bool {
	return this.DoPropSet(s, o.GetUProperties(), n, v)
}

func (this *editorHelper) DoPropSet(s *Session, props []*uprop.UProperty, n, v string) bool {
	ns := strings.Split(n, ":")
	_, pv := uprop.Find(props, ns)
	if pv != nil {
		err := pv.Commit(v)
		if err != nil {
			s.Writeln(fmt.Sprintf("ERROR: set '%s' fail - %s", n, err))
			return false
		}
		return true
	}
	s.Writeln("ERROR: unknow prop '" + n + "'")
	return false
}

func (this *editorHelper) DoPropAdd(s *Session, props []*uprop.UProperty, n string, vs []string) bool {
	ns := strings.Split(n, ":")
	p, _ := uprop.Find(props, ns)
	if p != nil {
		err := p.CallAdd(vs)
		if err != nil {
			s.Writeln(fmt.Sprintf("ERROR: add '%s' fail - %s", n, err))
			return false
		}
		return true
	}
	s.Writeln("ERROR: unknow prop '" + n + "'")
	return false
}

func (this *editorHelper) DoPropRemove(s *Session, props []*uprop.UProperty, n string, vs []string) bool {
	ns := strings.Split(n, ":")
	p, _ := uprop.Find(props, ns)
	if p != nil {
		err := p.CallRemove(vs)
		if err != nil {
			s.Writeln(fmt.Sprintf("ERROR: remove '%s' fail - %s", n, err))
			return false
		}
		return true
	}
	s.Writeln("ERROR: unknow prop '" + n + "'")
	return false
}
