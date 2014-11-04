package vmmsql

import (
	"bmautil/valutil"
	"bytes"
	"database/sql"
	"fmt"
	"golua"
)

func NewDBObject(vm *golua.VM, db *sql.DB) golua.VMTable {
	gos := vm.GetGoLua().CreateGoService("db", db, func() {
		db.Close()
	})
	return golua.NewGOO(gos, gooDB(0))
}

type gooDB int

func (gooDB) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if gos, ok := o.(*golua.GoService); ok {
		obj := gos.Data.(*sql.DB)
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
				gos.Close()
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
		case "ExecInsert":
			// ExecInsert(tableName:string, fv:table[, lastId:bool]) int
			return golua.NewGOF("DB.ExecInsert", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				tn, fv, lid, err1 := vm.API_pop3X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vtn := valutil.ToString(tn, "")
				if vtn == "" {
					return 0, fmt.Errorf("insert table name invalid(%v)", tn)
				}
				vfv := vm.API_toMap(fv)
				if len(vfv) == 0 {
					return 0, fmt.Errorf("insert fields invalid")
				}

				buf := bytes.NewBuffer(make([]byte, 0, 128))
				buf2 := bytes.NewBuffer(make([]byte, 0, 128))
				buf.WriteString("INSERT INTO ")
				buf.WriteString(vtn)
				buf.WriteString(" (")
				buf2.WriteString(" VALUES (")
				ps := make([]interface{}, 0)
				first := true
				for k, v := range vfv {
					if !first {
						buf.WriteString(",")
						buf2.WriteString(",")
					} else {
						first = false
					}
					buf.WriteString(k)
					buf2.WriteString("?")
					ps = append(ps, v)
				}
				buf.WriteString(")")
				buf2.WriteString(")")
				sql := buf.String() + buf2.String()

				rs, err2 := obj.Exec(sql, ps...)
				if err2 != nil {
					return 0, err2
				}
				ra, err3 := rs.RowsAffected()
				if err3 != nil {
					return 0, err3
				}
				vm.API_push(ra)
				if ra > 0 && valutil.ToBool(lid, false) {
					rid, err4 := rs.LastInsertId()
					if err4 != nil {
						return 1, err4
					}
					vm.API_push(rid)
				}
				return 2, nil
			}), nil
		case "ExecUpdate":
			// ExecUpdate(tableName:string, fv:table, tj:table) int
			return golua.NewGOF("DB.ExecUpdate", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(3)
				if err0 != nil {
					return 0, err0
				}
				tn, fv, tj, err1 := vm.API_pop3X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vtn := valutil.ToString(tn, "")
				if vtn == "" {
					return 0, fmt.Errorf("update table name invalid(%v)", tn)
				}
				vfv := vm.API_toMap(fv)
				if len(vfv) == 0 {
					return 0, fmt.Errorf("update fields invalid")
				}
				vtj := vm.API_toMap(tj)
				if len(vtj) == 0 {
					return 0, fmt.Errorf("update condition invalid")
				}

				buf := bytes.NewBuffer(make([]byte, 0, 128))
				buf.WriteString("UPDATE ")
				buf.WriteString(vtn)
				buf.WriteString(" SET ")
				ps := make([]interface{}, 0)
				first := true
				for k, v := range vfv {
					if !first {
						buf.WriteString(",")
					} else {
						first = false
					}
					buf.WriteString(k)
					buf.WriteString("=?")
					ps = append(ps, v)
				}
				buf.WriteString(" WHERE ")
				first = true
				for k, v := range vtj {
					if !first {
						buf.WriteString(" AND ")
					} else {
						first = false
					}
					buf.WriteString(k)
					buf.WriteString("=?")
					ps = append(ps, v)
				}
				rs, err2 := obj.Exec(buf.String(), ps...)
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
		case "ExecDelete":
			// ExecDelete(tableName:string, tj:table) int
			return golua.NewGOF("DB.ExecDelete", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				tn, tj, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vtn := valutil.ToString(tn, "")
				if vtn == "" {
					return 0, fmt.Errorf("delete table name invalid(%v)", tn)
				}
				vtj := vm.API_toMap(tj)
				if len(vtj) == 0 {
					return 0, fmt.Errorf("delete condition invalid")
				}

				buf := bytes.NewBuffer(make([]byte, 0, 128))
				buf.WriteString("DELETE FROM ")
				buf.WriteString(vtn)
				buf.WriteString(" WHERE ")
				ps := make([]interface{}, 0)
				first := true
				for k, v := range vtj {
					if !first {
						buf.WriteString(" AND ")
					} else {
						first = false
					}
					buf.WriteString(k)
					buf.WriteString("=?")
					ps = append(ps, v)
				}
				rs, err2 := obj.Exec(buf.String(), ps...)
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
