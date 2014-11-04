package vmmsql

import (
	"bmautil/valutil"
	"database/sql"
	"fmt"
	"golua"
	"logger"
)

const tag = "vmmsql"

func Module() *golua.VMModule {
	m := golua.NewVMModule("sql")
	m.Init("create", GOF_sql_create(0))
	m.Init("open", GOF_sql_open(0))
	return m
}

func InitGoLua(gl *golua.GoLua) {
	Module().Bind(gl)
	gl.SetObjectFactory("SmartDB", SmartDBFactory)
}

// sql.create(driver, dataSource[, maxConn, maxIdle]) DB:object
type GOF_sql_create int

func (this GOF_sql_create) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	dr, ds, mc, mi, err1 := vm.API_pop4X(-1, true)
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
	vmc := valutil.ToInt(mc, -1)
	vmi := valutil.ToInt(mi, -1)
	db, err2 := sql.Open(vdr, vds)
	if err2 != nil {
		return 0, err2
	}
	if vmc >= 0 {
		db.SetMaxOpenConns(vmc)
	}
	if vmi >= 0 {
		db.SetMaxIdleConns(vmi)
	}
	vm.API_push(NewDBObject(vm, db))
	return 1, nil
}

func (this GOF_sql_create) IsNative() bool {
	return true
}

func (this GOF_sql_create) String() string {
	return "GoFunc<sql.create>"
}

// sql.open(driver, dataSource[, maxConn, maxIdle]) DB:object
type GOF_sql_open int

func (this GOF_sql_open) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	dr, ds, mc, mi, err1 := vm.API_pop4X(-1, true)
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
	vmc := valutil.ToInt(mc, -1)
	vmi := valutil.ToInt(mi, -1)
	gl := vm.GetGoLua()
	sn := fmt.Sprintf("db_%s_!!%s", vdr, vds)
	gos, err2 := gl.SingletonService(sn, func() (interface{}, error) {
		logger.Debug(tag, "create new db(%s)", vdr)
		db, err := sql.Open(vdr, vds)
		if err != nil {
			return nil, err
		}
		if vmc >= 0 {
			db.SetMaxOpenConns(vmc)
		}
		if vmi >= 0 {
			db.SetMaxIdleConns(vmi)
		}
		gos := new(golua.GoService)
		gos.GL = gl
		gos.SID = sn
		gos.Data = db
		gos.CloseFunc = func() {
			db.Close()
		}
		return gos, nil
	})
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(golua.NewGOO(gos, gooDB(0)))
	return 1, nil
}

func (this GOF_sql_open) IsNative() bool {
	return true
}

func (this GOF_sql_open) String() string {
	return "GoFunc<sql.open>"
}
