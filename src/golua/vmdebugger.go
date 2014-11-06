package golua

import (
	"logger"
	"sync/atomic"
)

type Breakpoint struct {
	vmid      uint32
	chunkName string
	line      int
}

func (this *Breakpoint) Same(o *Breakpoint) bool {
	if this.vmid != o.vmid {
		return false
	}
	if this.chunkName != o.chunkName {
		return false
	}
	if this.line != o.line {
		return false
	}
	return true
}

func (this *GoLua) addDebuger(dg *VMDebugger) {
	this.dgMutex.Lock()
	defer this.dgMutex.Unlock()
	this.dgs[dg.vm.id] = dg
}

func (this *GoLua) GetDebugger(id uint32) *VMDebugger {
	this.dgMutex.RLock()
	defer this.dgMutex.RUnlock()
	return this.dgs[id]
}

func (this *GoLua) removeDebuger(dg *VMDebugger) {
	this.dgMutex.Lock()
	defer this.dgMutex.Unlock()
	delete(this.dgs, dg.vm.id)
}

func (this *GoLua) AddBreakpoint(o *Breakpoint) bool {
	this.dgMutex.Lock()
	defer this.dgMutex.Unlock()
	for old, _ := range this.breakpoints {
		if o.Same(old) {
			return false
		}
	}
	this.breakpoints[o] = true
	return true
}

func (this *GoLua) RemoveBreakpoint(o *Breakpoint) bool {
	this.dgMutex.Lock()
	defer this.dgMutex.Unlock()
	for old, _ := range this.breakpoints {
		if o.Same(old) {
			delete(this.breakpoints, old)
			return true
		}
	}
	return false
}

func (this *GoLua) ClearBreakpoint() {
	this.dgMutex.Lock()
	defer this.dgMutex.Unlock()
	for o, _ := range this.breakpoints {
		delete(this.breakpoints, o)
	}
}

func (this *GoLua) EnableBreakpoint(o *Breakpoint, v bool) bool {
	this.dgMutex.Lock()
	defer this.dgMutex.Unlock()
	for old, _ := range this.breakpoints {
		if o.Same(old) {
			this.breakpoints[old] = v
			return true
		}
	}
	return false
}

func (this *GoLua) checkBreakpoint(vm *VM) bool {
	n := vm.stack.chunkName
	l := vm.stack.line
	this.dgMutex.RLock()
	defer this.dgMutex.RUnlock()
	for o, e := range this.breakpoints {
		if !e {
			continue
		}
		if o.vmid != 0 {
			if vm.id != o.vmid {
				continue
			}
		}
		if n != o.chunkName {
			continue
		}
		if l != 0 {
			if l != o.line {
				continue
			}
		}
		return true
	}
	return false
}

func (this *VM) closeDebugger() {
	if this.debugger != nil {
		this.gl.removeDebuger(this.debugger)
	}
}

func (this *VM) DebugSet(flag int) {
	this.runMode = flag
	if flag > 0 {
		if this.debugger == nil {
			dg := newVMDebugger(this)
			this.gl.addDebuger(dg)
			this.debugger = dg
		}
	} else {
		if this.debugger == nil {
			this.gl.removeDebuger(this.debugger)
			this.debugger.close()
			this.debugger = nil
		}
	}
}

///// VMDebugger ///////////////////////////////////////
type debuggerCommand struct {
	title  string
	action func(vm *VM)
}
type VMDebugger struct {
	vm      *VM
	cmds    chan *debuggerCommand
	last    int
	st      *VMStack
	line    int
	waiting uint32
}

func newVMDebugger(vm *VM) *VMDebugger {
	r := new(VMDebugger)
	r.vm = vm
	r.last = 0
	r.cmds = make(chan *debuggerCommand, 8)
	return r
}

func (this *VMDebugger) mark() bool {
	l := this.vm.stack.line
	if l != 0 {
		this.st = this.vm.stack
		this.line = l
		// logger.Debug(tag, "%s mark at %s:%d", this.vm, this.st.chunkName, this.line)
		return true
	}
	return false
}

