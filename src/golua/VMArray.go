package golua

import "fmt"

// CommonVMArray
type CommonVMArray struct {
	data []interface{}
}

func NewVMArray(a []interface{}) VMArray {
	r := new(CommonVMArray)
	r.data = a
	return r
}

func (this *CommonVMArray) String() string {
	return fmt.Sprintf("@%v", this.data)
}

func (this *CommonVMArray) Get(vm *VM, idx int) (interface{}, error) {
	if idx >= 0 && idx < len(this.data) {
		return this.data[idx], nil
	}
	return nil, fmt.Errorf("index(%d) out range(%d)", idx, len(this.data))
}

func (this *CommonVMArray) Set(vm *VM, idx int, val interface{}) error {
	if idx >= 0 {
		if idx < len(this.data) {
			this.data[idx] = val
			return nil
		}
		if idx == len(this.data) {
			this.data = append(this.data, val)
			return nil
		}
	}
	return fmt.Errorf("index(%d) out range(%d)", idx, len(this.data))
}

func (this *CommonVMArray) Insert(vm *VM, idx int, val interface{}) error {
	if idx >= 0 {
		if idx < len(this.data) {
			this.data = append(this.data, nil)
			copy(this.data[idx+1:], this.data[idx:len(this.data)-1])
			this.data[idx] = val
			return nil
		}
	}
	return fmt.Errorf("index(%d) out range(%d)", idx, len(this.data))
}

func (this *CommonVMArray) Add(vm *VM, val interface{}) error {
	this.data = append(this.data, val)
	return nil
}

func (this *CommonVMArray) Delete(vm *VM, idx int) error {
	if idx >= 0 {
		if idx < len(this.data) {
			this.data[idx] = nil
			copy(this.data[idx:], this.data[idx+1:len(this.data)])
			this.data = this.data[0 : len(this.data)-1]
			return nil
		}
	}
	return fmt.Errorf("index(%d) out range(%d)", idx, len(this.data))
}

func (this *CommonVMArray) SubArray(start int, end int) ([]interface{}, error) {
	if start >= 0 && start < len(this.data) {
		if end > 0 && end <= len(this.data) {
			return this.data[start:end], nil
		}
		return nil, fmt.Errorf("index(%d) out range(%d)", end, len(this.data))
	}
	return nil, fmt.Errorf("index(%d) out range(%d)", start, len(this.data))
}

func (this *CommonVMArray) Len() int {
	return len(this.data)
}

func (this *CommonVMArray) ToArray() []interface{} {
	return this.data
}

func (this *VM) API_array(v interface{}) VMArray {
	if v == nil {
		return nil
	}
	if o, ok := v.(VMArray); ok {
		return o
	}
	if o, ok := v.([]interface{}); ok {
		r := new(CommonVMArray)
		r.data = o
		return r
	}
	return nil
}

func (this *VM) API_toSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	if o, ok := v.([]interface{}); ok {
		return o
	}
	if o, ok := v.(VMArray); ok {
		return o.ToArray()
	}
	return nil
}

func (this *VM) API_newarray() VMArray {
	r := new(CommonVMArray)
	r.data = make([]interface{}, 0)
	return r
}
