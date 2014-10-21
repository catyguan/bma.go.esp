package vmmsql

import (
	"database/sql"
	"golua"
)

type gooTx int

func (gooTx) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*sql.Tx); ok {
		switch key {
		case "Hours":
			return golua.NewGOF("Duration:Hours", func(vm *golua.VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj)
				return 1, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooTx) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooTx) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooTx) CanClose() bool {
	return false
}

func (gooTx) Close(o interface{}) {
}
