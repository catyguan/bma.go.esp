package golua

import (
	"bytes"
	"errors"
	"fmt"
	"golua/goluaparser"
)

type linfo struct {
	line int
}

type actionCol interface {
	AppendAction(act Action) error
}

type StackTracable interface {
	StackInfo() string
}

type ErrAction struct {
	s    string
	line int
}

func dumpActions(buf *bytes.Buffer, act Action, tab string, idx int) {
	for i := 0; i < idx; i++ {
		buf.WriteString(tab)
	}
	if act == nil {
		buf.WriteString("<nil>\n")
	} else {
		buf.WriteString(act.String())
		buf.WriteString("\n")
		for _, ch := range act.Children() {
			dumpActions(buf, ch, tab, idx+1)
		}
	}
}

func DumpActions(act Action, tab string) string {
	if act == nil {
		return "<nil action>"
	}
	buf := bytes.NewBuffer([]byte{})
	dumpActions(buf, act, tab, 0)
	return buf.String()
}

// ChunkAction
type ChunkAction struct {
	name  string
	block Action
}

func newChunkAction(n string) *ChunkAction {
	r := new(ChunkAction)
	r.name = n
	return r
}

func (this *ChunkAction) Line() int {
	if this.block != nil {
		return this.block.Line()
	}
	return -1
}

func (this *ChunkAction) Children() []Action {
	if this.block != nil {
		return []Action{this.block}
	}
	return nil
}

func (this *ChunkAction) StackInfo() string {
	return this.name
}

func (this *ChunkAction) String() string {
	return "Chunk"
}

func (this *ChunkAction) DoMerge(b *LuaBuilder, bn *goluaparser.SimpleNode) error {
	list := b.ChildrenAsAction(bn)
	if len(list) > 0 {
		b.Tracef("ChunkMerge ==> %s", list[0])
		this.block = list[0]
	}
	return nil
}

func (this *ChunkAction) Exec(vm *VM) (int, error) {
	top := vm.API_gettop()
	_, err := this.Process(vm)
	r := vm.API_gettop() - top
	return r, err
}

func (this *ChunkAction) Process(vm *VM) (ACTRES, error) {
	if this.block != nil {
		return this.block.Process(vm)
	}
	return ACTRES_NEXT, nil
}

// BlockAction
type BlockAction struct {
	actions []Action
}

func newBlockAction() *BlockAction {
	r := new(BlockAction)
	return r
}

func (this *BlockAction) Line() int {
	if len(this.actions) > 0 {
		return this.actions[0].Line()
	}
	return -1
}

func (this *BlockAction) Children() []Action {
	return this.actions
}

func (this *BlockAction) String() string {
	return "Block"
}

func (this *BlockAction) DoMerge(b *LuaBuilder, bn *goluaparser.SimpleNode) error {
	if this.actions != nil {
		return nil
	}
	list := b.ChildrenAsAction(bn)
	b.Tracef("BlockMerge ==> size(%d)", len(list))
	this.actions = list
	return nil
}

func (this *BlockAction) Process(vm *VM) (ACTRES, error) {
	for _, act := range this.actions {
		res, err := act.Process(vm)
		if err != nil {
			return ACTRES_ERROR, err
		}
		switch res {
		case ACTRES_NEXT:
		case ACTRES_BREAK, ACTRES_CONTINUE, ACTRES_RETURN:
			return res, nil
		}
	}
	return ACTRES_NEXT, nil
}

// LocalAction
type LocalAction struct {
	names  []string
	assign Action
	line   int
}

func newLocalAction() *LocalAction {
	r := new(LocalAction)
	return r
}

func (this *LocalAction) Line() int {
	return this.line
}

func (this *LocalAction) Children() []Action {
	if this.assign != nil {
		return []Action{this.assign}
	}
	return nil
}

func (this *LocalAction) String() string {
	r := bytes.NewBuffer([]byte{})
	r.WriteString("Local")
	if this.names != nil {
		r.WriteString(fmt.Sprintf("(%v)", this.names))
	}
	return r.String()
}

