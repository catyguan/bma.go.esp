package golua

import "fmt"

func (this *VM) API_toMap(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	if o, ok := v.(map[string]interface{}); ok {
		return o
	}
	if o, ok := v.(VMTable); ok {
		return o.ToMap()
	}
	return nil
}

func (this *VM) API_table(v interface{}) (VMTable, map[string]interface{}) {
	if v == nil {
		return nil, nil
	}
	if o, ok := v.(VMTable); ok {
		return o, nil
	}
	if m, ok := v.(map[string]interface{}); ok {
		return nil, m
	}
	return nil, nil
}

// objectVMTable
type objectVMTable struct {
	o interface{}
	p GoObject
}

func NewGOO(o interface{}, p GoObject) VMTable {
	r := new(objectVMTable)
	r.o = o
	r.p = p
	return r
}

func (this *objectVMTable) String() string {
	return fmt.Sprintf("@%v", this.o)
}

func (this *objectVMTable) Get(vm *VM, key string) (interface{}, error) {
	r, err := this.p.Get(vm, this.o, key)
	// fmt.Println("objectVMTable:Get", this.o, key, r, err)
	return r, err
}

func (this *objectVMTable) Rawget(key string) interface{} {
	return nil
}

func (this *objectVMTable) Set(vm *VM, key string, val interface{}) error {
	return this.p.Set(vm, this.o, key, val)
}

func (this *objectVMTable) Rawset(key string, val interface{}) {

}

func (this *objectVMTable) Delete(key string) {

}

func (this *objectVMTable) Len() int {
	return 0
}

func (this *objectVMTable) ToMap() map[string]interface{} {
	return this.p.ToMap(this.o)
}

func (this *VM) API_object(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	if o, ok := v.(*objectVMTable); ok {
		return o.o
	}
	return nil
}

func BaseData(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case bool, int, int8, int16, int32, uint, uint8, uint16, uint32, float32, int64, float64:
		return v
	case string, []byte, []interface{}:
		return v
	case map[string]interface{}:
		m := v.(map[string]interface{})
		rm := make(map[string]interface{})
		for k, v := range m {
			rm[k] = BaseData(v)
		}
		return rm
	case VMArray:
		a := v.(VMArray).ToArray()
		ra := make([]interface{}, len(a))
		for k, v := range a {
			ra[k] = BaseData(v)
		}
		return ra
	}
	return nil
}

func ScriptData(d interface{}) interface{} {
	if d == nil {
		return nil
	}
	switch ro := d.(type) {
	case map[string]interface{}:
		for k, v := range ro {
			ro[k] = ScriptData(v)
		}
		return ro
	case []interface{}:
		for k, v := range ro {
			ro[k] = ScriptData(v)
		}
		r := NewVMArray(ro)
		return r
	}
	return d
}
