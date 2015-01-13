package goyacc

import (
	"fmt"
	"strconv"
	"strings"
)

type OP uint8

const (
	OP_NONE         = OP(0)
	OP_VALUE        = OP(1)
	OP_NAME         = OP(2)
	OP_ANNO         = OP(3)
	OP_GOM          = OP(4)
	OP_FIELD        = OP(5)
	OP_TYPE         = OP(6)
	OP_STRUCT       = OP(7)
	OP_STRUCT_BODY  = OP(8)
	OP_SFIELD       = OP(9)
	OP_SERVICE      = OP(10)
	OP_SERVICE_BODY = OP(11)
	OP_SMETHOD      = OP(12)
	OP_SM_PARAMS    = OP(13)
	OP_SM_PARAM     = OP(14)
	OP_OBJECT       = OP(15)
	OP_OBJECT_BODY  = OP(16)
)

var OPNames = []string{
	"NONE", "VALUE",
	"NAME", "ANNO",
	"GOM",
	"FIELD", "TYPE",
	"STRUCT", "STRUCT_BODY", "SFIELD",
	"SERVICE", "SERVICE_BODY", "SMETHOD", "SM_PARAMS", "SM_PARAM",
	"OBJECT", "OBJECT_BODY",
}

func beNode(yylex yyLexer, lval *yySymType, val1 *yySymType) {
	n1, err1 := toNode(yylex, val1)
	if err1 != nil {
		yylex.Error(err1.Error())
		return
	}
	lval.Be(n1)
}

func toNode(yylex yyLexer, val *yySymType) (Node, error) {
	if val == nil {
		return nil, nil
	}
	if val.op == OP_VALUE {
		r := new(Node0)
		r.Bev(OP_VALUE, val)
		r.Value = val.value
		return r, nil
	}
	if val.value != nil {
		if n, ok := val.value.(Node); ok {
			return n, nil
		}
		if n, ok := val.value.(string); ok {
			r := new(Node0)
			r.Bev(OP_NAME, val)
			r.Value = n
			return r, nil
		}
	}
	if val.token.kind == NAME {
		r := new(Node0)
		r.Bev(OP_NAME, val)
		r.Value = val.token.image
		return r, nil
	}
	return nil, nil
	// return nil, fmt.Errorf("unknow node(%d)", val.token.kind)
}

func opN(yylex yyLexer, lval *yySymType, op OP, v1 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("opN >> ", op, v1)
	}
	var nlist []Node
	if v1 != nil {
		nlist = v1.value.([]Node)
	}
	r := new(NodeN)
	r.Bev(op, v1)
	r.Childs = nlist
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("opN end: ", nlist)
	}
}

func op2(yylex yyLexer, lval *yySymType, op OP, v1 *yySymType, v2 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("op2 >> ", op, v1, v2)
	}
	n1, err1 := toNode(yylex, v1)
	if err1 != nil {
		yylex.Error(err1.Error())
		return
	}
	n2, err2 := toNode(yylex, v2)
	if err2 != nil {
		yylex.Error(err2.Error())
		return
	}
	r := new(Node2)
	r.Bev2(op, v1, v2)
	r.Child1 = n1
	r.Child2 = n2
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("op2 end: ", r, n1, n2)
	}
}

func op3(yylex yyLexer, lval *yySymType, op OP, v1 *yySymType, v2 *yySymType, v3 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("op3 >> ", op, v1, v2, v3)
	}
	n1, err1 := toNode(yylex, v1)
	if err1 != nil {
		yylex.Error(err1.Error())
		return
	}
	n2, err2 := toNode(yylex, v2)
	if err2 != nil {
		yylex.Error(err2.Error())
		return
	}
	n3, err3 := toNode(yylex, v3)
	if err3 != nil {
		yylex.Error(err3.Error())
		return
	}
	r := new(Node3)
	r.Bev2(op, v1, v2)
	r.Child1 = n1
	r.Child2 = n2
	r.Child3 = n3
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("op3 end: ", r, n1, n2, n3)
	}
}

