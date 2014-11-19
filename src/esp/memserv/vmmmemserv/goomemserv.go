package vmmmemserv

import (
	"bmautil/valutil"
	"esp/memserv"
	"fmt"
	"golua"
)

const (
	tag = "vmmMemServ"
)

type MemServFactory struct {
	s *memserv.MemoryServ
}

func (this *MemServFactory) FactoryFunc(vm *golua.VM, n string) (interface{}, error) {
	gl := vm.GetGoLua()
	o, _ := gl.SingletonService("MemServ", func() (interface{}, error) {
		o := new(ObjectMemServ)
		o.s = this.s
		o.gl = gl
		o.appkeys = make(map[string]bool)
		return o, nil
	})
	return golua.NewGOO(o, gooMemServ(0)), nil
}

func API_toMemServ(vm *golua.VM, o interface{}) *ObjectMemServ {
	v := vm.API_object(o)
	if v == nil {
		return nil
	}
	r, ok := v.(*ObjectMemServ)
	if ok {
		return r
	}
	return nil
}

type gooMemServ int

func (gooMemServ) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*ObjectMemServ); ok {
		switch key {
		case "Close":
			return golua.NewGOF("MemServ.Close", func(vm *golua.VM, self interface{}) (int, error) {
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
				r, err2 := obj.Close(vm, vn)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(r)
				return 1, nil
			}), nil
		case "Create":
			return golua.NewGOF("MemServ.Create", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				n, cfg, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				vcfg := vm.API_toMap(cfg)
				if vcfg == nil {
					return 0, fmt.Errorf("Config invalid")
				}
				var co *memserv.MemGoConfig
				co = new(memserv.MemGoConfig)
				valutil.ToBean(vcfg, co)
				err1x := co.Valid()
				if err1x != nil {
					return 0, err1x
				}
				r, err2 := obj.Create(vm, vn, co)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(NewGOOMemGo(r))
				return 1, nil
			}), nil
		case "Get":
			return golua.NewGOF("MemServ.Get", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, cfg, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				vcfg := vm.API_toMap(cfg)
				var co *memserv.MemGoConfig
				if vcfg != nil {
					co = new(memserv.MemGoConfig)
					valutil.ToBean(vcfg, co)
					err1x := co.Valid()
					if err1x != nil {
						return 0, err1x
					}
				}
				r, err2 := obj.Get(vm, vn, co)
				if err2 != nil {
					return 0, err2
				}
				if r != nil {
					vm.API_push(NewGOOMemGo(r))
					return 1, nil
				}
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooMemServ) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooMemServ) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooMemServ) CanClose() bool {
	return false
}

func (gooMemServ) Close(o interface{}) {
}
