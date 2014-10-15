package golua

import (
	"fmt"
	"testing"
)

type GOF_test int

func (this GOF_test) Exec(vm *VM) (int, error) {
	v1, v2, err := vm.API_pop2X(-1, true)
	if err != nil {
		return 0, err
	}
	r := v1.(int) + v2.(int)
	vm.API_push(r)
	return 1, nil
}

func (this GOF_test) IsNative() bool {
	return true
}

func (this GOF_test) String() string {
	return "GoFunc<test>"
}

func T2estVMStack(t *testing.T) {
	safeCall()
	vmg := NewVMG("test")
	defer vmg.Close()
	vm, _ := vmg.CreateVM()

	if true {
		vm.API_insert(1, 1)
		fmt.Println("insert1", vm.DumpStack())
		vm.API_insert(1, 2)
		fmt.Println("insert2", vm.DumpStack())
		vm.API_remove(2)
		fmt.Println("remove", vm.DumpStack())
	}

}

func T2estVMAPI(t *testing.T) {
	safeCall()
	vmg := NewVMG("test")
	defer vmg.Close()
	vm, _ := vmg.CreateVM()

	if false {
		vm.API_insert(1, 1)
		fmt.Println(vm.DumpStack())
		vm.API_insert(1, 2)
		fmt.Println(vm.DumpStack())
		vm.API_remove(2)
		fmt.Println(vm.DumpStack())
	}

	if true {
		vm.API_push(1)
		vm.API_push(2)
		vm.API_push(3)
		vm.API_push(4)
		vm.API_push(5)
		fmt.Println(vm.DumpStack())
		r1, r2, r3, r4, err := vm.API_pop4X(-1, true)
		fmt.Println(vm.DumpStack())
		fmt.Println("result", r1, r2, r3, r4, err)
	}

	if false {
		f := GOF_print(0)
		vm.API_push(f)
		vm.API_push(1)
		vm.API_push("hello world")
		fmt.Println(vm.DumpStack())
		_, err1 := vm.Call(2, 1)
		if err1 != nil {
			t.Error(err1)
			return
		}
		fmt.Println(vm.DumpStack())
		r, err2 := vm.API_pop1X(-1, true)
		if err2 != nil {
			t.Error(err2)
			return
		}
		fmt.Println("call return ->", r)
		fmt.Println(vm.DumpStack())
	}

	if false {
		f := GOF_error(0)
		vm.API_push(f)
		vm.API_push("test error")
		fmt.Println(vm.DumpStack())
		_, err1 := vm.Call(1, 0)
		fmt.Println(vm.DumpStack())
		fmt.Println(err1)
	}

	if false {
		f := GOF_test(0)
		vm.API_push(f)
		vm.API_push(1)
		vm.API_push(2)
		vm.API_push(3)
		fmt.Println(vm.DumpStack())
		r, err1 := vm.Call(3, 1)
		fmt.Println(vm.DumpStack())
		fmt.Println(r, err1)
	}
}
