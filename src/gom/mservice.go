package gom

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
	"golua"
)

// MServiceMethodParam
type MServiceMethodParam struct {
	annos     *MAnnotations
	name      string
	paramType *MValType
}

func (this *MServiceMethodParam) String() string {
	return fmt.Sprintf("%s:%s", this.name, this.paramType)
}

// MServiceMethod
type MServiceMethod struct {
	annos      *MAnnotations
	name       string
	params     []*MServiceMethodParam
	returnType *MValType
	moo
}

func (this *MServiceMethod) String() string {
	buf := bytes.NewBuffer(make([]byte, 0, 16))
	buf.WriteString(this.name)
	buf.WriteString("(")
	for i, p := range this.params {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(p.String())
	}
	buf.WriteString("):")
	buf.WriteString(this.returnType.String())
	return buf.String()
}

func (this *MServiceMethod) Dump(buf *bytes.Buffer, prex string) {
	if this.annos != nil {
		this.annos.Dump(buf, prex)
	}
	bp := false
	for _, p := range this.params {
		if p.annos != nil && len(p.annos.list) > 0 {
			bp = true
			break
		}
	}
	if !bp {
		buf.WriteString(prex)
		buf.WriteString(this.String())
		return
	}

	buf.WriteString(prex)
	buf.WriteString(this.name)
	buf.WriteString("(")
	for i, p := range this.params {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
		p2 := prex + "\t"
		p.annos.Dump(buf, p2)
		buf.WriteString(p2)
		buf.WriteString(p.String())
	}
	buf.WriteString("\n")
	buf.WriteString(prex)
	buf.WriteString("):")
	buf.WriteString(this.returnType.String())
}

// MService
type MService struct {
	annos   *MAnnotations
	name    string
	methods []*MServiceMethod
	moo
}

func (this *MService) String() string {
	return fmt.Sprintf("service(%s)", this.name)
}

func (this *MService) GetMethod(n string) *MServiceMethod {
	for _, o := range this.methods {
		if o.name == n {
			return o
		}
	}
	return nil
}

func (this *MService) Dump(buf *bytes.Buffer, prex string) {
	if this.annos != nil {
		this.annos.Dump(buf, prex)
	}
	buf.WriteString(prex)
	buf.WriteString(fmt.Sprintf("service %s {", this.name))
	for i, m := range this.methods {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
		p2 := prex + "\t"
		m.Dump(buf, p2)
	}
	buf.WriteString("\n}\n")
}

//// vmm
func (this *MService) ToVMTable() golua.VMTable {
	return this
}
func (this *MService) Get(vm *golua.VM, key string) (interface{}, error) {
	if ok, v := gooCheckAnno(this.annos, key); ok {
		return v, nil
	}
	switch key {
	case "Name":
		return this.name, nil
	case "GetMethod":
		return golua.NewGOF("Service.Method", func(vm *golua.VM, self interface{}) (int, error) {
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

func (this *MService) Set(vm *golua.VM, key string, val interface{}) error {
	return nil
}
