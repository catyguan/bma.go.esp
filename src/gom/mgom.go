package gom

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
	"golua"
	gomcc "gom/goyacc"
	"logger"
	"strings"
)

type gomItem struct {
	value *MValue
	items map[string]*gomItem
	tmap  map[string]interface{}
	moo
}

type MGOM struct {
	items    map[string]*gomItem
	structs  []*MStruct
	services []*MService
	objects  []*MObject
	moo
}

func NewGOM() *MGOM {
	r := new(MGOM)
	r.items = make(map[string]*gomItem)
	r.structs = make([]*MStruct, 0)
	r.services = make([]*MService, 0)
	r.objects = make([]*MObject, 0)
	return r
}

func (this *MGOM) GetItem(name string) *gomItem {
	nlist := strings.Split(name, ".")
	var litem *gomItem
	h := this.items
	for _, n := range nlist {
		item, ok := h[n]
		if ok {
			litem = item
			h = item.items
		} else {
			return nil
		}
	}
	return litem
}

func (this *MGOM) GetStruct(name string) *MStruct {
	for _, o := range this.structs {
		if o.name == name {
			return o
		}
	}
	return nil
}

func (this *MGOM) GetService(name string) *MService {
	for _, o := range this.services {
		if o.name == name {
			return o
		}
	}
	return nil
}

func (this *MGOM) GetObject(name string) *MObject {
	for _, o := range this.objects {
		if o.name == name {
			return o
		}
	}
	return nil
}

func (this *MGOM) put(name string, v *MValue) {
	nlist := strings.Split(name, ".")
	var litem *gomItem
	h := this.items
	for _, n := range nlist {
		item, ok := h[n]
		if ok {
			litem = item
		} else {
			litem = new(gomItem)
			litem.items = make(map[string]*gomItem)
			h[n] = litem
		}
		h = litem.items
	}
	if litem != nil {
		litem.value = v
	}
}

func (this *MGOM) dumpItem(buf *bytes.Buffer, prex string, k string, item *gomItem) {
	if item.value != nil {
		if item.value.annos != nil {
			item.value.annos.Dump(buf, prex)
		}
	}
	buf.WriteString(prex)
	buf.WriteString(k)
	if item.value != nil {
		buf.WriteString(" = ")
		buf.WriteString(fmt.Sprintf("%v\n", item.value.value))
	} else {
		buf.WriteString("\n")
	}
	p2 := prex + "\t"
	for nk, nitem := range item.items {
		this.dumpItem(buf, p2, nk, nitem)
	}
}

func (this *MGOM) Dump(buf *bytes.Buffer, prex string) {
	for k, item := range this.items {
		this.dumpItem(buf, prex, k, item)
	}

	for _, o := range this.structs {
		o.Dump(buf, prex)
	}

	for _, o := range this.services {
		o.Dump(buf, prex)
	}

	for _, o := range this.objects {
		o.Dump(buf, prex)
	}
}

func (this *MGOM) toAnnos(node gomcc.Node) *MAnnotations {
	as := node.GetAnnotations()
	if len(as) == 0 {
		return nil
	}
	r := new(MAnnotations)
	r.list = make([]*MAnnotation, len(as))
	for i, a := range as {
		ma := new(MAnnotation)
		ma.Name = a.Name
		ma.Value = a.Value.(*gomcc.Node0).Value
		r.list[i] = ma
	}
	return r
}

func (this *MGOM) toValType(node gomcc.Node) *MValType {
	r := new(MValType)
	if node.GetOp() == gomcc.OP_NAME {
		r.name = node.(*gomcc.Node0).Value.(string)
		return r
	}

	n := node.(*gomcc.Node2)
	n1 := n.Child1.(*gomcc.Node0)
	r.name = n1.Value.(string)
	if n.Child2 != nil {
		nx := n.Child2.(*gomcc.Node2)
		r.innerType1 = this.toValType(nx.Child1)
		if nx.Child2 != nil {
			r.innerType2 = this.toValType(nx.Child2)
		}
	}
	return r
}

