package golua

import (
	"bmautil/valutil"
	"fmt"
	"strconv"
)

type LuaBuilder struct {
	root *SimpleNode
	name string
}

func NewLuaBuilder(node Node, n string) *LuaBuilder {
	r := new(LuaBuilder)
	r.root = node.(*SimpleNode)
	r.name = n
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
		if act != nil {
			fmt.Println("doMerge1", bn, "==>", act)
			return this.bindAction(bn, act)
		}
		return nil
	}
	var r Action
	var skip = -1
	for i, cn := range bn.children {
		if i <= skip {
			continue
		}
		cbn := cn.(*SimpleNode)
		act := cbn.action
		if op2, ok := act.(*Op2Action); ok {
			if op2.action1 == nil {
				n1 := bn.jjtGetChild(i - 1).(*SimpleNode)
				op2.action1 = n1.action
			}
			if op2.action2 == nil {
				if i+1 < bn.jjtGetNumChildren() {
					n2 := bn.jjtGetChild(i + 1).(*SimpleNode)
					op2.action2 = n2.action
					skip = i + 1
				}
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

func childToken(bn *SimpleNode) *Token {
	if bn.jjtGetNumChildren() == 1 {
		n := bn.jjtGetChild(0)
		cbn := n.(*SimpleNode)
		if cbn.id == JJTTOKEN {
			return cbn.token
		}
	}
	return nil
}

func (this *LuaBuilder) doCreate(bn *SimpleNode) (Action, error) {
	var r Action
	switch bn.id {
	case JJTASSIGN:
		r = newOp2Action(bn.id)
	case JJTBINOP:
		op2 := newOp2Action(bn.id)
		tk := childToken(bn)
		if tk != nil {
			op2.kind = tk.Kind
			return op2, nil
		}
	case JJTFIELDOP:
		op2 := newOp2Action(bn.id)
		tk := childToken(bn)
		if tk != nil && tk.Kind == NAME {
			op2.action2 = NewValueAction(tk.Image)
			return op2, nil
		}
	case JJTTOKEN:
		var err error
		r, err = this.token2action(bn.token)
		if err != nil {
			return nil, err
		}
	case JJTBLOCK:
		r = NewBlockAction(this.name)
	}
	err0 := this.doExpand(bn)
	if err0 != nil {
		return nil, err0
	}
	return r, nil
}

func (this *LuaBuilder) doBuild(bn *SimpleNode) error {
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
	case TRUE:
		return NewValueAction(true), nil
	case FALSE:
		return NewValueAction(true), nil
	case STRING, CHAR, CHARSTRING:
		return NewValueAction(tk.Image), nil
	case NUMBER:
		v32, err1 := strconv.ParseInt(tk.Image, 10, 32)
		if err1 == nil {
			return NewValueAction(int(v32)), nil
		}
		nerr := err1.(*strconv.NumError)
		if nerr.Err == strconv.ErrRange {
			v64, err2 := strconv.ParseInt(tk.Image, 10, 64)
			if err2 != nil {
				return nil, err2
			}
			return NewValueAction(v64), nil
		}
		if nerr.Err == strconv.ErrSyntax {
			f64, err3 := strconv.ParseFloat(tk.Image, 64)
			if err3 != nil {
				return nil, err3
			}
			return NewValueAction(f64), nil
		}
		return nil, err1
	default:
		v := valutil.ToBool(tk.Image, false)
		return NewValueAction(v), nil
	}
}