func nodeAppend(yylex yyLexer, lval *yySymType, val1 *yySymType, val2 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("nodeAppend", val1, val2)
	}
	r := make([]Node, 0, 2)
	if val1 != nil {
		n, ok := val1.value.(Node)
		if ok {
			r = append(r, n)
		} else {
			l := val1.value.([]Node)
			for _, n := range l {
				r = append(r, n)
			}
		}
	}
	if val2 != nil {
		n, err := toNode(yylex, val2)
		if err != nil {
			yylex.Error(err.Error())
			return
		}
		r = append(r, n)
	}
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("nodeAppend end: ", r)
	}
}

/////////////////////////////////////////////////////////
func endGOM(yylex yyLexer, lval *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("endGOM")
	}
}

func nameAppend(yylex yyLexer, lval *yySymType, val1 *yySymType, val2 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("nameAppend", val1, val2)
	}
	ns := make([]string, 0, 2)
	if val1 != nil {
		if val1.value == nil {
			ns = append(ns, val1.token.image)
		} else {
			ns = append(ns, val1.value.(string))
		}
	}
	if val2 != nil {
		if val2.value == nil {
			ns = append(ns, val2.token.image)
		} else {
			ns = append(ns, val2.value.(string))
		}
	}
	str := strings.Join(ns, ".")
	lval.Be(str)
	if yyDebug >= 2 {
		fmt.Println("nameAppend end: ", str)
	}
}

func annoAppend(yylex yyLexer, lval *yySymType, val1 *yySymType, val2 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("annoAppend", val1, val2)
	}
	r := make([]*Annotation, 0, 2)
	if val1 != nil {
		a, ok := val1.value.(*Annotation)
		if ok {
			r = append(r, a)
		} else {
			l := val1.value.([]*Annotation)
			for _, a := range l {
				r = append(r, a)
			}
		}
	}
	if val2 != nil {
		a, ok := val2.value.(*Annotation)
		if ok {
			r = append(r, a)
		} else {
			l := val2.value.([]*Annotation)
			for _, a := range l {
				r = append(r, a)
			}
		}
	}
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("annoAppend end: ", r)
	}
}

func defineAnnotation(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	n1, err1 := toNode(yylex, v1)
	if err1 != nil {
		yylex.Error(err1.Error())
		return
	}
	n2, err2 := toNode(yylex, v2)
	if err2 != nil {
		yylex.Error(err2.Error())
		return
	}
	a := new(Annotation)
	a.Name = n1.(*Node0).Value.(string)
	a.Value = n2
	lval.Be(a)
	if yyDebug >= 2 {
		fmt.Println("defineAnnotaion end: ", a)
	}
}

func commitValue(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	p := yylex.(*Parser)
	r := v1.value.(*Node2)
	if v2 != nil {
		r.AnnoList = Annotations(v2.value.([]*Annotation))
	}
	p.nodes = append(p.nodes, r)
	// lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("commitValue end: ", r, r.Child1, r.Child2)
	}
}

func commitStruct(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	p := yylex.(*Parser)
	r := v1.value.(*Node2)
	if v2 != nil {
		r.AnnoList = Annotations(v2.value.([]*Annotation))
	}
	p.nodes = append(p.nodes, r)
	if yyDebug >= 2 {
		fmt.Println("commitStruct end", r)
	}
}

func commitStructField(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	r := v1.value.(*Node2)
	if v2 != nil {
		r.AnnoList = Annotations(v2.value.([]*Annotation))
	}
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("commitStructField end")
	}
}

func commitService(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	p := yylex.(*Parser)
	r := v1.value.(*Node2)
	if v2 != nil {
		r.AnnoList = Annotations(v2.value.([]*Annotation))
	}
	p.nodes = append(p.nodes, r)
	if yyDebug >= 2 {
		fmt.Println("commitService end", r)
	}
}

