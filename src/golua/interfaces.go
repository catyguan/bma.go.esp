package golua

import (
	"errors"
	"fmt"
)

type GoFunction interface {
	Exec(vm *VM) (int, error)
	IsNative() bool
}

type supportFuncName interface {
	FuncName() (string, string)
}

type ER int

var (
	ER_ERROR    = ER(0)
	ER_NEXT     = ER(1)
	ER_BREAK    = ER(2)
	ER_CONTINUE = ER(3)
	ER_RETURN   = ER(4)
)

type VMVar interface {
	Get(vm *VM) (interface{}, error)
	Set(vm *VM, v interface{}) (bool, error)
}

const (
	METATABLE_INDEX    = "__index"
	METATABLE_NEWINDEX = "__newindex"
)

type VMTable interface {
	Get(vm *VM, key string) (interface{}, error)
	Rawget(key string) interface{}
	Set(vm *VM, key string, val interface{}) error
	Rawset(key string, val interface{})
	Delete(key string)
	Len() int
	ToMap() map[string]interface{}
}

func AssertNil(n string, v interface{}) error {
	if v == nil {
		if n != "" {
			return fmt.Errorf("%s null pointer", n)
		}
		return errors.New("null pointer")
	}
	return nil
}
