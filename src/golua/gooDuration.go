package golua

import "time"

type gooDuration int

func (gooDuration) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(time.Duration); ok {
		switch key {
		case "Hours":
			return NewGOF("Duration:Hours", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Hours())
				return 1, nil
			}), nil
		case "Minutes":
			return NewGOF("Duration:Minutes", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Minutes())
				return 1, nil
			}), nil
		case "Nanoseconds":
			return NewGOF("Duration:Nanoseconds", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Nanoseconds())
				return 1, nil
			}), nil
		case "Seconds":
			return NewGOF("Duration:Seconds", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Seconds())
				return 1, nil
			}), nil
		case "String":
			return NewGOF("Duration:String", func(vm *VM) (int, error) {
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
