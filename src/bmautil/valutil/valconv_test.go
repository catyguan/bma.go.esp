package valutil

import (
	"fmt"
	"reflect"
	"testing"
)

func doTest(t *testing.T, v interface{}, tt reflect.Type) interface{} {
	var s string
	s = fmt.Sprintf("%s TO %s =>", reflect.TypeOf(v), tt)
	v1, ok := Convert(v, tt)
	t.Log("--", s, v1, ok)
	if !ok {
		t.Error("ERROR", s, v1, ok)
		return nil
	}
	return v1
}

type testObj struct {
	Name string
	Age  int
}

func (this testObj) String() string {
	return "My name is " + this.Name
}

func (this *testObj) ValueConvert(n string) bool {
	this.Name = n
	return true
}

func T2estValconv(t *testing.T) {
	doTest(t, true, BaseType(reflect.String))
	doTest(t, true, BaseType(reflect.Uint32))
	doTest(t, -123, BaseType(reflect.Bool))
	doTest(t, uint(123), BaseType(reflect.Int))
	doTest(t, 123.2, BaseType(reflect.Int16))
	doTest(t, "123.2", BaseType(reflect.Float32))
	doTest(t, reflect.Complex128, BaseType(reflect.Float64))
	doTest(t, testObj{}, BaseType(reflect.String))

	if true {
		tt := reflect.SliceOf(BaseType(reflect.String))
		doTest(t, true, tt)
	}
	if true {
		tt := reflect.SliceOf(BaseType(reflect.String))
		v := [...]int{1, 2, 3, 4, 5}
		doTest(t, v, tt)
	}
	if true {
		tt := reflect.MapOf(BaseType(reflect.String), BaseType(reflect.Int))
		v := make(map[int]int)
		v[1] = 2
		v[2] = 3
		doTest(t, v, tt)
	}
	if true {
		tt := reflect.TypeOf(testObj{})
		v := make(map[string]interface{})
		v["Name"] = "Kitty"
		v["Age"] = 3
		nv := doTest(t, v, tt)
		if nv != nil {
			o := nv.(testObj)
			t.Log("testObj =>", o.Name, o.Age)
		}
	}
	if true {
		tt := reflect.TypeOf((map[string]interface{})(nil))
		v := testObj{"Hello", 5}
		doTest(t, v, tt)
	}
	if true {
		tt := reflect.TypeOf(testObj{})
		v := testObj{"Hello", 5}
		nv := doTest(t, v, tt)
		if nv != nil {
			o := nv.(testObj)
			t.Log("testObj =>", o.Name, o.Age)
		}
	}
	if true {
		tt := reflect.TypeOf(testObj{})
		v := "MyName"
		nv := doTest(t, v, tt)
		if nv != nil {
			o := nv.(testObj)
			t.Log("testObj =>", o.Name, o.Age)
		}
	}
	if true {
		v := 123
		var ptr *int = &v
		doTest(t, ptr, BaseType(reflect.String))
	}
}
