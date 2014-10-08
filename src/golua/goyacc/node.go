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

type Node interface {
	/** This method returns a child node.  The children are numbered
	  from zero, left to right. */
	GetChild(i int) Node

	/** Return the number of children the node has. */
	GetNumChildren() int

	GetOp() OP

	String() string
}

type baseNode struct {
	op OP
}

func (this *baseNode) GetOp() OP {
	return this.op
}

func (this *baseNode) String() string {
	return OPNames[this.op]
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
			buf.WriteString(prefix + " <nil>")
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

// NodeLocal
type NodeLocal struct {
	baseNode
	Names   []string
	ExpList Node
}

func (this *NodeLocal) String() string {
	return fmt.Sprintf("local(%v)", this.Names)
}

func (this *NodeLocal) GetChild(i int) Node {
	switch i {
	case 0:
		return this.ExpList
	}
	return nil
}

func (this *NodeLocal) GetNumChildren() int {
	if this.ExpList == nil {
		return 0
	}
	return 1
}

// NodeFunc
type NodeFunc struct {
	baseNode
	Name   string
	Params []string
	CVars  []string
	Block  Node
}

func (this *NodeFunc) String() string {
	return fmt.Sprintf("func(%v, %v)", this.Name, this.Params)
}

func (this *NodeFunc) GetChild(i int) Node {
	switch i {
	case 0:
		return this.Block
	}
	return nil
}

func (this *NodeFunc) GetNumChildren() int {
	if this.Block == nil {
		return 0
	}
	return 1
}

// NodeFor
type NodeFor struct {
	baseNode
	Names  []string
	ForExp Node
}

func (this *NodeFor) String() string {
	return fmt.Sprintf("%s(%v)", OPNames[this.op], this.Names)
}

func (this *NodeFor) GetChild(i int) Node {
	switch i {
	case 0:
		return this.ForExp
	}
	return nil
}

func (this *NodeFor) GetNumChildren() int {
	if this.ForExp == nil {
		return 0
	}
	return 1
}

// NodeIf
type NodeIf struct {
	baseNode
	Exp       Node
	Block     Node
	ElseBlock Node
}

func (this *NodeIf) GetChild(i int) Node {
	switch i {
	case 0:
		return this.Exp
	case 1:
		return this.Block
	case 2:
		return this.ElseBlock
	}
	return nil
}

func (this *NodeIf) GetNumChildren() int {
	return 3
}
