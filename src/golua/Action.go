package golua

import (
	"bytes"
	"fmt"
)

type linfo struct {
	line int
}

type Action interface {
	Line() int
	Children() []Action
	String() string
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

// BlockAction
type BlockAction struct {
	actions []Action
}

func NewBlockAction() *BlockAction {
	r := new(BlockAction)
	r.actions = make([]Action, 0)
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

// ValueAction
type ValueAction struct {
	value interface{}
}

func NewValueAction(v interface{}) *ValueAction {
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

// VarAction
type VarAction struct {
	name string
	linfo
}

func NewVarAction(field bool, tk *Token) *VarAction {
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

// Op2Action
type Op2Action struct {
	id      int
	action1 Action
	action2 Action
}

func newOp2Action(id int) *Op2Action {
	r := new(Op2Action)
	r.id = id
	return r
}

func NewOp2Action(id int, a1, a2 Action) *Op2Action {
	r := newOp2Action(id)
	r.action1 = a1
	r.action2 = a2
	return r
}

func (this *Op2Action) canMerge() bool {
	return this.action1 == nil
}

func (this *Op2Action) Line() int {
	return this.action1.Line()
}

func (this *Op2Action) Children() []Action {
	return []Action{this.action1, this.action2}
}

func (this *Op2Action) String() string {
	return fmt.Sprintf("%s", JJT_NODE_NAME[this.id])
}
