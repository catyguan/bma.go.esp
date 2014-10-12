package golua

import (
	"bmautil/syncutil"
	"bytes"
	"fmt"
	"logger"
	"sync/atomic"
)

const (
	tag = "golua"
)

type VM struct {
	id      uint32
	running int32
	vmg     *VMG
	stack   *VMStack
	syncutil.CloseState
	trace bool
	sdata []interface{}
}

func newVM(vmg *VMG, id uint32) *VM {
	vm := new(VM)
	vm.id = id
	vm.vmg = vmg
	vm.InitCloseState()
	vm.sdata = make([]interface{}, 0, 16)
	return vm
}

func (this *VM) initStack(st *VMStack) {
	this.stack = st
}

func (this *VM) Id() uint32 {
	return this.id
}

func (this *VM) GetVMG() *VMG {
	return this.vmg
}

func (this *VM) String() string {
	return fmt.Sprintf("VM(%s:%d)", this.vmg.name, this.id)
}

func (this *VM) Spawn(n string) (*VM, error) {
	vm2, err := this.vmg.newVM()
	if err != nil {
		return nil, err
	}
	st := newVMStack(this.stack)
	st.chunkName = n
	vm2.initStack(st)
	logger.Debug(tag, "%s spawn -> %s", this, vm2)
	return vm2, nil
}

func (this *VM) Destroy() {
	if this.IsRunning() {
		//therefore we are in a different goroutine
		this.AskClose()
		return
	}
	if this.IsClosed() {
		return
	}
	if !this.vmg.removeVM(this.id) {
		return
	}
	logger.Debug(tag, "%s destoryed", this)
	st := this.stack
	for st != nil {
		p := st.parent
		st.clear()
		st = p
	}
	this.stack = nil
	this.DoneClose()
}

func (this *VM) IsRunning() bool {
	return atomic.LoadInt32(&this.running) > 0
}

type StackTraceError struct {
	s []string
}

func (this *StackTraceError) String() string {
	return this.Error()
}

func (this *StackTraceError) Error() string {
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	for i, err := range this.s {
		if i != 0 {
			buf.WriteString("\nat ")
		}
		buf.WriteString(err)
	}
	return buf.String()
}

func (this *VM) EnableTrace(b bool) bool {
	old := this.trace
	this.trace = b
	return old
}

func (this *VM) Trace(format string, args ...interface{}) {
	if this.trace {
		logger.Info(tag, format, args...)
	}
}

func (this *VM) DumpStack() string {
	s := this.stack.Dump(this.sdata)
	return fmt.Sprintf("%sSDATA: %v\n", s, this.sdata)
}