// 0-normal 1-watch 2-wait 11-step 12-step.in 13-step.out
func (this *VMDebugger) Check() {
	switch this.vm.runMode {
	case 0:
		return
	case 2:
		this.enterWait()
	case 3:
		if this.vm.gl.checkBreakpoint(this.vm) {
			this.vm.runMode = 2
			this.enterWait()
		}
	case 11:
		if this.last != 11 {
			if this.mark() {
				this.last = 11
			}
		} else {
			end := false
			st := this.vm.stack
			if st != this.st {
				end = true
				st = st.parent
				for st != nil {
					if this.st == st {
						// in sub call
						return
					}
					st = st.parent
				}
			}
			if !end {
				l := this.vm.stack.line
				if l != 0 {
					end = this.line != l
				}
			}
			if end {
				this.last = 0
				this.vm.runMode = 2
				this.enterWait()
			}
		}
	case 12:
		if this.last != 12 {
			if this.mark() {
				this.last = 12
			}
		} else {
			end := false
			st := this.vm.stack
			if st != this.st {
				end = true
			}
			if !end {
				l := this.vm.stack.line
				if l != 0 {
					end = this.line != l
				}
			}
			if end {
				this.last = 0
				this.vm.runMode = 2
				this.enterWait()
			}
		}
	case 13:
		if this.last != 13 {
			if this.mark() {
				this.last = 13
			}
		} else {
			st := this.vm.stack
			for st != nil {
				if this.st == st {
					// in sub call
					return
				}
				st = st.parent
			}
			this.last = 0
			this.vm.runMode = 2
			this.enterWait()
		}
	default:
	}
}

func (this *VMDebugger) call(cmd *debuggerCommand) {
	defer func() {
		recover()
	}()
	this.cmds <- cmd
}

func (this *VMDebugger) DoRun() {
	cmd := new(debuggerCommand)
	cmd.action = func(vm *VM) {
		vm.DebugSet(0)
	}
	cmd.title = "RUN"
	this.cmds <- cmd
}

func (this *VMDebugger) DoStep() {
	cmd := new(debuggerCommand)
	cmd.action = func(vm *VM) {
		vm.DebugSet(11)
	}
	cmd.title = "STEP"
	this.cmds <- cmd
}

func (this *VMDebugger) DoStepIn() {
	cmd := new(debuggerCommand)
	cmd.action = func(vm *VM) {
		vm.DebugSet(12)
	}
	cmd.title = "STEP_IN"
	this.cmds <- cmd
}

func (this *VMDebugger) DoStepOut() {
	cmd := new(debuggerCommand)
	cmd.action = func(vm *VM) {
		vm.DebugSet(13)
	}
	cmd.title = "STEP_OUT"
	this.cmds <- cmd
}

func (this *VMDebugger) DoWatch() {
	cmd := new(debuggerCommand)
	cmd.action = func(vm *VM) {
		vm.DebugSet(3)
	}
	cmd.title = "WATCH"
	this.cmds <- cmd
}

func (this *VMDebugger) DoStackTrace() []string {
	out := make(chan []string, 1)
	defer close(out)

	cmd := new(debuggerCommand)
	cmd.action = func(vm *VM) {
		r := make([]string, 0)
		st := vm.stack
		for st != nil {
			r = append(r, st.String())
			st = st.parent
		}
		out <- r
	}
	cmd.title = "STACK_TRACE"
	this.cmds <- cmd
	r := <-out
	return r
}

func (this *VMDebugger) CanDo() bool {
	v := atomic.LoadUint32(&this.waiting)
	return v == 1
}

func (this *VMDebugger) enterWait() {
	atomic.StoreUint32(&this.waiting, 1)
	defer atomic.StoreUint32(&this.waiting, 0)

	this.vm.maxExecutionTime = 0 // don't timeout when DEBUG
	st := this.vm.stack
	for {
		if this.vm.runMode != 2 {
			return
		}
		logger.Debug(tag, "%s debugger waiting at %s ...", this.vm, st)
		cmd := <-this.cmds
		if cmd == nil {
			return
		}
		logger.Debug(tag, "%s debugger command -> %s", this.vm, cmd.title)
		func() {
			defer func() {
				recover()
			}()
			cmd.action(this.vm)
		}()
	}
}

func (this *VMDebugger) close() {
	close(this.cmds)
}
