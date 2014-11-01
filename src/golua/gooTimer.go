package golua

import (
	"bmautil/valutil"
	"fmt"
	"time"
)

func CreateGoTimer(tm *time.Timer, gos *GoService) VMTable {
	gos.CloseFunc = func() {
		tm.Stop()
	}
	gos.Data = tm
	goo := NewGOO(gos, gooTimer(0))
	return goo
}

type gooTimer int

func (gooTimer) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if gos, ok := o.(*GoService); ok {
		obj := gos.Data.(*time.Timer)
		switch key {
		case "Reset":
			return NewGOF("timer.Reset", func(vm *VM, self interface{}) (int, error) {
				err1 := vm.API_checkStack(1)
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
			return NewGOF("timer.Stop", func(vm *VM, self interface{}) (int, error) {
				gos.Close()
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

func (gooTimer) CanClose() bool {
	return false
}

func (gooTimer) Close(o interface{}) {
}

//////////////////////////// Ticker
func CreateGoTicker(tm *time.Ticker, gos *GoService) VMTable {
	gos.CloseFunc = func() {
		tm.Stop()
	}
	gos.Data = tm
	goo := NewGOO(gos, gooTicker(0))
	return goo
}

type gooTicker int

func (gooTicker) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if gos, ok := o.(*GoService); ok {
		obj := gos.Data.(*time.Ticker)
		switch key {
		case "Stop":
			return NewGOF("timer:Stop", func(vm *VM, self interface{}) (int, error) {
				if obj != nil {
					gos.Close()
				}
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

func (gooTicker) CanClose() bool {
	return false
}

func (gooTicker) Close(o interface{}) {
}
