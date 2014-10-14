package golua

import (
	"bmautil/syncutil"
	"bytes"
	"fmt"
	"logger"
	"sync/atomic"
	"time"
)

const (
	tag = "golua"
)

type VMConfig struct {
	MaxStack  int
	TimeLimit int // MS, 0 = nolimit
	TimeCheck int
}

func (this *VMConfig) Valid() error {
	if this.MaxStack <= 0 {
		return fmt.Errorf("MaxStack invalid")
	}
	if this.TimeLimit < 0 {
		return fmt.Errorf("TimeLimit invalid")
	}
	if this.TimeCheck <= 0 {
		return fmt.Errorf("TimeCheck invalid")
	}
	return nil
}

var (
	defaultConfig VMConfig
)

func init() {
	defaultConfig.MaxStack = 128
	defaultConfig.TimeLimit = 30 * 1000
	defaultConfig.TimeCheck = 1000
}

type VM struct {
	id         uint32
	name       string
	running    int32
	vmg        *VMG
	stack      *VMStack
	numOfStack int
	sdata      []interface{}
	syncutil.CloseState

	config           *VMConfig
	maxExecutionTime int
	numOfTime        int
	executeTime      time.Time
	trace            bool
}

func newVM(vmg *VMG, id uint32) *VM {
	vm := new(VM)
	vm.id = id
	vm.vmg = vmg
	vm.InitCloseState()
	vm.sdata = make([]interface{}, 0, 16)
	vm.config = &defaultConfig
	vm.maxExecutionTime = -1
	vm.executeTime = time.Now()
	return vm
}

func (this *VM) SetMaxExecutionTime(v int) {
	this.maxExecutionTime = v
}
func (this *VM) GetMaxExecutionTime() int {
	if this.maxExecutionTime < 0 {
		return this.config.TimeLimit
	}
	return this.maxExecutionTime
}
func (this *VM) ResetExecutionTime() {
	this.maxExecutionTime = -1
	this.executeTime = time.Now()
}

func (this *VM) Setup(cfg *VMConfig) {
	this.config = cfg
}

func (this *VM) initStack(st *VMStack) {
	this.stack = st
	this.numOfStack = 1
}

func (this *VM) Id() uint32 {
	return this.id
}

func (this *VM) GetVMG() *VMG {
	return this.vmg
}

func (this *VM) String() string {
	return fmt.Sprintf("VM(%s)", this.name)
}

func (this *VM) Spawn(n string, childStack bool) (*VM, error) {
	vm2, err := this.vmg.newVM()
	if err != nil {
		return nil, err
	}
	vm2.name = fmt.Sprintf("%s-%d", this.name, vm2.id)
	vm2.config = this.config
	vm2.trace = this.trace
	pst := this.stack
	if !childStack {
		pst = nil
	}
	st := newVMStack(pst)
	if n == "" {
		n = fmt.Sprintf("VM<%s>", vm2.name)
	}
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
	for i := 0; i < len(this.sdata); i++ {
		this.sdata[i] = nil
	}
	this.sdata = nil
	this.DoneClose()
}

func (this *VM) IsRunning() bool {
	return atomic.LoadInt32(&this.running) > 0
}

func (this *VM) PrepareRun(b bool) {
	if b {
		atomic.AddInt32(&this.running, 1)
	} else {
		atomic.AddInt32(&this.running, -1)
	}
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
		logger.Info(tag, this.name+": "+format, args...)
	}
}

func (this *VM) DumpStack() string {
	s := this.stack.Dump(this.sdata)
	return fmt.Sprintf("%sSDATA: %v\n", s, this.sdata)
}
