package golua

import (
	"bmautil/valutil"
	"fmt"
	"logger"
	"reflect"
	"runtime"
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
	m.Init("debug", GOF_go_log(logger.LEVEL_DEBUG))
	m.Init("info", GOF_go_log(logger.LEVEL_INFO))
	m.Init("warn", GOF_go_log(logger.LEVEL_WARN))
	m.Init("setGlobal", GOF_go_setGlobal(0))
	m.Init("getGlobal", GOF_go_getGlobal(0))
	m.Init("new", GOF_go_new(0))
	m.Init("invoke", GOF_go_invoke(0))
	m.Init("yield", GOF_go_yield(0))
	m.Init("lookup", GOF_go_lookup(0))
	return m
}

// go.run(func())
type GOF_go_run int

func (this GOF_go_run) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	f, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	if !vm.API_canCall(f) {
		return 0, fmt.Errorf("param1(%T) can't call", f)
	}
	vm2, err3 := vm.gl.GetVM()
	if err3 != nil {
		return 0, err3
	}
	vm.Trace("go.run %s", vm2)
	// vm2.PrepareRun(true)
	go func() {
		defer vm2.Finish()
		vm2.API_push(f)
		_, errX := vm2.Call(0, 0, nil)
		if errX != nil {
			logger.Debug(tag, "go.run %s fail - %s", vm2, errX)
		}
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

func (this GOF_go_defer) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	f, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	if canClose(f) {
		o := f
		f = NewGOF("deferClose", func(vm *VM, self interface{}) (int, error) {
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

// go.cleanDefer(xxx)
type GOF_go_cleanDefer int

func (this GOF_go_cleanDefer) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	f, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	err3 := vm.API_cleanDefer(f)
	return 0, err3
}

func (this GOF_go_cleanDefer) IsNative() bool {
	return true
}

func (this GOF_go_cleanDefer) String() string {
	return "GoFunc<go.cleanDefer>"
}

// go.chan(sz[, closeOnShutdown:bool])
type GOF_go_chan int

func (this GOF_go_chan) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	sz, cos, err2 := vm.API_pop2X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vsz := valutil.ToInt(sz, 0)
	vcos := valutil.ToBool(cos, false)
	if vsz <= 0 {
		return 0, fmt.Errorf("size invalid (%v)", sz)
	}
	ch := make(chan interface{}, vsz)
	vm.API_push(ch)
	if vcos {
		vm.GetGoLua().CreateGoService("chan", ch, func() {
			doClose(ch)
		})
	}
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

func (this GOF_go_write) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
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

// go.read(ch [,timeoutMS:int]) value,(isTimeout bool)
// go.read(array<ch> [,timeoutMS:int]) (idx, value,(isTimeout bool))
type GOF_go_read int

func (this GOF_go_read) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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
	if vtm <= 0 {
		vtm = 5000
	}
	if vch, ok := ch.(chan interface{}); ok {
		to := time.NewTimer(time.Duration(vtm) * time.Millisecond)
		var rv interface{}
		rb := false
		select {
		case <-to.C:
			rb = true
		case val := <-vch:
			to.Stop()
			rv = val
		}
		errC := vm.API_checkRun()
		if errC != nil {
			return 0, errC
		}
		vm.API_push(rv)
		vm.API_push(rb)
		return 2, nil
	} else {
		arr := vm.API_toSlice(ch)
		if arr != nil {
			var to *time.Timer
			c := len(arr) + 1
			to = time.NewTimer(time.Duration(vtm) * time.Millisecond)
			cases := make([]reflect.SelectCase, c)
			for i, tmp := range arr {
				cho, ok := tmp.(chan interface{})
				if !ok {
					return 0, fmt.Errorf("chan invalid at array[%d]", i)
				}
				cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cho)}
			}
			cases[c-1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(to.C)}

			var val interface{}
			rb := false
			chosen, value, ok := reflect.Select(cases)
			to.Stop()
			errC := vm.API_checkRun()
			if errC != nil {
				return 0, errC
			}
			if chosen == c-1 {
				rb = true
			}
			if ok && value.IsValid() {
				val = value.Interface()
			}
			vm.API_push(chosen)
			vm.API_push(val)
			vm.API_push(rb)
			return 3, nil
		} else {
			return 0, fmt.Errorf("chan invalid (%v)", ch)
		}
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

func (this GOF_go_select) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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
	case SupportClose:
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
		defer func() {
			recover()
		}()
		close(ro)
		return true
	case SupportClose:
		ro.Close()
		return true
	case *objectVMTable:
		ro.p.Close(ro.o)
		return true
	}
	return false
}

// go.close(obj)
//		obj - chan
type GOF_go_close int

func (this GOF_go_close) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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
		return 0, fmt.Errorf("invalid close object(%T)", ch)
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

func (this GOF_go_enableSafe) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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
		vm.API_push(o)
		return 1, nil
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

