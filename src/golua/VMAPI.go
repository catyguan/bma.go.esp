package golua

import (
	"bmautil/valutil"
	"context"
	"fmt"
	"time"
)

func (this *VM) API_checkExecuteTime() error {
	du := time.Since(this.executeTime).Seconds()
	if int(du*1000) > this.GetMaxExecutionTime() {
		return fmt.Errorf("max execute time(%f)", du)
	}
	return nil
}

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

func (this *VM) API_popAll() {
	this.API_popto(0)
}

func (this *VM) API_popN(n int, popval bool) ([]interface{}, error) {
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
		if popval {
			nv, err := this.API_value(r)
			if err != nil {
				return nil, err
			}
			r = nv
		}
		ra[n-1-i] = r
	}
	return ra, nil
}

func (this *VM) API_pop1X(c int, popval bool) (interface{}, error) {
	st := this.stack
	if c == -1 {
		c = st.stackTop
	}
	var r1 interface{}
	for i := c - 1; i >= 0; i-- {
		if st.stackTop >= 1 {
			st.stackTop--
			pos := st.stackBegin + st.stackTop
			// fmt.Println(st.stackBegin, st.stackTop, len(this.sdata), pos, i)
			if i == 0 {
				r1 = this.sdata[pos]
			}
			this.sdata[pos] = nil
		}
	}
	if popval {
		var err error
		r1, err = this.API_value(r1)
		if err != nil {
			return nil, err
		}
	}
	return r1, nil
}

func (this *VM) API_pop2X(c int, popval bool) (interface{}, interface{}, error) {
	st := this.stack
	if c == -1 {
		c = st.stackTop
	}
	var r1 interface{}
	var r2 interface{}
	for i := c - 1; i >= 0; i-- {
		if st.stackTop >= 1 {
			st.stackTop--
			pos := st.stackBegin + st.stackTop
			switch i {
			case 0:
				r1 = this.sdata[pos]
			case 1:
				r2 = this.sdata[pos]
			}
			this.sdata[pos] = nil
		}
	}
	if popval {
		var err error
		r1, err = this.API_value(r1)
		if err != nil {
			return nil, nil, err
		}
		r2, err = this.API_value(r2)
		if err != nil {
			return nil, nil, err
		}
	}
	return r1, r2, nil
}

