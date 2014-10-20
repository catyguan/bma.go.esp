package vmmacclog

import (
	"bmautil/valutil"
	"esp/acclog"
	"golua"
)

func Module() *golua.VMModule {
	m := golua.NewVMModule("acclog")
	m.Init("log", GOF_acclog_log(0))
	return m
}

// acclog.log(n:string, val[, notOverwrite:bool])
type GOF_acclog_log int

func (this GOF_acclog_log) Exec(vm *golua.VM) (int, error) {
	ctx := vm.API_getContext()
	if ctx != nil {
		adt, ok := acclog.AcclogDataFromContext(ctx)
		if ok {
			err0 := vm.API_checkstack(2)
			if err0 != nil {
				return 0, err0
			}
			n, v, not, err1 := vm.API_pop3X(-1, true)
			if err1 != nil {
				return 0, err1
			}
			vn := valutil.ToString(n, "")
			vnot := valutil.ToBool(not, false)

			if vnot {
				if _, ok2 := adt[vn]; ok2 {
					return 0, nil
				}
			}
			adt[vn] = v
		}
	}
	return 0, nil
}

func (this GOF_acclog_log) IsNative() bool {
	return true
}

func (this GOF_acclog_log) String() string {
	return "GoFunc<acclog.log>"
}
