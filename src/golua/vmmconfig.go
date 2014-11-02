package golua

import (
	"bmautil/valutil"
	"fmt"
)

func ConfigModule() *VMModule {
	m := NewVMModule("config")
	m.Init("set", GOF_config_set(0))
	m.Init("get", GOF_config_get(0))
	m.Init("query", GOF_config_query(0))
	m.Init("asTable", GOF_config_asTable(0))
	m.Init("parse", GOF_config_parse(0))
	return m
}

// config.set(n, v)
type GOF_config_set int

func (this GOF_config_set) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	n, v, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn := valutil.ToString(n, "")
	ok := vm.gl.SetConfig(vn, v)
	vm.API_push(ok)
	return 1, nil
}

func (this GOF_config_set) IsNative() bool {
	return true
}

func (this GOF_config_set) String() string {
	return "GoFunc<config.set>"
}

// config.get(n)
type GOF_config_get int

func (this GOF_config_get) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn := valutil.ToString(n, "")
	v, ok := vm.gl.GetConfig(vn)
	if ok {
		vm.API_push(v)
		return 1, nil
	} else {
		return 0, fmt.Errorf("invalid config('%s')", n)
	}
}

func (this GOF_config_get) IsNative() bool {
	return true
}

func (this GOF_config_get) String() string {
	return "GoFunc<config.get>"
}

// config.query(n, def)
type GOF_config_query int

func (this GOF_config_query) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, defv, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn := valutil.ToString(n, "")
	v, ok := vm.gl.GetConfig(vn)
	if ok {
		vm.API_push(v)
	} else {
		vm.API_push(defv)
	}
	return 1, nil
}

func (this GOF_config_query) IsNative() bool {
	return true
}

func (this GOF_config_query) String() string {
	return "GoFunc<config.query>"
}

// config.parse(v string) string
type GOF_config_parse int

func (this GOF_config_parse) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn := valutil.ToString(n, "")
	v, err2 := vm.gl.ParseConfig(vn)
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(v)
	return 1, nil
}

func (this GOF_config_parse) IsNative() bool {
	return true
}

func (this GOF_config_parse) String() string {
	return "GoFunc<config.parse>"
}

// config.asTable(prex)
type GOF_config_asTable int

func (this GOF_config_asTable) Exec(vm *VM, self interface{}) (int, error) {
	prex, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vprex := valutil.ToString(prex, "")
	r := new(configTable)
	r.prex = vprex
	vm.API_push(r)
	return 1, nil
}

func (this GOF_config_asTable) IsNative() bool {
	return true
}

func (this GOF_config_asTable) String() string {
	return "GoFunc<config.asTable>"
}

// configTable
type configTable struct {
	prex string
}

func (this *configTable) String() string {
	return fmt.Sprintf("@configTable[%s]", this.prex)
}

func (this *configTable) Get(vm *VM, key string) (interface{}, error) {
	n := key
	if this.prex != "" {
		n = this.prex + "." + n
	}
	v, ok := vm.gl.GetConfig(n)
	if ok {
		return v, nil
	}
	return nil, nil
}

func (this *configTable) Rawget(key string) interface{} {
	return nil
}

func (this *configTable) Set(vm *VM, key string, val interface{}) error {
	return nil
}

func (this *configTable) Rawset(key string, val interface{}) {

}

func (this *configTable) Delete(key string) {

}

func (this *configTable) Len() int {
	return 0
}

func (this *configTable) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	return r
}
