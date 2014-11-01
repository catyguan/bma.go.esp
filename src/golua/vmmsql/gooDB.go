package vmmsql

import (
	"bmautil/valutil"
	"database/sql"
	"fmt"
	"golua"
)

type gooDB int

func (gooDB) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*sql.DB); ok {
		switch key {
		case "Begin":
			return golua.NewGOF("DB.Begin", func(vm *golua.VM, self interface{}) (int, error) {
				vm.API_popAll()
				tx, err := obj.Begin()
				if err != nil {
					return 0, err
				}
				ro := golua.NewGOO(tx, gooTx(0))
				errl := vm.API_cleanDefer(ro)
				if errl != nil {
					return 0, errl
				}
				vm.API_push(ro)
				return 1, nil
			}), nil
		case "Close":
			return golua.NewGOF("DB.Close", func(vm *golua.VM, self interface{}) (int, error) {
				vm.API_popAll()
				obj.Close()
				return 0, nil
			}), nil
		case "Exec":
			return golua.NewGOF("DB.Exec", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				top := vm.API_gettop()
				ns, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				vsql := valutil.ToString(ns[0], "")
				if vsql == "" {
					return 0, fmt.Errorf("query string invalid(%v)", ns[0])
				}
				rs, err2 := obj.Exec(vsql, ns[1:]...)
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
			return golua.NewGOF("DB.ExecLastId", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				top := vm.API_gettop()
				ns, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				vsql := valutil.ToString(ns[0], "")
				if vsql == "" {
					return 0, fmt.Errorf("query string invalid(%v)", ns[0])
				}
				rs, err2 := obj.Exec(vsql, ns[1:]...)
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
		case "Ping":
			return golua.NewGOF("DB.Ping", func(vm *golua.VM, self interface{}) (int, error) {
				vm.API_popAll()
				err0 := obj.Ping()
				if err0 != nil {
					return 0, err0
				}
				return 0, nil
			}), nil
		case "Prepare":
			return golua.NewGOF("DB.Prepare", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				q, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vq := valutil.ToString(q, "")
				if vq == "" {
					return 0, fmt.Errorf("query string invalid(%v)", q)
				}
				st, err2 := obj.Prepare(vq)
				if err2 != nil {
					return 0, err2
				}
				ro := golua.NewGOO(st, gooStmt(0))
				errl := vm.API_cleanDefer(ro)
				if errl != nil {
					return 0, errl
				}
				vm.API_push(ro)
				return 1, nil
			}), nil
		case "Query":
			return golua.NewGOF("DB.Query", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				top := vm.API_gettop()
				ns, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				vsql := valutil.ToString(ns[0], "")
				if vsql == "" {
					return 0, fmt.Errorf("query string invalid(%v)", ns[0])
				}
				rs, err2 := obj.Query(vsql, ns[1:]...)
				if err2 != nil {
					return 0, err2
				}
				ro := golua.NewGOO(rs, gooRows(0))
				errl := vm.API_cleanDefer(ro)
				if errl != nil {
					return 0, errl
				}
				vm.API_push(ro)
				return 1, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooDB) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooDB) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooDB) CanClose() bool {
	return false
}

func (gooDB) Close(o interface{}) {
}