func (this *MGOM) build(node gomcc.Node) error {
	switch node.GetOp() {
	case gomcc.OP_GOM:
		n := node.(*gomcc.NodeN)
		for _, cn := range n.Childs {
			this.build(cn)
		}
	case gomcc.OP_VALUE:
		n := node.(*gomcc.Node2)
		n1 := n.Child1.(*gomcc.Node0)
		n2 := n.Child2.(*gomcc.Node0)
		v := new(MValue)
		v.annos = this.toAnnos(node)
		v.value = n2.Value
		name := n1.Value.(string)
		this.put(name, v)
	case gomcc.OP_STRUCT:
		m := new(MStruct)
		m.fields = make([]*MStructField, 0)
		m.annos = this.toAnnos(node)

		n := node.(*gomcc.Node2)
		n1 := n.Child1.(*gomcc.Node0)
		nx := n.Child2.(*gomcc.NodeN)
		m.name = n1.Value.(string)

		for _, cn := range nx.Childs {
			no := cn.(*gomcc.Node2)
			no1 := no.Child1.(*gomcc.Node0)

			mf := new(MStructField)
			mf.annos = this.toAnnos(cn)
			mf.name = no1.Value.(string)
			mf.valtype = this.toValType(no.Child2)
			m.fields = append(m.fields, mf)
		}
		this.structs = append(this.structs, m)
	case gomcc.OP_SERVICE:
		m := new(MService)
		m.methods = make([]*MServiceMethod, 0)
		m.annos = this.toAnnos(node)

		n := node.(*gomcc.Node2)
		n1 := n.Child1.(*gomcc.Node0)
		n2 := n.Child2.(*gomcc.NodeN)
		m.name = n1.Value.(string)

		for _, cn := range n2.Childs {
			smethod := cn.(*gomcc.Node3)

			mm := new(MServiceMethod)
			mm.params = make([]*MServiceMethodParam, 0)
			mm.annos = this.toAnnos(smethod)
			mm.name = smethod.Child1.(*gomcc.Node0).Value.(string)
			mm.returnType = this.toValType(smethod.Child3)
			sparams := smethod.Child2.(*gomcc.NodeN)
			for _, pn := range sparams.Childs {
				sparam := pn.(*gomcc.Node2)
				mp := new(MServiceMethodParam)
				mp.annos = this.toAnnos(sparam)
				mp.name = sparam.Child1.(*gomcc.Node0).Value.(string)
				mp.paramType = this.toValType(sparam.Child2)
				mm.params = append(mm.params, mp)
			}
			m.methods = append(m.methods, mm)
		}
		this.services = append(this.services, m)
	case gomcc.OP_OBJECT:
		m := new(MObject)
		m.fields = make([]*MStructField, 0)
		m.methods = make([]*MServiceMethod, 0)
		m.annos = this.toAnnos(node)

		n := node.(*gomcc.Node2)
		n1 := n.Child1.(*gomcc.Node0)
		n2 := n.Child2.(*gomcc.NodeN)
		m.name = n1.Value.(string)

		for _, cn := range n2.Childs {
			if cn.GetOp() == gomcc.OP_SMETHOD {
				smethod := cn.(*gomcc.Node3)

				mm := new(MServiceMethod)
				mm.params = make([]*MServiceMethodParam, 0)
				mm.annos = this.toAnnos(smethod)
				mm.name = smethod.Child1.(*gomcc.Node0).Value.(string)
				mm.returnType = this.toValType(smethod.Child3)
				sparams := smethod.Child2.(*gomcc.NodeN)
				for _, pn := range sparams.Childs {
					sparam := pn.(*gomcc.Node2)
					mp := new(MServiceMethodParam)
					mp.annos = this.toAnnos(sparam)
					mp.name = sparam.Child1.(*gomcc.Node0).Value.(string)
					mp.paramType = this.toValType(sparam.Child2)
					mm.params = append(mm.params, mp)
				}
				m.methods = append(m.methods, mm)
			} else {
				no := cn.(*gomcc.Node2)
				no1 := no.Child1.(*gomcc.Node0)

				mf := new(MStructField)
				mf.annos = this.toAnnos(cn)
				mf.name = no1.Value.(string)
				mf.valtype = this.toValType(no.Child2)
				m.fields = append(m.fields, mf)
			}
		}
		this.objects = append(this.objects, m)
	}
	return nil
}

