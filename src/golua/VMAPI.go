package golua

import "fmt"

func (this *VM) API_checkstack(n int) error {
	st := this.stack
	if st.stackTop < n {
		return fmt.Errorf("stack invalid, size=%d expect=%d", st.stackTop, n)
	}
	return nil
}

func (this *VM) API_validindex(idx int) (int, error) {
	v := this.API_absindex(idx)
	if v == 0 {
		return 0, fmt.Errorf("invalid pos %d(%d)", idx, v)
	}
	if v > this.stack.stackTop {
		return 0, fmt.Errorf("overflow pos %d(%d)", idx, v)
	}
	return v, nil
}

// Converts the acceptable index idx into an absolute index
// (that is, one that does not depend on the stack top).
func (this *VM) API_absindex(idx int) int {
	if idx >= 0 {
		return idx
	}
	p := this.stack.stackTop + idx + 1
	if p < 0 {
		p = 0
	}
	return p
}

func (this *VM) API_canCall(v interface{}) bool {
	if _, ok := v.(GoFunction); ok {
		return true
	}
	return false
}

func (this *VM) API_getglobal(name string) (interface{}, bool) {
	return this.vmg.GetGlobal(name)
}

// Returns the index of the top element in the stack.
// Because indices start at 1, this result is equal to the number of elements in the stack
// (and so 0 means an empty stack).
func (this *VM) API_gettop() int {
	return this.stack.stackTop
}

func (this *VM) API_insert(idx int, v interface{}) error {
	at := this.API_absindex(idx)
	if at == 0 {
		return fmt.Errorf("invalid pos %d(%d)", idx, at)
	}
	st := this.stack
	if at > st.stackTop+1 {
		return fmt.Errorf("overflow pos %d(%d)", idx, at)
	}
	if at == st.stackTop+1 {
		this.API_push(v)
		return nil
	}
	at = st.stackBegin + at - 1
	for i := st.stackBegin + st.stackTop - 1; i >= at; i-- {
		old := this.sdata[i]
		ni := i + 1
		if ni < len(this.sdata) {
			this.sdata[ni] = old
		} else {
			this.sdata = append(this.sdata, old)
		}
	}
	this.sdata[at] = v
	st.stackTop++

	return nil
}

func (this *VM) API_pop(n int) error {
	st := this.stack
	if n > st.stackTop {
		return fmt.Errorf("pop %d overflow", n)
	}
	for i := 0; i < n; i++ {
		st.stackTop--
		this.sdata[st.stackBegin+st.stackTop] = nil
	}
	return nil
}

func (this *VM) API_popto(pos int) {
	st := this.stack
	at := this.API_absindex(pos)
	for st.stackTop > at {
		st.stackTop--
		this.sdata[st.stackBegin+st.stackTop] = nil
	}
}

func (this *VM) API_popN(n int) ([]interface{}, error) {
	st := this.stack
	if n > st.stackTop {
		return nil, fmt.Errorf("pop %d overflow", n)
	}
	ra := make([]interface{}, n)
	for i := 0; i < n; i++ {
		st.stackTop--
		pos := st.stackBegin + st.stackTop
		r := this.sdata[pos]
		this.sdata[pos] = nil
		ra[n-1-i] = r
	}
	return ra, nil
}

func (this *VM) API_pop1() (interface{}, error) {
	st := this.stack
	if 1 > st.stackTop {
		return nil, fmt.Errorf("pop %d overflow", 1)
	}
	st.stackTop--
	pos := st.stackBegin + st.stackTop
	r1 := this.sdata[pos]
	this.sdata[pos] = nil
	return r1, nil
}

func (this *VM) API_pop1X(c int) interface{} {
	st := this.stack
	var r1 interface{}
	for i := 0; i < c; i++ {
		if st.stackTop >= 1 {
			st.stackTop--
			pos := st.stackBegin + st.stackTop
			if i == 0 {
				r1 = this.sdata[pos]
			}
			this.sdata[pos] = nil
		}
	}
	return r1
}

