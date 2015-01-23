package valutil

import (
	"fmt"
	"reflect"
	"strconv"
)

func keep_fmt() {
	fmt.Println("")
}

var (
	typeKinds map[reflect.Kind]reflect.Type = func() map[reflect.Kind]reflect.Type {
		r := make(map[reflect.Kind]reflect.Type)
		r[reflect.Bool] = reflect.TypeOf((*bool)(nil)).Elem()
		r[reflect.Int] = reflect.TypeOf((*int)(nil)).Elem()
		r[reflect.Int8] = reflect.TypeOf((*int8)(nil)).Elem()
		r[reflect.Int16] = reflect.TypeOf((*int16)(nil)).Elem()
		r[reflect.Int32] = reflect.TypeOf((*int32)(nil)).Elem()
		r[reflect.Int64] = reflect.TypeOf((*int64)(nil)).Elem()
		r[reflect.Uint] = reflect.TypeOf((*uint)(nil)).Elem()
		r[reflect.Uint8] = reflect.TypeOf((*uint8)(nil)).Elem()
		r[reflect.Uint16] = reflect.TypeOf((*uint16)(nil)).Elem()
		r[reflect.Uint32] = reflect.TypeOf((*uint32)(nil)).Elem()
		r[reflect.Uint64] = reflect.TypeOf((*uint64)(nil)).Elem()
		r[reflect.Float32] = reflect.TypeOf((*float32)(nil)).Elem()
		r[reflect.Float64] = reflect.TypeOf((*float64)(nil)).Elem()
		r[reflect.String] = reflect.TypeOf((*string)(nil)).Elem()
		return r
	}()
)

func bool01(v bool) int {
	if v {
		return 1
	}
	return 0
}

func nil_conv(tt reflect.Type) (interface{}, bool) {
	toKind := tt.Kind()
	switch toKind {
	case reflect.Bool:
		return false, true
	case reflect.Int:
		return int(0), true
	case reflect.Int8:
		return int8(0), true
	case reflect.Int16:
		return int16(0), true
	case reflect.Int32:
		return int32(0), true
	case reflect.Int64:
		return int64(0), true
	case reflect.Uint:
		return uint(0), true
	case reflect.Uint8:
		return uint8(0), true
	case reflect.Uint16:
		return uint16(0), true
	case reflect.Uint32:
		return uint32(0), true
	case reflect.Uint64:
		return uint64(0), true
	case reflect.Float32:
		return float32(0), true
	case reflect.Float64:
		return float64(0), true
	case reflect.String:
		return "<nil>", true
	default:
		return nil, true
	}
}

func bool_conv(iv interface{}, tt reflect.Type) (interface{}, bool) {
	toKind := tt.Kind()
	v := iv.(bool)
	switch toKind {
	case reflect.Bool:
		return v, true
	case reflect.Int:
		return int(bool01(v)), true
	case reflect.Int8:
		return int8(bool01(v)), true
	case reflect.Int16:
		return int16(bool01(v)), true
	case reflect.Int32:
		return int32(bool01(v)), true
	case reflect.Int64:
		return int64(bool01(v)), true
	case reflect.Uint:
		return uint(bool01(v)), true
	case reflect.Uint8:
		return uint8(bool01(v)), true
	case reflect.Uint16:
		return uint16(bool01(v)), true
	case reflect.Uint32:
		return uint32(bool01(v)), true
	case reflect.Uint64:
		return uint64(bool01(v)), true
	case reflect.Float32:
		return float32(bool01(v)), true
	case reflect.Float64:
		return float64(bool01(v)), true
	case reflect.String:
		if v {
			return "true", true
		}
		return "false", true
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.Struct:
	case reflect.Uintptr:
	case reflect.UnsafePointer:
		return nil, false
	}
	return nil, false
}

