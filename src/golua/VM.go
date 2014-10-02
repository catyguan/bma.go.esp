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
}

func newVM(vmg *VMG, id uint32) *VM {
	vm := new(VM)
	vm.id = id
	vm.vmg = vmg
	vm.InitCloseState()
	return vm
}

func (this *VM) initStack(st *VMStack) {
	this.stack = st
}

func (this *VM) Id() uint32 {
	return this.id
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
	st.name = n
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

func (this *VM) Call(nargs int, nresults int) (rerr error) {
	if this.IsClosing() {
		return fmt.Errorf("%s closed", this)
	}
	st := this.stack
	var nst *VMStack
	err := func(nargs int, nresults int) error {
		atomic.AddInt32(&this.running, 1)
		defer func() {
			atomic.AddInt32(&this.running, -1)
			if x := recover(); x != nil {
				logger.Warn(tag, "runtime panic: %v", x)
				if err, ok := x.(error); ok {
					rerr = err
				} else {
					rerr = fmt.Errorf("%v", x)
				}
			}
		}()
		n := nargs + 1
		err1 := this.API_checkstack(n)
		if err1 != nil {
			return err1
		}
		at := this.API_absindex(-n)
		f, err5 := this.API_peek(at)
		if err5 != nil {
			return err5
		}
		f, err5 = this.API_value(f)
		if err5 != nil {
			return err5
		}
		if !this.API_canCall(f) {
			return fmt.Errorf("can't call at '%v'", f)
		}
		nst = newVMStack(st)
		if tt, ok := f.(StackTracable); ok {
			nst.name = tt.StackInfo()
		}
		for i := 1; i <= nargs; i++ {
			v, err2 := this.API_peek(at + i)
			if err2 != nil {
				return err2
			}
			nst.stack = append(nst.stack, v)
			nst.stackTop++
		}
		this.API_pop(n)
		this.stack = nst

		if gof, ok := f.(GoFunction); ok {
			nst.gof = gof
			rc, err3 := gof.Exec(this)
			if err3 != nil {
				return err3
			}
			at = this.API_absindex(-rc)
			nres := nresults
			if nres < 0 {
				nres = rc
			}
			for i := 0; i < nres; i++ {
				var r interface{}
				if i < rc {
					v, err4 := this.API_peek(at + i)
					if err4 != nil {
						return err4
					}
					r = v
				} else {
					r = nil
				}
				if st.stackTop < len(st.stack) {
					st.stack[st.stackTop] = r
				} else {
					st.stack = append(st.stack, r)
				}
				st.stackTop++
			}
			logger.Debug(tag, "Call %s(%d,%d) -> %d", gof, nargs, nresults, rc)
		} else {
			panic(fmt.Errorf("unknow callable '%v'", f))
		}
		return nil
	}(nargs, nresults)

	if err != nil {
		if _, ok := err.(*StackTraceError); !ok {
			nerr := new(StackTraceError)
			nerr.s = make([]string, 0, 8)
			nerr.s = append(nerr.s, err.Error())
			p := this.stack
			for p != nil {
				nerr.s = append(nerr.s, p.String())
				p = p.parent
			}
			err = nerr
		}
	}

	if nst != nil {
		nst.clear()
	}
	this.stack = st

	return err
}

func (this *VM) DumpStack() string {
	return this.stack.Dump()
}
