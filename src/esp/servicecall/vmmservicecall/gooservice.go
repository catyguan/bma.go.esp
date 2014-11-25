package vmmservicecall

import (
	"bmautil/valutil"
	"esp/servicecall"
	"fmt"
	"golua"
	"time"
)

func ServiceCallGoLuaFactoryFunc(s *servicecall.Service) golua.GoObjectFactory {
	return func(vm *golua.VM, n string) (interface{}, error) {
		ns, errC := vm.GetGoLua().SingletonService("ServiceCall", func() (interface{}, error) {
			o := servicecall.NewService(n, s)
			return o, nil
		})
		if errC != nil {
			return nil, errC
		}
		return golua.NewGOO(ns, gooService(0)), nil
	}
}

func API_toService(vm *golua.VM, v interface{}) *servicecall.Service {
	if v == nil {
		return nil
	}
	if s, ok := v.(*servicecall.Service); ok {
		return s
	}
	return nil
}

type goluaServiceCaller struct {
	gl *golua.GoLua
	f  interface{}
}

func (this *goluaServiceCaller) Start() error {
	return nil
}

func (this *goluaServiceCaller) Ping() bool {
	if this.gl.IsClose() {
		return false
	}
	return true
}

func (this *goluaServiceCaller) Stop() {
}

func (this *goluaServiceCaller) Call(method string, params []interface{}, timeout time.Duration) (interface{}, error) {
	vm, err0 := this.gl.GetVM()
	if err0 != nil {
		return nil, err0
	}
	defer vm.Finish()
	vm.API_push(this.f)
	vm.API_push(method)
	vm.API_push(params)
	vm.API_push(int(timeout.Seconds() * 1000))
	rc, err2 := vm.Call(3, 1, nil)
	if err2 != nil {
		return nil, err2
	}
	rv, err3 := vm.API_pop1X(rc, true)
	if err3 != nil {
		return nil, err3
	}
	return rv, nil
}

type gooService int

func (gooService) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*servicecall.Service); ok {
		switch key {
		case "Bind":
			return golua.NewGOF("Service.Bind", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				n, f, ow, err1 := vm.API_pop3X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				if !vm.API_canCall(f) {
					return 0, fmt.Errorf("Func(%T) invalid", f)
				}
				vow := valutil.ToBool(ow, true)
				sc := new(goluaServiceCaller)
				sc.gl = vm.GetGoLua()
				sc.f = f
				rv := obj.SetServiceCall(vn, sc, vow)
				vm.API_push(rv)
				return 1, nil
			}), nil
		case "Assert":
			return golua.NewGOF("Service.Assert", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, tms, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				vtms := valutil.ToInt(tms, 0)
				if vtms <= 0 {
					return 0, fmt.Errorf("TimeoutMS invalid")
				}
				sc, err2 := obj.Assert(vn, time.Duration(vtms)*time.Millisecond)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(NewGOOServiceCall(sc))
				return 1, nil
			}), nil
		case "Get":
			return golua.NewGOF("Service.Get", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, tms, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				vtms := valutil.ToInt(tms, 0)
				if vtms <= 0 {
					return 0, fmt.Errorf("TimeoutMS invalid")
				}
				sc, err2 := obj.Get(vn, time.Duration(vtms)*time.Millisecond)
				if err2 != nil {
					return 0, err2
				}
				if sc == nil {
					vm.API_push(nil)
				} else {
					vm.API_push(NewGOOServiceCall(sc))
				}
				return 1, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooService) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooService) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooService) CanClose() bool {
	return false
}

func (gooService) Close(o interface{}) {
}
