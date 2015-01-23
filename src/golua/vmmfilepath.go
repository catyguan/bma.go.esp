package golua

import (
	"bmautil/valutil"
	"path"
	"path/filepath"
)

func FilePathModule() *VMModule {
	m := NewVMModule("filepath")
	m.Init("base", GOF_filepath_base(0))
	m.Init("clean", GOF_filepath_clean(0))
	m.Init("dir", GOF_filepath_dir(0))
	m.Init("ext", GOF_filepath_ext(0))
	m.Init("isAbs", GOF_filepath_isAbs(0))
	m.Init("join", GOF_filepath_join(0))
	m.Init("match", GOF_filepath_match(0))
	m.Init("abs", GOF_filepath_abs(0))
	m.Init("rel", GOF_filepath_rel(0))
	m.Init("changeExt", GOF_filepath_changeExt(0))
	return m
}

// filepath.base(s string) string
type GOF_filepath_base int

func (this GOF_filepath_base) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	rv := filepath.Base(vs)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_filepath_base) IsNative() bool {
	return true
}

func (this GOF_filepath_base) String() string {
	return "GoFunc<filepath.base>"
}

// filepath.clean(s string) string
type GOF_filepath_clean int

func (this GOF_filepath_clean) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	rv := path.Clean(vs)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_filepath_clean) IsNative() bool {
	return true
}

func (this GOF_filepath_clean) String() string {
	return "GoFunc<filepath.clean>"
}

// filepath.dir(s string) string
type GOF_filepath_dir int

func (this GOF_filepath_dir) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	rdir := filepath.Dir(vs)
	_, rfile := filepath.Split(vs)
	vm.API_push(rdir)
	vm.API_push(rfile)
	return 2, nil
}

func (this GOF_filepath_dir) IsNative() bool {
	return true
}

func (this GOF_filepath_dir) String() string {
	return "GoFunc<filepath.dir>"
}

// filepath.ext(s string) string
type GOF_filepath_ext int

func (this GOF_filepath_ext) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	rv := filepath.Ext(vs)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_filepath_ext) IsNative() bool {
	return true
}

func (this GOF_filepath_ext) String() string {
	return "GoFunc<filepath.ext>"
}

// filepath.isAbs(s string) bool
type GOF_filepath_isAbs int

func (this GOF_filepath_isAbs) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	rv := filepath.IsAbs(vs)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_filepath_isAbs) IsNative() bool {
	return true
}

func (this GOF_filepath_isAbs) String() string {
	return "GoFunc<filepath.isAbs>"
}

// filepath.join(s string...) string
type GOF_filepath_join int

func (this GOF_filepath_join) Exec(vm *VM, self interface{}) (int, error) {
	top := vm.API_gettop()
	ns, err1 := vm.API_popN(top, true)
	if err1 != nil {
		return 0, err1
	}
	slist := make([]string, len(ns))
	for i, v := range ns {
		slist[i] = valutil.ToString(v, "")
	}
	rv := filepath.Join(slist...)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_filepath_join) IsNative() bool {
	return true
}

func (this GOF_filepath_join) String() string {
	return "GoFunc<filepath.join>"
}

// filepath.match(pattern, name string) bool
type GOF_filepath_match int

func (this GOF_filepath_match) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	p, n, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vp := valutil.ToString(p, "")
	vn := valutil.ToString(n, "")
	rv, _ := filepath.Match(vp, vn)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_filepath_match) IsNative() bool {
	return true
}

func (this GOF_filepath_match) String() string {
	return "GoFunc<filepath.match>"
}

// filepath.abs(s string) (string,bool)
type GOF_filepath_abs int

func (this GOF_filepath_abs) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	rv, err := filepath.Abs(vs)
	vm.API_push(rv)
	vm.API_push(err == nil)
	return 1, nil
}

func (this GOF_filepath_abs) IsNative() bool {
	return true
}

func (this GOF_filepath_abs) String() string {
	return "GoFunc<filepath.abs>"
}

// filepath.rel(basepath,targpath string) (string, bool)
type GOF_filepath_rel int

func (this GOF_filepath_rel) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s1, s2, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs1 := valutil.ToString(s1, "")
	vs2 := valutil.ToString(s2, "")
	rv, err := filepath.Rel(vs1, vs2)
	if err != nil {
		vm.API_push(nil)
		vm.API_push(false)
	} else {
		vm.API_push(rv)
		vm.API_push(true)
	}
	return 2, nil
}

func (this GOF_filepath_rel) IsNative() bool {
	return true
}

func (this GOF_filepath_rel) String() string {
	return "GoFunc<filepath.rel>"
}

// filepath.changeExt(fn, ext string) string
type GOF_filepath_changeExt int

func (this GOF_filepath_changeExt) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s1, s2, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs1 := valutil.ToString(s1, "")
	vs2 := valutil.ToString(s2, "")

	var npath string
	ext := filepath.Ext(vs1)
	if ext != "" {
		npath = vs1[:len(vs1)-len(ext)]
	} else {
		npath = vs1
	}
	npath = npath + vs2

	vm.API_push(npath)
	return 1, nil
}

func (this GOF_filepath_changeExt) IsNative() bool {
	return true
}

func (this GOF_filepath_changeExt) String() string {
	return "GoFunc<filepath.changeExt>"
}
