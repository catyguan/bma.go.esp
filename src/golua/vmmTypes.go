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

// types.name(v[, extInfo:bool])
type GOF_types_name int

func (this GOF_types_name) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	v, ext, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vext := valutil.ToBool(ext, false)
	n := "<unknow>"
	einfo := ""

	if v == nil {
		n = "nil"
	} else {
		switch ro := v.(type) {
		case bool:
			n = "bool"
		case int, uint, int8, uint8, int16, uint16, int32, uint32:
			n = "int32"
		case int64, uint64:
			n = "int64"
		case float32, float64:
			n = "float"
		case string:
			n = "string"
		case map[string]interface{}:
			n = "table"
		case VMTable:
			n = "table"
			if tmp, ok := ro.(*objectVMTable); ok {
				n = "object"
				if vext {
					einfo = fmt.Sprintf("%T", tmp.o)
				}
			}
		case []interface{}:
			n = "array"
		case VMArray:
			n = "array"
			einfo = "VMArray"
		}
	}
	if einfo == "" && vext {
		einfo = fmt.Sprintf("%T", v)
	}

	r := 1
	vm.API_push(n)
	if vext {
		vm.API_push(einfo)
		r = 2
	}
	return r, nil
}

func (this GOF_types_name) IsNative() bool {
	return true
}

func (this GOF_types_name) String() string {
	return "GoFunc<types.name>"
}

func types_pop_v(vm *VM) (interface{}, interface{}, error) {
	c := vm.API_gettop()
	switch c {
	case 0:
	case 1:
		v, err1 := vm.API_pop1X(-1, true)
		if err1 != nil {
			return nil, nil, err1
		}
		return v, nil, nil
	case 2:
		v, dv, err1 := vm.API_pop2X(-1, true)
		if err1 != nil {
			return nil, nil, err1
		}
		return v, dv, nil
	default:
		o, m, dv, err1 := vm.API_pop3X(-1, true)
		if err1 != nil {
			return nil, nil, err1
		}
		v, err2 := vm.API_getMember(o, m)
		if err2 != nil {
			return nil, nil, err2
		}
		return v, dv, nil
	}
	return nil, nil, vm.API_checkStack(1)

}

// types.int32(v, defv)
type GOF_types_int32 int

func (this GOF_types_int32) Exec(vm *VM, self interface{}) (int, error) {
	v, dv, err1 := types_pop_v(vm)
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

func (this GOF_types_int64) Exec(vm *VM, self interface{}) (int, error) {
	v, dv, err1 := types_pop_v(vm)
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

func (this GOF_types_float) Exec(vm *VM, self interface{}) (int, error) {
	v, dv, err1 := types_pop_v(vm)
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

func (this GOF_types_string) Exec(vm *VM, self interface{}) (int, error) {
	v, dv, err1 := types_pop_v(vm)
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

func (this GOF_types_bool) Exec(vm *VM, self interface{}) (int, error) {
	v, dv, err1 := types_pop_v(vm)
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
