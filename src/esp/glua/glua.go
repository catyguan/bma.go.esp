package glua

import (
	"bmautil/goo"
	"boot"
	"bytes"
	"errors"
	"fmt"
	"logger"
	"lua51"
	"sync/atomic"
	"time"
)

const (
	tag = "glua"
)

// Config
type ConfigInfo struct {
	Paths    []string
	Preloads []string
}

func (this *ConfigInfo) Valid() error {
	return nil
}

func (this *ConfigInfo) Compare(old *ConfigInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	// compare Paths
	same := func() bool {
		if len(this.Paths) != len(old.Paths) {
			return false
		}
		tmp := make(map[string]bool)
		for _, s := range this.Paths {
			tmp[s] = true
		}
		for _, s := range old.Paths {
			if _, ok := tmp[s]; !ok {
				return false
			}
		}
		return true
	}()
	if !same {
		return boot.CCR_NEED_START
	}

	return boot.CCR_NONE
}

// GLua
type GLua struct {
	name   string
	config *ConfigInfo

	goo     goo.Goo
	l       *lua51.State
	plugins map[string]GLuaPlugin
	execId  uint32
	statis  StatisInfo
	context *Context
}

func NewGLua(n string, queueSize int, cfg *ConfigInfo) *GLua {
	r := new(GLua)
	r.name = n
	r.config = cfg
	r.goo.InitGoo(tag, queueSize, r.exitHandler)
	r.plugins = make(map[string]GLuaPlugin)
	return r
}

func (this *GLua) Add(pl GLuaPlugin) {
	this.plugins[pl.Name()] = pl
}

func (this *GLua) String() string {
	return this.name
}

func (this *GLua) exitHandler() {
	if this.l != nil {
		logger.Debug(tag, "'%s' close", this.name)
		for _, gp := range this.plugins {
			gp.OnCloseLua(this.l)
		}
		this.l.Close()
		this.l = nil
	}
}

func (this *GLua) Run() error {
	if this.goo.GetState() == goo.STATE_INIT {
		this.goo.Run()
	}
	return this.goo.DoSync(func() error {
		return this.doInitLua()
	})
}

func (this *GLua) Stop() {
	this.goo.Stop()
}

func (this *GLua) StopAndWait() {
	this.goo.StopAndWait()
}

func (this *GLua) doInitLua() error {
	logger.Debug(tag, "'%s' init", this.name)
	l := lua51.NewState()
	this.l = l
	l.OpenLibs()
	// set paths
	pathBuf := bytes.NewBuffer([]byte{})
	for _, s := range this.config.Paths {
		if pathBuf.Len() > 0 {
			pathBuf.WriteByte(';')
		}
		pathBuf.WriteString(s)
		pathBuf.WriteString("/?.lua")
	}
	path := pathBuf.String()
	logger.Debug(tag, "'%s' path='%s'", this.name, path)
	l.SetPath(path)

	this.initGoFunctions()

	printStr := "hostOrgPrint = print\n print = function(...)\n local msg = \"\"\n for i=1, arg.n do msg = msg .. tostring(arg[i])..\"\\t\" end\n glua_print(msg)\n end\n"
	l.Eval(printStr)

	for _, gp := range this.plugins {
		gp.OnInitLua(l)
	}

	for _, pl := range this.config.Preloads {
		logger.Debug(tag, "'%s' preload %s", this.name, pl)
		err := l.Eval(fmt.Sprintf("require(\"%s\")", pl))
		if err != nil {
			logger.Error(tag, "'%s' preload -> %s", this.name, err)
		}
	}

	return nil
}

func (this *GLua) NewContext(f string) *Context {
	id := atomic.AddUint32(&this.execId, 1)
	if id == 0 {
		id = atomic.AddUint32(&this.execId, 1)
	}
	r := new(Context)
	r.Id = id
	r.FuncName = f
	r.Timeout = 5 * time.Second
	return r
}

