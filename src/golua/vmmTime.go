package golua

import (
	"bmautil/valutil"
	"fmt"
	"time"
)

func TimeModule() *VMModule {
	m := NewVMModule("time")
	m.Init("parseDuration", GOF_time_parseDuration(0))
	m.Init("date", GOF_time_date(0))
	m.Init("now", GOF_time_now(0))
	m.Init("parse", GOF_time_parse(0))
	m.Init("unix", GOF_time_unix(0))
	return m
}

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

// time.parseDuration(s string) string
type GOF_time_parseDuration int

func (this GOF_time_parseDuration) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	ro, err2 := time.ParseDuration(vs)
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(NewGOO(ro, gooDuration(0)))
	return 1, nil
}

func (this GOF_time_parseDuration) IsNative() bool {
	return true
}

func (this GOF_time_parseDuration) String() string {
	return "GoFunc<time.contains>"
}

// time.date(year int, month, day, hour, min, sec int[, loc string]) Time:object
type GOF_time_date int

func (this GOF_time_date) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(6)
	if err0 != nil {
		return 0, err0
	}
	top := vm.API_gettop()
	vs, err1 := vm.API_popN(top, true)
	if err1 != nil {
		return 0, err1
	}
	y := valutil.ToInt(vs[0], 0)
	m := valutil.ToInt(vs[1], 0)
	d := valutil.ToInt(vs[2], 0)
	h := valutil.ToInt(vs[3], 0)
	n := valutil.ToInt(vs[4], 0)
	s := valutil.ToInt(vs[5], 0)
	loc := ""
	if top > 6 {
		loc = valutil.ToString(vs[6], "")
	}
	var vloc *time.Location
	if loc != "" {
		vloc, err1 = time.LoadLocation(loc)
		if err1 != nil {
			return 0, err1
		}
	} else {
		vloc = time.Local
	}
	ro := time.Date(y, time.Month(m), d, h, n, s, 0, vloc)
	vm.API_push(NewGOO(&ro, gooTime(0)))
	return 1, nil
}

func (this GOF_time_date) IsNative() bool {
	return true
}

func (this GOF_time_date) String() string {
	return "GoFunc<time.date>"
}

// time.now() Time:object
type GOF_time_now int

func (this GOF_time_now) Exec(vm *VM) (int, error) {
	vm.API_popAll()
	ro := time.Now()
	vm.API_push(NewGOO(&ro, gooTime(0)))
	return 1, nil
}

func (this GOF_time_now) IsNative() bool {
	return true
}

func (this GOF_time_now) String() string {
	return "GoFunc<time.now>"
}

// time.parse(layout, value string) Time:object
type GOF_time_parse int

func (this GOF_time_parse) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	s1, s2, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs1 := valutil.ToString(s1, "")
	vs2 := valutil.ToString(s2, "")
	if vs2 == "" {
		vs2 = vs1
		vs1 = "2006-01-02 15:04:05"
	}
	ro, err2 := time.ParseInLocation(vs1, vs2, time.Local)
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(NewGOO(&ro, gooTime(0)))
	return 1, nil
}

func (this GOF_time_parse) IsNative() bool {
	return true
}

func (this GOF_time_parse) String() string {
	return "GoFunc<time.parse>"
}

// time.unix(v int64) Time:object
type GOF_time_unix int

func (this GOF_time_unix) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	v, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vv := valutil.ToInt64(v, 0)
	ro := time.Unix(vv, 0)
	vm.API_push(NewGOO(&ro, gooTime(0)))
	return 1, nil
}

func (this GOF_time_unix) IsNative() bool {
	return true
}

func (this GOF_time_unix) String() string {
	return "GoFunc<time.unix>"
}