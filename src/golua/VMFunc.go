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
	for i, n := range ns {
		var val interface{}
		val = nil
		if i < len(vs) {
			val = vs[i]
		}
		vm.API_createLocal(n, val)
	}
	r1, _, err1 := vm.runCode(this.node.Block)
	return r1, err1
}