func (this *GLua) ExecuteSync(ctx *Context) {
	ev := make(chan bool, 1)
	defer close(ev)
	cb := func(rctx *Context) {
		ev <- true
	}
	if err := this.ExecuteNow(ctx, cb); err != nil {
		ctx.End(err)
		return
	}
	<-ev
}

func (this *GLua) ExecuteNow(ctx *Context, cb ExecuteCallback) error {
	ctx.callback = func(rctx *Context) {
		this.statis.Active = this.statis.Active - 1
		cb(rctx)
	}
	if ctx.Result == nil {
		ctx.Result = make(map[string]interface{})
	}
	this.statis.Active = this.statis.Active + 1
	return this.goo.DoNow(func() {
		this.doExecute(ctx)
	})
}

func (this *GLua) doExecute(ctx *Context) {
	err := this.processExecute(ctx.FuncName, ctx)
	if err != nil {
		logger.Error(tag, "'%s' [%s] execute fail -> %s", this.name, ctx, err)
		ctx.End(err)
	}
	if !ctx.IsEnd() {
		time.AfterFunc(ctx.Timeout, func() {
			this.timeout(ctx)
		})
	}
}

func (this *GLua) processExecute(f string, ctx *Context) error {
	this.context = ctx
	defer func() {
		this.context = nil
	}()
	ctx.Step = ctx.Step + 1
	logger.Debug(tag, "'%s' [%s] execute(%s) start", this.name, ctx, f)
	if ctx.IsEnd() {
		logger.Debug(tag, "'%s' [%s] is end, skip", this.name, ctx)
		return nil
	}
	if f == "" {
		return fmt.Errorf("miss func name")
	}
	l := this.l
	l.GetGlobal(f)
	if !l.IsFunction(-1) {
		l.Pop(1)
		return fmt.Errorf("func '%s' not exists", f)
	}
	defer l.ClearGValues()
	l.PushGValue(ctx.Data)
	l.PushGValue(ctx.Result)
	if l.PCall(2, 1, 0) != 0 {
		err := fmt.Errorf("run(%s) fail %s", f, l.ToString(-1))
		l.Pop(1)
		return err
	}
	if l.IsBoolean(-1) {
		r := l.ToBoolean(-1)
		l.Pop(1)
		if r {
			logger.Debug(tag, "'%s' [%s] execute(%s) done", this.name, ctx, f)
			ctx.End(nil)
		}
	} else if l.IsString(-1) {
		r := l.ToString(-1)
		l.Pop(1)
		logger.Debug(tag, "'%s' [%s] execute(%s) fail -> %s", this.name, ctx, f, r)
		ctx.End(errors.New(r))
	}
	return nil
}

func (this *GLua) timeout(ctx *Context) {
	err := this.goo.DoNow(func() {
		logger.Debug(tag, "'%s' [%s] timeout", this.name, ctx)
		ctx.End(errors.New("timeout"))
	})
	if err != nil {
		ctx.End(err)
	}
}

func (this *GLua) StartTask(n string, ctx *Context, req map[string]interface{}, next string, cb func(n string, cu ContextUpdater)) error {
	pl, ok0 := this.plugins[n]
	if !ok0 {
		return fmt.Errorf("task '%s' not exists", n)
	}
	task := new(PluginTask)
	task.Context = ctx
	task.Next = next
	task.Request = req
	task.Service = this
	task.State = this.l
	task.cb = cb
	logger.Debug(tag, "'%s' [%s] task[%s] start", this.name, task.Context, n)
	err := pl.Execute(task)
	if err != nil {
		logger.Debug(tag, "'%s' [%s] task[%s] fail - %s", this.name, ctx, n, err)
		return err
	}
	return nil
}

func (this *GLua) TaskCallback(n string, f string, ctx *Context, cu ContextUpdater) {
	logger.Debug(tag, "'%s' [%s] task[%s] end", this.name, ctx, n)
	err := this.goo.DoNow(func() {
		if cu != nil {
			cu(ctx)
		}
		if f != "" {
			this.processExecute(f, ctx)
		}
	})
	if err != nil {
		ctx.End(err)
	}
}
