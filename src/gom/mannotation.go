package gom

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
	"golua"
)

type MAnnotation struct {
	Name  string
	Value interface{}
}

func (this *MAnnotation) String() string {
	return fmt.Sprintf("%s=%v", this.Name, this.Value)
}

type MAnnotations struct {
	list []*MAnnotation
	moo
}

func (this *MAnnotations) Has(n string) bool {
	for _, a := range this.list {
		if a.Name == n {
			return true
		}
	}
	return false
}

func (this *MAnnotations) GetAnnotation(n string) interface{} {
	for _, a := range this.list {
		if a.Name == n {
			return a.Value
		}
	}
	return nil
}

func (this *MAnnotations) All() []*MAnnotation {
	return this.list
}

func (this *MAnnotations) List(n string) []interface{} {
	r := make([]interface{}, 0)
	for _, a := range this.list {
		if a.Name == n {
			r = append(r, a.Value)
		}
	}
	return r
}

func (this *MAnnotations) Dump(buf *bytes.Buffer, prex string) {
	for _, a := range this.list {
		buf.WriteString(prex)
		buf.WriteString(fmt.Sprintf("@%s\n", a))
	}
}

//// vmm
func (this *MAnnotations) ToVMTable() golua.VMTable {
	return this
}
func (this *MAnnotations) funcGet() interface{} {
	return golua.NewGOF("MAnnotations.GetAnnotation", func(vm *golua.VM, self interface{}) (int, error) {
		err0 := vm.API_checkStack(1)
		if err0 != nil {
			return 0, err0
		}
		n, err1 := vm.API_pop1X(-1, true)
		if err1 != nil {
			return 0, err1
		}
		vn := valutil.ToString(n, "")
		v := this.GetAnnotation(vn)
		vm.API_push(v)
		return 1, nil
	})
}
func (this *MAnnotations) Get(vm *golua.VM, key string) (interface{}, error) {
	switch key {
	case "All":
		r := make([]interface{}, len(this.list))
		for i, a := range this.list {
			r[i] = map[string]interface{}{"name": a.Name, "value": a.Value}
		}
		return r, nil
	case "Get":
		return this.funcGet(), nil
	case "List":
		return golua.NewGOF("MAnnotations.List", func(vm *golua.VM, self interface{}) (int, error) {
			err0 := vm.API_checkStack(1)
			if err0 != nil {
				return 0, err0
			}
			n, err1 := vm.API_pop1X(-1, true)
			if err1 != nil {
				return 0, err1
			}
			vn := valutil.ToString(n, "")
			v := this.List(vn)
			vm.API_push(v)
			return 1, nil
		}), nil
	}
	return nil, nil
}

func (this *MAnnotations) Set(vm *golua.VM, key string, val interface{}) error {
	// return this.p.Set(vm, this.o, key, val)
	return nil
}

func gooCheckAnno(anns *MAnnotations, key string) (bool, interface{}) {
	switch key {
	case "Annotation":
		if anns == nil {
			return true, nil
		}
		return true, anns.funcGet()
	case "Annotations":
		return true, anns
	}
	return false, nil
}
