package vmmsql

import (
	"bmautil/valutil"
	"database/sql"
	"fmt"
	"golua"
)

const tag = "vmmsql"

func Module() *golua.VMModule {
	m := golua.NewVMModule("sql")
	m.Init("open", GOF_sql_open(0))
	return m
}

func NewDBObject(db *sql.DB) golua.VMTable {
	return golua.NewGOO(db, gooDB(0))
}

// sql.open(driver, dataSource) DB:object
type GOF_sql_open int

func (this GOF_sql_open) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	dr, ds, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vdr := valutil.ToString(dr, "")
	if vdr == "" {
		return 0, fmt.Errorf("driver invalid(%v)", dr)
	}
	vds := valutil.ToString(ds, "")
	if vds == "" {
		return 0, fmt.Errorf("dataSource invalid(%v)", ds)
	}
	db, err2 := sql.Open(vdr, vds)
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(NewDBObject(db))
	return 1, nil
}

func (this GOF_sql_open) IsNative() bool {
	return true
}

func (this GOF_sql_open) String() string {
	return "GoFunc<sql.open>"
}
