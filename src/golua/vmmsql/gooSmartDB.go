package vmmsql

import (
	"bmautil/valutil"
	"fmt"
	"golua"
	"logger"
)

func SmartDBFactory(n string) (interface{}, error) {
	return golua.NewGOO(0, gooSmartDB(0)), nil
}

type gooSmartDB int

func (gooSmartDB) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	gos, errC := vm.GetGoLua().SingletonService("SmartDB", createSmartDB)
	if errC != nil {
		return nil, errC
	}
	if obj, ok := gos.(*smartDB); ok {
		switch key {
		case "Add":
			return golua.NewGOF("SmartDB.Add", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				o, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				mo := vm.API_toMap(o)
				dbi := new(dbInfo)
				if !valutil.ToBean(mo, dbi) {
					return 0, fmt.Errorf("convert dbInfo fail - %v", mo)
				}
				if dbi.Name == "" {
					return 0, fmt.Errorf("dbInfo.Name invalid")
				}
				if dbi.Driver == "" {
					return 0, fmt.Errorf("dbInfo.Driver invalid")
				}
				if dbi.DataSource == "" {
					return 0, fmt.Errorf("dbInfo.DataSource invalid")
				}
				obj.Add(dbi)
				return 0, nil
			}), nil
		case "Select":
			return golua.NewGOF("SmartDB.Select", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				tbn, write, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vtbn := valutil.ToString(tbn, "")
				if vtbn == "" {
					return 0, fmt.Errorf("TableName invalid")
				}
				vwrite := valutil.ToBool(write, false)
				dbi := obj.Select(vtbn, vwrite)
				if dbi == nil {
					return 0, fmt.Errorf("Select(%s) fail", vtbn)
				}
				logger.Debug(tag, "select '%s' => %s", vtbn, dbi)
				vm.API_push(dbi.Driver)
				vm.API_push(dbi.DataSource)
				vm.API_push(dbi.MaxOpenConns)
				vm.API_push(dbi.MaxIdleConns)
				return GOF_sql_open(0).Exec(vm, self)
			}), nil
		case "Refresh":
			return golua.NewGOF("SmartDB.Refresh", func(vm *golua.VM, self interface{}) (int, error) {
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooSmartDB) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooSmartDB) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooSmartDB) CanClose() bool {
	return false
}

func (gooSmartDB) Close(o interface{}) {
}
