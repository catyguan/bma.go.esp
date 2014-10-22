package vmmsql

import (
	"database/sql"
	"golua"
)

type gooStmt int

func (gooStmt) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*sql.Stmt); ok {
		switch key {
		case "Close":
			return golua.NewGOF("Stmt:Close", func(vm *golua.VM) (int, error) {
				vm.API_popAll()
				obj.Close()
				return 0, nil
			}), nil
		case "Exec":
			return golua.NewGOF("Stmt:Exec", func(vm *golua.VM) (int, error) {
				top := vm.API_gettop()
				ns, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				rs, err2 := obj.Exec(ns[1:]...)
				if err2 != nil {
					return 0, err2
				}
				ra, err3 := rs.RowsAffected()
				if err3 != nil {
					return 0, err3
				}
				vm.API_push(ra)
				return 1, nil
			}), nil
		case "ExecLastId":
			return golua.NewGOF("Stmt:ExecLastId", func(vm *golua.VM) (int, error) {
				top := vm.API_gettop()
				ns, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				rs, err2 := obj.Exec(ns[1:]...)
				if err2 != nil {
					return 0, err2
				}
				ra, err3 := rs.RowsAffected()
				if err3 != nil {
					return 0, err3
				}
				vm.API_push(ra)
				if ra > 0 {
					rid, err4 := rs.LastInsertId()
					if err4 != nil {
						return 1, err4
					}
					vm.API_push(rid)
				} else {
					vm.API_push(0)
				}
				return 2, nil
			}), nil
		case "Query":
			return golua.NewGOF("Stmt:Query", func(vm *golua.VM) (int, error) {
				top := vm.API_gettop()
				ns, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				rs, err2 := obj.Query(ns[1:]...)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(golua.NewGOO(rs, gooRows(0)))
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
	return true
}

func (gooStmt) Close(o interface{}) {
	if obj, ok := o.(*sql.Stmt); ok {
		obj.Close()
	}
}
