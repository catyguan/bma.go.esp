package golua

import (
	"bytes"
	"fmt"
)

type VMStack struct {
	parent     *VMStack
	chunkName  string
	funcName   string
	line       int
	gof        GoFunction
	local      map[string]VMVar
	stackBegin int
	stackTop   int
	defers     []interface{}
}

func newVMStack(p *VMStack) *VMStack {
	r := new(VMStack)
	r.parent = p
	if p != nil {
		r.stackBegin = p.stackBegin + p.stackTop
	} else {
		r.stackBegin = 0
	}
	r.stackTop = 0
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

func (this *VMStack) Dump(sdata []interface{}) string {
	buf := bytes.NewBuffer([]byte{})
	st := this
	for st != nil {
		buf.WriteString(fmt.Sprintf("%s\n", st.String()))
		buf.WriteString(fmt.Sprintf("\tLOCAL: %v\n", st.local))
		buf.WriteString(fmt.Sprintf("\tSTACK: %d:%d, %v\n", st.stackBegin, st.stackTop, sdata[st.stackBegin:st.stackBegin+st.stackTop]))
		st = st.parent
	}
	return buf.String()
}

func (this *VMStack) clear() {
	this.parent = nil
	for k, _ := range this.local {
		delete(this.local, k)
	}
	this.local = nil
	for i, _ := range this.defers {
		this.defers[i] = nil
	}
	this.defers = nil
}

func (this *VMStack) createLocal(vm *VM, n string, val interface{}) {
	if this.local == nil {
		this.local = make(map[string]VMVar)
	}
	va, ok := this.local[n]
	if !ok {
		if va == nil {
			va = new(localVar)
		}
		this.local[n] = va
	}
	va.Set(vm, val)
}
