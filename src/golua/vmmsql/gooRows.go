package vmmsql

import (
	"bmautil/sqlutil"
	"database/sql"
	"golua"
)

type gooRows int

func (gooRows) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*sql.Rows); ok {
		switch key {
		case "Fetch":
			return golua.NewGOF("Rows:Fetch", func(vm *golua.VM) (int, error) {
				vm.API_popAll()
				r, err0 := sqlutil.FetchRow(obj)
				if err0 != nil {
					return 0, err0
				}
				vm.API_push(r != nil)
				vm.API_push(r)
				return 2, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooRows) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooRows) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooRows) CanClose() bool {
	return true
}

func (gooRows) Close(o interface{}) {
	if obj, ok := o.(*sql.Rows); ok {
		obj.Close()
	}
}
