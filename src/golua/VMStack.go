package golua

import (
	"bytes"
	"fmt"
)

type VMStack struct {
	parent    *VMStack
	chunkName string
	funcName  string
	line      int
	gof       GoFunction
	local     map[string]VMVar
	stack     []interface{}
	stackTop  int
}

func newVMStack(p *VMStack) *VMStack {
	r := new(VMStack)
	r.parent = p
	r.local = make(map[string]VMVar)
	r.stack = make([]interface{}, 0, 8)
	return r
}

func (this *VMStack) String() string {
	buf := bytes.NewBuffer(make([]byte, 0, 64))
	if this.chunkName != "" {
		buf.WriteString(this.chunkName)
		if this.line > 0 {
			buf.WriteString(fmt.Sprintf(":%d", this.line))
		}
		if this.funcName != "" {
			buf.WriteString(" ")
			buf.WriteString(this.funcName)
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
	for i := 0; i < len(this.stack); i++ {
		this.stack[i] = nil
	}
	this.stack = nil
	for k, _ := range this.local {
		delete(this.local, k)
	}
	this.local = nil
}

func (this *VMStack) createLocal(n string, val interface{}) {
	va, ok := this.local[n]
	if !ok {
		va = new(localVar)
		this.local[n] = va
	}
	va.Set(val)
}
