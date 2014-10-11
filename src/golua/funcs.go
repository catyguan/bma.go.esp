package golua

import (
	"bytes"
	"errors"
	"fmt"
)

type GOF_print int

func (this GOF_print) Exec(vm *VM) (int, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	top := vm.API_gettop()
	for i := 1; i <= top; i++ {
		v, err := vm.API_peek(i)
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

type GOF_error int

func (this GOF_error) Exec(vm *VM) (int, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	top := vm.API_gettop()
	for i := 1; i <= top; i++ {
		v, err := vm.API_peek(i)
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
