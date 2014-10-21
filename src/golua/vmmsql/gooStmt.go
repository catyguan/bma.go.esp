package vmmsql

import (
	"database/sql"
	"golua"
)

type gooStmt int

func (gooStmt) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*sql.Stmt); ok {
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

func (gooStmt) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooStmt) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooStmt) CanClose() bool {
	return false
}

func (gooStmt) Close(o interface{}) {
}