func (this *LocalAction) DoCombine(b *LuaBuilder, bn *goluaparser.SimpleNode, clist []goluaparser.Node, pos int) error {
	if this.names != nil {
		return nil
	}
	b.Tracef("doCombine %s", clist[pos])
	if pos+1 < len(clist) {
		cbn := clist[pos+1].(*goluaparser.SimpleNode)
		// NAMELIST
		if cbn.Id != goluaparser.JJTNAMELIST {
			return errors.New("invalid local syntax")
		}
		this.names = make([]string, 0, len(cbn.Children))
		for _, ccn := range cbn.Children {
			ccbn := ccn.(*goluaparser.SimpleNode)
			if ccbn.Token != nil {
				this.names = append(this.names, ccbn.Token.Image)
			}
		}
	} else {
		return errors.New("invalid local syntax")
	}
	if pos+2 < len(clist) {
		cbn := clist[pos+2].(*goluaparser.SimpleNode)
		// init value
		if cbn.Data != nil {
			this.assign = cbn.Data.(Action)
			cbn.Data = nil
		}
	}
	return nil
}

func (this *LocalAction) Process(vm *VM) (ACTRES, error) {
	var v interface{}
	if this.assign != nil {
		top := vm.API_gettop()
		r, err := this.assign.Process(vm)
		if err != nil {
			return r, err
		}
		if r != ACTRES_NEXT {
			return r, nil
		}
		if vm.API_gettop()-top > 0 {
			v, err = vm.API_peek(top + 1)
			if err != nil {
				return ACTRES_ERROR, vm.ActionError(this, err)
			}
			v, err = vm.API_value(v)
			if err != nil {
				return ACTRES_ERROR, vm.ActionError(this, err)
			}
			vm.API_popto(top)
		}
	}
	for _, n := range this.names {
		if vv, ok := vm.stack.local[n]; ok {
			if vv != nil {
				if vvv, ok2 := vv.(VMVar); ok2 {
					vvv.Set(v)
				}
			}
		} else {
			vm.stack.local[n] = &localVar{value: v}
		}
	}
	fmt.Println(vm.stack.Dump())
	return ACTRES_NEXT, nil
}

// ValueAction
type ValueAction struct {
	value interface{}
}

func newValueAction(v interface{}) *ValueAction {
	r := new(ValueAction)
	r.value = v
	return r
}

func (this *ValueAction) Line() int {
	return -1
}

func (this *ValueAction) Children() []Action {
	return nil
}

func (this *ValueAction) String() string {
	return fmt.Sprintf("Value(%v)", this.value)
}

func (this *ValueAction) Process(vm *VM) (ACTRES, error) {
	vm.API_push(this.value)
	return ACTRES_NEXT, nil
}

// VarAction
type VarAction struct {
	name string
	linfo
}

func newVarAction(field bool, tk *goluaparser.Token) *VarAction {
	r := new(VarAction)
	r.name = tk.Image
	r.line = tk.BeginLine
	return r
}

func (this *VarAction) Line() int {
	return this.line
}

func (this *VarAction) Children() []Action {
	return nil
}

func (this *VarAction) String() string {
	return fmt.Sprintf("Var(%s)", this.name)
}

func (this *VarAction) Process(vm *VM) (ACTRES, error) {
	vv := vm.API_var(this.name)
	vm.API_push(vv)
	return ACTRES_NEXT, nil
}

// CallAction
type CallAction struct {
	member bool
	args   []Action
	line   int
}

func (this *CallAction) Line() int {
	return this.line
}

func (this *CallAction) Children() []Action {
	return this.args
}

func (this *CallAction) String() string {
	return fmt.Sprintf("Call(%v)", this.member)
}

func (this *CallAction) DoMerge(b *LuaBuilder, bn *goluaparser.SimpleNode) error {
	if this.args != nil {
		return nil
	}
	this.args = b.ChildrenAsAction(bn)
	return nil
}