func (this *MGOM) Compile(content, fname string) error {
	// bs, errr := ioutil.ReadFile(fname)
	// if errr != nil {
	// 	return nil, errr
	// }
	// // compile
	// content := string(bs)
	p := gomcc.NewParser(fname, content)
	node, err2 := p.Parse()
	if err2 != nil {
		err0 := fmt.Errorf("compile '%s' fail - %s", fname, err2)
		return err0
	}
	logger.Debug(tag, "compile('%s') done", fname)
	err3 := this.build(node)
	if err3 != nil {
		return err3
	}
	return nil
}

//// vmm
func (this *gomItem) ToVMTable() golua.VMTable {
	return this
}
func (this *gomItem) Get(vm *golua.VM, key string) (interface{}, error) {
	if this.value != nil {
		ok, v := gooCheckAnno(this.value.annos, key)
		if ok {
			return v, nil
		}
	}
	if this.items != nil {
		if item, ok := this.items[key]; ok {
			return item, nil
		}
	}
	switch key {
	case "Kind":
		return "Item", nil
	case "Value", "_Value":
		if this.value == nil {
			return nil, nil
		}
		return this.value.value, nil
	}
	return nil, nil
}

func (this *gomItem) Set(vm *golua.VM, key string, val interface{}) error {
	return nil
}

func (this *gomItem) ToMap() map[string]interface{} {
	if this.tmap != nil {
		return this.tmap
	}
	r := make(map[string]interface{})
	for k, o := range this.items {
		r[k] = o
	}
	this.tmap = r
	return r
}

func (this *MGOM) ToVMTable() golua.VMTable {
	return this
}
func (this *MGOM) funcGet(getval bool) interface{} {
	n := "GOM.GetItem"
	if getval {
		n = "GOM.GetValue"
	}
	return golua.NewGOF(n, func(vm *golua.VM, self interface{}) (int, error) {
		err0 := vm.API_checkStack(1)
		if err0 != nil {
			return 0, err0
		}
		n, err1 := vm.API_pop1X(-1, true)
		if err1 != nil {
			return 0, err1
		}
		vn := valutil.ToString(n, "")
		item := this.GetItem(vn)
		if item == nil {
			vm.API_push(nil)
		} else {
			if getval {
				if item.value != nil {
					vm.API_push(item.value.value)
				} else {
					vm.API_push(nil)
				}
			} else {
				vm.API_push(item.ToVMTable())
			}
		}
		return 1, nil
	})
}
func (this *MGOM) Get(vm *golua.VM, key string) (interface{}, error) {
	switch key {
	case "Kind":
		return "GOM", nil
	case "Items":
		r := make(map[string]interface{})
		for k, o := range this.items {
			r[k] = o
		}
		return r, nil
	case "GetItem":
		return this.funcGet(false), nil
	case "GetValue":
		return this.funcGet(true), nil
	case "GetStruct":
		return golua.NewGOF("GOM.GetStruct", func(vm *golua.VM, self interface{}) (int, error) {
			err0 := vm.API_checkStack(1)
			if err0 != nil {
				return 0, err0
			}
			n, err1 := vm.API_pop1X(-1, true)
			if err1 != nil {
				return 0, err1
			}
			vn := valutil.ToString(n, "")
			st := this.GetStruct(vn)
			vm.API_push(st)
			return 1, nil
		}), nil
	case "GetService":
		return golua.NewGOF("GOM.GetService", func(vm *golua.VM, self interface{}) (int, error) {
			err0 := vm.API_checkStack(1)
			if err0 != nil {
				return 0, err0
			}
			n, err1 := vm.API_pop1X(-1, true)
			if err1 != nil {
				return 0, err1
			}
			vn := valutil.ToString(n, "")
			st := this.GetService(vn)
			vm.API_push(st)
			return 1, nil
		}), nil
	case "GetObject":
		return golua.NewGOF("GOM.GetObject", func(vm *golua.VM, self interface{}) (int, error) {
			err0 := vm.API_checkStack(1)
			if err0 != nil {
				return 0, err0
			}
			n, err1 := vm.API_pop1X(-1, true)
			if err1 != nil {
				return 0, err1
			}
			vn := valutil.ToString(n, "")
			st := this.GetObject(vn)
			vm.API_push(st)
			return 1, nil
		}), nil
	}
	return nil, nil
}

func (this *MGOM) Set(vm *golua.VM, key string, val interface{}) error {
	// return this.p.Set(vm, this.o, key, val)
	return nil
}
