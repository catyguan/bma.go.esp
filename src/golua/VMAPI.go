package golua

import

// Converts the acceptable index idx into an absolute index
// (that is, one that does not depend on the stack top).
"fmt"

func (this *VM) API_validindex(idx int) (int, error) {
	v := this.API_absindex(idx)
	if v == 0 {
		return 0, fmt.Errorf("invalid pos %d(%d)", idx, v)
	}
	if v > this.stackTop {
		return 0, fmt.Errorf("overflow pos %d(%d)", idx, v)
	}
	return v, nil
}

func (this *VM) API_absindex(idx int) int {
	if idx >= 0 {
		return idx
	}
	p := this.stackTop + idx + 1
	if p < 0 {
		p = 0
	}
	return p
}

func (this *VM) API_canCall(v interface{}) bool {
	if _, ok := v.(GoFunction); ok {
		return true
	}
	if _, ok := v.(Action); ok {
		return true
	}
	return false
}

func (this *VM) API_getglobal(name string) (interface{}, bool) {
	sh := this.sh
	for {
		if sh == nil {
			break
		}
		sh.RLock()
		v, ok := sh.heap[name]
		sh.RUnlock()
		if ok {
			return v, true
		}
		sh = sh.parent
	}
	return nil, false
}

// Returns the index of the top element in the stack.
// Because indices start at 1, this result is equal to the number of elements in the stack
// (and so 0 means an empty stack).
func (this *VM) API_gettop() int {
	return this.stackTop
}

func (this *VM) API_insert(idx int, v interface{}) error {
	at := this.API_absindex(idx)
	if at == 0 {
		return fmt.Errorf("invalid pos %d(%d)", idx, at)
	}
	if at > this.stackTop+1 {
		return fmt.Errorf("overflow pos %d(%d)", idx, at)
	}
	if at == this.stackTop+1 {
		this.API_push(v)
		return nil
	}
	at--
	result := make([]interface{}, this.stackTop+1, 8)
	copy(result, this.stack[:at])
	result[at] = v
	copy(result[at+1:], this.stack[at:this.stackTop])
	this.stack = result
	this.stackTop++

	return nil
}

func (this *VM) API_pop(n int) error {
	if n > this.stackTop {
		return fmt.Errorf("pop %d overflow", n)
	}
	for i := 0; i < n; i++ {
		this.stack[this.stackTop-n] = nil
	}
	this.stackTop -= n
	return nil
}

func (this *VM) API_push(v interface{}) {
	if this.stackTop < len(this.stack) {
		this.stack[this.stackTop] = v
	} else {
		this.stack = append(this.stack, v)
	}
	this.stackTop++
}

func (this *VM) API_remove(idx int) error {
	at, err := this.API_validindex(idx)
	if err != nil {
		return err
	}
	at--
	copy(this.stack[at:], this.stack[at+1:this.stackTop])
	this.stackTop--
	this.stack[this.stackTop] = nil
	return nil
}

func (this *VM) API_replace(idx int, v interface{}) error {
	at, err := this.API_validindex(idx)
	if err != nil {
		return err
	}
	this.stack[at] = v
	return nil
}

func (this *VM) API_setglobal(n string, v interface{}) {

}

// Accepts any index, or 0, and sets the stack top to this index.
// If the new top is larger than the old one, then the new elements are filled with nil.
// If index is 0, then all stack elements are removed
func (this *VM) API_settop(idx int) {

}

func (this *VM) API_var(n string) (VMVar, bool) {
	return nil, false
}

func (this *VM) API_peek(idx int) (interface{}, bool) {
	return nil, false
}
