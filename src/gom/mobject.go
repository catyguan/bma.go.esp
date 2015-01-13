package gom

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
	"golua"
)

// MObject
type MObject struct {
	annos   *MAnnotations
	name    string
	fields  []*MStructField
	methods []*MServiceMethod
	moo
}

func (this *MObject) String() string {
	return fmt.Sprintf("object(%s)", this.name)
}

func (this *MObject) GetField(n string) *MStructField {
	for _, o := range this.fields {
		if o.name == n {
			return o
		}
	}
	return nil
}

func (this *MObject) GetMethod(n string) *MServiceMethod {
	for _, o := range this.methods {
		if o.name == n {
			return o
		}
	}
	return nil
}

func (this *MObject) Dump(buf *bytes.Buffer, prex string) {
	if this.annos != nil {
		this.annos.Dump(buf, prex)
	}
	buf.WriteString(prex)
	buf.WriteString(fmt.Sprintf("object %s {", this.name))
	i := 0
	for _, f := range this.fields {
		if i != 0 {
			buf.WriteString(",")
		}
		i++
		buf.WriteString("\n")
		p2 := prex + "\t"
		f.Dump(buf, p2)
	}
	for _, m := range this.methods {
		if i != 0 {
			buf.WriteString(",")
		}
		i++
		buf.WriteString("\n")
		p2 := prex + "\t"
		m.Dump(buf, p2)
	}
	buf.WriteString("\n}\n")
}

//// vmm
func (this *MObject) ToVMTable() golua.VMTable {
	return this
}
func (this *MObject) Get(vm *golua.VM, key string) (interface{}, error) {
	if ok, v := gooCheckAnno(this.annos, key); ok {
		return v, nil
	}
	switch key {
	case "Name":
		return this.name, nil
	case "GetField":
		return golua.NewGOF("Object.Field", func(vm *golua.VM, self interface{}) (int, error) {
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
	case "GetMethod":
		return golua.NewGOF("Object.Method", func(vm *golua.VM, self interface{}) (int, error) {
			err0 := vm.API_checkStack(1)
			if err0 != nil {
				return 0, err0
			}
			n, err1 := vm.API_pop1X(-1, true)
			if err1 != nil {
				return 0, err1
			}
			vn := valutil.ToString(n, "")
			v := this.GetMethod(vn)
			vm.API_push(v)
			return 1, nil
		}), nil
	}
	return nil, nil
}

func (this *MObject) Set(vm *golua.VM, key string, val interface{}) error {
	return nil
}
