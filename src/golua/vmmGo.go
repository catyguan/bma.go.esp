package golua

import (
	"bmautil/valutil"
	"fmt"
	"logger"
	"reflect"
	"time"
)

func GoModule() *VMModule {
	m := NewVMModule("go")
	m.Init("run", GOF_go_run(0))
	m.Init("defer", GOF_go_defer(0))
	m.Init("chan", GOF_go_chan(0))
	m.Init("write", GOF_go_write(0))
	m.Init("read", GOF_go_read(0))
	m.Init("select", GOF_go_select(0))
	m.Init("close", GOF_go_close(0))
	m.Init("deferClose", GOF_go_deferClose(0))
	return m
}

// go.run(func)
type GOF_go_run int

func (this GOF_go_run) Exec(vm *VM) (int, error) {
	n := ""
	c := vm.API_gettop()
	if c > 1 {
		v, err1 := vm.API_pop1(true)
		if err1 != nil {
			return 0, err1
		}
		n = valutil.ToString(v, "")
	}
	f, err2 := vm.API_pop1(true)
	if err2 != nil {
		return 0, err2
	}
	if !vm.API_canCall(f) {
		return 0, fmt.Errorf("param1(%T) can't call", f)
	}
	vm2, err3 := vm.Spawn(n, false)
	if err3 != nil {
		return 0, err3
	}
	vm.Trace("go.run %s", vm2)
	vm2.PrepareRun(true)
	go func() {
		vm2.API_push(f)
		_, errX := vm2.Call(0, 0)
		if errX != nil {
			logger.Debug(tag, "go.run %s fail - %s", vm2, errX)
		}
		vm2.PrepareRun(false)
		vm2.Destroy()
	}()
	return 0, nil
}

func (this GOF_go_run) IsNative() bool {
	return true
}

func (this GOF_go_run) String() string {
	return "GoFunc<go.run>"
}

// go.defer(func)
type GOF_go_defer int

func (this GOF_go_defer) Exec(vm *VM) (int, error) {
	f, err2 := vm.API_pop1(true)
	if err2 != nil {
		return 0, err2
	}
	err3 := vm.API_defer(f, true)
	return 0, err3
}

func (this GOF_go_defer) IsNative() bool {
	return true
}

func (this GOF_go_defer) String() string {
	return "GoFunc<go.defer>"
}

// go.chan(sz)
type GOF_go_chan int

func (this GOF_go_chan) Exec(vm *VM) (int, error) {
	sz, err2 := vm.API_pop1(true)
	if err2 != nil {
		return 0, err2
	}
	vsz := valutil.ToInt(sz, 0)
	if vsz <= 0 {
		return 0, fmt.Errorf("size invalid (%v)", sz)
	}
	ch := make(chan interface{}, vsz)
	vm.API_push(ch)
	return 1, nil
}

func (this GOF_go_chan) IsNative() bool {
	return true
}

func (this GOF_go_chan) String() string {
	return "GoFunc<go.chan>"
}

// go.write(chan, v) bool
type GOF_go_write int

func (this GOF_go_write) Exec(vm *VM) (int, error) {
	ch, val, err1 := vm.API_pop2(true)
	if err1 != nil {
		return 0, err1
	}
	if vch, ok := ch.(chan interface{}); ok {
		r := true
		func() {
			defer func() {
				x := recover()
				if x != nil {
					r = false
				}
			}()
			vch <- val
		}()
		vm.API_push(r)
		return 1, nil
	} else {
		return 0, fmt.Errorf("chan invalid (%v)", ch)
	}
}

func (this GOF_go_write) IsNative() bool {
	return true
}

func (this GOF_go_write) String() string {
	return "GoFunc<go.write>"
}

// go.read(ch [,timeoutMS:int]) value
// go.read(array<ch> [,timeoutMS:int]) (idx, value)
type GOF_go_read int

