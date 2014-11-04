package valutil

import (
	// 	"fmt"
	"encoding/hex"
	"encoding/json"
	"reflect"
)

func ToBool(v interface{}, defv bool) bool {
	if v == nil {
		return false
	}
	if r, ok := Convert(v, reflect.TypeOf(defv)); ok {
		return r.(bool)
	}
	return defv
}

func ToBoolNil(v interface{}) (bool, bool) {
	if r, ok := Convert(v, reflect.TypeOf(false)); ok {
		return r.(bool), true
	}
	return false, false
}

func ToInt64(v interface{}, defv int64) int64 {
	if r, ok := Convert(v, reflect.TypeOf(defv)); ok {
		return r.(int64)
	}
	return defv
}

func ToInt(v interface{}, defv int) int {
	return int(ToInt64(v, int64(defv)))
}

func ToInt8(v interface{}, defv int8) int8 {
	return int8(ToInt64(v, int64(defv)))
}

func ToInt16(v interface{}, defv int16) int16 {
	return int16(ToInt64(v, int64(defv)))
}

func ToInt32(v interface{}, defv int32) int32 {
	return int32(ToInt64(v, int64(defv)))
}

func ToUint64(v interface{}, defv uint64) uint64 {
	if r, ok := Convert(v, reflect.TypeOf(defv)); ok {
		return r.(uint64)
	}
	return defv
}

func ToByte(v interface{}, defv byte) byte {
	return byte(ToUint64(v, uint64(defv)))
}

func ToUint8(v interface{}, defv uint8) uint8 {
	return uint8(ToUint64(v, uint64(defv)))
}

func ToUint16(v interface{}, defv uint16) uint16 {
	return uint16(ToUint64(v, uint64(defv)))
}

func ToUint32(v interface{}, defv uint32) uint32 {
	return uint32(ToUint64(v, uint64(defv)))
}

func ToFloat64(v interface{}, defv float64) float64 {
	if r, ok := Convert(v, reflect.TypeOf(defv)); ok {
		return r.(float64)
	}
	return defv
}

func ToFloat32(v interface{}, defv float32) float32 {
	if r, ok := Convert(v, reflect.TypeOf(defv)); ok {
		return r.(float32)
	}
	return defv
}

func ToString(v interface{}, defv string) string {
	if r, ok := Convert(v, reflect.TypeOf(defv)); ok {
		if rs, ok := r.(string); ok {
			return rs
		}
	}
	return defv
}

func ToArray(v interface{}) []interface{} {
	switch v.(type) {
	case []interface{}:
		return v.([]interface{})
	}
	return nil
}

func ToSlice(v interface{}, elemType reflect.Type) interface{} {
	atype := reflect.SliceOf(elemType)
	if r, ok := Convert(v, atype); ok {
		return r
	}
	return nil
}

func ToStringMap(v interface{}) map[string]interface{} {
	switch v.(type) {
	case map[string]interface{}:
		return v.(map[string]interface{})
	}
	return nil
}

func ToMap(v interface{}, keyType reflect.Type, elemType reflect.Type) interface{} {
	atype := reflect.MapOf(keyType, elemType)
	if r, ok := Convert(v, atype); ok {
		return r
	}
	return nil
}

func ToBean(m map[string]interface{}, beanPtr interface{}) bool {
	if m == nil {
		return false
	}
	val := reflect.ValueOf(beanPtr)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return false
	}
	val = val.Elem()

	for k, v := range m {
		if field := val.FieldByName(k); field.IsValid() {
			et := field.Type()
			if nval, ok := Convert(v, et); ok {
				field.Set(reflect.ValueOf(nval))
			} else {
				field.Set(reflect.Zero(et))
			}
		}
	}
	return true
}

func BeanToMap(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	b, _ := json.Marshal(v)
	r := make(map[string]interface{})
	json.Unmarshal(b, &r)
	return r
}

func ToBytes(s string) []byte {
	r, err := hex.DecodeString(s)
	if err != nil {
		return []byte{}
	}
	return r
}