func int_conv(v int64, tt reflect.Type) (interface{}, bool) {
	toKind := tt.Kind()
	switch toKind {
	case reflect.Bool:
		if v != 0 {
			return true, true
		}
		return false, true
	case reflect.Int:
		return int(v), true
	case reflect.Int8:
		return int8(v), true
	case reflect.Int16:
		return int16(v), true
	case reflect.Int32:
		return int32(v), true
	case reflect.Int64:
		return v, true
	case reflect.Uint:
		return uint(v), true
	case reflect.Uint8:
		return uint8(v), true
	case reflect.Uint16:
		return uint16(v), true
	case reflect.Uint32:
		return uint32(v), true
	case reflect.Uint64:
		return uint64(v), true
	case reflect.Float32:
		return float32(v), true
	case reflect.Float64:
		return float64(v), true
	case reflect.String:
		return strconv.FormatInt(v, 10), true
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.Struct:
	case reflect.Uintptr:
	case reflect.UnsafePointer:
		return nil, false
	}
	return nil, false
}

func uint_conv(v uint64, tt reflect.Type) (interface{}, bool) {
	toKind := tt.Kind()
	switch toKind {
	case reflect.Bool:
		if v != 0 {
			return true, true
		}
		return false, true
	case reflect.Int:
		return int(v), true
	case reflect.Int8:
		return int8(v), true
	case reflect.Int16:
		return int16(v), true
	case reflect.Int32:
		return int32(v), true
	case reflect.Int64:
		return int64(v), true
	case reflect.Uint:
		return uint(v), true
	case reflect.Uint8:
		return uint8(v), true
	case reflect.Uint16:
		return uint16(v), true
	case reflect.Uint32:
		return uint32(v), true
	case reflect.Uint64:
		return v, true
	case reflect.Float32:
		return float32(v), true
	case reflect.Float64:
		return float64(v), true
	case reflect.String:
		return strconv.FormatUint(v, 10), true
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.Struct:
	case reflect.Uintptr:
	case reflect.UnsafePointer:
		return nil, false
	}
	return nil, false
}

func float_conv(v float64, tt reflect.Type) (interface{}, bool) {
	toKind := tt.Kind()
	switch toKind {
	case reflect.Bool:
		if v != 0 {
			return true, true
		}
		return false, true
	case reflect.Int:
		return int(v), true
	case reflect.Int8:
		return int8(v), true
	case reflect.Int16:
		return int16(v), true
	case reflect.Int32:
		return int32(v), true
	case reflect.Int64:
		return int64(v), true
	case reflect.Uint:
		return uint(v), true
	case reflect.Uint8:
		return uint8(v), true
	case reflect.Uint16:
		return uint16(v), true
	case reflect.Uint32:
		return uint32(v), true
	case reflect.Uint64:
		return uint64(v), true
	case reflect.Float32:
		return float32(v), true
	case reflect.Float64:
		return v, true
	case reflect.String:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.Struct:
	case reflect.Uintptr:
	case reflect.UnsafePointer:
		return nil, false
	}
	return nil, false
}

func string_int64(v string) (int64, bool) {
	o, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return o, true
	}
	return 0, false
}

func string_uint64(v string) (uint64, bool) {
	o, err := strconv.ParseUint(v, 10, 64)
	if err == nil {
		return o, true
	}
	return 0, false
}

func string_float64(v string) (float64, bool) {
	o, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return o, true
	}
	return 0, false
}

func string_conv(vv interface{}, tt reflect.Type) (interface{}, bool) {
	toKind := tt.Kind()
	v := vv.(string)
	switch toKind {
	case reflect.Bool:
		o, err := strconv.ParseBool(v)
		if err == nil {
			return o, true
		}
		return nil, false
	case reflect.Int:
		if r, ok := string_int64(v); ok {
			return int(r), true
		}
		return nil, false
	case reflect.Int8:
		if r, ok := string_int64(v); ok {
			return int8(r), true
		}
		return nil, false
	case reflect.Int16:
		if r, ok := string_int64(v); ok {
			return int16(r), true
		}
		return nil, false
	case reflect.Int32:
		if r, ok := string_int64(v); ok {
			return int32(r), true
		}
		return nil, false
	case reflect.Int64:
		return string_int64(v)
	case reflect.Uint:
		if r, ok := string_uint64(v); ok {
			return uint(r), true
		}
	case reflect.Uint8:
		if r, ok := string_uint64(v); ok {
			return uint8(r), true
		}
	case reflect.Uint16:
		if r, ok := string_uint64(v); ok {
			return uint16(r), true
		}
	case reflect.Uint32:
		if r, ok := string_uint64(v); ok {
			return uint32(r), true
		}
	case reflect.Uint64:
		return string_uint64(v)
	case reflect.Float32:
		if r, ok := string_float64(v); ok {
			return float32(r), true
		}
	case reflect.Float64:
		return string_float64(v)
	case reflect.String:
		return v, true
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.Struct:
	case reflect.Uintptr:
	case reflect.UnsafePointer:
	}
	return nil, false
}

