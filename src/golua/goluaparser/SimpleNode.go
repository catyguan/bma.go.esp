package goluaparser

import "fmt"

type SimpleNode struct {
	Parent   Node
	Children []Node
	Id       int
	Token    *Token
	line     int
	Data     interface{}
}

var (
	counter = 0
)

func NewSimpleNode(Id int) *SimpleNode {
	// if Id == 2 {
	// 	counter++
	// 	if counter > 10 {
	// 		panic("fuck")
	// 	}
	// }
	// fmt.Println("new node", Id)
	this := new(SimpleNode)
	this.Id = Id
	return this
}

func NewSimpleNodeT(tk *Token) *SimpleNode {
	r := NewSimpleNode(JJTTOKEN)
	r.Token = tk
	return r
}

func (this *SimpleNode) Line() int {
	if this.Token != nil {
		return this.Token.BeginLine
	}
	return this.line
}

func (this *SimpleNode) jjtOpen() {
}

func (this *SimpleNode) jjtClose() {
}

func (this *SimpleNode) jjtSetParent(n Node) {
	this.Parent = n
}

func (this *SimpleNode) jjtGetParent() Node {
	return this.Parent
}

func (this *SimpleNode) jjtAddChild(n Node, i int) {
	if this.Children == nil {
		this.Children = make([]Node, i+1)
	} else if i >= len(this.Children) {
		c := make([]Node, i+1)
		copy(c, this.Children)
		this.Children = c
	}
	this.Children[i] = n
}

func (this *SimpleNode) jjtGetChild(i int) Node {
	return this.Children[i]
}

func (this *SimpleNode) jjtGetNumChildren() int {
	if this.Children == nil {
		return 0
	}
	return len(this.Children)
}

/* You can overrIde these two methods in subclasses of SimpleNode to
   customize the way the node appears when the tree is dumped.  If
   your output uses more than one line you should overrIde
   toString(String), otherwise overrIding toString() is probably all
   you need to do. */

func (this *SimpleNode) String() string {
	s := JJT_NODE_NAME[this.Id]
	if this.Id == JJTTOKEN && this.Token != nil {
		ts := this.Token.String()
		s += "(" + ts + ")"
	}
	return s
}

func (this *SimpleNode) toString(prefix string) string {
	return prefix + this.String()
}

func (this *SimpleNode) getId() int {
	return this.Id
}

/* OverrIde this method if you want to customize how the node dumps
   out its Children. */
func (this *SimpleNode) dump(prefix string) {
	fmt.Println(this.toString(prefix))
	if this.Children != nil {
		for _, n := range this.Children {
			if sn, ok := n.(*SimpleNode); ok {
				sn.dump(prefix + " ")
			}
		}
	}
}
