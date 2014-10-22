package golua

import (
	"boot"
	"context"
	"fileloader"
	"fmt"
	"golua/goyacc"
	"logger"
	"sync"
	"time"
)

// GoLua
type GoLua struct {
	ss  fileloader.FileLoader
	cfg *VMConfig
	vmg *VMG

	mux   sync.RWMutex
	codes map[string]*ChunkCode
}

func NewGoLua(n string, ss fileloader.FileLoader, gm VMGInitor, cfg *VMConfig) *GoLua {
	r := new(GoLua)
	r.vmg = NewVMG(n)
	r.vmg.gl = r
	r.ss = ss
	gm(r.vmg)
	r.cfg = cfg
	r.codes = make(map[string]*ChunkCode)
	return r
}

func (this *GoLua) String() string {
	return this.vmg.name
}

func (this *GoLua) ResetCodes() {
	this.mux.Lock()
	defer this.mux.Unlock()
	for k, _ := range this.codes {
		delete(this.codes, k)
	}
}

func (this *GoLua) Close() {
	this.vmg.Close()
	this.vmg.gl = nil
}

func (this *GoLua) Load(script string, save bool) (*ChunkCode, error) {
	// compile
	bs, err := this.ss.Load(script)
	if err != nil {
		err0 := fmt.Errorf("load '%s' fail - %s", script, err)
		logger.Debug(tag, "%s: %s", this, err0)
		return nil, err0
	}
	if bs == nil {
		err0 := fmt.Errorf("can't locate '%s'", script)
		logger.Debug(tag, "%s: %s", this, err0)
		return nil, err0
	}
	logger.Debug(tag, "%s: load('%s') done", this, script)

	content := string(bs)
	p := goyacc.NewParser(script, content)
	node, err2 := p.Parse()
	if err2 != nil {
		err0 := fmt.Errorf("compile '%s' fail - %s", script, err)
		logger.Debug(tag, "%s: %s", this, err0)
		return nil, err0
	}
	logger.Debug(tag, "%s: compile('%s') done", this, script)
	r := NewChunk(script, node)

	if save {
		this.mux.Lock()
		this.codes[script] = r
		this.mux.Unlock()
		logger.Debug(tag, "%s: update('%s')", this, script)
	}
	return r, nil
}

func (this *GoLua) Execute(ctx context.Context) (interface{}, error) {
	req, ok := RequestFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("nil request")
	}
	err0 := req.Valid()
	if err0 != nil {
		return nil, err0
	}

	var cc *ChunkCode
	this.mux.RLock()
	if boot.DevMode {
		for k, _ := range this.codes {
			delete(this.codes, k)
		}
	} else {
		cc = this.codes[req.Script]
	}
	this.mux.RUnlock()

	if cc == nil {
		var err2 error
		cc, err2 = this.Load(req.Script, true)
		if err2 != nil {
			return nil, err2
		}
	}

	// build up env
	locals := make(map[string]interface{})
	locals[KEY_OBJECT_CONTEXT] = ctx
	locals[KEY_CONTEXT] = NewVMTable(req.Context)
	locals[KEY_REQUEST] = NewVMTable(req.Data)

	vm, err3 := this.vmg.CreateVM()
	if err3 != nil {
		return nil, fmt.Errorf("create vm error - %s", err3)
	}
	defer func() {
		vm.CleanDefer()
		vm.Destroy()
	}()

	if this.cfg != nil {
		vm.config = this.cfg
	}
	vm.EnableTrace(req.Trace)
	tm, ok := ctx.Deadline()
	if ok {
		du := tm.Sub(time.Now())
		vm.ResetExecutionTime()
		vm.SetMaxExecutionTime(int(du.Seconds() * 1000))
	}
	vm.context = ctx

	vm.API_push(cc)
	_, err4 := vm.Call(0, 1, locals)
	if err4 != nil {
		logger.Debug(tag, "%s: execute %s fail - %s", this, req, err4)
		return nil, err4
	}
	r, err5 := vm.API_pop1X(-1, true)
	if err5 != nil {
		logger.Debug(tag, "%s: execute %s fail - %s", this, req, err5)
		return nil, err5
	}
	r = GoData(r)
	if logger.EnableDebug(tag) {
		execId, _ := context.ExecIdFromContext(ctx)
		logger.Debug(tag, "%s:[%d] execute %s done -> %v", this, execId, req, r)
	}
	return r, nil
}

func (this *GoLua) Require(pvm *VM, n string) error {
	var cc *ChunkCode
	this.mux.RLock()
	cc = this.codes[n]
	this.mux.RUnlock()

	if cc != nil {
		logger.Debug(tag, "%s: require '%s' exists", this, n)
		return nil
	}

	var err2 error
	cc, err2 = this.Load(n, true)
	if err2 != nil {
		return err2
	}

	var vm *VM
	if pvm != nil {
		vm = pvm
	} else {
		vm, err2 = this.vmg.CreateVM()
		if err2 != nil {
			return fmt.Errorf("create vm error - %s", err2)
		}
		defer vm.Destroy()
	}

	if this.cfg != nil {
		vm.config = this.cfg
	}
	vm.ResetExecutionTime()

	vm.API_push(cc)
	_, err4 := vm.Call(0, 0, nil)
	if err4 != nil {
		logger.Debug(tag, "%s: require '%s' fail - %s", this, n, err4)
		return err4
	}
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "%s: require '%s' done", this, n)
	}
	return nil
}