func commitServiceMethod(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	r := v1.value.(*Node3)
	if v2 != nil {
		r.AnnoList = Annotations(v2.value.([]*Annotation))
	}
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("commitServiceMethod end")
	}
}

func commitMethodParam(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	r := v1.value.(*Node2)
	if v2 != nil {
		r.AnnoList = Annotations(v2.value.([]*Annotation))
	}
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("commitMethodParam end")
	}
}

func commitObject(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	p := yylex.(*Parser)
	r := v1.value.(*Node2)
	if v2 != nil {
		r.AnnoList = Annotations(v2.value.([]*Annotation))
	}
	p.nodes = append(p.nodes, r)
	if yyDebug >= 2 {
		fmt.Println("commitObject end", r)
	}
}

func opValue(yylex yyLexer, lval *yySymType) {
	err := func() error {
		lval.op = OP_VALUE
		switch lval.token.kind {
		case NIL:
			lval.value = nil
		case TRUE:
			lval.value = true
		case FALSE:
			lval.value = false
		case STRING:
			lval.value = lval.token.image
		case NUMBER:
			s := lval.token.image
			if !strings.Contains(s, ".") {
				v32, err1 := strconv.ParseInt(s, 10, 32)
				if err1 == nil {
					lval.value = int32(v32)
					break
				}
				nerr := err1.(*strconv.NumError)
				if nerr.Err == strconv.ErrRange {
					v64, err2 := strconv.ParseInt(s, 10, 64)
					if err2 != nil {
						return err2
					}
					lval.value = v64
					break
				}
				if nerr.Err != strconv.ErrSyntax {
					return err1
				}
			}
			f64, err3 := strconv.ParseFloat(s, 64)
			if err3 != nil {
				return err3
			}
			lval.value = f64
			break
		}
		return nil
	}()
	if err != nil {
		yylex.Error(err.Error())
	}
	if yyDebug >= 2 {
		fmt.Println("opValue end: ", lval)
	}
}

func defineArray(yylex yyLexer, lval *yySymType, v1 *yySymType) {
	v := make([]interface{}, 0)
	if v1 != nil {
		if n, ok := v1.value.(Node); ok {
			v = append(v, n)
		} else {
			nlist := v1.value.([]Node)
			for _, n := range nlist {
				v = append(v, n)
			}
		}
	}
	r := new(Node0)
	r.Bev(OP_VALUE, v1)
	r.Value = v
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("defineArray end: ", r)
	}
}

func defineTable(yylex yyLexer, lval *yySymType, v1 *yySymType) {
	v := make(map[string]interface{})
	if v1 != nil {
		nlist := v1.value.([]Node)
		for _, n := range nlist {
			n2 := n.(*Node2)
			nf1 := n2.Child1.(*Node0)
			nf2 := n2.Child2.(*Node0)
			v[nf1.Value.(string)] = nf2.Value
		}
	}
	r := new(Node0)
	r.Bev(OP_VALUE, v1)
	r.Value = v
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("defineTable end: ", r)
	}
}

func defineField(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("defineField >> ", v1, v2)
	}
	n1 := new(Node0)
	n1.Bev(OP_NAME, v1)
	n1.Value = v1.token.image

	n2, err2 := toNode(yylex, v2)
	if err2 != nil {
		yylex.Error(err2.Error())
		return
	}
	r := new(Node2)
	r.Bev2(OP_FIELD, v1, v2)
	r.Child1 = n1
	r.Child2 = n2
	lval.Be(r)
	if yyDebug >= 2 {
		fmt.Println("defineField end: ", r, n1, n2)
	}
}

////////////////
func walk(node Node, f func(n Node) bool) bool {
	if node == nil {
		return true
	}
	c := node.GetNumChildren()
	for i := 0; i < c; i++ {
		cn := node.GetChild(i)
		if cn != nil {
			if !f(cn) {
				return false
			}
		}
	}
	return true
}

func execOptimize(node Node) bool {
	if node == nil {
		return true
	}
	walk(node, execOptimize)
	return true
}
