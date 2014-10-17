package golua

import (
	"context"
	"fmt"
	"golua/goyacc"
	"logger"
	"sync"
	"time"
)

// GoLua
type GoLua struct {
	sr  GoSourceRepository
	cfg *VMConfig
	vmg *VMG

	mux   sync.RWMutex
	codes map[string]*ChunkCode
}

func NewGoLua(n string, sr GoSourceRepository, gm VMGInitor, cfg *VMConfig) *GoLua {
	r := new(GoLua)
	r.vmg = NewVMG(n)
	r.sr = sr
	gm(r.vmg)
	r.cfg = cfg
	r.codes = make(map[string]*ChunkCode)
	return r
}

func (this *GoLua) String() string {
	return this.vmg.name
}

func (this *GoLua) Close() {
	this.vmg.Close()
}

func (this *GoLua) Load(script string, reload bool, save bool) (*ChunkCode, error) {
	// compile
	cok, content, err := this.sr.Load(script, reload)
	if err != nil {
		err0 := fmt.Errorf("load '%s' fail - %s", script, err)
		logger.Debug(tag, "%s: %s", this, err0)
		return nil, err0
	}
	if !cok {
		err0 := fmt.Errorf("can't locate '%s'", script)
		logger.Debug(tag, "%s: %s", this, err0)
		return nil, err0
	}
	logger.Debug(tag, "%s: load('%s',%v) done", this, script, reload)

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

	var cc *ChunkCode
	if !req.Reload {
		this.mux.RLock()
		cc = this.codes[req.Script]
		this.mux.RUnlock()
	}
	if cc == nil {
		var err2 error
		cc, err2 = this.Load(req.Script, req.Reload, true)
		if err2 != nil {
			return nil, err2
		}
	}

	vm, err3 := this.vmg.CreateVM()
	if err3 != nil {
		return nil, fmt.Errorf("create vm error - %s", err3)
	}
	defer vm.Destroy()

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

	// build up env
	locals := make(map[string]interface{})
	locals[KEY_OBJECT_CONTEXT] = ctx
	locals[KEY_CONTEXT] = vm.API_table(req.Context)
	locals[KEY_REQUEST] = vm.API_table(req.Data)

	vm.API_push(cc)
	_, err4 := vm.CallX(0, 1, locals)
	if req.DumpStack {
		fmt.Println(vm.DumpStack())
	}
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
	logger.Debug(tag, "%s: execute %s done -> %v", this, req, r)
	return r, nil
}