func array_conv(vv reflect.Value, tt reflect.Type) (interface{}, bool) {
	if tt.Kind() != reflect.Slice {
		return nil, false
	}
	if !vv.IsValid() {
		return nil, false
	}
	sz := vv.Len()
	r := reflect.MakeSlice(tt, sz, sz)
	et := tt.Elem()
	for i := 0; i < sz; i++ {
		sval := vv.Index(i)
		if nval, ok := conv(sval, et); ok {
			r.Index(i).Set(reflect.ValueOf(nval))
		} else {
			r.Index(i).Set(reflect.Zero(et))
		}
	}
	return r.Interface(), true
}

func map_conv(vv reflect.Value, tt reflect.Type) (interface{}, bool) {
	switch tt.Kind() {
	case reflect.Map:
		return map2map(vv, tt)
	case reflect.Struct:
		if r, ok := map2struct(vv, tt); ok {
			return r.Interface(), true
		}
	case reflect.Ptr:
		if tt.Elem().Kind() == reflect.Struct {
			if r, ok := map2struct(vv, tt.Elem()); ok {
				return r.Addr().Interface(), true
			}
		}
	}
	return nil, false
}

func map2map(vv reflect.Value, tt reflect.Type) (interface{}, bool) {
	if !vv.IsValid() {
		return nil, false
	}
	r := reflect.MakeMap(tt)
	kt := tt.Key()
	et := tt.Elem()

	mkeys := vv.MapKeys()
	for _, k := range mkeys {
		if nkey, ok := conv(k, kt); ok {
			sval := vv.MapIndex(k)
			if nval, ok := conv(sval, et); ok {
				r.SetMapIndex(reflect.ValueOf(nkey), reflect.ValueOf(nval))
			} else {
				r.SetMapIndex(reflect.ValueOf(nkey), reflect.Zero(et))
			}
		}
	}
	return r.Interface(), true
}

func map2struct(vv reflect.Value, tt reflect.Type) (reflect.Value, bool) {
	if !vv.IsValid() {
		return reflect.Zero(tt), false
	}
	if vv.Type().Key().Kind() != reflect.String {
		return reflect.Zero(tt), false
	}

	r := reflect.New(tt).Elem()

	mkeys := vv.MapKeys()
	for _, k := range mkeys {
		name := k.String()
		if field := r.FieldByName(name); field.IsValid() {
			sval := vv.MapIndex(k)
			et := field.Type()
			if nval, ok := conv(sval, et); ok {
				field.Set(reflect.ValueOf(nval))
			} else {
				field.Set(reflect.Zero(et))
			}
		}
	}
	return r, true
}

func struct_conv(vv reflect.Value, tt reflect.Type) (interface{}, bool) {
	switch tt.Kind() {
	case reflect.Map:
		return struct2map(vv, tt)
	case reflect.Struct:
		if r, ok := struct2struct(vv, tt); ok {
			return r.Interface(), true
		}
	case reflect.Ptr:
		if tt.Elem().Kind() == reflect.Struct {
			if r, ok := struct2struct(vv, tt.Elem()); ok {
				return r.Addr().Interface(), true
			}
		}
	}
	return nil, false
}

func struct2map(vv reflect.Value, tt reflect.Type) (interface{}, bool) {
	if !vv.IsValid() {
		return nil, false
	}
	r := reflect.MakeMap(tt)
	kt := tt.Key()
	et := tt.Elem()

	vt := vv.Type()
	sz := vt.NumField()
	for i := 0; i < sz; i++ {
		tfield := vt.Field(i)
		if nkey, ok := Convert(tfield.Name, kt); ok {
			sval := vv.Field(i)
			if nval, ok := conv(sval, et); ok {
				r.SetMapIndex(reflect.ValueOf(nkey), reflect.ValueOf(nval))
			} else {
				r.SetMapIndex(reflect.ValueOf(nkey), reflect.Zero(et))
			}
		}
	}
	return r.Interface(), true
}

