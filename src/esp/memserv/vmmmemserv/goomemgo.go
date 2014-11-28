package vmmmemserv

import (
	"bmautil/valutil"
	"esp/memserv"
	"fmt"
	"golua"
	"logger"
	"strings"
	"time"
)

func NewGOOMemGo(mg *memserv.MemGo) golua.VMTable {
	return golua.NewGOO(mg, gooMemGo(0))
}

type gooMemGo int

func (gooMemGo) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*memserv.MemGo); ok {
		switch key {
		case "Call":
			return golua.NewGOF("MemGo.Call", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				top := vm.API_gettop()
				ps, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				f := ps[0]
				if !vm.API_canCall(f) {
					return 0, fmt.Errorf("invalid Call func")
				}
				vm2, err3 := vm.GetGoLua().GetVM()
				if err3 != nil {
					return 0, err3
				}
				defer vm2.Finish()
				var rs []interface{}
				err2 := obj.DoSync(func(mgi *memserv.MemGoI) error {
					goomgi := golua.NewGOO(mgi, gooMemGoI(0))
					vm2.API_push(f)
					vm2.API_push(goomgi)
					for _, v := range ps[1:] {
						vm2.API_push(v)
					}
					rc, err := vm2.Call(top, -1, nil)
					if err != nil {
						return err
					}
					rs, err = vm2.API_popN(rc, true)
					return err
				})
				if err2 != nil {
					return 0, err2
				}
				for _, v := range rs {
					vm.API_push(v)
				}
				return len(rs), nil
			}), nil
		case "Post":
			return golua.NewGOF("MemGo.Post", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				top := vm.API_gettop()
				ps, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				f := ps[0]
				if !vm.API_canCall(f) {
					return 0, fmt.Errorf("invalid Post func")
				}
				vm2, err3 := vm.GetGoLua().GetVM()
				if err3 != nil {
					return 0, err3
				}
				err2 := obj.DoNow(func(mgi *memserv.MemGoI) error {
					defer vm2.Finish()
					goomgi := golua.NewGOO(mgi, gooMemGoI(0))
					vm2.API_push(f)
					vm2.API_push(goomgi)
					for _, v := range ps[1:] {
						vm2.API_push(v)
					}
					_, err := vm2.Call(top, 0, nil)
					if err != nil {
						logger.Warn(tag, "MemGo.Post fail - %s", err)
						return nil
					}
					vm2.API_popAll()
					return nil
				})
				if err2 != nil {
					vm2.Finish()
					return 0, err2
				}
				return 0, nil
			}), nil
		case "Scan":
			// reutrn isEnd, array()
			return golua.NewGOF("MemGo.Scan", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				a, n, count, err1 := vm.API_pop3X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				va := valutil.ToString(a, "")
				if va == "" {
					return 0, fmt.Errorf("Action invalid")
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("ScanName invalid")
				}
				switch strings.ToLower(va) {
				case "begin":
					err2 := obj.BeginScan(vn)
					if err2 != nil {
						return 0, err2
					}
					return 0, nil
				case "end":
					err2 := obj.EndScan(vn)
					if err2 != nil {
						return 0, err2
					}
					return 0, nil
				case "do":
					ra := vm.API_newarray()
					vcount := valutil.ToInt(count, 10)
					isEnd, err2 := obj.Scan(vn, vcount, func(k string, v interface{}) {
						m := make(map[string]interface{})
						m["Key"] = k
						m["Value"] = golua.ScriptData(v)
						ra.Add(nil, m)
					})
					if err2 != nil {
						return 0, err2
					}
					vm.API_push(isEnd)
					vm.API_push(ra)
					return 2, nil
				default:
					return 0, fmt.Errorf("unknow scan action(%s)", va)
				}
			}), nil
		case "Size":
			return golua.NewGOF("MemServ.Size", func(vm *golua.VM, self interface{}) (int, error) {
				vm.API_popAll()
				c, sz := obj.Size()
				vm.API_push(c)
				vm.API_push(sz)
				return 2, nil
			}), nil
		case "Incr":
			return golua.NewGOF("MemGo.Incr", func(vm *golua.VM, self interface{}) (int, error) {
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
				var vr interface{}
				err2 := obj.DoSync(func(mgi *memserv.MemGoI) error {
					r, err2 := mgi.Incr(vn, vv, vtm)
					if err2 != nil {
						return err2
					}
					vr = r
					return nil
				})
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(vr)
				return 1, nil
			}), nil
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
				err2 := obj.DoSync(func(mgi *memserv.MemGoI) error {
					return mgi.Set(vn, vv, vtm)
				})
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "Put":
			return golua.NewGOF("MemGo.Put", func(vm *golua.VM, self interface{}) (int, error) {
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
				var vr interface{}
				err2 := obj.DoSync(func(mgi *memserv.MemGoI) error {
					rb, err2 := mgi.Put(vn, vv, vtm)
					vr = rb
					return err2
				})
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(vr)
				return 1, nil
			}), nil
		case "Touch":
			return golua.NewGOF("MemGo.Touch", func(vm *golua.VM, self interface{}) (int, error) {
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
				err2 := obj.DoSync(func(mgi *memserv.MemGoI) error {
					return mgi.Touch(vn, vtm)
				})
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
				var vr interface{}
				err2 := obj.DoSync(func(mgi *memserv.MemGoI) error {
					_, r, err2 := mgi.Get(vn, ptm)
					vr = r
					return err2
				})
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(golua.ScriptData(vr))
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
				err2 := obj.DoSync(func(mgi *memserv.MemGoI) error {
					return mgi.Remove(vn)
				})
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
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
