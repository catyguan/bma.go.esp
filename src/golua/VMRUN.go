package golua

import (
	"bmautil/valutil"
	"fmt"
	"golua/goyacc"
	"logger"
	"runtime"
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

func (this *ChunkCode) Exec(vm *VM, self interface{}) (int, error) {
	return vm.runChunk(this, self)
}

func (this *ChunkCode) IsNative() bool {
	return false
}

func (this *ChunkCode) String() string {
	return "Chunk<" + this.name + ">"
}

func (this *VM) Call(nargs int, nresults int, locals map[string]interface{}) (rint int, rerr error) {
	if this.IsClosing() {
		return 0, fmt.Errorf("%s closed", this)
	}
	this.numOfStack++
	if this.numOfStack >= this.config.MaxStack {
		return 0, fmt.Errorf("stack overflow %d", this.numOfStack)
	}
	st := this.stack
	var nst *VMStack
	r, err := func(nargs int, nresults int) (rint int, rerr error) {
		atomic.AddInt32(&this.running, 1)
		defer func() {
			if x := recover(); x != nil {
				trace := make([]byte, 1024)
				runtime.Stack(trace, true)
				logger.Warn(tag, "runtime panic: %v\n%s", x, trace)
				if err, ok := x.(error); ok {
					rerr = err
				} else {
					rerr = fmt.Errorf("%v", x)
				}
			}
			if this.stack.defers != nil {
				l := len(this.stack.defers)
				for i := l - 1; i >= 0; i-- {
					f := this.stack.defers[i]
					this.API_push(f)
					_, errX := this.Call(0, 0, nil)
					if errX != nil {
						if errX != nil {
							logger.Debug(tag, "%s defer %s fail - %s", this, f, errX)
						}
					}
				}
			}
			atomic.AddInt32(&this.running, -1)
		}()
		n := nargs + 1
		err1 := this.API_checkstack(n)
		if err1 != nil {
			return 0, err1
		}
		// function
		at := this.API_absindex(-n)
		at = this.stack.stackBegin + at - 1
		f := this.sdata[at]
		var self interface{}
		if mvar, ok := f.(*memberVar); ok {
			nf, err1x := mvar.Get(this)
			if err1x != nil {
				return 0, err1x
			}
			self = mvar.obj
			f = nf

			// fmt.Println("self = ", self)
		}
		f, err1 = this.API_value(f)
		if err1 != nil {
			return 0, err1
		}
		if !this.API_canCall(f) {
			return 0, fmt.Errorf("can't call at '%v'", f)
		}
		this.sdata[at] = nil

		nst = newVMStack(st)
		nst.stackBegin = at + 1
		nst.stackTop = nargs
		this.stack.stackTop -= n
		this.stack = nst

		for lk, lv := range locals {
			nst.createLocal(this, lk, lv)
		}

		if gof, ok := f.(GoFunction); ok {
			nst.gof = gof
			if sfn, ok := f.(supportFuncName); ok {
				nst.chunkName, nst.funcName = sfn.FuncName()
			}
			rc, err3 := gof.Exec(this, self)
			if err3 != nil {
				return rc, err3
			}
			at = this.API_absindex(-rc)
			at = this.stack.stackBegin + at - 1
			nres := nresults
			if nres < 0 {
				nres = rc
			}
			// if nres > 0 {
			// 	fmt.Printf("BeforeReturn %v AT=%d, OB=%d, OTOP=%d\n", this.sdata, at, st.stackBegin, st.stackTop)
			// }
			for i := 0; i < nres; i++ {
				var r interface{}
				pos := at + i
				if i < rc {
					r = this.sdata[pos]
				} else {
					r = nil
				}
				npos := st.stackBegin + st.stackTop
				this.sdata[npos] = r
				st.stackTop++
				// fmt.Printf("Return%d %v AT=%d, OB=%d, OTOP=%d\n", i+1, this.sdata, pos, st.stackBegin, st.stackTop)
			}
			sttop := st.stackBegin + st.stackTop
			for i := nst.stackBegin; i < nst.stackBegin+nst.stackTop; i++ {
				if i >= sttop {
					this.sdata[i] = nil
				}
			}
			// if nres > 0 {
			// 	fmt.Printf("AfterReturn %v AT=%d, OB=%d, OTOP=%d\n", this.sdata, at, st.stackBegin, st.stackTop)
			// }
			this.Trace("Call %s(%d,%d) -> %d", gof, nargs, nresults, rc)
			return nres, nil
		} else {
			panic(fmt.Errorf("unknow callable '%v'", f))
		}
	}(nargs, nresults)
	this.numOfStack--

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

func (this *VM) runChunk(cc *ChunkCode, self interface{}) (int, error) {
	st := this.stack
	st.chunkName = cc.name
	if self != nil {
		st.createLocal(this, "self", self)
	}
	r, _, err := this.runCode(cc.node)
	return r, err
}

func (this *VM) runCode(node goyacc.Node) (int, ER, error) {
	if this.trace {
		this.Trace(">>> %v", node)
	}
	this.numOfTime++
	if this.numOfTime > this.config.TimeCheck {
		this.numOfTime = 0
		err := this.API_checkExecuteTime()
		if err != nil {
			return 0, ER_ERROR, err
		}
	}

	if node == nil {
		return 0, ER_NEXT, nil
	}
	op := node.GetOp()
	vline := node.GetLine()
	if vline > 0 {
		this.stack.line = vline
	}
	switch n := node.(type) {
	case *goyacc.Node0:
		switch op {
		case goyacc.OP_BREAK:
			return 0, ER_BREAK, nil
		case goyacc.OP_CONTINUE:
			return 0, ER_CONTINUE, nil
		case goyacc.OPF_CLOSURE:
			return 0, ER_NEXT, nil
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
		case goyacc.OP_NOT:
			r1, er1, err1 := this.runCode(n.Child)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			v, err2 := this.API_pop1X(r1, true)
			if err2 != nil {
				return 0, ER_ERROR, err2
			}
			nv := !valutil.ToBool(v, false)
			this.API_push(nv)
			return 1, ER_NEXT, nil
		case goyacc.OP_NSIGN:
			r1, er1, err1 := this.runCode(n.Child)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			v, err2 := this.API_pop1X(r1, true)
			if err2 != nil {
				return 0, ER_ERROR, err2
			}
			ok3, val3, err3 := goyacc.ExecOp2(goyacc.OP_SUB, 0, v)
			if err3 != nil {
				return 0, ER_ERROR, err3
			}
			if !ok3 {
				return 0, ER_ERROR, fmt.Errorf("invalid -%v", v)
			}
			this.API_push(val3)
			return 1, ER_NEXT, nil
		case goyacc.OP_LEN:
			r1, er1, err1 := this.runCode(n.Child)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			v, err2 := this.API_pop1X(r1, true)
			if err2 != nil {
				return 0, ER_ERROR, err2
			}
			var nv interface{}
			nv = 0
			if v != nil {
				switch rv := v.(type) {
				case string:
					nv = len(rv)
				case []interface{}:
					nv = len(rv)
				case map[string]interface{}:
					nv = len(rv)
				case VMTable:
					nv = rv.Len()
				case VMArray:
					nv = rv.Len()
				}
			}
			this.API_push(nv)
			return 1, ER_NEXT, nil
		case goyacc.OP_TABLE:
			tb := this.API_newtable()
			this.API_push(tb)
			r1, er1, err1 := this.runCode(n.Child)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			return 1, ER_NEXT, nil
		case goyacc.OP_ARRAY:
			r1, er1, err1 := this.runCode(n.Child)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			val, err2 := this.API_popN(r1, true)
			if err2 != nil {
				return 0, ER_ERROR, err2
			}
			this.API_push(this.API_array(val))
			return 1, ER_NEXT, nil
		case goyacc.OP_RETURN:
			r1, er1, err1 := this.runCode(n.Child)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			if r1 > 0 {
				pos := this.API_absindex(-r1)
				for i := 0; i < r1; i++ {
					v, err2 := this.API_peek(pos+i, true)
					if err2 != nil {
						return r1, ER_ERROR, err2
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
			v1, err12 := this.API_pop1X(r1, true)
			if err12 != nil {
				return 0, ER_ERROR, err12
			}

			r2, er2, err2 := this.runCode(n.Child2)
			if er2 != ER_NEXT {
				return r2, er2, err2
			}
			v2, err22 := this.API_pop1X(r2, true)
			if err22 != nil {
				return 0, ER_ERROR, err22
			}
			_, rv, err := goyacc.ExecOp2(op, v1, v2)
			if err != nil {
				return 0, ER_ERROR, err
			}
			this.API_push(rv)
			return 1, ER_NEXT, nil
		case goyacc.OP_CALL:
			r1, er1, err1 := this.runCode(n.Child1)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			// fmt.Println("CALL 1", r1, this.DumpStack())
			r2, er2, err2 := this.runCode(n.Child2)
			if er2 != ER_NEXT {
				return r1 + r2, er2, err2
			}
			// fmt.Println("CALL 2", r1, this.DumpStack())
			calll := node.GetLine()
			if calll > 0 {
				this.stack.line = node.GetLine()
			}
			r0, err0 := this.Call(r2, -1, nil)
			if err0 != nil {
				return r0, ER_ERROR, err0
			}
			// fmt.Println("CALL End", r1, this.DumpStack())
			return r0, ER_NEXT, nil
		case goyacc.OP_ASSIGN:
			r1, er1, err1 := this.runCode(n.Child1)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			vas, err4 := this.API_popN(r1, false)
			if err4 != nil {
				return 0, ER_ERROR, err4
			}

			r2, er2, err2 := this.runCode(n.Child2)
			if er2 != ER_NEXT {
				return r2, er2, err2
			}
			vs, err3 := this.API_popN(r2, true)
			if err3 != nil {
				return 0, ER_ERROR, err3
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
						_, err5 := vao.Set(this, v)
						if err5 != nil {
							return 0, ER_ERROR, err5
						}
					} else {
						return 0, ER_ERROR, fmt.Errorf("invalid var(%T)", va)
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
				this.API_pop(r1)

				r2, er2, err2 := this.runCode(n.Child2)
				if er2 != ER_NEXT {
					return r2, er2, err2
				}
				v2, err22 := this.API_pop1X(r2, true)
				if err22 != nil {
					return 0, ER_ERROR, err22
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
				v1, err12 := this.API_pop1X(r1, true)
				if err12 != nil {
					return 0, ER_ERROR, err12
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
				this.API_pop(r2)
			}
			return 0, ER_NEXT, nil
		case goyacc.OP_FIELD:
			r1, er1, err1 := this.runCode(n.Child1)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			v1, err12 := this.API_pop1X(r1, true)
			if err12 != nil {
				return 0, ER_ERROR, err12
			}
			s1 := valutil.ToString(v1, "")

			r2, er2, err2 := this.runCode(n.Child2)
			if er2 != ER_NEXT {
				return r2, er2, err2
			}
			v2, err22 := this.API_pop1X(r2, true)
			if err22 != nil {
				return 0, ER_ERROR, err22
			}

			tb, err3 := this.API_peek(-1, true)
			if err3 != nil {
				return 0, ER_ERROR, err3
			}
			err4 := tb.(VMTable).Set(this, s1, v2)
			if err4 != nil {
				return 0, ER_ERROR, err4
			}
			return 0, ER_NEXT, nil
		case goyacc.OP_MEMBER:
			r1, er1, err1 := this.runCode(n.Child1)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			v1, err12 := this.API_pop1X(r1, true)
			if err12 != nil {
				return 0, ER_ERROR, err12
			}
			if v1 == nil {
				return 0, ER_ERROR, fmt.Errorf("null pointer")
			}

			r2, er2, err2 := this.runCode(n.Child2)
			if er2 != ER_NEXT {
				return r2, er2, err2
			}
			v2, err22 := this.API_pop1X(r2, true)
			if err22 != nil {
				return 0, ER_ERROR, err22
			}

			mvar := new(memberVar)
			mvar.obj = v1
			mvar.key = v2

			this.API_push(mvar)
			return 1, ER_NEXT, nil
		}
	case *goyacc.NodeFor:
		if op == goyacc.OP_FOR {
			// for var=exp1,exp2,exp3 do
			// var从exp1变化到exp2，每次变化以exp3为步长递增var，并执行一次“执行体”。
			// exp3是可选的，如果不指定，默认为1
			name := n.Names[0]
			r1, er1, err1 := this.runCode(n.ForExp)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			if r1 < 2 {
				return r1, ER_ERROR, fmt.Errorf("miss exp1 or exp2(for var=exp1,exp2[,exp3])")
			}
			expl, _ := this.API_popN(r1, true)
			v1 := expl[0]
			v2 := expl[1]

			var v3 interface{}
			if len(expl) > 2 {
				v3 = expl[2]
			} else {
				v3 = 1
			}

			st := this.stack
			st.createLocal(this, name, v1)

			opx := goyacc.OP_LTEQ
			ok7, expv7, err7 := goyacc.ExecOp2(goyacc.OP_LT, v3, 0)
			if err7 != nil {
				return 0, ER_ERROR, err7
			}
			if !ok7 {
				return 0, ER_ERROR, fmt.Errorf("invalid for exp3(%v)", v3)
			}
			if valutil.ToBool(expv7, false) {
				opx = goyacc.OP_GTEQ
			}

			for {
				va := st.local[name]
				val1, err3 := va.Get(this)
				if err3 != nil {
					return 0, ER_ERROR, err3
				}
				ok, expv, err4 := goyacc.ExecOp2(opx, val1, v2)
				if err4 != nil {
					return 0, ER_ERROR, err4
				}
				if !ok {
					return 0, ER_ERROR, fmt.Errorf("invalid exp(%v == %v)", val1, v2)
				}
				if !valutil.ToBool(expv, false) {
					return 0, ER_NEXT, nil
				}

				r5, er5, err5 := this.runCode(n.Block)
				switch er5 {
				case ER_NEXT, ER_CONTINUE:
				case ER_BREAK:
					return r5, ER_NEXT, err5
				default:
					return r5, er5, err5
				}
				this.API_pop(r5)

				ok6, nval, err6 := goyacc.ExecOp2(goyacc.OP_ADD, val1, v3)
				if err6 != nil {
					return 0, ER_ERROR, err6
				}
				if !ok6 {
					return 0, ER_ERROR, fmt.Errorf("invalid exp(%v + %v)", val1, v3)
				}
				_, same, _ := goyacc.ExecOp2(goyacc.OP_EQ, val1, nval)
				if valutil.ToBool(same, true) {
					return 0, ER_ERROR, fmt.Errorf("deadloop exp3(%v + %v)", val1, v3)
				}
				va.Set(this, nval)
			}
		} else {
			// op==OP_FORIN
			// FOR NameList IN ExpList
			name1 := n.Names[0]
			name2 := ""
			if len(n.Names) > 1 {
				name2 = n.Names[1]
			}
			r1, er1, err1 := this.runCode(n.ForExp)
			if er1 != ER_NEXT {
				return r1, er1, err1
			}
			var vlist []interface{}
			var mlist map[string]interface{}
			if r1 == 1 {
				v1, err2 := this.API_pop1X(1, true)
				if err2 != nil {
					return 0, ER_ERROR, err2
				}
				if v1 != nil {
					if tmp, ok := v1.([]interface{}); ok {
						vlist = tmp
					} else {
						if tmp, ok := v1.(VMArray); ok {
							vlist = tmp.ToArray()
						}
					}
				}
				if vlist == nil {
					tb := this.API_table(v1)
					if tb == nil {
						vlist = []interface{}{v1}
					} else {
						mlist = tb.ToMap()
					}
				}
			} else {
				vlist, err1 = this.API_popN(r1, true)
				if err1 != nil {
					return 0, ER_ERROR, err1
				}
			}

			st := this.stack
			st.createLocal(this, name1, nil)
			if name2 != "" {
				st.createLocal(this, name2, nil)
			}

			if mlist != nil {
				for k, val := range mlist {
					va1 := st.local[name1]
					var va2 VMVar
					if name2 != "" {
						va2 = st.local[name2]
						va1.Set(this, k)
						va2.Set(this, val)
					} else {
						va2 = nil
						va1.Set(this, val)
					}

					r5, er5, err5 := this.runCode(n.Block)
					switch er5 {
					case ER_NEXT, ER_CONTINUE:
					case ER_BREAK:
						return r5, ER_NEXT, err5
					default:
						return r5, er5, err5
					}
					this.API_pop(r5)
				}
			} else {
				for i, val := range vlist {
					va1 := st.local[name1]
					var va2 VMVar
					if name2 != "" {
						va2 = st.local[name2]
						va1.Set(this, i)
						va2.Set(this, val)
					} else {
						va2 = nil
						va1.Set(this, val)
					}

					r5, er5, err5 := this.runCode(n.Block)
					switch er5 {
					case ER_NEXT, ER_CONTINUE:
					case ER_BREAK:
						return r5, ER_NEXT, err5
					default:
						return r5, er5, err5
					}
					this.API_pop(r5)
				}
			}
			return 0, ER_NEXT, nil
		}
	case *goyacc.NodeFunc:
		fo := new(VMFunc)
		fo.chunk = this.stack.chunkName
		fo.node = n
		for _, name := range n.CVars {
			va := this.API_findVar(name)
			if va != nil {
				if fo.closures == nil {
					fo.closures = make(map[string]VMVar)
				}
				fo.closures[name] = va
			}
		}
		this.API_push(fo)
		return 1, ER_NEXT, nil
	case *goyacc.NodeIf:
		r1, er1, err1 := this.runCode(n.Exp)
		if er1 != ER_NEXT {
			return r1, er1, err1
		}
		v1, err12 := this.API_pop1X(r1, true)
		if err12 != nil {
			return 0, ER_ERROR, err12
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
			vs, _ = this.API_popN(r0, true)
		}
		st := this.stack
		for i, name := range ns {
			var v interface{}
			if i < len(vs) {
				v = vs[i]
			} else {
				v = nil
			}
			st.createLocal(this, name, v)
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
	return 0, ER_ERROR, fmt.Errorf("unknow op(%d, %s)", op, node)
}

// func (this *VM) runCode(node goyacc.Node) (int, ER, error) {
// 	// this.Trace(">>> %v", node)
// 	r1, r2, r3 := this._runCode(node)
// 	this.Trace(">>> %v -> %d, %d, %v", node, r1, r2, r3)
// 	return r1, r2, r3
// }
