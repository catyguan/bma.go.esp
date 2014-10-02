package golua

import (
	"fmt"
	"testing"
)

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
		f := GOF_print(0)
		vm.API_push(f)
		vm.API_push(1)
		vm.API_push("hello world")
		fmt.Println(vm.DumpStack())
		err1 := vm.Call(2, 1)
		if err1 != nil {
			t.Error(err1)
			return
		}
		fmt.Println(vm.DumpStack())
		r, err2 := vm.API_pop1()
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
		err1 := vm.Call(1, 0)
		fmt.Println(vm.DumpStack())
		fmt.Println(err1)
	}
}
