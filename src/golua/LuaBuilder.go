package golua

import (
	"bmautil/valutil"
	"fmt"
	"golua/goluaparser"
	"strconv"
)

type LuaBuilder struct {
	Tracable bool
	name     string
}

type LuaDoBuild interface {
	DoBuild(b *LuaBuilder) error
}

type LuaDoMerge interface {
	DoMerge(b *LuaBuilder, bn *goluaparser.SimpleNode) error
}

type LuaDoCombine interface {
	DoCombine(b *LuaBuilder, bn *goluaparser.SimpleNode, plist []goluaparser.Node, pos int) error
}

func NewLuaBuilder() *LuaBuilder {
	r := new(LuaBuilder)
	return r
}

func (this *LuaBuilder) Tracef(format string, args ...interface{}) {
	if this.Tracable {
		s := fmt.Sprintf(format, args...)
		fmt.Println(s)
	}
}

func (this *LuaBuilder) Build(node goluaparser.Node, name string) (Action, error) {
	this.name = name
	bn := node.(*goluaparser.SimpleNode)
	err := this.DoBuild(bn)
	if err != nil {
		return nil, err
	}
	if bn.Data != nil {
		return bn.Data.(Action), nil
	}
	return nil, nil
}

func (this *LuaBuilder) DoNode(node goluaparser.Node) error {
	bn := node.(*goluaparser.SimpleNode)
	return this.DoBuild(bn)
}

func (this *LuaBuilder) DoBuild(bn *goluaparser.SimpleNode) error {
	act, err1 := this.DoCreate(bn)
	if err1 != nil {
		return err1
	}
	bn.Data = act
	if o, ok := act.(LuaDoBuild); ok {
		return o.DoBuild(this)
	}

	switch bn.Id {
	case goluaparser.JJTNAMELIST:
		return nil
	}
	err2 := this.DoChildren(bn)
	if err2 != nil {
		return err2
	}
	return this.DoMerge(bn)
}

func (this *LuaBuilder) DoCreate(bn *goluaparser.SimpleNode) (Action, error) {
	var r Action
	switch bn.Id {
	case goluaparser.JJTASSIGN:
		r = newOp2Action(bn.Id)
	case goluaparser.JJTTOKEN:
		var err error
		r, err = this.token2action(bn.Token)
		if err != nil {
			return nil, err
		}
	case goluaparser.JJTBLOCK:
		r = newBlockAction()
	case goluaparser.JJTCHUNK:
		r = newChunkAction(this.name)
	case goluaparser.JJTFUNCOP:
		act := new(CallAction)
		act.line = bn.Line()
		r = act
	}
	return r, nil
}

func (this *LuaBuilder) DoChildren(bn *goluaparser.SimpleNode) error {
	for _, n := range bn.Children {
		err := this.DoNode(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *LuaBuilder) DoMerge(bn *goluaparser.SimpleNode) error {
	this.Tracef("doMerge %s:%v start", bn, bn.Data)
	if bn.Data != nil {
		act := bn.Data.(Action)
		if o, ok := act.(LuaDoMerge); ok {
			return o.DoMerge(this, bn)
		}
	}

	c := len(bn.Children)
	if c == 0 {
		this.Tracef("skip %s", bn)
		return nil
	}
	if c == 1 {
		cbn := bn.Children[0].(*goluaparser.SimpleNode)
		if cbn.Data != nil {
			act := cbn.Data.(Action)
			return this.BindAction(bn, act)
		}
		return nil
	}
	alist := this.ChildrenAsAction(bn)
	if bn.Data == nil && len(alist) > 1 {
		block := newBlockAction()
		block.actions = alist
		bn.Data = block
		return nil
	} else {
		r := Action(nil)
		if len(alist) > 0 {
			r = alist[len(alist)-1]
			this.Tracef("doMerge %s N ==> %s", bn, r)
			return this.BindAction(bn, r)
		}
		this.Tracef("skip %s", bn)
		return nil
	}
}

func (this *LuaBuilder) BindAction(bn *goluaparser.SimpleNode, act Action) error {
	if bn.Data == nil {
		this.Tracef("doBind %s ==> %s", bn, act)
		bn.Data = act
	}
	return nil
}

func (this *LuaBuilder) ChildrenAsAction(bn *goluaparser.SimpleNode) []Action {
	r := make([]Action, 0, len(bn.Children))
	for i, cn := range bn.Children {
		cbn := cn.(*goluaparser.SimpleNode)
		if cbn.Data != nil {
			act := cbn.Data.(Action)
			if o, ok := act.(LuaDoCombine); ok {
				o.DoCombine(this, bn, bn.Children, i)
			}
		}
	}
	for _, cn := range bn.Children {
		cbn := cn.(*goluaparser.SimpleNode)
		if cbn.Data != nil {
			act := cbn.Data.(Action)
			r = append(r, act)
		}
	}
	return r
}

func (this *LuaBuilder) ChildToken(bn *goluaparser.SimpleNode) *goluaparser.Token {
	if len(bn.Children) == 1 {
		n := bn.Children[0]
		cbn := n.(*goluaparser.SimpleNode)
		if cbn.Id == goluaparser.JJTTOKEN {
			return cbn.Token
		}
	}
	return nil
}

func (this *LuaBuilder) token2action(tk *goluaparser.Token) (Action, error) {
	switch tk.Kind {
	case goluaparser.LOCAL:
		r := newLocalAction()
		r.line = tk.BeginLine
		return r, nil
	case goluaparser.NAME:
		return newVarAction(false, tk), nil
	case goluaparser.TRUE:
		return newValueAction(true), nil
	case goluaparser.FALSE:
		return newValueAction(true), nil
	case goluaparser.STRING, goluaparser.CHAR, goluaparser.CHARSTRING:
		return newValueAction(tk.Image), nil
	case goluaparser.NUMBER:
		v32, err1 := strconv.ParseInt(tk.Image, 10, 32)
		if err1 == nil {
			return newValueAction(int(v32)), nil
		}
		nerr := err1.(*strconv.NumError)
		if nerr.Err == strconv.ErrRange {
			v64, err2 := strconv.ParseInt(tk.Image, 10, 64)
			if err2 != nil {
				return nil, err2
			}
			return newValueAction(v64), nil
		}
		if nerr.Err == strconv.ErrSyntax {
			f64, err3 := strconv.ParseFloat(tk.Image, 64)
			if err3 != nil {
				return nil, err3
			}
			return newValueAction(f64), nil
		}
		return nil, err1
	default:
		v := valutil.ToBool(tk.Image, false)
		return newValueAction(v), nil
	}
}
