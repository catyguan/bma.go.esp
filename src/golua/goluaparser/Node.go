package goluaparser

type Node interface {

	/** This method is called after the node has been made the current
	  node.  It indicates that child nodes can now be added to it. */
	jjtOpen()

	/** This method is called after all the child nodes have been
	  added. */
	jjtClose()

	/** This pair of methods are used to inform the node of its
	  parent. */
	jjtSetParent(n Node)
	jjtGetParent() Node

	/** This method tells the node to add its argument to the node's
	  list of children.  */
	jjtAddChild(n Node, i int)

	/** This method returns a child node.  The children are numbered
	  from zero, left to right. */
	jjtGetChild(i int) Node

	/** Return the number of children the node has. */
	jjtGetNumChildren() int

	getId() int

	dump(prefix string)
}

func DumpNode(n Node, p string) {
	n.dump(p)
}
