package golua

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
)

func TableModule() *VMModule {
	m := NewVMModule("table")
	m.Init("concat", GOF_table_concat(0))
	m.Init("insert", GOF_table_insert(0))
	m.Init("remove", GOF_table_remove(0))
	m.Init("subtable", GOF_table_subtable(0))
	m.Init("newArray", GOF_table_newArray(0))
	return m
}

// table.concat(array[, sep,  start, end])
type GOF_table_concat int

func (this GOF_table_concat) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	arr, sep, start, end, err1 := vm.API_pop4X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	varr := vm.API_array(arr)
	if varr == nil {
		return 0, fmt.Errorf("array invalid(%T)", arr)
	}
	c := varr.Len()
	vsep := valutil.ToString(sep, "")
	vstart := valutil.ToInt(start, 0)
	vend := valutil.ToInt(end, c)
	buf := bytes.NewBuffer(make([]byte, 0, 16))
	for i := vstart; i < vend && i < c; i++ {
		if i != vstart {
			buf.WriteString(vsep)
		}
		av, err2 := varr.Get(vm, i)
		if err2 != nil {
			return 0, err2
		}
		buf.WriteString(valutil.ToString(av, ""))
	}

	vm.API_push(buf.String())
	return 1, nil
}

func (this GOF_table_concat) IsNative() bool {
	return true
}

func (this GOF_table_concat) String() string {
	return "GoFunc<table.concat>"
}

// table.insert(table[, pos], value)
type GOF_table_insert int

func (this GOF_table_insert) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	top := vm.API_gettop()
	arr, p1, p2, err1 := vm.API_pop3X(top, true)
	if err1 != nil {
		return 0, err1
	}
	varr := vm.API_array(arr)
	if varr == nil {
		return 0, fmt.Errorf("array invalid(%T)", arr)
	}
	c := varr.Len()
	var vpos int
	var vval interface{}
	if top == 2 {
		vpos = c
		vval = p1
	} else {
		vpos = valutil.ToInt(p1, c)
		vval = p2
	}
	if vpos >= c {
		varr.Add(vm, vval)
	} else {
		varr.Insert(vm, vpos, vval)
	}
	return 0, nil
}

func (this GOF_table_insert) IsNative() bool {
	return true
}

func (this GOF_table_insert) String() string {
	return "GoFunc<table.insert>"
}

// table.remove(table, pos) value
type GOF_table_remove int

func (this GOF_table_remove) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	o, pos, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	varr := vm.API_array(o)
	if varr != nil {
		vpos := valutil.ToInt(pos, -1)
		if vpos < 0 {
			return 0, fmt.Errorf("pos invalid(%v)", pos)
		}
		rv, err2 := varr.Get(vm, vpos)
		if err2 != nil {
			return 0, err2
		}
		err3 := varr.Delete(vm, vpos)
		if err3 != nil {
			return 0, err3
		}
		vm.API_push(rv)
		return 1, nil
	}
	vtb := vm.API_table(o)
	if vtb != nil {
		vkey := valutil.ToString(pos, "")
		rv := vtb.Rawget(vkey)
		vtb.Delete(vkey)
		vm.API_push(rv)
		return 1, nil
	}
	return 0, fmt.Errorf("table invalid(%T)", o)
}

func (this GOF_table_remove) IsNative() bool {
	return true
}

func (this GOF_table_remove) String() string {
	return "GoFunc<table.insert>"
}

// table.sub(table , int $start [, int $length ] ) table
type GOF_table_subtable int

func (this GOF_table_subtable) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	o, start, l, err1 := vm.API_pop3X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vo := vm.API_array(o)
	if vo == nil {
		return 0, fmt.Errorf("invalid array(%T)", o)
	}
	vstart := valutil.ToInt(start, 0)
	vlen := valutil.ToInt(l, -1)
	if vstart < 0 || vstart >= vo.Len() {
		return 0, fmt.Errorf("invalid start(%v) on array(%d)", vstart, vo.Len())
	}
	if vlen < 0 {
		vlen = vo.Len() - vstart
	}
	if vstart+vlen > vo.Len() {
		return 0, fmt.Errorf("invalid len(%v, %v) on array(%d)", vstart, l, vo.Len())
	}
	rv, err2 := vo.SubArray(vstart, vstart+vlen)
	if err2 != nil {
		return 0, err2
	}
	tmp := make([]interface{}, len(rv))
	copy(tmp, rv)
	vm.API_push(vm.API_array(tmp))
	return 1, nil
}

func (this GOF_table_subtable) IsNative() bool {
	return true
}

func (this GOF_table_subtable) String() string {
	return "GoFunc<table.subtable>"
}

// table.newArray(...) table
type GOF_table_newArray int

func (this GOF_table_newArray) Exec(vm *VM, self interface{}) (int, error) {
	top := vm.API_gettop()
	ns, err1 := vm.API_popN(top, true)
	if err1 != nil {
		return 0, err1
	}
	vm.API_push(vm.API_array(ns))
	return 1, nil
}

func (this GOF_table_newArray) IsNative() bool {
	return true
}

func (this GOF_table_newArray) String() string {
	return "GoFunc<table.newArray>"
}
