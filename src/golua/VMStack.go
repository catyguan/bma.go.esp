package golua

import (
	"bytes"
	"fmt"
)

type VMStack struct {
	parent   *VMStack
	name     string
	line     int
	gof      GoFunction
	local    map[string]interface{}
	stack    []interface{}
	stackTop int
}

func newVMStack(p *VMStack) *VMStack {
	r := new(VMStack)
	r.parent = p
	r.local = make(map[string]interface{})
	r.stack = make([]interface{}, 0, 8)
	return r
}

func (this *VMStack) String() string {
	buf := bytes.NewBuffer(make([]byte, 0, 64))
	if this.name != "" {
		buf.WriteString(this.name)
		if this.line > 0 {
			buf.WriteString(fmt.Sprintf(":%d", this.line))
		}
	} else if this.gof != nil {
		buf.WriteString(this.gof.String())
	} else {
		buf.WriteString("<unknow>")
	}
	return buf.String()
}

func (this *VMStack) Dump() string {
	buf := bytes.NewBuffer([]byte{})
	st := this
	for st != nil {
		buf.WriteString(fmt.Sprintf("%s\n", st.String()))
		buf.WriteString(fmt.Sprintf("\tLOCAL: %v\n", st.local))
		buf.WriteString(fmt.Sprintf("\tSTACK: %d, %v\n", st.stackTop, st.stack))
		st = st.parent
	}
	return buf.String()
}

func (this *VMStack) clear() {
	this.parent = nil
	this.stack = nil
}
