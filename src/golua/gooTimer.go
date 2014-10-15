package golua

import (
	"bmautil/valutil"
	"fmt"
	"time"
)

type gooTimer int

func (gooTimer) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*time.Timer); ok {
		switch key {
		case "Reset":
			return NewGOF("timer:Reset", func(vm *VM) (int, error) {
				err1 := vm.API_checkstack(1)
				if err1 != nil {
					return 0, err1
				}
				tm, err2 := vm.API_pop1X(-1, true)
				if err2 != nil {
					return 0, err2
				}
				vtm := valutil.ToInt(tm, -1)
				if vtm < 0 {
					return 0, fmt.Errorf("invalid timer time(%v)", tm)
				}
				r := obj.Reset(time.Duration(vtm) * time.Millisecond)
				vm.API_push(r)
				return 1, nil
			}), nil
		case "Stop":
			return NewGOF("timer:Stop", func(vm *VM) (int, error) {
				obj.Stop()
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooTimer) Set(vm *VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooTimer) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

type gooTicker int

func (gooTicker) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*time.Ticker); ok {
		switch key {
		case "Stop":
			return NewGOF("timer:Stop", func(vm *VM) (int, error) {
				obj.Stop()
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooTicker) Set(vm *VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooTicker) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}
