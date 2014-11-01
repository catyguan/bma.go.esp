package vmmsql

import (
	"bmautil/valutil"
	"database/sql"
	"fmt"
	"golua"
)

type gooRows int

func (gooRows) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*sql.Rows); ok {
		switch key {
		case "Fetch":
			return golua.NewGOF("Rows.Fetch", func(vm *golua.VM, self interface{}) (int, error) {
				errX := vm.API_checkStack(1)
				if errX != nil {
					return 0, errX
				}
				va, desc, err1 := vm.API_pop2X(-1, false)
				if err1 != nil {
					return 0, err1
				}
				if va == nil {
					return 0, fmt.Errorf("fetch var nil")
				}
				vva, ok := va.(golua.VMVar)
				if !ok {
					return 0, fmt.Errorf("fetch var invalid(%T)", va)
				}
				desc, err1 = vm.API_value(desc)
				if err1 != nil {
					return 0, err1
				}
				mdesc := make(map[string]string)
				if desc != nil {
					m := vm.API_toMap(desc)
					if m != nil {
						for k, v := range m {
							mdesc[k] = valutil.ToString(v, "")
						}
					}
				}
				if obj.Next() {
					r, err2 := FetchRow(obj, mdesc)
					if err2 != nil {
						return 0, err2
					}
					vva.Set(vm, r)
					vm.API_push(true)
				} else {
					vm.API_push(false)
				}
				return 1, nil
			}), nil
		case "Close":
			return golua.NewGOF("Rows.Close", func(vm *golua.VM, self interface{}) (int, error) {
				vm.API_popAll()
				obj.Close()
				return 0, nil
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
