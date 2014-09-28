package golua

import "fmt"

type SimpleNode struct {
	parent   Node
	children []Node
	id       int
	token    *Token
	action   Action
}

var (
	counter = 0
)

func NewSimpleNode(id int) *SimpleNode {
	// if id == 2 {
	// 	counter++
	// 	if counter > 10 {
	// 		panic("fuck")
	// 	}
	// }
	// fmt.Println("new node", id)
	this := new(SimpleNode)
	this.id = id
	return this
}

func NewSimpleNodeT(tk *Token) *SimpleNode {
	r := NewSimpleNode(JJTTOKEN)
	r.token = tk
	return r
}

func (this *SimpleNode) jjtOpen() {
}

func (this *SimpleNode) jjtClose() {
}

func (this *SimpleNode) jjtSetParent(n Node) {
	this.parent = n
}

func (this *SimpleNode) jjtGetParent() Node {
	return this.parent
}

func (this *SimpleNode) jjtAddChild(n Node, i int) {
	if this.children == nil {
		this.children = make([]Node, i+1)
	} else if i >= len(this.children) {
		c := make([]Node, i+1)
		copy(c, this.children)
		this.children = c
	}
	this.children[i] = n
}

func (this *SimpleNode) jjtGetChild(i int) Node {
	return this.children[i]
}

func (this *SimpleNode) jjtGetNumChildren() int {
	if this.children == nil {
		return 0
	}
	return len(this.children)
}

/* You can override these two methods in subclasses of SimpleNode to
   customize the way the node appears when the tree is dumped.  If
   your output uses more than one line you should override
   toString(String), otherwise overriding toString() is probably all
   you need to do. */

func (this *SimpleNode) String() string {
	s := JJT_NODE_NAME[this.id]
	if this.id == JJTTOKEN && this.token != nil {
		ts := this.token.String()
		s += "(" + ts + ")"
	}
	return s
}

func (this *SimpleNode) toString(prefix string) string {
	return prefix + this.String()
}

func (this *SimpleNode) getId() int {
	return this.id
}

/* Override this method if you want to customize how the node dumps
   out its children. */
func (this *SimpleNode) dump(prefix string) {
	fmt.Println(this.toString(prefix))
	if this.children != nil {
		for _, n := range this.children {
			if sn, ok := n.(*SimpleNode); ok {
				sn.dump(prefix + " ")
			}
		}
	}
}
