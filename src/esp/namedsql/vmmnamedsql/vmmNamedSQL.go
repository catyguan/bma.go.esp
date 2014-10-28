package vmmnamedsql

import (
	"bmautil/valutil"
	"esp/namedsql"
	"fmt"
	"golua"
	"golua/vmmsql"
)

const tag = "vmmsql"

func Module(s *namedsql.Service) *golua.VMModule {
	m := vmmsql.Module()
	m.Init("open", &GOF_sql_named{s})
	return m
}

// sql.named(name) DB:object
type GOF_sql_named struct {
	s *namedsql.Service
}

func (this GOF_sql_named) Exec(vm *golua.VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	n, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn := valutil.ToString(n, "")
	if vn == "" {
		return 0, fmt.Errorf("namedSQL name invalid(%v)", n)
	}
	db, err2 := this.s.Get(vn)
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(vmmsql.NewDBObject(db))
	return 1, nil
}

func (this GOF_sql_named) IsNative() bool {
	return true
}

func (this GOF_sql_named) String() string {
	return "GoFunc<sql.named>"
}