func (this GOF_go_mutex) Exec(vm *VM, self interface{}) (int, error) {
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

func (this GOF_go_sleep) Exec(vm *VM, self interface{}) (int, error) {
	tm, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vtm := valutil.ToInt(tm, -1)
	if vtm < 0 {
		return 0, fmt.Errorf("invalid sleep time(%v)", tm)
	}
	time.Sleep(time.Duration(vtm) * time.Millisecond)
	errC := vm.API_checkRun()
	if errC != nil {
		return 0, errC
	}
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

func (this GOF_go_timer) Exec(vm *VM, self interface{}) (int, error) {
	err1 := vm.API_checkStack(2)
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
	gos := vm.GetGoLua().CreateGoService("timer", nil, nil)
	timer := time.AfterFunc(time.Duration(vtm)*time.Millisecond, func() {
		gos.Close()
		vm2, err3 := gos.GL.GetVM()
		if err3 != nil {
			logger.Debug(tag, "go.timer start fail - %s", err3)
			return
		}
		defer vm2.Finish()

		vm2.API_push(f)
		_, errX := vm2.Call(0, 0, nil)
		if errX != nil {
			logger.Debug(tag, "go.timer %s call fail - %s", vm2, errX)
		}
	})
	r := CreateGoTimer(timer, gos)
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

func (this GOF_go_ticker) Exec(vm *VM, self interface{}) (int, error) {
	err1 := vm.API_checkStack(2)
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
		return 0, fmt.Errorf("ticker func(%T) can't call", f)
	}
	gos := vm.GetGoLua().CreateGoService("ticker", nil, nil)
	ticker := time.NewTicker(time.Duration(vtm) * time.Millisecond)
	go func() {
		for {
			_, ok := <-ticker.C
			if !ok {
				break
			}
			vm2, err := gos.GL.GetVM()
			if err != nil {
				logger.Debug(tag, "go.ticker start fail - %s", err)
				return
			}
			vm2.API_push(f)
			_, errX := vm2.Call(0, 0, nil)
			if errX != nil {
				logger.Debug(tag, "go.ticker %s call fail - %s", vm2, errX)
			}
			vm2.Finish()
		}
	}()
	r := CreateGoTicker(ticker, gos)
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

func (this GOF_go_exec) Exec(vm *VM, self interface{}) (int, error) {
	err1 := vm.API_checkStack(1)
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

	gl := vm.gl
	cc, err3 := gl.ChunkLoad(vm, vsn, true, nil)
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

// go.debug|info|warn(tag, str, ....)
type GOF_go_log int

func (this GOF_go_log) Exec(vm *VM, self interface{}) (int, error) {
	err1 := vm.API_checkStack(2)
	if err1 != nil {
		return 0, err1
	}
	top := vm.API_gettop()
	vs, err2 := vm.API_popN(top, true)
	if err2 != nil {
		return 0, err2
	}
	t := vs[0]
	sf := vs[1]
	vt := valutil.ToString(t, "")
	if vt == "" {
		vt = "log"
	}
	vt = "golua-" + vt
	vsf := valutil.ToString(sf, "")
	if vsf == "" {
		vsf = "<empty log format>"
	}
	logger.DoLog(vt, int(this), vsf, vs[2:]...)
	return 0, nil
}

func (this GOF_go_log) IsNative() bool {
	return true
}

func (this GOF_go_log) String() string {
	return "GoFunc<go.log>"
}

// go.setGlobal(n, v)
type GOF_go_setGlobal int

func (this GOF_go_setGlobal) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	n, v, err2 := vm.API_pop2X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vn := valutil.ToString(n, "")
	vm.API_setglobal(vn, v)
	return 0, nil
}

func (this GOF_go_setGlobal) IsNative() bool {
	return true
}

func (this GOF_go_setGlobal) String() string {
	return "GoFunc<go.setGlobal>"
}

// go.getGlobal(n)
type GOF_go_getGlobal int

func (this GOF_go_getGlobal) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vn := valutil.ToString(n, "")
	v, _ := vm.API_getglobal(vn)
	vm.API_push(v)
	return 1, nil
}

func (this GOF_go_getGlobal) IsNative() bool {
	return true
}

func (this GOF_go_getGlobal) String() string {
	return "GoFunc<go.getGlobal>"
}

// go.new(n)
type GOF_go_new int

func (this GOF_go_new) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vn := valutil.ToString(n, "")
	o, err3 := vm.API_newObject(vn)
	if err3 != nil {
		return 0, err3
	}
	vm.API_push(o)
	return 1, nil
}

func (this GOF_go_new) IsNative() bool {
	return true
}

func (this GOF_go_new) String() string {
	return "GoFunc<go.new>"
}

// go.invoke(self, f, ...) ...
type GOF_go_invoke int

func (this GOF_go_invoke) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	top := vm.API_gettop()
	ns, err2 := vm.API_popN(top, true)
	if err2 != nil {
		return 0, err2
	}
	se := ns[0]
	f := ns[1]
	if se == nil {
		vm.API_push(f)
	} else {
		vm.API_pushMemberCall(se, f)
	}
	for _, v := range ns[2:] {
		vm.API_push(v)
	}
	return vm.Call(top-2, -1, nil)
}

func (this GOF_go_invoke) IsNative() bool {
	return true
}

func (this GOF_go_invoke) String() string {
	return "GoFunc<go.invoke>"
}

// go.yield()
type GOF_go_yield int

func (this GOF_go_yield) Exec(vm *VM, self interface{}) (int, error) {
	runtime.Gosched()
	return 0, nil
}

func (this GOF_go_yield) IsNative() bool {
	return true
}

func (this GOF_go_yield) String() string {
	return "GoFunc<go.yield>"
}

// go.lookup(n)
type GOF_go_lookup int

func (this GOF_go_lookup) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vn := valutil.ToString(n, "")
	vv := vm.API_findVar(vn)
	if vv != nil {
		r, err3 := vv.Get(vm)
		if err3 != nil {
			return 0, err3
		}
		vm.API_push(r)
	} else {
		vm.API_push(nil)
	}
	return 1, nil
}

func (this GOF_go_lookup) IsNative() bool {
	return true
}

func (this GOF_go_lookup) String() string {
	return "GoFunc<go.lookup>"
}
