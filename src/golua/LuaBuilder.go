package golua

import (
	"bmautil/valutil"
	"fmt"
)

type LuaBuilder struct {
	root *SimpleNode
}

func NewLuaBuilder(node Node) *LuaBuilder {
	r := new(LuaBuilder)
	r.root = node.(*SimpleNode)
	return r
}

func (this *LuaBuilder) Build() (Action, error) {
	err := this.doBuild(this.root)
	if err != nil {
		return nil, err
	}
	return this.root.action, nil
}

func (this *LuaBuilder) doNode(p *SimpleNode, node Node) error {
	bn := node.(*SimpleNode)
	return this.doBuild(bn)
}

func (this *LuaBuilder) doExpand(bn *SimpleNode) error {
	for _, n := range bn.children {
		err := this.doNode(bn, n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *LuaBuilder) bindAction(bn *SimpleNode, act Action) error {
	if bn.action == nil {
		bn.action = act
	} else if op2, ok := bn.action.(*Op2Action); ok {
		op2.action2 = act
	} else {
		return fmt.Errorf("unknow bindAction(%s, %s)", bn, act)
	}
	return nil
}

func (this *LuaBuilder) doMerge(bn *SimpleNode) error {
	fmt.Println("doMerge", bn, "start")
	c := len(bn.children)
	if c == 0 {
		fmt.Println("doMerge0", bn)
		return nil
	}
	if c == 1 {
		n := bn.children[0]
		act := n.(*SimpleNode).action
		fmt.Println("doMerge1", bn, "==>", act)
		err := this.bindAction(bn, act)
		return err
	}
	var r Action
	for i, cn := range bn.children {
		cbn := cn.(*SimpleNode)
		act := cbn.action
		if op2, ok := act.(*Op2Action); ok {
			if op2.canMerge() {
				n1 := bn.jjtGetChild(i - 1).(*SimpleNode)
				op2.action1 = n1.action
			}
		}
		if act != nil {
			fmt.Println("doMergeN", bn, i, act)
			r = act
		}
	}
	fmt.Println("doMergeN", bn, "==>", r)
	return this.bindAction(bn, r)
}

func (this *LuaBuilder) doOp2(bn *SimpleNode, idx int) (Action, error) {
	fmt.Println("doOp2", bn, idx)
	p := bn.parent
	n1 := p.(*SimpleNode).jjtGetChild(idx - 1)
	a1 := n1.(*SimpleNode).action
	a2 := bn.action
	r := NewOp2Action(bn.id, a1, a2)
	bn.action = r
	return r, nil
}

func (this *LuaBuilder) doCreate(bn *SimpleNode) (Action, error) {
	switch bn.id {
	case JJTASSIGN, JJTFIELDOP:
		return newOp2Action(bn.id), nil
	case JJTTOKEN:
		return this.token2action(bn.token)
	case JJTBLOCK:
		return NewBlockAction(), nil
	}
	return nil, nil
}

func (this *LuaBuilder) doBuild(bn *SimpleNode) error {
	err0 := this.doExpand(bn)
	if err0 != nil {
		return err0
	}
	act, err1 := this.doCreate(bn)
	if err1 != nil {
		return err1
	}
	bn.action = act

	switch bn.id {
	case JJTBLOCK:
		p := act.(*BlockAction)
		for _, cbn := range bn.children {
			act := cbn.(*SimpleNode).action
			if act != nil {
				p.actions = append(p.actions, act)
			}
		}
	default:
		return this.doMerge(bn)
	}
	return nil
}

func (this *LuaBuilder) token2action(tk *Token) (Action, error) {
	switch tk.Kind {
	case NAME:
		return NewVarAction(false, tk), nil
	case NUMBER:
		v := valutil.ToInt(tk.Image, 0)
		return NewValueAction(v), nil
	}
	return nil, nil
}
