package golua

import (
	"fmt"
	"sync"
)

// voidVar
type voidVar struct {
	name string
}

func (this *voidVar) Get(vm *VM) (interface{}, error) {
	return nil, nil
}

func (this *voidVar) Set(vm *VM, v interface{}) (bool, error) {
	return true, nil
}

func (this *voidVar) String() string {
	return "voidVar"
}

var (
	VoidVar voidVar
)

// globalVar
type globalVar struct {
	name string
}

func (this *globalVar) Get(vm *VM) (interface{}, error) {
	vmg := vm.GetVMG()
	v, ok := vmg.GetGlobal(this.name)
	if ok {
		return v, nil
	}
	return nil, nil
}

func (this *globalVar) Set(vm *VM, v interface{}) (bool, error) {
	vmg := vm.GetVMG()
	vmg.SetGlobal(this.name, v)
	return true, nil
}

func (this *globalVar) EnableSafe() {

}

func (this *globalVar) String() string {
	return fmt.Sprintf("globalVar(%s)", this.name)
}

// localVar
type localVar struct {
	value interface{}
	mux   *sync.RWMutex
}

func (this *localVar) Get(vm *VM) (interface{}, error) {
	if this.mux != nil {
		this.mux.RLock()
		defer this.mux.RUnlock()
	}
	return this.value, nil
}

func (this *localVar) Set(vm *VM, v interface{}) (bool, error) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	this.value = v
	return true, nil
}

func (this *localVar) EnableSafe() {
	if this.mux == nil {
		this.mux = new(sync.RWMutex)
	}
}

func (this *localVar) String() string {
	return fmt.Sprintf("localVar(%v)", this.value)
}

// memberVar
type memberVar struct {
	obj interface{}
	key interface{}
}

func (this *memberVar) Get(vm *VM) (interface{}, error) {
	if this.obj == nil {
		return nil, nil
	}
	return vm.API_getMember(this.obj, this.key)
}

func (this *memberVar) Set(vm *VM, v interface{}) (bool, error) {
	if this.obj == nil {
		return false, nil
	}
	return vm.API_setMember(this.obj, this.key, v)
}

func (this *memberVar) String() string {
	return fmt.Sprintf("memberVar(%T, %v)", this.obj, this.key)
}
