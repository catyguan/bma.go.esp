package memserv

import (
	"fmt"
	"strings"
)

func MemGoData(v interface{}) (interface{}, int) {
	if v == nil {
		return nil, 0
	}
	switch v.(type) {
	case bool, int, int8, int16, int32, uint, uint8, uint16, uint32, float32:
		return v, 4
	case int64, float64:
		return v, 8
	case string:
		rv := v.(string)
		return v, 4 + len(rv)
	case []byte:
		rv := v.([]byte)
		return v, 4 + len(rv)
	case []interface{}:
		sz := 8
		a := v.([]interface{})
		for i, av := range a {
			nv, ns := MemGoData(av)
			if ns > 0 {
				sz += ns
			}
			a[i] = nv
		}
		return v, sz
	case map[string]interface{}:
		sz := 4
		m := v.(map[string]interface{})
		for k, mv := range m {
			sz += 4 + len(k)
			nv, ns := MemGoData(mv)
			if ns > 0 {
				sz += ns
			}
			m[k] = nv
		}
		return v, sz
	}
	return nil, -1
}

func Key(typ, n string) string {
	return fmt.Sprintf("%s-%s", typ, n)
}

func SplitTypeName(key string) (string, string) {
	p := strings.SplitN(key, "-", 2)
	if len(p) > 1 {
		return p[0], p[1]
	} else {
		return key, ""
	}
}
