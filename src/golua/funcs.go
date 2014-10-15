package golua

import (
	"bmautil/valutil"
	"bytes"
	"errors"
	"fmt"
)

// print(...)
type GOF_print int

func (this GOF_print) Exec(vm *VM) (int, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	top := vm.API_gettop()
	for i := 1; i <= top; i++ {
		v, err := vm.API_peek(i, true)
		if err != nil {
			return 0, err
		}
		v, err = vm.API_value(v)
		if err != nil {
			return 0, err
		}
		if i != 1 {
			buf.WriteString("\t")
		}
		buf.WriteString(fmt.Sprintf("%v", v))
	}
	fmt.Println(buf.String())
	vm.API_pop(top)
	return 0, nil
}

func (this GOF_print) IsNative() bool {
	return true
}

func (this GOF_print) String() string {
	return "GoFunc<print>"
}

// error(...)
type GOF_error int

func (this GOF_error) Exec(vm *VM) (int, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	top := vm.API_gettop()
	for i := 1; i <= top; i++ {
		v, err := vm.API_peek(i, true)
		if err != nil {
			return 0, err
		}
		v, err = vm.API_value(v)
		if err != nil {
			return 0, err
		}
		if i != 1 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("%v", v))
	}
	vm.API_pop(top)
	return 0, errors.New(buf.String())
}

func (this GOF_error) IsNative() bool {
	return true
}

func (this GOF_error) String() string {
	return "GoFunc<error>"
}

// setmetatable(table, metatable)
type GOF_setmetatable int

func (this GOF_setmetatable) Exec(vm *VM) (int, error) {
	err00 := vm.API_checkstack(2)
	if err00 != nil {
		return 0, err00
	}
	t1, mt, err0 := vm.API_pop2X(-1, true)
	if err0 != nil {
		return 0, err0
	}
	err0 = AssertNil("table", t1)
	if err0 != nil {
		return 0, err0
	}
	err0 = AssertNil("metatable", mt)
	if err0 != nil {
		return 0, err0
	}
	vt, ok := t1.(*CommonVMTable)
	if !ok {
		return 0, fmt.Errorf("invalid table for setmetatable")
	}
	vmt := vm.API_table(mt)
	if vmt == nil {
		return 0, fmt.Errorf("invalid metatable for setmetatable")
	}
	vt.SetMetaTable(vmt)
	return 0, nil
}

func (this GOF_setmetatable) IsNative() bool {
	return true
}

func (this GOF_setmetatable) String() string {
	return "GoFunc<setmetatable>"
}

// getmetatable(table) metatable
type GOF_getmetatable int

func (this GOF_getmetatable) Exec(vm *VM) (int, error) {
	err00 := vm.API_checkstack(1)
	if err00 != nil {
		return 0, err00
	}
	t1, err0 := vm.API_pop1X(-1, true)
	if err0 != nil {
		return 0, err0
	}
	err0 = AssertNil("table", t1)
	if err0 != nil {
		return 0, err0
	}
	vt, ok := t1.(*CommonVMTable)
	if !ok {
		return 0, fmt.Errorf("invalid table for setmetatable")
	}
	vmt := vt.GetMetaTable()
	vm.API_push(vmt)
	return 1, nil
}

func (this GOF_getmetatable) IsNative() bool {
	return true
}

func (this GOF_getmetatable) String() string {
	return "GoFunc<getmetatable>"
}

// rawget(table, key) value
type GOF_rawget int

func (this GOF_rawget) Exec(vm *VM) (int, error) {
	err00 := vm.API_checkstack(2)
	if err00 != nil {
		return 0, err00
	}
	t1, key, err0 := vm.API_pop2X(-1, true)
	if err0 != nil {
		return 0, err0
	}
	err0 = AssertNil("table", t1)
	if err0 != nil {
		return 0, err0
	}
	err0 = AssertNil("key", key)
	if err0 != nil {
		return 0, err0
	}
	vt := vm.API_table(t1)
	if vt == nil {
		return 0, fmt.Errorf("invalid table for rawget")
	}
	k2 := valutil.ToString(key, "")
	r := vt.Rawget(k2)
	vm.API_push(r)
	return 1, nil
}

func (this GOF_rawget) IsNative() bool {
	return true
}

func (this GOF_rawget) String() string {
	return "GoFunc<rawget>"
}

// rawset(table, key[, value])
type GOF_rawset int

func (this GOF_rawset) Exec(vm *VM) (int, error) {
	err00 := vm.API_checkstack(2)
	if err00 != nil {
		return 0, err00
	}
	t1, key, val, err0 := vm.API_pop3X(-1, true)
	if err0 != nil {
		return 0, err0
	}
	err0 = AssertNil("table", t1)
	if err0 != nil {
		return 0, err0
	}
	err0 = AssertNil("key", key)
	if err0 != nil {
		return 0, err0
	}

	vt := vm.API_table(t1)
	if vt == nil {
		return 0, fmt.Errorf("invalid table for rawget")
	}
	k2 := valutil.ToString(key, "")
	vt.Rawset(k2, val)

	return 0, nil
}

func (this GOF_rawset) IsNative() bool {
	return true
}

func (this GOF_rawset) String() string {
	return "GoFunc<rawset>"
}

// core module
func CoreModule(vmg *VMG) {
	vmg.SetGlobal("print", GOF_print(0))
	vmg.SetGlobal("error", GOF_error(0))
	vmg.SetGlobal("setmetatable", GOF_setmetatable(0))
	vmg.SetGlobal("getmetatable", GOF_getmetatable(0))
	vmg.SetGlobal("rawget", GOF_rawget(0))
	vmg.SetGlobal("rawset", GOF_rawset(0))
}
