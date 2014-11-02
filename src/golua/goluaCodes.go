package golua

import (
	"boot"
	"context"
	"fileloader"
	"fmt"
	"golua/goyacc"
	"io/ioutil"
	"logger"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (this *GoLua) ResetCodes() {
	this.codeMutex.Lock()
	defer this.codeMutex.Unlock()
	for k, _ := range this.codes {
		delete(this.codes, k)
	}
}

func (this *GoLua) FileLoad(vm *VM, fn string) ([]byte, error) {
	fn = this.ParseScriptName(vm, fn)
	bs, err := this.ss.Load(fn)
	if err != nil {
		err0 := fmt.Errorf("fileLoad '%s' fail - %s", fn, err)
		logger.Debug(tag, "%s: %s", this, err0)
		return nil, err0
	}
	return bs, nil
}

func (this *GoLua) ChunkLoad(vm *VM, script string, save bool, spp ScriptPreprocess) (*ChunkCode, error) {
	n := this.ParseScriptName(vm, script)
	var cc *ChunkCode
	this.codeMutex.RLock()
	cc = this.codes[script]
	this.codeMutex.RUnlock()
	if cc == nil {
		return this.compile(n, save, spp)
	}
	return cc, nil
}

func (this *GoLua) compile(script string, save bool, spp ScriptPreprocess) (*ChunkCode, error) {
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
	if spp != nil {
		str, err0 := spp(content)
		if err0 != nil {
			logger.Debug(tag, "%s: preprocess('%s') fail - %s", this, script, err0)
			return nil, err0
		}
		content = str

		if boot.DevMode {
			m, file := fileloader.SplitModuleScript(script)
			fn, _ := filepath.Abs("tmp/" + m + "/" + file)
			dir := filepath.Dir(fn)
			os.MkdirAll(dir, os.ModePerm)
			errF2 := ioutil.WriteFile(fn, []byte(content), os.ModePerm)
			if errF2 != nil {
				logger.Debug(tag, "write %s fail - %s", fn, errF2)
			} else {
				logger.Debug(tag, "'%s' preprocess -> %s", script, fn)
			}
		}
	}

	p := goyacc.NewParser(script, content)
	node, err2 := p.Parse()
	if err2 != nil {
		err0 := fmt.Errorf("compile '%s' fail - %s", script, err2)
		logger.Debug(tag, "%s: %s", this, err0)
		return nil, err0
	}
	logger.Debug(tag, "%s: compile('%s') done", this, script)
	r := NewChunk(script, node)

	if save {
		this.codeMutex.Lock()
		this.codes[script] = r
		this.codeMutex.Unlock()
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

	script := this.ParseScriptName(nil, req.Script)

	var cc *ChunkCode
	this.codeMutex.RLock()
	if boot.DevMode {
		for k, _ := range this.codes {
			delete(this.codes, k)
		}
	} else {
		cc = this.codes[script]
	}
	this.codeMutex.RUnlock()

	if cc == nil {
		var err2 error
		cc, err2 = this.compile(script, true, nil)
		if err2 != nil {
			return nil, err2
		}
	}

	// build up env
	locals := make(map[string]interface{})
	locals[KEY_OBJECT_CONTEXT] = ctx
	locals[KEY_CONTEXT] = NewVMTable(req.Context)
	locals[KEY_REQUEST] = NewVMTable(req.Data)

	vm, err3 := this.GetVM()
	if err3 != nil {
		return nil, fmt.Errorf("create vm error - %s", err3)
	}
	defer func() {
		vm.Finish()
	}()

	if this.cfg != nil {
		vm.config = this.cfg
	}
	vm.EnableTrace(req.Trace)
	vm.ResetExecutionTime()

	tm, ok := ctx.Deadline()
	if ok {
		du := tm.Sub(time.Now())
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

func (this *GoLua) ParseScriptName(vm *VM, n string) string {
	change := false
	m, f := fileloader.SplitModuleScript(n)
	if strings.HasPrefix(m, "_") && vm != nil {
		cn := vm.stack.chunkName
		m2, _ := fileloader.SplitModuleScript(cn)
		m = m2 + strings.TrimPrefix(m, "_")
		change = true
	}
	if strings.HasPrefix(f, "/") {
		f = f[1:]
		change = true
	}
	if change {
		return fileloader.BuildModuleScript(m, f)
	}
	return n
}

func (this *GoLua) Require(pvm *VM, n string) error {
	n = this.ParseScriptName(pvm, n)

	var cc *ChunkCode
	this.codeMutex.RLock()
	cc = this.codes[n]
	this.codeMutex.RUnlock()

	if cc != nil {
		logger.Debug(tag, "%s: require '%s' exists", this, n)
		return nil
	}

	var err2 error
	cc, err2 = this.compile(n, true, nil)
	if err2 != nil {
		return err2
	}

	var vm *VM
	if pvm != nil {
		vm = pvm
	} else {
		vm, err2 = this.GetVM()
		if err2 != nil {
			return fmt.Errorf("create vm error - %s", err2)
		}
		defer func() {
			vm.Finish()
		}()
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
