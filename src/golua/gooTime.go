package golua

import (
	"bmautil/valutil"
	"fmt"
	"time"
)

func CreateGoTime(tm *time.Time) VMTable {
	return NewGOO(tm, gooTime(0))
}

type gooTime int

func (this gooTime) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*time.Time); ok {
		switch key {
		case "Year":
			return NewGOF("Time.Year", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Day())
				return 1, nil
			}), nil
		case "YearDay":
			return NewGOF("Time.YearDay", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.YearDay())
				return 1, nil
			}), nil
		case "Month":
			return NewGOF("Time.Month", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Day())
				return 1, nil
			}), nil
		case "Day":
			return NewGOF("Time.Day", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Day())
				return 1, nil
			}), nil
		case "Weekday":
			return NewGOF("Time.Weekday", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Weekday())
				return 1, nil
			}), nil
		case "Hour":
			return NewGOF("Time.Hour", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Hour())
				return 1, nil
			}), nil
		case "Minute":
			return NewGOF("Time.Minute", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Minute())
				return 1, nil
			}), nil
		case "Nanosecond":
			return NewGOF("Time.Nanosecond", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Nanosecond())
				return 1, nil
			}), nil
		case "Second":
			return NewGOF("Time.Second", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Second())
				return 1, nil
			}), nil
		case "String":
			return NewGOF("Time.String", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.String())
				return 1, nil
			}), nil
		case "ToMap":
			return NewGOF("Time.ToMap", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(vm.API_table(this.ToMap(obj)))
				return 1, nil
			}), nil
		case "Clock":
			return NewGOF("Time.Clock", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				h, m, s := obj.Clock()
				vm.API_push(h)
				vm.API_push(m)
				vm.API_push(s)
				return 3, nil
			}), nil
		case "Date":
			return NewGOF("Time.Clock", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				y, m, d := obj.Date()
				vm.API_push(y)
				vm.API_push(int(m))
				vm.API_push(d)
				return 3, nil
			}), nil
		case "Locate":
			return NewGOF("Time.Locate", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				loc := obj.Location()
				vm.API_push(loc.String())
				return 1, nil
			}), nil
		case "Add":
			return NewGOF("Time.Add", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkstack(1)
				if err0 != nil {
					return 0, err0
				}
				v, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				du, err2 := ToDuration(v)
				if err2 != nil {
					return 0, err2
				}
				tm := obj.Add(du)
				vm.API_push(NewGOO(&tm, gooTime(0)))
				return 1, nil
			}), nil
		case "AddDate":
			return NewGOF("Time.AddDate", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkstack(3)
				if err0 != nil {
					return 0, err0
				}
				y, m, d, err1 := vm.API_pop3X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vy := valutil.ToInt(y, 0)
				vmo := valutil.ToInt(m, 0)
				vd := valutil.ToInt(d, 0)
				tm := obj.AddDate(vy, vmo, vd)
				vm.API_push(NewGOO(&tm, gooTime(0)))
				return 1, nil
			}), nil
		case "After", "Before", "Equal":
			return NewGOF("Time.After", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkstack(1)
				if err0 != nil {
					return 0, err0
				}
				ctm, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vo := vm.API_object(ctm)
				if vo != nil {
					if vctm, ok := vo.(*time.Time); ok {
						var rv bool
						if key == "After" {
							rv = obj.After(*vctm)
						} else if key == "Before" {
							rv = obj.Before(*vctm)
						} else {
							rv = obj.Equal(*vctm)
						}
						vm.API_push(rv)
						return 1, nil
					}
				}
				return 0, fmt.Errorf("checkTime invalid(%v)", ctm)
			}), nil
		case "Format":
			return NewGOF("Time.Format", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkstack(1)
				if err0 != nil {
					return 0, err0
				}
				s, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vs := valutil.ToString(s, "")
				rv := obj.Format(vs)
				vm.API_push(rv)
				return 1, nil
			}), nil
		case "In":
			return NewGOF("Time.In", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkstack(1)
				if err0 != nil {
					return 0, err0
				}
				s, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vs := valutil.ToString(s, "")
				loc, err2 := time.LoadLocation(vs)
				if err2 != nil {
					return 0, err2
				}
				tm := obj.In(loc)
				vm.API_push(NewGOO(&tm, gooTime(0)))
				return 1, nil
			}), nil
		case "Local":
			return NewGOF("Time.Local", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				tm := obj.Local()
				vm.API_push(NewGOO(&tm, gooTime(0)))
				return 1, nil
			}), nil
		case "Round":
			return NewGOF("Time.Round", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkstack(1)
				if err0 != nil {
					return 0, err0
				}
				du, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vdu, err2 := ToDuration(du)
				if err2 != nil {
					return 0, err2
				}
				tm := obj.Round(vdu)
				vm.API_push(NewGOO(&tm, gooTime(0)))
				return 1, nil
			}), nil
		case "Sub":
			return NewGOF("Time.Sub", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkstack(1)
				if err0 != nil {
					return 0, err0
				}
				stm, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vo := vm.API_object(stm)
				if vo != nil {
					if vstm, ok := vo.(*time.Time); ok {
						du := obj.Sub(*vstm)
						vm.API_push(CreateDuration(du))
						return 1, nil
					}
				}
				return 0, fmt.Errorf("subTime invalid(%v)", stm)
			}), nil
		case "Truncate":
			return NewGOF("Time.Truncate", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkstack(1)
				if err0 != nil {
					return 0, err0
				}
				du, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vdu, err2 := ToDuration(du)
				if err2 != nil {
					return 0, err2
				}
				tm := obj.Truncate(vdu)
				vm.API_push(NewGOO(&tm, gooTime(0)))
				return 1, nil
			}), nil
		case "UTC":
			return NewGOF("Time.UTC", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				tm := obj.UTC()
				vm.API_push(NewGOO(&tm, gooTime(0)))
				return 1, nil
			}), nil
		case "Unix":
			return NewGOF("Time.Unix", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				v := obj.Unix()
				vm.API_push(v)
				return 1, nil
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
		y, m, d := obj.Date()
		h, n, s := obj.Clock()
		r["Year"] = y
		r["YearDay"] = obj.YearDay()
		r["Month"] = int(m)
		r["Day"] = d
		r["Weekday"] = obj.Weekday()
		r["Hour"] = h
		r["Minute"] = n
		r["Second"] = s
		r["Nanosecond"] = obj.Nanosecond()
	}
	return r
}

func (gooTime) CanClose() bool {
	return false
}

func (gooTime) Close(o interface{}) {
}
