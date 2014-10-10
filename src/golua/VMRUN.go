package golua

import (
	"bmautil/valutil"
	"errors"
	"fmt"
	"golua/goyacc"
	"logger"
	"sync/atomic"
)

type ChunkCode struct {
	name string
	node goyacc.Node
}

func NewChunk(name string, node goyacc.Node) *ChunkCode {
	r := new(ChunkCode)
	r.name = name
	r.node = node
	return r
}

func (this *ChunkCode) Exec(vm *VM) (int, error) {
	return vm.runChunk(this)
}

func (this *ChunkCode) String() string {
	return "Chunk<" + this.name + ">"
}

func (this *VM) Call(nargs int, nresults int) (rint int, rerr error) {
	if this.IsClosing() {
		return 0, fmt.Errorf("%s closed", this)
	}
	st := this.stack
	var nst *VMStack
	r, err := func(nargs int, nresults int) (rint int, rerr error) {
		atomic.AddInt32(&this.running, 1)
		defer func() {
			atomic.AddInt32(&this.running, -1)
			if x := recover(); x != nil {
				logger.Warn(tag, "runtime panic: %v", x)
				if err, ok := x.(error); ok {
					rerr = err
				} else {
					rerr = fmt.Errorf("%v", x)
				}
			}
		}()
		n := nargs + 1
		err1 := this.API_checkstack(n)
		if err1 != nil {
			return 0, err1
		}
		at := this.API_absindex(-n)
		f, err5 := this.API_peek(at)
		if err5 != nil {
			return 0, err5
		}
		f, err5 = this.API_value(f)
		if err5 != nil {
			return 0, err5
		}
		if !this.API_canCall(f) {
			return 0, fmt.Errorf("can't call at '%v'", f)
		}
		nst = newVMStack(st)
		// if tt, ok := f.(StackTracable); ok {
		// 	nst.name = tt.StackInfo()
		// }
		for i := 1; i <= nargs; i++ {
			v, err2 := this.API_peek(at + i)
			if err2 != nil {
				return 0, err2
			}
			nst.stack = append(nst.stack, v)
			nst.stackTop++
		}
		this.API_pop(n)
		this.stack = nst

		if gof, ok := f.(GoFunction); ok {
			nst.gof = gof
			if sfn, ok := f.(supportFuncName); ok {
				nst.funcName = sfn.FuncName()
			}
			rc, err3 := gof.Exec(this)
			if err3 != nil {
				return rc, err3
			}
			at = this.API_absindex(-rc)
			nres := nresults
			if nres < 0 {
				nres = rc
			}
			for i := 0; i < nres; i++ {
				var r interface{}
				if i < rc {
					v, err4 := this.API_peek(at + i)
					if err4 != nil {
						return 0, err4
					}
					r = v
				} else {
					r = nil
				}
				if st.stackTop < len(st.stack) {
					st.stack[st.stackTop] = r
				} else {
					st.stack = append(st.stack, r)
				}
				st.stackTop++
			}
			logger.Debug(tag, "Call %s(%d,%d) -> %d", gof, nargs, nresults, rc)
			return nres, nil
		} else {
			panic(fmt.Errorf("unknow callable '%v'", f))
		}
	}(nargs, nresults)

	if err != nil {
		if _, ok := err.(*StackTraceError); !ok {
			nerr := new(StackTraceError)
			nerr.s = make([]string, 0, 8)
			nerr.s = append(nerr.s, err.Error())
			p := this.stack
			for p != nil {
				nerr.s = append(nerr.s, p.String())
				p = p.parent
			}
			err = nerr
		}
	}

	if nst != nil {
		nst.clear()
	}
	this.stack = st

	return r, err
}

func (this *VM) codeErr(node goyacc.Node, msg string) error {
	return this.codeError(node, errors.New(msg))
}

func (this *VM) codeError(node goyacc.Node, err error) error {
	this.stack.line = node.GetLine()
	return err
}

