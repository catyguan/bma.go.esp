package golua

type bytesObject struct {
	data []byte
}

type gooBytes int

func ToBytes(v interface{}) []byte {
	if bo, ok := v.(*bytesObject); ok {
		return bo.data
	}
	if bs, ok := v.([]byte); ok {
		return bs
	}
	return nil
}

func CreateGoBytes(bs []byte) VMTable {
	buf := new(bytesObject)
	buf.data = bs
	return NewGOO(buf, gooBytes(0))
}

func (gooBytes) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*bytesObject); ok {
		switch key {
		case "Length":
			return len(obj.data), nil
		case "Size", "Len":
			return NewGOF("Bytes:Len", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(len(obj.data))
				return 1, nil
			}), nil
		case "String":
			return NewGOF("Bytes:String", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(string(obj.data))
				return 1, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooBytes) Set(vm *VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooBytes) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooBytes) CanClose() bool {
	return false
}

func (gooBytes) Close(o interface{}) {
}
