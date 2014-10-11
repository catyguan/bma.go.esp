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
	if this.gof != nil && this.gof.IsNative() {
		buf.WriteString(fmt.Sprintf("%v", this.gof))
	} else {
		if this.chunkName != "" {
			buf.WriteString(this.chunkName)
			if this.line > 0 {
				buf.WriteString(fmt.Sprintf(":%d", this.line))
			}
			if this.funcName != "" {
				buf.WriteString(fmt.Sprintf(" %s(...)", this.funcName))
			}
		} else {
			buf.WriteString("<unknow>")
		}
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

func (this *VMStack) createLocal(vm *VM, n string, val interface{}) {
	va, ok := this.local[n]
	if !ok {
		if va == nil {
			va = new(localVar)
		}
		this.local[n] = va
	}
	va.Set(vm, val)
}
