package golua

import (
	"bmautil/valutil"
	"bytes"
	"errors"
	"fmt"
)

// print(...)
type GOF_print int

func (this GOF_print) Exec(vm *VM, self interface{}) (int, error) {
	// fmt.Println("PRINT", vm.DumpStack())
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
	vm.API_popAll()
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

func (this GOF_error) Exec(vm *VM, self interface{}) (int, error) {
	top := vm.API_gettop()
	if top == 0 {
		return 0, errors.New("<unknow error>")
	}
	ns, err1 := vm.API_popN(top, true)
	if err1 != nil {
		return 0, err1
	}
	errf := valutil.ToString(ns[0], "")
	if top > 1 {
		return 0, fmt.Errorf(errf, ns[1:]...)
	}
	return 0, errors.New(errf)
}

func (this GOF_error) IsNative() bool {
	return true
}

func (this GOF_error) String() string {
	return "GoFunc<error>"
}

// rawget(table, key) value
type GOF_rawget int

func (this GOF_rawget) Exec(vm *VM, self interface{}) (int, error) {
	err00 := vm.API_checkStack(2)
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
	k2 := valutil.ToString(key, "")
	vt, m := vm.API_table(t1)
	if m != nil {
		r := m[k2]
		vm.API_push(r)
		return 1, nil
	}
	if vt != nil {
		r := vt.Rawget(k2)
		vm.API_push(r)
		return 1, nil
	}
	return 0, fmt.Errorf("invalid table for rawget")
}

func (this GOF_rawget) IsNative() bool {
	return true
}

func (this GOF_rawget) String() string {
	return "GoFunc<rawget>"
}

// rawset(table, key[, value])
type GOF_rawset int

func (this GOF_rawset) Exec(vm *VM, self interface{}) (int, error) {
	err00 := vm.API_checkStack(2)
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
	k2 := valutil.ToString(key, "")
	vt, m := vm.API_table(t1)
	if m != nil {
		if val == nil {
			delete(m, k2)
		} else {
			m[k2] = val
		}
		return 0, nil
	}
	if vt != nil {
		vt.Rawset(k2, val)
		return 0, nil
	}
	return 0, fmt.Errorf("invalid table for rawget")
}

func (this GOF_rawset) IsNative() bool {
	return true
}

func (this GOF_rawset) String() string {
	return "GoFunc<rawset>"
}

// pcall(f, ...) true, ... | false, error
type GOF_pcall int

func (this GOF_pcall) Exec(vm *VM, self interface{}) (int, error) {
	top := vm.API_gettop()
	if top == 0 {
		vm.API_push(true)
		return 1, nil
	}
	f, err0 := vm.API_peek(-top, true)
	if err0 != nil {
		vm.API_popAll()
		vm.API_push(false)
		vm.API_push(err0)
		return 2, nil
	}
	if !vm.API_canCall(f) {
		err1 := fmt.Errorf("pcall func(%T) can't call", f)
		vm.API_popAll()
		vm.API_push(false)
		vm.API_push(err1)
		return 2, nil
	}
	r, err2 := vm.Call(top-1, -1, nil)
	if err2 != nil {
		vm.API_popAll()
		vm.API_push(false)
		vm.API_push(err2)
		return 2, nil
	}
	if r == 0 {
		vm.API_push(true)
		return 1, nil
	} else {
		vm.API_insert(-r, true)
		return r + 1, nil
	}
}

func (this GOF_pcall) IsNative() bool {
	return true
}

func (this GOF_pcall) String() string {
	return "GoFunc<pcall>"
}

// require(scriptName)
type GOF_require int

func (this GOF_require) Exec(vm *VM, self interface{}) (int, error) {
	n, err0 := vm.API_pop1X(-1, true)
	if err0 != nil {
		return 0, err0
	}
	vn := valutil.ToString(n, "")
	if vn == "" {
		return 0, fmt.Errorf("require script invalid(%v)", n)
	}
	err2 := vm.API_require(vn)
	if err2 != nil {
		return 0, err2
	}
	return 0, nil
}

func (this GOF_require) IsNative() bool {
	return true
}

func (this GOF_require) String() string {
	return "GoFunc<require>"
}

// core module
func CoreModule(gl *GoLua) {
	gl.SetGlobal("print", GOF_print(0))
	gl.SetGlobal("error", GOF_error(0))
	gl.SetGlobal("rawget", GOF_rawget(0))
	gl.SetGlobal("rawset", GOF_rawset(0))
	gl.SetGlobal("pcall", GOF_pcall(0))
	gl.SetGlobal("require", GOF_require(0))
}

func InitCoreLibs(gl *GoLua) {
	CoreModule(gl)
	GoModule().Bind(gl)
	TypesModule().Bind(gl)
	TableModule().Bind(gl)
	InitGoLuaStringsModule(gl)
	TimeModule().Bind(gl)
	ConfigModule().Bind(gl)
	FilePathModule().Bind(gl)
	OSModule().Bind(gl)
}