func (this *VM) API_pop3X(c int, popval bool) (interface{}, interface{}, interface{}, error) {
	st := this.stack
	if c == -1 {
		c = st.stackTop
	}
	var r1 interface{}
	var r2 interface{}
	var r3 interface{}
	for i := c - 1; i >= 0; i-- {
		if st.stackTop >= 1 {
			st.stackTop--
			pos := st.stackBegin + st.stackTop
			switch i {
			case 0:
				r1 = this.sdata[pos]
			case 1:
				r2 = this.sdata[pos]
			case 2:
				r3 = this.sdata[pos]
			}
			this.sdata[pos] = nil
		}
	}
	if popval {
		var err error
		r1, err = this.API_value(r1)
		if err != nil {
			return nil, nil, nil, err
		}
		r2, err = this.API_value(r2)
		if err != nil {
			return nil, nil, nil, err
		}
		r3, err = this.API_value(r3)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return r1, r2, r3, nil
}

func (this *VM) API_pop4X(c int, popval bool) (interface{}, interface{}, interface{}, interface{}, error) {
	st := this.stack
	if c == -1 {
		c = st.stackTop
	}
	var r1 interface{}
	var r2 interface{}
	var r3 interface{}
	var r4 interface{}
	for i := c - 1; i >= 0; i-- {
		if st.stackTop >= 1 {
			st.stackTop--
			pos := st.stackBegin + st.stackTop
			switch i {
			case 0:
				r1 = this.sdata[pos]
			case 1:
				r2 = this.sdata[pos]
			case 2:
				r3 = this.sdata[pos]
			case 3:
				r4 = this.sdata[pos]
			}
			this.sdata[pos] = nil
		}
	}
	if popval {
		var err error
		r1, err = this.API_value(r1)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		r2, err = this.API_value(r2)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		r3, err = this.API_value(r3)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		r4, err = this.API_value(r4)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
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
	if n == "_" {
		return &VoidVar
	}
	st := this.stack
	if st.local != nil {
		if v, ok := st.local[n]; ok {
			if vv, ok2 := v.(VMVar); ok2 {
				return vv
			}
		}
	}
	return &globalVar{n}
}

func (this *VM) API_peek(idx int, peekval bool) (interface{}, error) {
	at, err := this.API_validindex(idx)
	if err != nil {
		return nil, err
	}
	at = this.stack.stackBegin + at - 1
	v := this.sdata[at]
	if peekval {
		return this.API_value(v)
	}
	return v, nil
}

func (this *VM) API_createLocal(n string, val interface{}) {
	this.stack.createLocal(this, n, val)
}

func (this *VM) API_findVar(n string) VMVar {
	st := this.stack
	for st != nil {
		if st.local != nil {
			if va, ok := st.local[n]; ok {
				return va
			}
		}
		st = st.parent
	}
	return nil
}

func (this *VM) API_defer(f interface{}, parentStack bool) error {
	if !this.API_canCall(f) {
		return fmt.Errorf("defer func(%T) invalid", f)
	}
	st := this.stack
	if parentStack {
		st = st.parent
		if st == nil {
			return fmt.Errorf("parent stack nil when defer")
		}
	}
	if st.defers == nil {
		st.defers = make([]interface{}, 0)
	}
	st.defers = append(st.defers, f)
	return nil
}

func (this *VM) API_getMember(obj interface{}, key interface{}) (interface{}, error) {
	switch o := obj.(type) {
	case []interface{}:
		i := valutil.ToInt(key, -1)
		if i < 0 || i >= len(o) {
			return nil, fmt.Errorf("index(%d) out of range(%d)", i, len(o))
		}
		return o[i], nil
	case VMArray:
		i := valutil.ToInt(key, -1)
		return o.Get(this, i)
	case map[string]interface{}:
		s := valutil.ToString(key, "")
		v := o[s]
		return v, nil
	case VMTable:
		s := valutil.ToString(key, "")
		return o.Get(this, s)

	}
	return nil, fmt.Errorf("unknow memberObject(%T)", obj)
}

func (this *VM) API_setMember(obj interface{}, key interface{}, v interface{}) (bool, error) {
	switch o := obj.(type) {
	case []interface{}:
		i := valutil.ToInt(key, -1)
		if i < 0 || i >= len(o) {
			return false, fmt.Errorf("index(%d) out of range(%d)", i, len(o))
		}
		o[i] = v
		return true, nil
	case VMArray:
		i := valutil.ToInt(key, -1)
		err := o.Set(this, i, v)
		if err != nil {
			return false, err
		}
		return true, nil
	case map[string]interface{}:
		s := valutil.ToString(key, "")
		o[s] = v
		return true, nil
	case VMTable:
		s := valutil.ToString(key, "")
		err := o.Set(this, s, v)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, fmt.Errorf("unknow memberObject(%T)", obj)
}

func (this *VM) API_getContext() context.Context {
	return this.context
}

func (this *VM) API_setContext(ctx context.Context) {
	this.context = ctx
}

func (this *VM) API_cleanDefer(f interface{}) error {
	if canClose(f) {
		o := f
		f = NewGOF("deferClose", func(vm *VM) (int, error) {
			doClose(o)
			return 0, nil
		})
	}
	if !this.API_canCall(f) {
		return fmt.Errorf("clean defer func(%T) invalid", f)
	}
	if this.defers == nil {

	}
	if this.defers == nil {
		this.defers = make([]interface{}, 0)
	}
	this.defers = append(this.defers, f)
	return nil
}