func (this *VM) API_pop2() (interface{}, interface{}, error) {
	st := this.stack
	if 2 > st.stackTop {
		return nil, nil, fmt.Errorf("pop %d overflow", 2)
	}
	st.stackTop--
	pos := st.stackBegin + st.stackTop
	r2 := this.sdata[pos]
	this.sdata[pos] = nil
	st.stackTop--
	pos = st.stackBegin + st.stackTop
	r1 := this.sdata[pos]
	this.sdata[pos] = nil
	return r1, r2, nil
}

func (this *VM) API_pop3() (interface{}, interface{}, interface{}, error) {
	st := this.stack
	if 3 > st.stackTop {
		return nil, nil, nil, fmt.Errorf("pop %d overflow", 3)
	}
	st.stackTop--
	pos := st.stackBegin + st.stackTop
	r3 := this.sdata[pos]
	this.sdata[pos] = nil
	st.stackTop--
	pos = st.stackBegin + st.stackTop
	r2 := this.sdata[pos]
	this.sdata[pos] = nil
	st.stackTop--
	pos = st.stackBegin + st.stackTop
	r1 := this.sdata[pos]
	this.sdata[pos] = nil
	return r1, r2, r3, nil
}

func (this *VM) API_pop4() (interface{}, interface{}, interface{}, interface{}, error) {
	st := this.stack
	if 4 > st.stackTop {
		return nil, nil, nil, nil, fmt.Errorf("pop %d overflow", 4)
	}
	st.stackTop--
	pos := st.stackBegin + st.stackTop
	r4 := this.sdata[pos]
	this.sdata[pos] = nil
	st.stackTop--
	pos = st.stackBegin + st.stackTop
	r3 := this.sdata[pos]
	this.sdata[pos] = nil
	st.stackTop--
	pos = st.stackBegin + st.stackTop
	r2 := this.sdata[pos]
	this.sdata[pos] = nil
	st.stackTop--
	pos = st.stackBegin + st.stackTop
	r1 := this.sdata[pos]
	this.sdata[pos] = nil
	return r1, r2, r3, r4, nil
}

func (this *VM) API_push(v interface{}) {
	st := this.stack
	pos := st.stackBegin + st.stackTop
	if pos < len(this.sdata) {
		this.sdata[pos] = v
	} else {
		this.sdata = append(this.sdata, v)
	}
	st.stackTop++
}

func (this *VM) API_remove(idx int) error {
	at, err := this.API_validindex(idx)
	if err != nil {
		return err
	}
	st := this.stack
	at = st.stackBegin + at - 1
	copy(this.sdata[at:], this.sdata[at+1:st.stackBegin+st.stackTop])
	st.stackTop--
	this.sdata[st.stackBegin+st.stackTop] = nil
	return nil
}

func (this *VM) API_replace(idx int, v interface{}) error {
	at, err := this.API_validindex(idx)
	if err != nil {
		return err
	}
	at = this.stack.stackBegin + at - 1
	this.sdata[at] = v
	return nil
}

func (this *VM) API_setglobal(n string, v interface{}) {
	this.vmg.SetGlobal(n, v)
}

func (this *VM) API_value(v interface{}) (interface{}, error) {
	nv := v
	for {
		if nv == nil {
			return nil, nil
		}
		if a, ok := nv.(VMVar); ok {
			var err error
			nv, err = a.Get(this)
			if err != nil {
				return nil, err
			}
			continue
		}
		return nv, nil
	}
}

func (this *VM) API_var(n string) VMVar {
	st := this.stack
	if v, ok := st.local[n]; ok {
		if vv, ok2 := v.(VMVar); ok2 {
			return vv
		}
	}
	return &globalVar{n}
}

func (this *VM) API_peek(idx int) (interface{}, error) {
	at, err := this.API_validindex(idx)
	if err != nil {
		return nil, err
	}
	at = this.stack.stackBegin + at - 1
	return this.sdata[at], nil
}

func (this *VM) API_createLocal(n string, val interface{}) {
	this.stack.createLocal(this, n, val)
}

func (this *VM) API_findVar(n string) VMVar {
	st := this.stack
	for st != nil {
		if va, ok := st.local[n]; ok {
			return va
		}
		st = st.parent
	}
	return nil
}
