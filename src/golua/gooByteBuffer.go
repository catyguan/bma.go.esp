package golua

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
)

func ByteBufferFactory(vm *VM, n string) (interface{}, error) {
	return CreateGoByteBuffer(nil), nil
}

type gooByteBuffer int

func CreateGoByteBuffer(bs []byte) VMTable {
	if bs == nil {
		bs = []byte{}
	}
	buf := bytes.NewBuffer(bs)
	return CreateGoByteBuffer2(buf)
}

func CreateGoByteBuffer2(bb *bytes.Buffer) VMTable {
	return NewGOO(bb, gooByteBuffer(0))
}

func (gooByteBuffer) Get(vm *VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*bytes.Buffer); ok {
		switch key {
		case "Bytes":
			return NewGOF("ByteBuffer:Bytes", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(CreateGoBytes(obj.Bytes()))
				return 1, nil
			}), nil
		case "Grow":
			return NewGOF("ByteBuffer:Grow", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToInt(n, 0)
				obj.Grow(vn)
				return 0, nil
			}), nil
		case "Size", "Len":
			return NewGOF("ByteBuffer:Len", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.Len())
				return 1, nil
			}), nil
		case "Next", "Read":
			return NewGOF("ByteBuffer:Read", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToInt(n, 0)
				bs := obj.Next(vn)
				vm.API_push(CreateGoBytes(bs))
				return 1, nil
			}), nil
		case "ReadByte":
			return NewGOF("ByteBuffer:ReadByte", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				b, err2 := obj.ReadByte()
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(b)
				return 1, nil
			}), nil
		case "ReadRune":
			return NewGOF("ByteBuffer:ReadRune", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				b, n, err2 := obj.ReadRune()
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(int(b))
				vm.API_push(n)
				return 2, nil
			}), nil
		case "Reset":
			return NewGOF("ByteBuffer:Reset", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				obj.Reset()
				return 0, nil
			}), nil
		case "String":
			return NewGOF("ByteBuffer:String", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				vm.API_push(obj.String())
				return 1, nil
			}), nil
		case "Truncate":
			return NewGOF("ByteBuffer:Truncate", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToInt(n, 0)
				obj.Truncate(vn)
				return 0, nil
			}), nil
		case "UnreadByte":
			return NewGOF("ByteBuffer:UnreadByte", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				err2 := obj.UnreadByte()
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "UnreadRune":
			return NewGOF("ByteBuffer:UnreadRune", func(vm *VM, self interface{}) (int, error) {
				vm.API_popAll()
				err2 := obj.UnreadRune()
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "Write", "WriteString":
			return NewGOF("ByteBuffer:Write", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				v, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				bs := ToBytes(v)
				if bs == nil {
					return 0, fmt.Errorf("invalid write data(%T)", v)
				}
				n, err2 := obj.Write(bs)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(n)
				return 1, nil
			}), nil
		case "WriteByte":
			return NewGOF("ByteBuffer:WriteByte", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				v, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				bs := valutil.ToByte(v, 0)
				err2 := obj.WriteByte(bs)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(1)
				return 1, nil
			}), nil
		case "WriteRune":
			return NewGOF("ByteBuffer:WriteRune", func(vm *VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				v, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				bs := valutil.ToInt(v, 0)
				n, err2 := obj.WriteRune(rune(bs))
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(n)
				return 1, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooByteBuffer) Set(vm *VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooByteBuffer) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooByteBuffer) CanClose() bool {
	return false
}

func (gooByteBuffer) Close(o interface{}) {
}