func (this GOF_go_read) Exec(vm *VM) (int, error) {
	c := vm.API_gettop()
	vtm := 0
	if c > 1 {
		tm, err0 := vm.API_pop1(true)
		if err0 != nil {
			return 0, err0
		}
		vtm = valutil.ToInt(tm, 0)
	}
	ch, err1 := vm.API_pop1(true)
	if err1 != nil {
		return 0, err1
	}
	if vch, ok := ch.(chan interface{}); ok {
		if vtm <= 0 {
			val := <-vch
			vm.API_push(val)
			return 1, nil
		}
		to := time.NewTimer(time.Duration(vtm) * time.Millisecond)
		select {
		case <-to.C:
			return 0, fmt.Errorf("read tiemout(%dms)", vtm)
		case val := <-vch:
			to.Stop()
			vm.API_push(val)
			return 1, nil
		}
	} else if arr, ok := ch.([]interface{}); ok {
		var to *time.Timer
		c := len(arr)
		if vtm > 0 {
			c++
			to = time.NewTimer(time.Duration(vtm) * time.Millisecond)
		}
		cases := make([]reflect.SelectCase, c)
		for i, tmp := range arr {
			cho, ok := tmp.(chan interface{})
			if !ok {
				return 0, fmt.Errorf("chan invalid at array[%d]", i)
			}
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cho)}
		}
		if vtm > 0 {
			cases[c-1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(to.C)}
		}
		chosen, value, ok := reflect.Select(cases)
		to.Stop()
		if vtm > 0 {
			if chosen == c-1 {
				return 0, fmt.Errorf("read tiemout(%dms)", vtm)
			}
		}
		var val interface{}
		if ok && value.IsValid() {
			val = value.Interface()
		}
		vm.API_push(chosen)
		vm.API_push(val)
		return 2, nil
	} else {
		return 0, fmt.Errorf("chan invalid (%v)", ch)
	}
}

func (this GOF_go_read) IsNative() bool {
	return true
}

func (this GOF_go_read) String() string {
	return "GoFunc<go.read>"
}

// go.select(...) (bool, value)
type GOF_go_select int

func (this GOF_go_select) Exec(vm *VM) (int, error) {
	ch, err1 := vm.API_pop1(true)
	if err1 != nil {
		return 0, err1
	}
	if vch, ok := ch.(chan interface{}); ok {
		select {
		case val := <-vch:
			vm.API_push(true)
			vm.API_push(val)
		default:
			vm.API_push(false)
			vm.API_push(nil)
		}
		return 2, nil
	} else {
		return 0, fmt.Errorf("chan invalid (%v)", ch)
	}
}

func (this GOF_go_select) IsNative() bool {
	return true
}

func (this GOF_go_select) String() string {
	return "GoFunc<go.select>"
}

// go.close(obj)
//		obj - chan
type GOF_go_close int

func doClose(o interface{}) bool {
	if o == nil {
		return true
	}
	if vch, ok := o.(chan interface{}); ok {
		close(vch)
		return true
	}
	return false
}

func (this GOF_go_close) Exec(vm *VM) (int, error) {
	ch, err1 := vm.API_pop1(true)
	if err1 != nil {
		return 0, err1
	}
	if doClose(ch) {
		return 0, nil
	} else {
		return 0, fmt.Errorf("invalid close object(%v)", ch)
	}
}

func (this GOF_go_close) IsNative() bool {
	return true
}

func (this GOF_go_close) String() string {
	return "GoFunc<go.close>"
}

// obj4deferClose
type obj4deferClose struct {
	o interface{}
}

func (this *obj4deferClose) Exec(vm *VM) (int, error) {
	doClose(this.o)
	return 0, nil
}

func (this *obj4deferClose) IsNative() bool {
	return true
}

func (this *obj4deferClose) String() string {
	return "GoFunc<obj4deferClose>"
}

// go.deferClose(obj)
type GOF_go_deferClose int

func (this GOF_go_deferClose) Exec(vm *VM) (int, error) {
	o, err1 := vm.API_pop1(true)
	if err1 != nil {
		return 0, err1
	}
	vm.API_defer(&obj4deferClose{o}, true)
	return 0, nil
}

func (this GOF_go_deferClose) IsNative() bool {
	return true
}

func (this GOF_go_deferClose) String() string {
	return "GoFunc<go.deferClose>"
}
