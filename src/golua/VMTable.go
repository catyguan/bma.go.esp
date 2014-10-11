package golua

import (
	"fmt"
	"sync"
)

type VMTable struct {
	Data map[string]interface{}
	mux  *sync.RWMutex
}

func (this *VMTable) String() string {
	return fmt.Sprintf("@%v", this.Data)
}

func (this *VMTable) Get(vm *VM, key string) (interface{}, error) {
	v := this.Rawget(key)
	return v, nil
}

func (this *VMTable) Rawget(key string) interface{} {
	if this.mux != nil {
		this.mux.RLock()
		defer this.mux.RUnlock()
	}
	return this.Data[key]
}

func (this *VMTable) Set(key string, val interface{}) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	this.Data[key] = val
}

func (this *VM) API_table(v interface{}) *VMTable {
	if v == nil {
		return nil
	}
	if o, ok := v.(*VMTable); ok {
		return o
	}
	if o, ok := v.(map[string]interface{}); ok {
		r := new(VMTable)
		r.Data = o
		return r
	}
	return nil
}

func (this *VM) API_newtable() *VMTable {
	r := new(VMTable)
	r.Data = make(map[string]interface{})
	return r
}

func (this *VM) API_newarray() []interface{} {
	return make([]interface{}, 0)
}