func (this *CallAction) Process(vm *VM) (ACTRES, error) {
	top := vm.API_gettop()
	if this.args != nil {
		for _, a := range this.args {
			res, err := a.Process(vm)
			if err != nil {
				return res, err
			}
			if res != ACTRES_NEXT {
				return res, nil
			}
		}
	}
	n := vm.API_gettop() - top
	old := vm.stack.line
	vm.stack.line = this.line
	err := vm.Call(n, -1)
	// fmt.Print("action call => ", vm.DumpStack())
	if err != nil {
		return ACTRES_ERROR, err
	}
	vm.stack.line = old
	return ACTRES_NEXT, nil
}

// Op2Action
type Op2Action struct {
	id      int
	kind    int
	action1 Action
	action2 Action
}

func newOp2Action(id int) *Op2Action {
	r := new(Op2Action)
	r.id = id
	return r
}

func (this *Op2Action) Line() int {
	return this.action1.Line()
}

func (this *Op2Action) Children() []Action {
	return []Action{this.action1, this.action2}
}

func (this *Op2Action) String() string {
	if this.kind != 0 {
		return fmt.Sprintf("%s(%s)", goluaparser.JJT_NODE_NAME[this.id], goluaparser.TokenImage[this.kind])
	}
	return fmt.Sprintf("%s", goluaparser.JJT_NODE_NAME[this.id])
}

func (this *Op2Action) DoMerge(b *LuaBuilder, bn *goluaparser.SimpleNode) error {
	if this.action2 != nil {
		return nil
	}
	list := b.ChildrenAsAction(bn)
	if len(list) > 0 {
		b.Tracef("Op2Merge %s =2=> %s", bn, list[0])
		this.action2 = list[0]
	}
	return nil
}

func (this *Op2Action) DoCombine(b *LuaBuilder, bn *goluaparser.SimpleNode, clist []goluaparser.Node, pos int) error {
	if this.action1 == nil {
		if pos > 0 {
			cbn := clist[pos-1].(*goluaparser.SimpleNode)
			if cbn.Data == nil {
				return errors.New("invalid op2 action1 syntax")
			}
			this.action1 = cbn.Data.(Action)
			cbn.Data = nil
			b.Tracef("doCombine %s =1=> %s", bn, this.action1)
		} else {
			return errors.New("invalid local syntax")
		}
	}
	if this.action2 == nil {
		if pos+1 < len(clist) {
			cbn := clist[pos+1].(*goluaparser.SimpleNode)
			if cbn.Data != nil {
				this.action2 = cbn.Data.(Action)
				cbn.Data = nil
				b.Tracef("doCombine %s =2=> %s", bn, this.action2)
			}
		}
	}
	return nil
}

func (this *Op2Action) Process(vm *VM) (ACTRES, error) {
	if this.action1 != nil {
		res, err := this.action1.Process(vm)
		if res != ACTRES_NEXT {
			return res, err
		}
	}
	if this.action2 != nil {
		fmt.Println("process", this.action2)
		res, err := this.action2.Process(vm)
		if res != ACTRES_NEXT {
			return res, err
		}
	}
	switch this.id {
	case goluaparser.JJTASSIGN:
		v1, v2, err1 := vm.API_pop2()
		if err1 != nil {
			return ACTRES_ERROR, vm.ActionError(this, err1)
		}
		vv, ok := v1.(VMVar)
		if !ok {
			return ACTRES_ERROR, vm.ActionError(this, fmt.Errorf("invalid var assign"))
		}
		v3, err2 := vm.API_value(v2)
		if err2 != nil {
			return ACTRES_ERROR, vm.ActionError(this, err2)
		}
		_, err3 := vv.Set(v3)
		if err3 != nil {
			return ACTRES_ERROR, vm.ActionError(this, err3)
		}
	default:
	}
	return ACTRES_NEXT, nil
}
