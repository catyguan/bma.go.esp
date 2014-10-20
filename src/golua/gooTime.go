package golua

import (
	"bmautil/valutil"
	"fmt"
	"time"
)

type gooTime int

func (this gooTime) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*time.Time); ok {
		switch key {
		case "Year":
			return NewGOF("Time:Year", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Day())
				return 1, nil
			}), nil
		case "Month":
			return NewGOF("Time:Month", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Day())
				return 1, nil
			}), nil
		case "Day":
			return NewGOF("Time:Day", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Day())
				return 1, nil
			}), nil
		case "Weekday":
			return NewGOF("Time:Weekday", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Weekday())
				return 1, nil
			}), nil
		case "Hour":
			return NewGOF("Time:Hour", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Hour())
				return 1, nil
			}), nil
		case "Minute":
			return NewGOF("Time:Minute", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Minute())
				return 1, nil
			}), nil
		case "Nanosecond":
			return NewGOF("Time:Nanosecond", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Nanosecond())
				return 1, nil
			}), nil
		case "Second":
			return NewGOF("Time:Second", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Second())
				return 1, nil
			}), nil
		case "String":
			return NewGOF("Time:String", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.String())
				return 1, nil
			}), nil
		case "ToMap":
			return NewGOF("Time:ToMap", func(vm *VM) (int, error) {
				vm.API_popAll()
				vm.API_push(vm.API_table(this.ToMap(obj)))
				return 1, nil
			}), nil
		case "Add":
			return NewGOF("Time:Add", func(vm *VM) (int, error) {
				err0 := vm.API_checkstack(2)
				if err0 != nil {
					return 0, err0
				}
				_, v, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				switch v.(type) {
				case string:
					rv := v.(string)
					du, err2 := time.ParseDuration(rv)
					if err2 != nil {
						return 0, err2
					}
					tm := obj.Add(du)
					vm.API_push(&tm)
					return 1, nil
				case int, uint, int8, uint8, int16, uint16, int32, int64, float32, float64:
					rv := valutil.ToInt64(v, 0)
					du := time.Duration(rv) * time.Millisecond
					tm := obj.Add(du)
					vm.API_push(&tm)
					return 1, nil
				case *objectVMTable:
					o := v.(*objectVMTable).o
					if du, ok := o.(time.Duration); ok {
						tm := obj.Add(du)
						vm.API_push(&tm)
						return 1, nil
					}
				}
				return 0, fmt.Errorf("duration invalid(%T)", v)
			}), nil
		}
	}
	return nil, nil
}

func (gooTime) Set(vm *VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooTime) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	if obj, ok := o.(*time.Time); ok {
		r["Year"] = obj.Year()
		r["Month"] = int(obj.Month())
		r["Day"] = obj.Day()
		r["Weekday"] = obj.Weekday()
		r["Hour"] = obj.Hour()
		r["Minute"] = obj.Minute()
		r["Second"] = obj.Second()
		r["Nanosecond"] = obj.Nanosecond()
	}
	return r
}

func (gooTime) CanClose() bool {
	return false
}

func (gooTime) Close(o interface{}) {
}
