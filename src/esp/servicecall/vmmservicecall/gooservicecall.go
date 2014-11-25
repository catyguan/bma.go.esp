package vmmservicecall

import (
	"bmautil/valutil"
	"esp/servicecall"
	"fmt"
	"golua"
	"logger"
	"time"
)

func NewGOOServiceCall(sc servicecall.ServiceCaller) golua.VMTable {
	return golua.NewGOO(sc, gooServiceCall(0))
}

type gooServiceCall int

func (gooServiceCall) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(servicecall.ServiceCaller); ok {
		switch key {
		case "Call":
			return golua.NewGOF("ServiceCall.Call", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				top := vm.API_gettop()
				vlist, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(vlist[0], "")
				if vn == "" {
					return 0, fmt.Errorf("Method invalid")
				}
				for i, v := range vlist {
					if i < 1 {
						continue
					}
					vlist[i] = golua.BaseData(v)
				}
				rv, err2 := obj.Call(vn, vlist[1:], time.Duration(0))
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(rv)
				return 1, nil
			}), nil
		case "CallTimeout":
			return golua.NewGOF("ServiceCall.CallTimeout", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				top := vm.API_gettop()
				vlist, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(vlist[0], "")
				if vn == "" {
					return 0, fmt.Errorf("Method invalid")
				}
				vtms := valutil.ToInt(vlist[1], 0)
				if vtms <= 0 {
					return 0, fmt.Errorf("TimeoutMS invalid")
				}
				for i, v := range vlist {
					if i < 2 {
						continue
					}
					vlist[i] = golua.BaseData(v)
				}
				rv, err2 := obj.Call(vn, vlist[2:], time.Duration(vtms)*time.Millisecond)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(rv)
				return 1, nil
			}), nil
		case "Post":
			return golua.NewGOF("ServiceCall.Post", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				top := vm.API_gettop()
				vlist, err1 := vm.API_popN(top, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(vlist[0], "")
				if vn == "" {
					return 0, fmt.Errorf("Method invalid")
				}
				for i, v := range vlist {
					if i < 1 {
						continue
					}
					vlist[i] = golua.BaseData(v)
				}
				go func() {
					_, err2 := obj.Call(vn, vlist[1:], time.Duration(0))
					if err2 != nil {
						logger.Warn(tag, "ServiceCall.Post(%s) fail - %s", vn, err2)
					}
				}()
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooServiceCall) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooServiceCall) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooServiceCall) CanClose() bool {
	return false
}

func (gooServiceCall) Close(o interface{}) {
}