func (this *VM) runChunk(cc *ChunkCode) (int, error) {
	st := this.stack
	st.chunkName = cc.name
	r, _, err := this.runCode(cc.node)
	return r, err
}

func (this *VM) _runCode(node goyacc.Node) (int, ER, error) {
	if node == nil {
		return 0, ER_NEXT, nil
	}
	op := node.GetOp()
	switch n := node.(type) {
	case *goyacc.Node0:
		switch op {
		case goyacc.OP_VALUE:
			this.API_push(n.Value)
			return 1, ER_NEXT, nil
		case goyacc.OP_VAR:
			s := n.Value.(string)
			va := this.API_var(s)
			this.API_push(va)
			return 1, ER_NEXT, nil
		default:
		}
	case *goyacc.Node1:
		switch op {
		case goyacc.OP_RETURN:
			r1, er1, err1 := this.runCode(n.Child)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			if r1 > 0 {
				pos := this.API_absindex(-r1)
				for i := 0; i < r1; i++ {
					v, _ := this.API_peek(pos + i)
					v, err1 = this.API_value(v)
					if err1 != nil {
						return r1, ER_ERROR, this.codeError(node, err1)
					}
					this.API_replace(pos+i, v)
				}
			}
			return r1, ER_RETURN, nil
		}
	case *goyacc.Node2:
		switch op {
		case goyacc.OP_ADD, goyacc.OP_SUB, goyacc.OP_MUL, goyacc.OP_DIV,
			goyacc.OP_PMUL, goyacc.OP_MOD,
			goyacc.OP_LT, goyacc.OP_GT, goyacc.OP_LTEQ, goyacc.OP_GTEQ, goyacc.OP_EQ, goyacc.OP_NOTEQ,
			goyacc.OP_STRADD, goyacc.OP_AND, goyacc.OP_OR:
			r1, er1, err1 := this.runCode(n.Child1)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			v1, err12 := this.API_value(this.API_pop1X(r1))
			if err12 != nil {
				return 0, ER_ERROR, this.codeError(node, err12)
			}

			r2, er2, err2 := this.runCode(n.Child2)
			if er2 != ER_NEXT {
				return r2, er2, err2
			}
			v2, err22 := this.API_value(this.API_pop1X(r2))
			if err22 != nil {
				return 0, ER_ERROR, this.codeError(node, err22)
			}
			_, rv, err := goyacc.ExecOp2(op, v1, v2)
			if err != nil {
				return 0, ER_ERROR, this.codeError(node, err)
			}
			this.API_push(rv)
			return 1, ER_NEXT, nil
		case goyacc.OP_CALL:
			r1, er1, err1 := this.runCode(n.Child1)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			r2, er2, err2 := this.runCode(n.Child2)
			if er2 != ER_NEXT {
				return r1 + r2, er2, err2
			}
			r0, err0 := this.Call(r2, -1)
			if err0 != nil {
				return r0, ER_ERROR, err0
			}
			return r0, ER_NEXT, nil
		case goyacc.OP_ASSIGN:
			r1, er1, err1 := this.runCode(n.Child1)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}

			r2, er2, err2 := this.runCode(n.Child2)
			if er2 != ER_NEXT {
				return r1 + r2, er2, err2
			}

			vs, err3 := this.API_popN(r2)
			if err3 != nil {
				return 0, ER_ERROR, this.codeError(node, err3)
			}
			vas, err4 := this.API_popN(r1)
			if err4 != nil {
				return 0, ER_ERROR, this.codeError(node, err4)
			}

			for i, va := range vas {
				var v interface{}
				if i < len(vs) {
					v = vs[i]
				} else {
					v = nil
				}
				if va != nil {
					if vao, ok := va.(VMVar); ok {
						_, err5 := vao.Set(v)
						if err5 != nil {
							return 0, ER_ERROR, this.codeError(node, err5)
						}
					} else {
						return 0, ER_ERROR, this.codeErr(node, fmt.Sprintf("invalid var(%T)", va))
					}
				}
			}
			return 0, ER_NEXT, nil
		case goyacc.OP_UNTIL:
			for {
				r1, er1, err1 := this.runCode(n.Child1)
				switch er1 {
				case ER_NEXT, ER_CONTINUE:
				case ER_BREAK:
					return r1, ER_NEXT, err1
				default:
					return r1, er1, err1
				}
				this.API_popN(r1)

				r2, er2, err2 := this.runCode(n.Child2)
				if er2 != ER_NEXT {
					return r2, er2, err2
				}
				v2, err22 := this.API_value(this.API_pop1X(r2))
				if err22 != nil {
					return 0, ER_ERROR, this.codeError(node, err22)
				}
				if valutil.ToBool(v2, true) {
					break
				}
			}
			return 0, ER_NEXT, nil
		case goyacc.OP_WHILE:
			for {
				r1, er1, err1 := this.runCode(n.Child1)
				if er1 != ER_NEXT {
					return r1, er1, err1
				}
				v1, err12 := this.API_value(this.API_pop1X(r1))
				if err12 != nil {
					return 0, ER_ERROR, this.codeError(node, err12)
				}
				if !valutil.ToBool(v1, false) {
					break
				}

				r2, er2, err2 := this.runCode(n.Child2)
				switch er2 {
				case ER_NEXT, ER_CONTINUE:
				case ER_BREAK:
					return r2, ER_NEXT, err2
				default:
					return r2, er1, err2
				}
				this.API_popN(r2)
			}
			return 0, ER_NEXT, nil
		default:
		}
	case *goyacc.NodeFor:
	case *goyacc.NodeFunc:
	case *goyacc.NodeIf:
		r1, er1, err1 := this.runCode(n.Exp)
		if er1 != ER_NEXT {
			return r1, er1, err1
		}
		v1, err12 := this.API_value(this.API_pop1X(r1))
		if err12 != nil {
			return 0, ER_ERROR, this.codeError(node, err12)
		}
		b := valutil.ToBool(v1, false)
		if b {
			return this.runCode(n.Block)
		} else {
			return this.runCode(n.ElseBlock)
		}
	case *goyacc.NodeLocal:
		r0 := 0
		ns := n.Names
		if n.ExpList != nil {
			r, er, err := this.runCode(n.ExpList)
			if er != ER_NEXT {
				return r, er, err
			}
			r0 = r
		}
		var vs []interface{}
		if r0 > 0 {
			vs, _ = this.API_popN(r0)
		}
		st := this.stack
		for i, name := range ns {
			var v interface{}
			if i < len(vs) {
				v = vs[i]
			} else {
				v = nil
			}
			st.createLocal(name, v)
		}
		return 0, ER_NEXT, nil
	case *goyacc.NodeN:
		switch op {
		case goyacc.OP_BLOCK:
			for _, cn := range n.Childs {
				r, er, err := this.runCode(cn)
				if err != nil {
					return 0, ER_ERROR, err
				}
				switch er {
				case ER_BREAK, ER_CONTINUE:
					this.API_pop(r)
					return 0, er, nil
				case ER_RETURN:
					return r, er, nil
				default:
					this.API_pop(r)
				}
			}
			return 0, ER_NEXT, nil
		case goyacc.OP_EXPLIST:
			r0 := 0
			for _, cn := range n.Childs {
				r, er, err := this.runCode(cn)
				if err != nil {
					return 0, ER_ERROR, err
				}
				switch er {
				case ER_BREAK, ER_CONTINUE:
					return r0, er, nil
				case ER_RETURN:
					return r0, er, nil
				default:
					r0 += r
				}
			}
			return r0, ER_NEXT, nil
		}
	}
	return 0, ER_NEXT, nil
}

func (this *VM) runCode(node goyacc.Node) (int, ER, error) {
	// this.Trace(">>> %v", node)
	r1, r2, r3 := this._runCode(node)
	this.Trace(">>> %v -> %d, %d, %v", node, r1, r2, r3)
	return r1, r2, r3
}
