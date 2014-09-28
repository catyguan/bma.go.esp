package golua

type JJTLuaParserState struct {
	nodes        *list
	marks        *list
	sp           int // number of nodes on stack
	mk           int // current mark
	node_created bool
	// openNode     Node
}

func newJJTLuaParserState() *JJTLuaParserState {
	r := new(JJTLuaParserState)
	r.nodes = newList()
	r.marks = newList()
	r.sp = 0
	r.mk = 0
	return r
}

func (this *JJTLuaParserState) nodeCreated() bool {
	return this.node_created
}

func (this *JJTLuaParserState) reset() {
	this.nodes.clear()
	this.marks.clear()
	this.sp = 0
	this.mk = 0
}

func (this *JJTLuaParserState) rootNode() Node {
	o := this.nodes.get(0)
	return o.(Node)
}

/* Pushes a node on to the stack. */
func (this *JJTLuaParserState) pushNode(n Node) {
	this.nodes.add(n)
	this.sp = this.sp + 1
}

/* Returns the node on the top of the stack, and remove it from the
   stack.  */
func (this *JJTLuaParserState) popNode() Node {
	this.sp = this.sp - 1
	if this.sp < this.mk {
		v := this.marks.remove(this.marks.size() - 1)
		if v != nil {
			this.mk = v.(int)
		}
	}
	v2 := this.nodes.remove(this.nodes.size() - 1)
	return v2.(Node)
}

/* Returns the node currently on the top of the stack. */
func (this *JJTLuaParserState) peekNode() Node {
	v := this.nodes.get(this.nodes.size() - 1)
	return v.(Node)
}

/* Returns the number of children on the stack in the current node
   scope. */
func (this *JJTLuaParserState) nodeArity() int {
	return this.sp - this.mk
}

func (this *JJTLuaParserState) clearNodeScope(n Node) {
	for this.sp > this.mk {
		this.popNode()
	}
	v := this.marks.remove(this.marks.size() - 1)
	this.mk = v.(int)
}

func (this *JJTLuaParserState) openNodeScope(n Node) {
	this.marks.add(this.mk)
	this.mk = this.sp
	n.jjtOpen()
}

/* A definite node is constructed from a specified number of
   children.  That number of nodes are popped from the stack and
   made the children of the definite node.  Then the definite node
   is pushed on to the stack. */
// func (this *JJTLuaParserState) closeNodeScopeI(n Node, num int) {
// 	v := this.marks.remove(this.marks.size() - 1)
// 	this.mk = v.(int)
// 	for num-1 > 0 {
// 		num = num - 1
// 		c := this.popNode()
// 		c.jjtSetParent(n)
// 		n.jjtAddChild(c, num)
// 	}
// 	n.jjtClose()
// 	this.pushNode(n)
// 	this.node_created = true
// }

/* A conditional node is constructed if its condition is true.  All
   the nodes that have been pushed since the node was opened are
   made children of the conditional node, which is then pushed
   on to the stack.  If the condition is false the node is not
   constructed and they are left on the stack. */
func (this *JJTLuaParserState) closeNodeScopeB(n Node, condition bool) {
	if condition {
		a := this.nodeArity()
		v := this.marks.remove(this.marks.size() - 1)
		this.mk = v.(int)
		// fmt.Println("closeNodeScope", n, a)
		for {
			if a <= 0 {
				break
			}
			a--
			c := this.popNode()
			c.jjtSetParent(n)
			n.jjtAddChild(c, a)
			// fmt.Println("closeNodeScope addChild", c, a)
		}
		n.jjtClose()
		this.pushNode(n)
		this.node_created = true
	} else {
		v := this.marks.remove(this.marks.size() - 1)
		this.mk = v.(int)
		this.node_created = false
	}
}
