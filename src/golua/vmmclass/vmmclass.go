package vmmclass

import (
	"bmautil/valutil"
	"fmt"
	"golua"
	"sync"
)

const tag = "vmmclass"

type classes struct {
	lock sync.RWMutex
	data map[string]*classVMTable
}

func (this *classes) Get(n string) *classVMTable {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.data[n]
}

func (this *classes) Put(n string, o *classVMTable) *classVMTable {
	this.lock.Lock()
	defer this.lock.Unlock()
	if c, ok := this.data[n]; ok {
		return c
	}
	this.data[n] = o
	return o
}

func (this *classes) Close() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, c := range this.data {
		delete(this.data, k)
		c.Clear()
	}
}

func getClasses(vm *golua.VM) (*classes, error) {
	css, _ := vm.GetGoLua().GetService("vmmclass")
	if css == nil {
		return nil, fmt.Errorf("vmmclass not init")
	}
	if o, ok := css.(*classes); ok {
		return o, nil
	}
	return nil, fmt.Errorf("vmmclass service invalid(%T)", css)
}

func InitGoLua(gl *golua.GoLua) {
	Module().Bind(gl)
	o := new(classes)
	o.data = make(map[string]*classVMTable)
	gl.SetService("vmmclass", o)
}

func Module() *golua.VMModule {
	m := golua.NewVMModule("class")
	m.Init("forName", GOF_class_forName(0))
	m.Init("define", GOF_class_define(0))
	m.Init("new", GOF_class_new(0))
	m.Init("is", GOF_class_is(0))
	m.Init("check", GOF_class_check(0))
	return m
}

// class.forName(name string) class
type GOF_class_forName int

func (this GOF_class_forName) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	css, err1x := getClasses(vm)
	if err1x != nil {
		return 0, err1x
	}
	vn := valutil.ToString(n, "")
	if vn == "" {
		return 0, fmt.Errorf("className empty")
	}
	ro := css.Get(vn)
	vm.API_push(ro)
	return 1, nil
}

func (this GOF_class_forName) IsNative() bool {
	return true
}

func (this GOF_class_forName) String() string {
	return "GoFunc<class.forName>"
}

// class.define(name string [,supers []string]) class
type GOF_class_define int

func (this GOF_class_define) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, sp, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	css, err1x := getClasses(vm)
	if err1x != nil {
		return 0, err1x
	}
	vn := valutil.ToString(n, "")
	if vn == "" {
		return 0, fmt.Errorf("className empty")
	}
	var spl []interface{}
	if sp != nil {
		if str, ok := sp.(string); ok {
			spl = []interface{}{str}
		} else {
			spl = vm.API_toSlice(sp)
			if spl == nil {
				return 0, fmt.Errorf("superName not array(%T)", sp)
			}
		}
	}
	ro := css.Get(vn)
	if ro == nil {
		clist := make([]*classVMTable, len(spl))
		if len(spl) > 0 {
			for i, sn := range spl {
				vsn := valutil.ToString(sn, "")
				if vsn == "" {
					return 0, fmt.Errorf("superName empty")
				}
				cls := css.Get(vsn)
				if cls == nil {
					return 0, fmt.Errorf("super class(%s) not exists", vsn)
				}
				if cls.HasBase(vn) {
					return 0, fmt.Errorf("'%s' is super class of '%s'", vn, vsn)
				}
				clist[i] = cls
			}
		}
		ro = css.Put(vn, newClassVMTable(vn, clist))
	}
	vm.API_push(ro)
	return 1, nil
}

func (this GOF_class_define) IsNative() bool {
	return true
}

func (this GOF_class_define) String() string {
	return "GoFunc<class.define>"
}

// class.new(name string) class
type GOF_class_new int

func (this GOF_class_new) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	css, err1x := getClasses(vm)
	if err1x != nil {
		return 0, err1x
	}
	vn := valutil.ToString(n, "")
	if vn == "" {
		return 0, fmt.Errorf("className empty")
	}
	cls := css.Get(vn)
	if cls == nil {
		return 0, fmt.Errorf("class(%s) not exists", vn)
	}
	ro, err3 := cls.New(vm, nil)
	if err3 != nil {
		return 0, err3
	}
	vm.API_push(ro)
	return 1, nil
}

func (this GOF_class_new) IsNative() bool {
	return true
}

func (this GOF_class_new) String() string {
	return "GoFunc<class.new>"
}

// class.is(v object or class or className, className string) bool
type GOF_class_is int

func (this GOF_class_is) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	x, n, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	css, err1x := getClasses(vm)
	if err1x != nil {
		return 0, err1x
	}
	vn := valutil.ToString(n, "")
	if vn == "" {
		return 0, fmt.Errorf("className empty")
	}
	b := false
	switch v := x.(type) {
	case *ciVMTable:
		b = v.cls.HasBase(vn)
	case *classVMTable:
		b = v.HasBase(vn)
	case string:
		cls := css.Get(v)
		if cls == nil {
			return 0, fmt.Errorf("class(%s) not exists", v)
		}
		b = cls.HasBase(vn)
	default:
		return 0, fmt.Errorf("invalid param(%v)", x)
	}
	vm.API_push(b)
	return 1, nil
}

func (this GOF_class_is) IsNative() bool {
	return true
}

func (this GOF_class_is) String() string {
	return "GoFunc<class.is>"
}

// class.check(o object, className string) bool
type GOF_class_check int

func (this GOF_class_check) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	x, n, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	css, err1x := getClasses(vm)
	if err1x != nil {
		return 0, err1x
	}
	vn := valutil.ToString(n, "")
	if vn == "" {
		return 0, fmt.Errorf("className empty")
	}
	o, ok := x.(golua.VMTable)
	if !ok {
		return 0, fmt.Errorf("invalid param(%v)", o)
	}
	cls := css.Get(vn)
	if cls == nil {
		return 0, fmt.Errorf("class(%s) not exists", vn)
	}
	b := true
	ss := make([]string, 0)
	for _, k := range cls.Props() {
		v, errX := o.Get(vm, k)
		if errX != nil {
			return 0, errX
		}
		if v == nil {
			b = false
			ss = append(ss, k)
		}
	}
	vm.API_push(b)
	vm.API_push(ss)
	return 2, nil
}

func (this GOF_class_check) IsNative() bool {
	return true
}

func (this GOF_class_check) String() string {
	return "GoFunc<class.check>"
}
