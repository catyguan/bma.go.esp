package gom

import (
	"bytes"
	"fmt"
	gomcc "gom/goyacc"
	"io/ioutil"
	"logger"
	"strings"
)

type gomItem struct {
	value *MValue
	items map[string]*gomItem
}

type MGOM struct {
	items    map[string]*gomItem
	structs  []*MStruct
	services []*MService
}

func (this *MGOM) inito() {
	this.items = make(map[string]*gomItem)
	this.structs = make([]*MStruct, 0)
	this.services = make([]*MService, 0)
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
	}
	return nil
}

func Compile(fname string) (*MGOM, error) {
	bs, errr := ioutil.ReadFile(fname)
	if errr != nil {
		return nil, errr
	}
	// compile
	content := string(bs)
	p := gomcc.NewParser(fname, content)
	node, err2 := p.Parse()
	if err2 != nil {
		err0 := fmt.Errorf("compile '%s' fail - %s", fname, err2)
		return nil, err0
	}
	logger.Debug(tag, "compile('%s') done", fname)
	r := new(MGOM)
	r.inito()
	err3 := r.build(node)
	if err3 != nil {
		return nil, err3
	}
	return r, nil
}
