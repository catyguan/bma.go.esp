package golua

import (
	"fmt"
	"logger"
	"sync"
	"sync/atomic"
)

type VMG struct {
	name    string
	mux     sync.RWMutex
	vmid    uint32
	vms     map[uint32]*VM
	globals map[string]interface{}
	closed  uint32
}

func NewVMG(n string) *VMG {
	r := new(VMG)
	r.vms = make(map[uint32]*VM)
	r.globals = make(map[string]interface{})
	r.name = n
	return r
}

func (this *VMG) String() string {
	return fmt.Sprintf("VMG[%s]", this.name)
}

func (this *VMG) newVM() (*VM, error) {
	if this.IsClose() {
		return nil, fmt.Errorf("%s closed", this)
	}
	this.mux.Lock()
	defer this.mux.Unlock()
	this.vmid++
	vm := newVM(this, this.vmid)
	this.vms[vm.id] = vm
	return vm, nil
}

func (this *VMG) CreateVM() (*VM, error) {
	vm, err := this.newVM()
	if err != nil {
		return nil, err
	}
	st := newVMStack(nil)
	st.name = fmt.Sprintf("VM<%d>", vm.id)
	vm.initStack(st)
	logger.Debug(tag, "createVM -> %s", vm)
	return vm, nil
}

func (this *VMG) removeVM(id uint32) bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	_, ok := this.vms[id]
	delete(this.vms, id)
	return ok
}

func (this *VMG) GetGlobal(n string) (interface{}, bool) {
	this.mux.RLock()
	defer this.mux.RUnlock()
	v, ok := this.globals[n]
	return v, ok
}

func (this *VMG) SetGlobal(n string, v interface{}) interface{} {
	this.mux.Lock()
	defer this.mux.Unlock()
	old := this.globals[n]
	this.globals[n] = v
	return old
}

func (this *VMG) IsClose() bool {
	return atomic.LoadUint32(&this.closed) == 1
}

func (this *VMG) Close() {
	if !atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		return
	}
	tmp := make(map[uint32]*VM)
	this.mux.Lock()
	for k, vm := range this.vms {
		tmp[k] = vm
	}
	this.mux.Unlock()
	for _, vm := range tmp {
		vm.Destroy()
	}
}
