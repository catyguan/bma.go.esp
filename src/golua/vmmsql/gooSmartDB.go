package vmmsql

import (
	"bmautil/valutil"
	"database/sql"
	"fmt"
	"golua"
	"logger"
)

func SmartDBFactory(n string) (interface{}, error) {
	return golua.NewGOO(0, gooSmartDB(0)), nil
}

func API_select(vm *golua.VM, tableName string, write bool) (*sql.DB, error) {
	gos, errC := vm.GetGoLua().SingletonService("SmartDB", createSmartDB)
	if errC != nil {
		return nil, errC
	}
	obj, ok := gos.(*smartDB)
	if !ok {
		return nil, fmt.Errorf("invalid SmartDB")
	}
	dbi := obj.Select(tableName, write)
	if dbi == nil {
		return nil, nil
	}
	logger.Debug(tag, "API_select '%s' => %s", tableName, dbi)
	vm.API_push(dbi.Driver)
	vm.API_push(dbi.DataSource)
	vm.API_push(dbi.MaxOpenConns)
	vm.API_push(dbi.MaxIdleConns)
	c, err2 := GOF_sql_open(0).Exec(vm, nil)
	if err2 != nil {
		return nil, err2
	}
	dbobj, err3 := vm.API_pop1X(c, true)
	if err3 != nil {
		return nil, err3
	}
	gos2 := vm.API_object(dbobj)
	return API_toDB(gos2), nil
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
			// SmartDB.Add(dbi:table, syncRefresh:bool)
			return golua.NewGOF("SmartDB.Add", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				o, sr, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				mo := vm.API_toMap(o)
				vsr := valutil.ToBool(sr, false)

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
				err2 := obj.Add(dbi, vsr)
				if err2 != nil {
					return 0, err2
				}
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
			// case "X":
			// 	return golua.NewGOF("SmartDB.X", func(vm *golua.VM, self interface{}) (int, error) {
			// 		db, err := API_select(vm, "tEST4", false)
			// 		fmt.Println("API_select", db, err)
			// 		return 0, nil
			// 	}), nil
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
