package golua

import "golua/goyacc"

type VMFunc struct {
	chunk    string
	node     *goyacc.NodeFunc
	closures map[string]VMVar
}

func (this *VMFunc) String() string {
	return this.node.String()
}

func (this *VMFunc) FuncName() (string, string) {
	return this.chunk, this.node.Name
}

func (this *VMFunc) IsNative() bool {
	return false
}

func (this *VMFunc) Exec(vm *VM) (int, error) {
	top := vm.API_gettop()
	vs, _ := vm.API_popN(top, true)
	ns := this.node.Params

	st := vm.stack
	for n, val := range this.closures {
		st.local[n] = val
	}
	more := false
	for i, n := range ns {
		if n == KEYWORD_MORE {
			more = true
			break
		}
		var val interface{}
		val = nil
		if i < len(vs) {
			val = vs[i]
		}
		vm.API_createLocal(n, val)
	}
	if more {
		a := make([]interface{}, 0)
		idx := len(ns) - 1
		for i := idx; i < len(vs); i++ {
			a = append(a, vs[i])
		}
		vm.API_createLocal(KEYWORD_MORE, a)
	}
	r1, _, err1 := vm.runCode(this.node.Block)
	return r1, err1
}
