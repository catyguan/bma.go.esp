package golua

import (
	"context"
	"fmt"
	"golua/goyacc"
	"sync"
	"time"
)

// GoLua
type GoLua struct {
	sr  GoSourceRepository
	cfg *VMConfig
	vmg *VMG

	mux   sync.RWMutex
	codes map[string]goyacc.Node
}

func NewGoLua(n string, sr GoSourceRepository, gm VMGInitor, cfg *VMConfig) *GoLua {
	r := new(GoLua)
	r.vmg = NewVMG(n)
	r.sr = sr
	gm(r.vmg)
	r.cfg = cfg
	r.codes = make(map[string]goyacc.Node)
	return r
}

func (this *GoLua) String() string {
	return this.vmg.name
}

func (this *GoLua) Close() {
	this.vmg.Close()
}

func (this *GoLua) Execute(ctx context.Context) (interface{}, error) {
	req, ok := RequestFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("nil request")
	}

	var runNode goyacc.Node
	if !req.Reload {
		this.mux.RLock()
		runNode = this.codes[req.Script]
		this.mux.RUnlock()
	}
	if runNode == nil {
		// compile
		cok, content, err := this.sr.Load(req.Script, req.Reload)
		if err != nil {
			return nil, fmt.Errorf("load '%s' fail - %s", req.Script, err)
		}
		if !cok {
			return nil, fmt.Errorf("can't locate '%s'", req.Script)
		}

		chunkName := req.Script
		p := goyacc.NewParser(chunkName, content)
		node, err2 := p.Parse()
		if err2 != nil {
			return nil, fmt.Errorf("compile '%s' fail - %s", req.Script, err)
		}
		this.mux.Lock()
		this.codes[req.Script] = node
		this.mux.Unlock()
		runNode = node
	}

	chunk := NewChunk(req.Script, runNode)

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

	vm.API_push(chunk)
	_, err4 := vm.CallX(0, 1, locals)
	if req.DumpStack {
		fmt.Println(vm.DumpStack())
	}
	if err4 != nil {
		return nil, err4
	}
	r, err5 := vm.API_pop1X(-1, true)
	if err5 != nil {
		return nil, err5
	}
	return GoData(r), nil
}
