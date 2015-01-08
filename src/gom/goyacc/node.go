package goyacc

import (
	"bytes"
	"fmt"
)

type yyToken struct {
	kind   int
	line   int
	column int

	image string
}

type yySymType struct {
	yys   int
	token yyToken
	op    OP
	value interface{}
}

func (this *yySymType) Be(val interface{}) {
	this.op = OP_NONE
	this.value = val
}

type Node interface {
	/** This method returns a child node.  The children are numbered
	  from zero, left to right. */
	GetChild(i int) Node

	/** Return the number of children the node has. */
	GetNumChildren() int

	GetOp() OP

	GetLine() int

	String() string
}

type baseNode struct {
	op       OP
	Line     int
	AnnoList Annotations
}

func (this *baseNode) Be(op OP, line int) {
	this.op = op
	this.Line = line
}

func (this *baseNode) Bev(op OP, v *yySymType) {
	this.op = op
	if v != nil {
		this.Line = v.token.line
	}
}

func (this *baseNode) Bev2(op OP, v1 *yySymType, v2 *yySymType) {
	this.op = op
	if v1 != nil {
		this.Line = v1.token.line
	} else if v2 != nil {
		this.Line = v2.token.line
	}
}

func (this *baseNode) GetLine() int {
	return this.Line
}

func (this *baseNode) GetOp() OP {
	return this.op
}

func (this *baseNode) String() string {
	s := OPNames[this.op]
	if len(this.AnnoList) > 0 {
		s += "@{\n"
		for _, a := range this.AnnoList {
			s += a.String()
		}
		s += "}"
	}
	return s
}

func dumpNode(buf *bytes.Buffer, prefix string, n Node) {
	buf.WriteString(prefix)
	if n == nil {
		buf.WriteString("<nil node>\n")
		return
	}
	buf.WriteString(n.String())
	buf.WriteString("\n")
	c := n.GetNumChildren()
	for i := 0; i < c; i++ {
		cn := n.GetChild(i)
		if cn != nil {
			dumpNode(buf, prefix+" ", cn)
		} else {
			buf.WriteString(prefix + " <nil>\n")
		}
	}
}

func DumpNode(prefix string, n Node) string {
	buf := bytes.NewBuffer(make([]byte, 0, 16))
	dumpNode(buf, prefix, n)
	return buf.String()
}

// Node0
type Node0 struct {
	baseNode
	Value interface{}
}

func (this *Node0) GetChild(i int) Node {
	return nil
}

func (this *Node0) GetNumChildren() int {
	return 0
}

func (this *Node0) String() string {
	if this.Value != nil {
		return fmt.Sprintf("%s(%v,%T)", this.baseNode.String(), this.Value, this.Value)
	}
	return this.baseNode.String()
}

// Node1
type Node1 struct {
	baseNode
	Child Node
}

func (this *Node1) GetChild(i int) Node {
	switch i {
	case 0:
		return this.Child
	}
	return nil
}

func (this *Node1) GetNumChildren() int {
	if this.Child == nil {
		return 0
	}
	return 1
}

// Node2
type Node2 struct {
	baseNode
	Child1 Node
	Child2 Node
}

func (this *Node2) GetChild(i int) Node {
	switch i {
	case 0:
		return this.Child1
	case 1:
		return this.Child2
	}
	return nil
}

func (this *Node2) GetNumChildren() int {
	return 2
}

// Node3
type Node3 struct {
	baseNode
	Child1 Node
	Child2 Node
	Child3 Node
}

func (this *Node3) GetChild(i int) Node {
	switch i {
	case 0:
		return this.Child1
	case 1:
		return this.Child2
	case 2:
		return this.Child3
	}
	return nil
}

func (this *Node3) GetNumChildren() int {
	return 3
}

// NodeN
type NodeN struct {
	baseNode
	Childs []Node
}

func (this *NodeN) GetChild(i int) Node {
	if i < len(this.Childs) {
		return this.Childs[i]
	}
	return nil
}

func (this *NodeN) GetNumChildren() int {
	return len(this.Childs)
}

func (this *NodeN) AddChild(n Node) {
	if this.Childs == nil {
		this.Childs = []Node{n}
	} else {
		this.Childs = append(this.Childs, n)
	}
}
