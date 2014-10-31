package golua

import (
	"fmt"
	"sync"
)

type VMModule struct {
	name    string
	funcs   map[string]GoFunction
	mutex   sync.RWMutex
	members map[string]interface{}
}

func NewVMModule(n string) *VMModule {
	r := new(VMModule)
	r.name = n
	r.funcs = make(map[string]GoFunction)
	return r
}

func (this *VMModule) Name() string {
	return this.name
}

func (this *VMModule) Bind(gl *GoLua) {
	gl.SetGlobal(this.name, this)
}

func (this *VMModule) Init(key string, f GoFunction) {
	this.funcs[key] = f
}

func (this *VMModule) String() string {
	return fmt.Sprintf("Module<%s>", this.name)
}

func (this *VMModule) Get(vm *VM, key string) (interface{}, error) {
	return this.Rawget(key), nil
}
func (this *VMModule) Rawget(key string) interface{} {
	if f, ok := this.funcs[key]; ok {
		return f
	}
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	if this.members != nil {
		if v, ok := this.members[key]; ok {
			return v
		}
	}
	return nil
}
func (this *VMModule) Set(vm *VM, key string, val interface{}) error {
	this.Rawset(key, val)
	return nil
}
func (this *VMModule) Rawset(key string, val interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.members == nil {
		this.members = make(map[string]interface{})
	}
	this.members[key] = val
}
func (this *VMModule) Delete(key string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.members != nil {
		delete(this.members, key)
	}
}
func (this *VMModule) Len() int {
	return len(this.funcs) + len(this.members)
}
func (this *VMModule) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	for k, v := range this.members {
		r[k] = v
	}
	return r
}
