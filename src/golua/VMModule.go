package golua

import (
	"fmt"
)

type VMModule struct {
	name  string
	funcs map[string]GoFunction
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

func (this *VMModule) Bind(vmg *VMG) {
	vmg.SetGlobal(this.name, this)
}

func (this *VMModule) Init(key string, f GoFunction) {
	this.funcs[key] = f
}

func (this *VMModule) String() string {
	return fmt.Sprintf("Module<%s>", this.name)
}

func (this *VMModule) Get(vm *VM, key string) (interface{}, error) {
	return this.funcs[key], nil
}
func (this *VMModule) Rawget(key string) interface{} {
	return this.funcs[key]
}
func (this *VMModule) Set(vm *VM, key string, val interface{}) error {
	return nil
}
func (this *VMModule) Rawset(key string, val interface{}) {

}
func (this *VMModule) Delete(key string) {

}
func (this *VMModule) Len() int {
	return len(this.funcs)
}
func (this *VMModule) ToMap() map[string]interface{} {
	return make(map[string]interface{})
}
