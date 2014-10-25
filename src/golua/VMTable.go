package golua

import (
	"fmt"
	"sync"
)

// CommonVMTable
type CommonVMTable struct {
	mux       *sync.RWMutex
	data      map[string]interface{}
	metaTable VMTable
}

func NewVMTable(o map[string]interface{}) VMTable {
	r := new(CommonVMTable)
	if o == nil {
		r.data = make(map[string]interface{})
	} else {
		r.data = o
	}
	return r
}

func (this *CommonVMTable) String() string {
	return fmt.Sprintf("@%v", this.data)
}

func (this *CommonVMTable) EnableSafe() {
	if this.mux == nil {
		this.mux = new(sync.RWMutex)
	}
}

func (this *CommonVMTable) GetMetaTable() VMTable {
	if this.mux != nil {
		this.mux.RLock()
		defer this.mux.RUnlock()
	}
	return this.metaTable
}

func (this *CommonVMTable) SetMetaTable(mt VMTable) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	this.metaTable = mt
}

func (this *CommonVMTable) Get(vm *VM, key string) (interface{}, error) {
	v, mt := this._rawget(key)
	if v == nil {
		if mt != nil {
			f := mt.Rawget(METATABLE_INDEX)
			if f != nil {
				vm.API_push(f)
				vm.API_push(this)
				vm.API_push(key)
				r0, err := vm.Call(2, 1, nil)
				if err != nil {
					return nil, err
				}
				v, err = vm.API_pop1X(r0, false)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return v, nil
}

func (this *CommonVMTable) _rawget(key string) (interface{}, VMTable) {
	if this.mux != nil {
		this.mux.RLock()
		defer this.mux.RUnlock()
	}
	return this.data[key], this.metaTable
}

func (this *CommonVMTable) Rawget(key string) interface{} {
	v, _ := this._rawget(key)
	return v
}

func (this *CommonVMTable) Set(vm *VM, key string, val interface{}) error {
	ok, mt := this._rawset(key, val, false)
	if !ok {
		if mt != nil {
			f := mt.Rawget(METATABLE_NEWINDEX)
			if f != nil {
				vm.API_push(f)
				vm.API_push(this)
				vm.API_push(key)
				vm.API_push(val)
				_, err := vm.Call(3, 0, nil)
				if err != nil {
					return err
				}
				return nil
			}
		}
		this._rawset(key, val, true)
	}
	return nil
}

func (this *CommonVMTable) _rawset(key string, val interface{}, force bool) (bool, VMTable) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	if val == nil {
		delete(this.data, key)
		return true, nil
	} else {
		if force {
			this.data[key] = val
			return true, nil
		} else {
			_, ok := this.data[key]
			if ok {
				this.data[key] = val
				return true, nil
			}
			return false, this.metaTable
		}
	}
}

func (this *CommonVMTable) Rawset(key string, val interface{}) {
	this._rawset(key, val, true)
}

func (this *CommonVMTable) Delete(key string) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	delete(this.data, key)
}

func (this *CommonVMTable) Len() int {
	return len(this.data)
}

func (this *CommonVMTable) ToMap() map[string]interface{} {
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

func (this *VM) API_toMap(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	if o, ok := v.(map[string]interface{}); ok {
		return o
	}
	if o, ok := v.(VMTable); ok {
		return o.ToMap()
	}
	return nil
}

func (this *VM) API_table(v interface{}) VMTable {
	if v == nil {
		return nil
	}
	if o, ok := v.(VMTable); ok {
		return o
	}
	if o, ok := v.(map[string]interface{}); ok {
		r := new(CommonVMTable)
		r.data = o
		return r
	}
	return nil
}

func (this *VM) API_newtable() VMTable {
	r := new(CommonVMTable)
	r.data = make(map[string]interface{})
	return r
}

// objectVMTable
type objectVMTable struct {
	o interface{}
	p GoObject
}

func NewGOO(o interface{}, p GoObject) VMTable {
	r := new(objectVMTable)
	r.o = o
	r.p = p
	return r
}

func (this *objectVMTable) String() string {
	return fmt.Sprintf("@%v", this.o)
}

func (this *objectVMTable) Get(vm *VM, key string) (interface{}, error) {
	r, err := this.p.Get(vm, this.o, key)
	// fmt.Println("objectVMTable:Get", this.o, key, r, err)
	return r, err
}

func (this *objectVMTable) Rawget(key string) interface{} {
	return nil
}

func (this *objectVMTable) Set(vm *VM, key string, val interface{}) error {
	return this.p.Set(vm, this.o, key, val)
}

func (this *objectVMTable) Rawset(key string, val interface{}) {

}

func (this *objectVMTable) Delete(key string) {

}

func (this *objectVMTable) Len() int {
	return 0
}

func (this *objectVMTable) ToMap() map[string]interface{} {
	return this.p.ToMap(this.o)
}

func (this *VM) API_object(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	if o, ok := v.(*objectVMTable); ok {
		return o.o
	}
	return nil
}

func GoData(d interface{}) interface{} {
	if d == nil {
		return nil
	}
	switch ro := d.(type) {
	case VMTable:
		m := ro.ToMap()
		rm := make(map[string]interface{})
		for k, v := range m {
			rm[k] = GoData(v)
		}
		return rm
	case VMArray:
		a := ro.ToArray()
		ra := make([]interface{}, len(a))
		for k, v := range a {
			ra[k] = GoData(v)
		}
		return ra
	}
	return d
}

func GoluaData(d interface{}) interface{} {
	if d == nil {
		return nil
	}
	switch ro := d.(type) {
	case map[string]interface{}:
		for k, v := range ro {
			ro[k] = GoluaData(v)
		}
		r := new(CommonVMTable)
		r.data = ro
		return r
	case []interface{}:
		for k, v := range ro {
			ro[k] = GoluaData(v)
		}
		r := new(CommonVMArray)
		r.data = ro
		return r
	}
	return d
}
