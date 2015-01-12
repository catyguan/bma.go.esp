package golua

import (
	"bmautil/valutil"
	"fmt"
	"io/ioutil"
	"os"
)

func CreateGoFile(f *os.File) VMTable {
	return NewGOO(f, gooFile(0))
}

func NewSafeGoFile(vm *VM, f *os.File) VMTable {
	gos := vm.GetGoLua().CreateGoService("file", f, func() {
		f.Close()
	})
	return NewGOO(gos, gooFile(0))
}

func ToFile(o interface{}) *os.File {
	if o == nil {
		return nil
	}
	if f, ok := o.(*os.File); ok {
		return f
	}
	if gos, ok := o.(*GoService); ok {
		if obj, ok2 := gos.Data.(*os.File); ok2 {
			return obj
		}
	}
	return nil
}

type gooFile int

func (this gooFile) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	obj := ToFile(o)
	if obj != nil {
		switch key {
		case "Read":
			return NewGOF("File.Read", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToInt(n, -1)
				if vn < 0 {
					return 0, fmt.Errorf("invalid Read size(%v)", n)
				}
				bs := make([]byte, vn)
				r, err2 := obj.Read(bs)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(CreateGoBytes(bs[:r]))
				vm.API_push(n)
				return 2, nil
			}), nil
		case "ReadAll":
			return NewGOF("File.ReadAll", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				bs, err2 := ioutil.ReadAll(obj)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(CreateGoBytes(bs))
				return 1, nil
			}), nil
		case "Write":
			return NewGOF("File.Write", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				b, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				bs := ToBytes(b)
				if bs == nil {
					return 0, fmt.Errorf("invalid bytes (%T)", b)
				}
				n, err2 := obj.Write(bs)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(n)
				return 1, nil
			}), nil
		case "Close":
			return NewGOF("File.Close", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				if gos, ok := o.(*GoService); ok {
					gos.Close()
				} else {
					obj.Close()
				}
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooFile) Set(vm *VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooFile) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooFile) CanClose() bool {
	return true
}

func (gooFile) Close(o interface{}) {
	if gos, ok := o.(*GoService); ok {
		gos.Close()
		return
	}
	if f, ok := o.(*os.File); ok {
		f.Close()
	}
}
