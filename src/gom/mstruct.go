package gom

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
	"golua"
)

// MStructField
type MStructField struct {
	annos   *MAnnotations
	name    string
	valtype *MValType
	moo
}

func (this *MStructField) String() string {
	return fmt.Sprintf("%s:%s", this.name, this.valtype)
}

func (this *MStructField) Dump(buf *bytes.Buffer, prex string) {
	if this.annos != nil {
		this.annos.Dump(buf, prex)
	}
	buf.WriteString(prex)
	buf.WriteString(this.String())
}

// MStruct
type MStruct struct {
	annos  *MAnnotations
	name   string
	fields []*MStructField
	moo
}

func (this *MStruct) String() string {
	return fmt.Sprintf("struct(%s)", this.name)
}

func (this *MStruct) GetField(n string) *MStructField {
	if this.fields == nil {
		return nil
	}
	for _, f := range this.fields {
		if f.name == n {
			return f
		}
	}
	return nil
}

func (this *MStruct) Dump(buf *bytes.Buffer, prex string) {
	if this.annos != nil {
		this.annos.Dump(buf, prex)
	}
	buf.WriteString(prex)
	buf.WriteString(fmt.Sprintf("struct %s {", this.name))
	for i, f := range this.fields {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
		p2 := prex + "\t"
		f.Dump(buf, p2)
	}
	buf.WriteString("\n}\n")
}

//// vmm
func (this *MStructField) ToVMTable() golua.VMTable {
	return this
}
func (this *MStructField) Get(vm *golua.VM, key string) (interface{}, error) {
	if ok, v := gooCheckAnno(this.annos, key); ok {
		return v, nil
	}
	switch key {
	case "Name":
		return this.name, nil
	case "Kind":
		return "Field", nil
	case "Type":
		if this.valtype != nil {
			return this.valtype.String(), nil
		}
		return "", nil
	case "TypeObject":
		return this.valtype, nil
	}
	return nil, nil
}

func (this *MStructField) Set(vm *golua.VM, key string, val interface{}) error {
	return nil
}

func (this *MStruct) ToVMTable() golua.VMTable {
	return this
}
func (this *MStruct) Get(vm *golua.VM, key string) (interface{}, error) {
	if ok, v := gooCheckAnno(this.annos, key); ok {
		return v, nil
	}
	switch key {
	case "Name":
		return this.name, nil
	case "Kind":
		return "Struct", nil
	case "Fields":
		r := make([]interface{}, len(this.fields))
		for i, o := range this.fields {
			r[i] = o
		}
		return r, nil
	case "Field":
		return golua.NewGOF("Struct.Field", func(vm *golua.VM, self interface{}) (int, error) {
			err0 := vm.API_checkStack(1)
			if err0 != nil {
				return 0, err0
			}
			n, err1 := vm.API_pop1X(-1, true)
			if err1 != nil {
				return 0, err1
			}
			vn := valutil.ToString(n, "")
			v := this.GetField(vn)
			vm.API_push(v)
			return 1, nil
		}), nil
	}
	return nil, nil
}

func (this *MStruct) Set(vm *golua.VM, key string, val interface{}) error {
	return nil
}
