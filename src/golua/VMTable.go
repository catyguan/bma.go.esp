package golua

import (
	"fmt"
	"sync"
)

type CommonVMTable struct {
	Data map[string]interface{}
	mux  *sync.RWMutex
}

func (this *CommonVMTable) String() string {
	return fmt.Sprintf("@%v", this.Data)
}

func (this *CommonVMTable) Get(vm *VM, key string) (interface{}, error) {
	v := this.Rawget(key)
	return v, nil
}

func (this *CommonVMTable) Rawget(key string) interface{} {
	if this.mux != nil {
		this.mux.RLock()
		defer this.mux.RUnlock()
	}
	return this.Data[key]
}

func (this *CommonVMTable) Set(key string, val interface{}) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	if val == nil {
		delete(this.Data, key)
	} else {
		this.Data[key] = val
	}
}

func (this *CommonVMTable) Delete(key string) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	delete(this.Data, key)
}

func (this *CommonVMTable) Len() int {
	return len(this.Data)
}

func (this *CommonVMTable) ToMap() map[string]interface{} {
	return this.Data
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
		r.Data = o
		return r
	}
	return nil
}

func (this *VM) API_newtable() VMTable {
	r := new(CommonVMTable)
	r.Data = make(map[string]interface{})
	return r
}

func (this *VM) API_newarray() []interface{} {
	return make([]interface{}, 0)
}
