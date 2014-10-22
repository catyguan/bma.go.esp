package vmmsql

import (
	"bmautil/sqlutil"
	"bmautil/valutil"
	"database/sql"
	"fmt"
	"golua"
	"time"
)

type gooRows int

func (gooRows) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*sql.Rows); ok {
		switch key {
		case "Fetch":
			return golua.NewGOF("Rows:Fetch", func(vm *golua.VM) (int, error) {
				errX := vm.API_checkstack(2)
				if errX != nil {
					return 0, errX
				}
				_, va, desc, err1 := vm.API_pop3X(-1, false)
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
					switch rdesc := desc.(type) {
					case map[string]interface{}:
						for k, v := range rdesc {
							mdesc[k] = valutil.ToString(v, "")
						}
					case golua.VMTable:
						m := rdesc.ToMap()
						for k, v := range m {
							mdesc[k] = valutil.ToString(v, "")
						}
					}
				}
				// fmt.Println("desc", desc)
				r, err2 := sqlutil.FetchRow(obj, mdesc)
				if err2 != nil {
					return 0, err2
				}
				for k, v := range r {
					if v != nil {
						if tm, ok := v.(time.Time); ok {
							v = golua.CreateGoTime(&tm)
							r[k] = v
						}
					}
				}
				vva.Set(vm, r)
				vm.API_push(r != nil)
				return 1, nil
			}), nil
		case "Close":
			return golua.NewGOF("Rows:Close", func(vm *golua.VM) (int, error) {
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
