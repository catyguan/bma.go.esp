package golua

import (
	"fmt"
	"sync"
)

type gooRMutex int

func (gooRMutex) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if mux, ok := o.(*sync.RWMutex); ok {
		switch key {
		case "RLocker":
			return NewGOF("rmutex:RLocker", func(vm *VM) (int, error) {
				nobj := NewGOO(mux.RLocker(), gooLocker(0))
				vm.API_push(nobj)
				return 1, nil
			}), nil
		}
	}
	return gooLocker(0).Get(vm, o, key)
}

func (gooRMutex) Set(vm *VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooRMutex) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

type gooLocker int

func (gooLocker) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if mux, ok := o.(sync.Locker); ok {
		switch key {
		case "Lock":
			return NewGOF("locker:Lock", func(vm *VM) (int, error) {
				mux.Lock()
				return 0, nil
			}), nil
		case "Unlock":
			return NewGOF("locker:Unlock", func(vm *VM) (int, error) {
				mux.Unlock()
				return 0, nil
			}), nil
		case "Sync":
			return NewGOF("locker:Sync", func(vm *VM) (int, error) {
				// o:Sync(f)
				_, f, err0 := vm.API_pop2X(-1, true)
				if err0 != nil {
					return 0, err0
				}
				if !vm.API_canCall(f) {
					return 0, fmt.Errorf("sync func(%T) can't call", f)
				}
				mux.Lock()
				defer mux.Unlock()
				vm.API_push(f)
				r, err2 := vm.Call(0, -1)
				if err2 != nil {
					return r, err2
				}
				return r, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooLocker) Set(vm *VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooLocker) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}