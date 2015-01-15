package gom

import (
	"bytes"
	"golua"
)

type MValType struct {
	name       string
	innerType1 *MValType
	innerType2 *MValType
	moo
}

func (this *MValType) String() string {
	buf := bytes.NewBuffer(make([]byte, 0, 16))
	buf.WriteString(this.name)
	if this.innerType1 != nil {
		buf.WriteByte('<')
		buf.WriteString(this.innerType1.String())
		if this.innerType2 != nil {
			buf.WriteByte(',')
			buf.WriteString(this.innerType2.String())
		}
		buf.WriteByte('>')
	}
	return buf.String()
}

type MValue struct {
	annos *MAnnotations
	value interface{}
}

//// vmm
func (this *MValType) ToVMTable() golua.VMTable {
	return this
}
func (this *MValType) Get(vm *golua.VM, key string) (interface{}, error) {
	switch key {
	case "Name":
		return this.name, nil
	case "Kind":
		return "Type", nil
	case "String":
		return this.String(), nil
	case "InnerType", "InnerType1":
		return this.innerType1, nil
	case "InnerType2":
		return this.innerType2, nil
	case "InnerTypes":
		r := make([]interface{}, 2)
		r[0] = this.innerType1
		r[1] = this.innerType2
		return r, nil
	}
	return nil, nil
}

func (this *MValType) Set(vm *golua.VM, key string, val interface{}) error {
	return nil
}
