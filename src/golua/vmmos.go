package golua

import (
	"bmautil/valutil"
	"os"
)

func OSModule() *VMModule {
	m := NewVMModule("os")
	m.Init("mkdir", GOF_os_mkdir(0))
	m.Init("remove", GOF_os_remove(0))
	m.Init("rename", GOF_os_rename(0))
	m.Init("createFile", GOF_os_createFile(0))
	m.Init("openFile", GOF_os_openFile(0))
	m.Init("fileExists", GOF_os_fileExists(0))
	return m
}

// os.mkdir(dir string, bool all) bool
type GOF_os_mkdir int

func (this GOF_os_mkdir) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	dir, all, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vdir := valutil.ToString(dir, "")
	vall := valutil.ToBool(all, false)
	if vall {
		err := os.MkdirAll(vdir, os.ModePerm)
		if err != nil {
			return 0, err
		}
	} else {
		err := os.Mkdir(vdir, os.ModePerm)
		if err != nil {
			return 0, err
		}
	}
	return 0, nil
}

func (this GOF_os_mkdir) IsNative() bool {
	return true
}

func (this GOF_os_mkdir) String() string {
	return "GoFunc<os.mkdir>"
}

// os.delete(name string, bool all) bool
type GOF_os_remove int

func (this GOF_os_remove) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, all, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn := valutil.ToString(n, "")
	vall := valutil.ToBool(all, false)
	if vall {
		err := os.RemoveAll(vn)
		if err != nil {
			return 0, err
		}
	} else {
		err := os.Remove(vn)
		if err != nil {
			return 0, err
		}
	}
	return 0, nil
}

func (this GOF_os_remove) IsNative() bool {
	return true
}

func (this GOF_os_remove) String() string {
	return "GoFunc<os.remove>"
}

// os.rename(name string, newName string) bool
type GOF_os_rename int

func (this GOF_os_rename) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	n1, n2, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn1 := valutil.ToString(n1, "")
	vn2 := valutil.ToString(n2, "")
	err := os.Rename(vn1, vn2)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (this GOF_os_rename) IsNative() bool {
	return true
}

func (this GOF_os_rename) String() string {
	return "GoFunc<os.rename>"
}

// os.createFile(name string) File
type GOF_os_createFile int

func (this GOF_os_createFile) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n1, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn1 := valutil.ToString(n1, "")
	f, err := os.Create(vn1)
	if err != nil {
		return 0, err
	}
	vm.API_push(NewSafeGoFile(vm, f))
	return 1, nil
}

func (this GOF_os_createFile) IsNative() bool {
	return true
}

func (this GOF_os_createFile) String() string {
	return "GoFunc<os.createFile>"
}

// os.openFile(name string) File
type GOF_os_openFile int

func (this GOF_os_openFile) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n1, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn1 := valutil.ToString(n1, "")
	f, err := os.OpenFile(vn1, os.O_RDWR, os.ModePerm)
	if err != nil {
		return 0, err
	}
	vm.API_push(NewSafeGoFile(vm, f))
	return 1, nil
}

func (this GOF_os_openFile) IsNative() bool {
	return true
}

func (this GOF_os_openFile) String() string {
	return "GoFunc<os.openFile>"
}

// os.fileExists(name string) bool
type GOF_os_fileExists int

func (this GOF_os_fileExists) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n1, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	r := true
	vn1 := valutil.ToString(n1, "")
	f, err := os.Open(vn1)
	if err != nil {
		if os.IsNotExist(err) {
			r = false
		} else {
			return 0, err
		}
	}
	if f != nil {
		f.Close()
	}
	vm.API_push(r)
	return 1, nil
}

func (this GOF_os_fileExists) IsNative() bool {
	return true
}

func (this GOF_os_fileExists) String() string {
	return "GoFunc<os.fileExists>"
}