func struct2struct(vv reflect.Value, tt reflect.Type) (reflect.Value, bool) {
	if !vv.IsValid() {
		return reflect.Zero(tt), false
	}

	r := reflect.New(tt).Elem()
	vt := vv.Type()
	sz := vt.NumField()
	for i := 0; i < sz; i++ {
		tfield := vt.Field(i)
		name := tfield.Name
		if field := r.FieldByName(name); field.IsValid() {
			sval := vv.Field(i)
			et := field.Type()
			if nval, ok := conv(sval, et); ok {
				field.Set(reflect.ValueOf(nval))
			} else {
				field.Set(reflect.Zero(et))
			}
		}
	}
	return r, true
}

func base_convert(vv reflect.Value, tt reflect.Type) (interface{}, bool) {
	switch vv.Kind() {
	case reflect.Bool:
		return bool_conv(vv.Bool(), tt)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int_conv(vv.Int(), tt)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint_conv(vv.Uint(), tt)
	case reflect.Float32, reflect.Float64:
		return float_conv(vv.Float(), tt)
	case reflect.String:
		return string_conv(vv.String(), tt)
	case reflect.Array, reflect.Slice:
		return array_conv(vv, tt)
	case reflect.Map:
		return map_conv(vv, tt)
	case reflect.Ptr:
		return conv(vv.Elem(), tt)
	case reflect.Struct:
		return struct_conv(vv, tt)
	case reflect.Interface:
		return Convert(vv.Interface(), tt)
	}
	return nil, false
}

func slice_maker(iv reflect.Value, tt reflect.Type) (interface{}, bool) {
	if nv, ok := conv(iv, tt.Elem()); ok {
		r := reflect.MakeSlice(tt, 1, 1)
		r.Index(0).Set(reflect.ValueOf(nv))
		return r.Interface(), true
	}
	return nil, false
}

func BaseType(k reflect.Kind) reflect.Type {
	if r, ok := typeKinds[k]; ok {
		return r
	}
	return nil
}

func Convert(v interface{}, toType reflect.Type) (interface{}, bool) {
	if v == nil {
		return nil, false
	}
	val := reflect.ValueOf(v)
	return conv(val, toType)
}

func conv(val reflect.Value, toType reflect.Type) (interface{}, bool) {
	if !val.IsValid() {
		return nil_conv(toType)
	}
	valType := val.Type()
	if valType.AssignableTo(toType) {
		// fmt.Println(1)
		return val.Interface(), true
	}
	if valType.ConvertibleTo(toType) {
		// fmt.Println(2)
		if toType.Kind() != reflect.String {
			return val.Convert(toType).Interface(), true
		}
	}

	if r, ok := base_convert(val, toType); ok {
		// fmt.Println(3)
		return r, true
	}
	kind := toType.Kind()
	switch kind {
	case reflect.String:
		if method := val.MethodByName("String"); method.IsValid() {
			methodType := method.Type()
			if methodType.NumIn() == 0 && methodType.NumOut() == 1 && methodType.Out(0).Kind() == reflect.String {
				return method.Call([]reflect.Value{}), true
			}
		}
	case reflect.Slice:
		return slice_maker(val, toType)
	case reflect.Struct:
		m, ok := reflect.PtrTo(toType).MethodByName("ValueConvert")
		if ok {
			methodType := m.Type
			if methodType.NumIn() == 2 && val.Type().AssignableTo(methodType.In(1)) && methodType.NumOut() == 1 && methodType.Out(0).Kind() == reflect.Bool {
				obj := reflect.New(toType)
				if method := obj.MethodByName("ValueConvert"); method.IsValid() {
					isDone := method.Call([]reflect.Value{val})[0]
					if isDone.Bool() {
						return obj.Elem().Interface(), true
					}
				}
			}
		}
	}
	return nil, false
}
