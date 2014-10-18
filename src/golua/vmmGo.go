package golua

import (
	"bmautil/valutil"
	"fmt"
	"logger"
	"reflect"
	"sync"
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
	m.Init("enableSafe", GOF_go_enableSafe(0))
	m.Init("mutex", GOF_go_mutex(0))
	m.Init("sleep", GOF_go_sleep(0))
	m.Init("timer", GOF_go_timer(0))
	m.Init("ticker", GOF_go_ticker(0))
	m.Init("exec", GOF_go_exec(0))
	return m
}

// go.run(func() [,vmname string])
type GOF_go_run int

func (this GOF_go_run) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	n := ""
	f, nv, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	if nv != nil {
		n = valutil.ToString(nv, "")
	}
	if !vm.API_canCall(f) {
		return 0, fmt.Errorf("param1(%T) can't call", f)
	}
	vm2, err3 := vm.Spawn(n)
	if err3 != nil {
		return 0, err3
	}
	vm.Trace("go.run %s", vm2)
	// vm2.PrepareRun(true)
	go func() {
		vm2.API_push(f)
		_, errX := vm2.Call(0, 0, nil)
		if errX != nil {
			logger.Debug(tag, "go.run %s fail - %s", vm2, errX)
		}
		// vm2.PrepareRun(false)
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

// go.defer(xxx)
//	xxx -- func(), ClosableObject(chan,)
type GOF_go_defer int

func (this GOF_go_defer) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	f, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	if canClose(f) {
		o := f
		f = NewGOF("deferClose", func(vm *VM) (int, error) {
			doClose(o)
			return 0, nil
		})
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
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	sz, err2 := vm.API_pop1X(-1, true)
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
	err0 := vm.API_checkstack(2)
	if err0 != nil {
		return 0, err0
	}
	ch, val, err1 := vm.API_pop2X(-1, true)
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
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	vtm := 0
	ch, tm, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	if tm != nil {
		vtm = valutil.ToInt(tm, 0)
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

// go.select(ch) (bool, value)
type GOF_go_select int

func (this GOF_go_select) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	ch, err1 := vm.API_pop1X(-1, true)
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

// close func
func canClose(o interface{}) bool {
	if o == nil {
		return true
	}
	switch ro := o.(type) {
	case chan interface{}:
		return true
	case *objectVMTable:
		return ro.p.CanClose()
	}
	return false
}

func doClose(o interface{}) bool {
	if o == nil {
		return true
	}
	switch ro := o.(type) {
	case chan interface{}:
		close(ro)
	case *objectVMTable:
		ro.p.Close(ro.o)
	}
	return false
}

// go.close(obj)
//		obj - chan
type GOF_go_close int

func (this GOF_go_close) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	ch, err1 := vm.API_pop1X(-1, true)
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

// go.enableSafe(obj[, val bool])
//		obj - var
type GOF_go_enableSafe int

func (this GOF_go_enableSafe) Exec(vm *VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	o, ival, err1 := vm.API_pop2X(-1, false)
	if err1 != nil {
		return 0, err1
	}
	if valutil.ToBool(ival, false) {
		o, err1 = vm.API_value(o)
		if err1 != nil {
			return 0, err1
		}
	}
	if so, ok := o.(supportSafe); ok {
		so.EnableSafe()
		return 0, nil
	} else {
		return 0, fmt.Errorf("invalid safe(%T)", o)
	}
}

func (this GOF_go_enableSafe) IsNative() bool {
	return true
}

func (this GOF_go_enableSafe) String() string {
	return "GoFunc<go.enableSafe>"
}

// go.mutex([rw:bool])
type GOF_go_mutex int

func (this GOF_go_mutex) Exec(vm *VM) (int, error) {
	rw, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vrw := valutil.ToBool(rw, false)
	var r interface{}
	if vrw {
		o := new(sync.RWMutex)
		r = NewGOO(o, gooRMutex(0))
	} else {
		o := new(sync.Mutex)
		r = NewGOO(o, gooLocker(0))
	}
	vm.API_push(r)
	return 1, nil
}

func (this GOF_go_mutex) IsNative() bool {
	return true
}

func (this GOF_go_mutex) String() string {
	return "GoFunc<go.mutex>"
}

// go.sleep(timeMS:int)
type GOF_go_sleep int

func (this GOF_go_sleep) Exec(vm *VM) (int, error) {
	tm, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vtm := valutil.ToInt(tm, -1)
	if vtm < 0 {
		return 0, fmt.Errorf("invalid sleep time(%v)", tm)
	}
	time.Sleep(time.Duration(vtm) * time.Millisecond)
	return 0, nil
}

func (this GOF_go_sleep) IsNative() bool {
	return true
}

func (this GOF_go_sleep) String() string {
	return "GoFunc<go.sleep>"
}

// go.timer(timeMS:int, func())
type GOF_go_timer int

func (this GOF_go_timer) Exec(vm *VM) (int, error) {
	err1 := vm.API_checkstack(2)
	if err1 != nil {
		return 0, err1
	}
	tm, f, err2 := vm.API_pop2X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vtm := valutil.ToInt(tm, -1)
	if vtm < 0 {
		return 0, fmt.Errorf("invalid timer time(%v)", tm)
	}
	if !vm.API_canCall(f) {
		return 0, fmt.Errorf("timer func(%T) can't call", f)
	}
	vm2, err3 := vm.Spawn("")
	if err3 != nil {
		return 0, err3
	}
	timer := time.AfterFunc(time.Duration(vtm)*time.Millisecond, func() {
		if vm2.IsClosing() {
			return
		}
		vm2.API_push(f)
		_, errX := vm2.Call(0, 0, nil)
		if errX != nil {
			logger.Debug(tag, "go.timer %s fail - %s", vm2, errX)
		}
		vm2.Destroy()
	})
	r := NewGOO(timer, gooTimer(0))
	vm.API_push(r)
	return 1, nil
}

func (this GOF_go_timer) IsNative() bool {
	return true
}

func (this GOF_go_timer) String() string {
	return "GoFunc<go.timer>"
}

// go.ticker(timeMS:int, func())
type GOF_go_ticker int

func (this GOF_go_ticker) Exec(vm *VM) (int, error) {
	err1 := vm.API_checkstack(2)
	if err1 != nil {
		return 0, err1
	}
	tm, f, err2 := vm.API_pop2X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vtm := valutil.ToInt(tm, -1)
	if vtm < 0 {
		return 0, fmt.Errorf("invalid timer time(%v)", tm)
	}
	if !vm.API_canCall(f) {
		return 0, fmt.Errorf("timer func(%T) can't call", f)
	}
	vm2, err3 := vm.Spawn("")
	if err3 != nil {
		return 0, err3
	}
	ticker := time.NewTicker(time.Duration(vtm) * time.Millisecond)
	go func() {
		defer vm2.Destroy()
		for {
			_, ok := <-ticker.C
			if !ok {
				break
			}
			if vm2.IsClosing() {
				break
			}
			vm2.API_push(f)
			_, errX := vm2.Call(0, 0, nil)
			if errX != nil {
				logger.Debug(tag, "go.ticker %s fail - %s", vm2, errX)
			}
		}
	}()
	r := NewGOO(ticker, gooTicker(0))
	vm.API_push(r)
	return 1, nil
}

func (this GOF_go_ticker) IsNative() bool {
	return true
}

func (this GOF_go_ticker) String() string {
	return "GoFunc<go.ticker>"
}

// go.exec(scriptName:string [, locals:map])
type GOF_go_exec int

func (this GOF_go_exec) Exec(vm *VM) (int, error) {
	err1 := vm.API_checkstack(1)
	if err1 != nil {
		return 0, err1
	}
	sn, locals, err2 := vm.API_pop2X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vsn := valutil.ToString(sn, "")
	if vsn == "" {
		return 0, fmt.Errorf("invalid ScriptName(%v)", sn)
	}
	var vtb map[string]interface{}
	if locals != nil {
		switch tmp := locals.(type) {
		case map[string]interface{}:
			vtb = tmp
		case VMTable:
			vtb = tmp.ToMap()
		}
		if vtb == nil {
			return 0, fmt.Errorf("invalid Locals(%v)", locals)
		}
	}

	cc, err3 := vm.vmg.gl.Load(vsn, true)
	if err3 != nil {
		return 0, err3
	}
	vm.API_push(cc)
	rc, err4 := vm.Call(0, -1, vtb)
	if err4 != nil {
		return rc, err4
	}
	return rc, nil
}

func (this GOF_go_exec) IsNative() bool {
	return true
}

func (this GOF_go_exec) String() string {
	return "GoFunc<go.exec>"
}
