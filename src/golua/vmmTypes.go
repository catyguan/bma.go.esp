package golua

import (
	"bmautil/valutil"
	"fmt"
)

func TypesModule() *VMModule {
	m := NewVMModule("types")
	m.Init("name", GOF_types_name(0))
	m.Init("int32", GOF_types_int32(0))
	m.Init("int", GOF_types_int32(0))
	m.Init("int64", GOF_types_int64(0))
	m.Init("float", GOF_types_float(0))
	m.Init("string", GOF_types_string(0))
	m.Init("bool", GOF_types_bool(0))
	return m
}

// types.name(v)
type GOF_types_name int

func (this GOF_types_name) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	v, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	n := "<unknow>"
	n = fmt.Sprintf("%T", v)
	vm.API_push(n)
	return 1, nil
}

func (this GOF_types_name) IsNative() bool {
	return true
}

func (this GOF_types_name) String() string {
	return "GoFunc<types.name>"
}

// types.int32(v, defv)
type GOF_types_int32 int

func (this GOF_types_int32) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	v, dv, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vdv := valutil.ToInt32(dv, 0)
	vv := valutil.ToInt32(v, vdv)
	vm.API_push(vv)
	return 1, nil
}

func (this GOF_types_int32) IsNative() bool {
	return true
}

func (this GOF_types_int32) String() string {
	return "GoFunc<types.int32>"
}

// types.int64(v, defv)
type GOF_types_int64 int

func (this GOF_types_int64) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	v, dv, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vdv := valutil.ToInt64(dv, 0)
	vv := valutil.ToInt64(v, vdv)
	vm.API_push(vv)
	return 1, nil
}

func (this GOF_types_int64) IsNative() bool {
	return true
}

func (this GOF_types_int64) String() string {
	return "GoFunc<types.int64>"
}

// types.float(v, defv)
type GOF_types_float int

func (this GOF_types_float) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	v, dv, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vdv := valutil.ToFloat64(dv, 0)
	vv := valutil.ToFloat64(v, vdv)
	vm.API_push(vv)
	return 1, nil
}

func (this GOF_types_float) IsNative() bool {
	return true
}

func (this GOF_types_float) String() string {
	return "GoFunc<types.float>"
}

// types.string(v, defv)
type GOF_types_string int

func (this GOF_types_string) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	v, dv, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vdv := valutil.ToString(dv, "")
	vv := valutil.ToString(v, vdv)
	vm.API_push(vv)
	return 1, nil
}

func (this GOF_types_string) IsNative() bool {
	return true
}

func (this GOF_types_string) String() string {
	return "GoFunc<types.string>"
}

// types.bool(v, defv)
type GOF_types_bool int

func (this GOF_types_bool) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	v, dv, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vdv := valutil.ToBool(dv, false)
	vv := valutil.ToBool(v, vdv)
	vm.API_push(vv)
	return 1, nil
}

func (this GOF_types_bool) IsNative() bool {
	return true
}

func (this GOF_types_bool) String() string {
	return "GoFunc<types.bool>"
}
