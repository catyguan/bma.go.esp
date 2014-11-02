package vmmclass

import (
	"fmt"
	"golua"
	"sync"
)

// superVMTable
type superVMTable struct {
	supers []*classVMTable
}

func newSuperVMTable(sp []*classVMTable) *superVMTable {
	r := new(superVMTable)
	r.supers = sp
	return r
}

func (this *superVMTable) String() string {
	return fmt.Sprintf("@super")
}

func (this *superVMTable) Get(vm *golua.VM, key string) (interface{}, error) {
	for _, sp := range this.supers {
		return sp.Get(vm, key)
	}
	return nil, nil
}

func (this *superVMTable) Rawget(key string) interface{} {
	return nil
}

func (this *superVMTable) Set(vm *golua.VM, key string, val interface{}) error {
	return nil
}

func (this *superVMTable) Rawset(key string, val interface{}) {

}

func (this *superVMTable) Delete(key string) {

}

func (this *superVMTable) Len() int {
	return len(this.supers)
}

func (this *superVMTable) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	for _, s := range this.supers {
		r[s.name] = s
	}
	return r
}

// classVMTable
type classVMTable struct {
	name   string
	lock   sync.RWMutex
	def    map[string]interface{}
	supers []*classVMTable
}

func newClassVMTable(n string, sp []*classVMTable) *classVMTable {
	r := new(classVMTable)
	r.name = n
	r.supers = sp
	r.def = make(map[string]interface{})
	return r
}

func (this *classVMTable) String() string {
	return fmt.Sprintf("@class(%s)", this.name)
}

func (this *classVMTable) Get(vm *golua.VM, key string) (interface{}, error) {
	switch key {
	case "Super":
		return newSuperVMTable(this.supers), nil
	case "Class":
		return this, nil
	case "ClassName":
		return this.name, nil
	case "New":
		return golua.NewGOF("Class.New", func(vm *golua.VM, self interface{}) (int, error) {
			top := vm.API_gettop()
			ns, err := vm.API_popN(top, true)
			if err != nil {
				return 0, err
			}
			o, err2 := this.New(vm, ns)
			if err2 != nil {
				return 0, err2
			}
			vm.API_push(o)
			return 1, nil
		}), nil
	}
	this.lock.RLock()
	v, ok := this.def[key]
	this.lock.RUnlock()
	if ok {
		return v, nil
	}
	for _, sp := range this.supers {
		return sp.Get(vm, key)
	}
	return nil, nil
}

func (this *classVMTable) Rawget(key string) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if v, ok := this.def[key]; ok {
		return v
	}
	return nil
}

func (this *classVMTable) Set(vm *golua.VM, key string, val interface{}) error {
	this.Rawset(key, val)
	return nil
}

func (this *classVMTable) Rawset(key string, val interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.def[key] = val
}

func (this *classVMTable) Delete(key string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.def, key)
}

func (this *classVMTable) Len() int {
	return len(this.def)
}

func (this *classVMTable) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	this.lock.RLock()
	defer this.lock.RUnlock()
	for k, v := range this.def {
		r[k] = v
	}
	return r
}

func (this *classVMTable) Props() []string {
	r := make([]string, 0)
	for _, cls := range this.supers {
		p := cls.Props()
		if p != nil {
			r = append(r, p...)
		}
	}
	this.lock.RLock()
	defer this.lock.RUnlock()
	for k, _ := range this.def {
		r = append(r, k)
	}
	return r
}

func (this *classVMTable) HasBase(n string) bool {
	for _, cls := range this.supers {
		if cls.name == n {
			return true
		}
		if cls.HasBase(n) {
			return true
		}
	}
	return false
}

func (this *classVMTable) ctor(vm *golua.VM, o *ciVMTable, args []interface{}) error {
	for _, cls := range this.supers {
		errX := cls.ctor(vm, o, args)
		if errX != nil {
			return errX
		}
	}
	ctor, err := this.Get(vm, "ctor")
	if err != nil {
		return err
	}
	if ctor != nil {
		err1 := vm.API_pushMemberCall(o, ctor)
		if err1 != nil {
			return err1
		}
		for _, arg := range args {
			vm.API_push(arg)
		}
		_, err2 := vm.Call(len(args), 0, nil)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func (this *classVMTable) New(vm *golua.VM, args []interface{}) (golua.VMTable, error) {
	r := new(ciVMTable)
	r.data = make(map[string]interface{})
	r.cls = this
	err0 := this.ctor(vm, r, args)
	if err0 != nil {
		return nil, err0
	}
	return r, nil
}

func (this *classVMTable) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, _ := range this.def {
		delete(this.def, k)
	}
	this.supers = nil
}

// ciVMTable
type ciVMTable struct {
	mux  *sync.RWMutex
	data map[string]interface{}
	cls  *classVMTable
}

func (this *ciVMTable) String() string {
	return fmt.Sprintf("@%s%v", this.cls.name, this.data)
}

func (this *ciVMTable) EnableSafe() {
	if this.mux == nil {
		this.mux = new(sync.RWMutex)
	}
}

func (this *ciVMTable) Get(vm *golua.VM, key string) (interface{}, error) {
	v := this.Rawget(key)
	if v == nil {
		return this.cls.Get(vm, key)
	}
	return v, nil
}

func (this *ciVMTable) Rawget(key string) interface{} {
	if this.mux != nil {
		this.mux.RLock()
		defer this.mux.RUnlock()
	}
	return this.data[key]
}

func (this *ciVMTable) Set(vm *golua.VM, key string, val interface{}) error {
	this.Rawset(key, val)
	return nil
}

func (this *ciVMTable) Rawset(key string, val interface{}) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	this.data[key] = val
}

func (this *ciVMTable) Delete(key string) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	delete(this.data, key)
}

func (this *ciVMTable) Len() int {
	return len(this.data)
}

func (this *ciVMTable) ToMap() map[string]interface{} {
	if this.mux != nil {
		this.mux.RLock()
		defer this.mux.RUnlock()
		r := make(map[string]interface{})
		for k, v := range this.data {
			r[k] = v
		}
		return r
	} else {
		return this.data
	}
}
