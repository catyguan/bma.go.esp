package vmmmemserv

import (
	"bmautil/valutil"
	"esp/memserv"
	"fmt"
	"golua"
	"time"
)

type gooMemGoI int

func (gooMemGoI) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*memserv.MemGoI); ok {
		switch key {
		case "Incr":
			return golua.NewGOF("MemGoI.Incr", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				n, v, tm, err1 := vm.API_pop3X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				vv := valutil.ToInt64(v, 0)
				vtm := valutil.ToInt(tm, 0)
				r, err2 := obj.Incr(vn, vv, vtm)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(r)
				return 1, nil
			}), nil
		case "Set":
			return golua.NewGOF("MemGoI.Set", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				n, v, tm, err1 := vm.API_pop3X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				vv := golua.BaseData(v)
				vtm := valutil.ToInt(tm, 0)
				err2 := obj.Set(vn, vv, vtm)
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "Put":
			return golua.NewGOF("MemGoI.Put", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				n, v, tm, err1 := vm.API_pop3X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				vv := golua.BaseData(v)
				vtm := valutil.ToInt(tm, 0)
				rb, err2 := obj.Put(vn, vv, vtm)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(rb)
				return 1, nil
			}), nil
		case "Touch":
			return golua.NewGOF("MemGoI.Touch", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, tm, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				vtm := valutil.ToInt(tm, 0)
				err2 := obj.Touch(vn, vtm)
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "Get":
			return golua.NewGOF("MemGoI.Get", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, tm, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				var ptm *time.Time
				if tm != nil {
					if o := vm.API_object(tm); o != nil {
						if tmp, ok := o.(*time.Time); ok {
							ptm = tmp
						}
					}
					if ptm == nil {
						if _, ok := tm.(bool); ok {
							tmp := time.Now()
							ptm = &tmp
						}
					}
				}
				_, r, err2 := obj.Get(vn, ptm)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(golua.ScriptData(r))
				return 1, nil
			}), nil
		case "Remove":
			return golua.NewGOF("MemGoI.Remove", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				err2 := obj.Remove(vn)
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooMemGoI) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooMemGoI) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooMemGoI) CanClose() bool {
	return false
}

func (gooMemGoI) Close(o interface{}) {
}
