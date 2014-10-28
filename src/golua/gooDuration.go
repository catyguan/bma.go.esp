package golua

import (
	"bmautil/valutil"
	"fmt"
	"time"
)

func ToDuration(v interface{}) (time.Duration, error) {
	switch v.(type) {
	case string:
		rv := v.(string)
		return time.ParseDuration(rv)
	case int, uint, int8, uint8, int16, uint16, int32, int64, float32, float64:
		rv := valutil.ToInt64(v, 0)
		return time.Duration(rv) * time.Millisecond, nil
	case *objectVMTable:
		o := v.(*objectVMTable).o
		if du, ok := o.(time.Duration); ok {
			return du, nil
		}
	}
	return 0, fmt.Errorf("duration invalid(%v)", v)
}

type gooDuration int

func (gooDuration) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(time.Duration); ok {
		switch key {
		case "Hours":
			return NewGOF("Duration:Hours", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Hours())
				return 1, nil
			}), nil
		case "Minutes":
			return NewGOF("Duration:Minutes", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Minutes())
				return 1, nil
			}), nil
		case "Nanoseconds":
			return NewGOF("Duration:Nanoseconds", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Nanoseconds())
				return 1, nil
			}), nil
		case "Seconds":
			return NewGOF("Duration:Seconds", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Seconds())
				return 1, nil
			}), nil
		case "String":
			return NewGOF("Duration:String", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.String())
				return 1, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooDuration) Set(vm *VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooDuration) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooDuration) CanClose() bool {
	return false
}

func (gooDuration) Close(o interface{}) {
}
