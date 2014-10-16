package golua

import (
	"bmautil/valutil"
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
	switch o := this.obj.(type) {
	case []interface{}:
		i := valutil.ToInt(this.key, -1)
		if i < 0 || i >= len(o) {
			return nil, fmt.Errorf("index(%d) out of range(%d)", i, len(o))
		}
		return o[i], nil
	case VMArray:
		i := valutil.ToInt(this.key, -1)
		return o.Get(vm, i)
	case map[string]interface{}:
		s := valutil.ToString(this.key, "")
		v := o[s]
		return v, nil
	case VMTable:
		s := valutil.ToString(this.key, "")
		return o.Get(vm, s)

	}
	return nil, fmt.Errorf("unknow memberVar(%t)", this.obj)
}

func (this *memberVar) Set(vm *VM, v interface{}) (bool, error) {
	if this.obj == nil {
		return false, nil
	}
	switch o := this.obj.(type) {
	case []interface{}:
		i := valutil.ToInt(this.key, -1)
		if i < 0 || i >= len(o) {
			return false, fmt.Errorf("index(%d) out of range(%d)", i, len(o))
		}
		o[i] = v
		return true, nil
	case VMArray:
		i := valutil.ToInt(this.key, -1)
		err := o.Set(vm, i, v)
		if err != nil {
			return false, err
		}
		return true, nil
	case map[string]interface{}:
		s := valutil.ToString(this.key, "")
		o[s] = v
		return true, nil
	case VMTable:
		s := valutil.ToString(this.key, "")
		err := o.Set(vm, s, v)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, fmt.Errorf("unknow memberVar(%t)", this.obj)
}

func (this *memberVar) String() string {
	return fmt.Sprintf("memberVar(%T, %v)", this.obj, this.key)
}

// selfm
type selfm struct {
	self   interface{}
	mvalue interface{}
}
