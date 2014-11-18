package vmmmemserv

import (
	"bmautil/valutil"
	"esp/memserv"
	"fmt"
	"golua"
	"time"
)

func NewGOOMemGo(mg *memserv.MemGo) golua.VMTable {
	return golua.NewGOO(mg, gooMemGo(0))
}

type gooMemGo int

func (gooMemGo) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*memserv.MemGo); ok {
		switch key {
		case "Set":
			return golua.NewGOF("MemGo.Set", func(vm *golua.VM, self interface{}) (int, error) {
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
		case "Get":
			return golua.NewGOF("MemGo.Get", func(vm *golua.VM, self interface{}) (int, error) {
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
				r, err2 := obj.Get(vn, ptm)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(r)
				return 1, nil
			}), nil
		case "Remove":
			return golua.NewGOF("MemGo.Remove", func(vm *golua.VM, self interface{}) (int, error) {
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
		case "Size":
			return golua.NewGOF("MemServ.Size", func(vm *golua.VM, self interface{}) (int, error) {
				vm.API_popAll()
				c, sz := obj.Size()
				vm.API_push(c)
				vm.API_push(sz)
				return 2, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooMemGo) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooMemGo) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooMemGo) CanClose() bool {
	return false
}

func (gooMemGo) Close(o interface{}) {
}
