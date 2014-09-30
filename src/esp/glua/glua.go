package glua

import (
	"bmautil/goo"
	"boot"
	"bytes"
	"context"
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
	QueueSize int
	Paths     []string
	Preloads  []string
}

func (this *ConfigInfo) Valid() error {
	if this.QueueSize <= 0 {
		this.QueueSize = 128
	}
	return nil
}

func (this *ConfigInfo) Compare(old *ConfigInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.QueueSize != old.QueueSize {
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

var (
	execId uint32
)

// GLua
type GLua struct {
	name   string
	config *ConfigInfo

	goo     goo.Goo
	l       *lua51.State
	plugins map[string]GLuaPlugin
	statis  StatisInfo
	context context.Context

	service *Service
}

func NewGLua(n string, cfg *ConfigInfo) *GLua {
	r := new(GLua)
	r.name = n
	r.config = cfg
	r.goo.InitGoo(tag, cfg.QueueSize, r.exitHandler)
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
	l.OpenJson()
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

func (this *GLua) NewContext(title string, acclog bool) context.Context {
	id := atomic.AddUint32(&execId, 1)
	if id == 0 {
		id = atomic.AddUint32(&execId, 1)
	}
	einfo := new(ContextExecuteInfo)
	einfo.Step = 0
	einfo.Title = title
	einfo.Result = make(map[string]interface{})

	var ainfo *ContextAcclogInfo
	if acclog {
		ainfo = new(ContextAcclogInfo)
		ainfo.StartTime = time.Now()
		ainfo.Logdata = make(map[string]interface{})
	}

	return GLuaContext.New(id, einfo, ainfo)
}

func (this *GLua) ExecuteSync(ctx context.Context) error {
	ev := make(chan error, 1)
	defer close(ev)
	cb := func(ctx context.Context, err error) {
		ev <- err
	}
	this.ExecuteNow(ctx, cb)
	select {
	case err := <-ev:
		return err
	case <-ctx.Done():
		err := ctx.Err()
		return err
	}
}

func (this *GLua) ExecuteNow(ctx context.Context, cb GLuaCallback) {
	atomic.AddInt32(&this.statis.Active, 1)
	lcb := func(ctx context.Context, err error) {
		atomic.AddInt32(&this.statis.Active, -1)
		cb(ctx, err)
	}
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	if einfo != nil {
		if einfo.Result == nil {
			einfo.Result = make(map[string]interface{})
		}
		einfo.callback = lcb
	}
	err := this.goo.DoNow(func() {
		this.doExecute(ctx)
	})
	if err != nil {
		cb(ctx, err)
	}
}

func (this *GLua) doExecute(ctx context.Context) {
	err := this.processExecute(ctx, nil)
	if err != nil {
		if logger.EnableDebug(tag) {
			s := GLuaContext.String(ctx)
			logger.Error(tag, "'%s' [%s] execute fail -> %s", this.name, s, err)
		}
		GLuaContext.End(ctx, err)
	}
}

func (this *GLua) processExecute(ctx context.Context, linfo *ContextLuaInfo) error {
	logstr := ""
	if logger.EnableDebug(tag) {
		logstr = GLuaContext.String(ctx)
	}
	if GLuaContext.IsEnd(ctx) {
		logger.Debug(tag, "'%s' [%s] is end, skip", this.name, logstr)
		return nil
	}

	this.context = ctx
	defer func() {
		this.context = nil
	}()	
	if linfo == nil {
		linfo = GLuaContext.Lua(ctx)
	} else {
		if linfo.Script=="" {
			ctxlinfo := GLuaContext.Lua(ctx)
			linfo.Script = ctxlinfo.Script
		}
	}
	if linfo == nil {
		return logger.Error(tag, "processExecute(%s) miss LuaInfo", logstr)
	}
	if linfo.FuncName == "" {
		return fmt.Errorf("processsExecute(%s) miss func name", logstr)
	}

	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	if einfo == nil {
		return logger.Error(tag, "processExecute(%s) miss ExecuteInfo", linfo)
	}
	einfo.Step = einfo.Step + 1
	if logstr != "" {
		logstr = GLuaContext.String(ctx)
	}

	ainfo, _ := GLuaContext.AcclogInfo(ctx)
	var accData map[string]interface{}
	if ainfo != nil {
		accData = ainfo.Logdata
	}

	logger.Debug(tag, "'%s' [%s] execute(%s) start", this.name, logstr, linfo)

	GLuaContext.DoAccessLog(ctx, "execute", nil)

	l := this.l
	if linfo.Script != "" {
		if linfo.Reload {
			err1 := l.Eval(fmt.Sprintf("package.loaded[\"%s\"] = nil", linfo.Script))
			if err1 != nil {
				return err1
			}
		}
		err2 := l.Eval(fmt.Sprintf("require(\"%s\")", linfo.Script))
		if err2 != nil {
			return err2
		}
		logger.Debug(tag, "'%s' load script '%s' done", this.name, linfo.Script)
	}
	l.GetGlobal(linfo.FuncName)
	if !l.IsFunction(-1) {
		l.Pop(1)
		return fmt.Errorf("func '%s' not exists", linfo.FuncName)
	}
	defer l.ClearGValues()
	l.PushGValue(einfo.Data)
	l.PushGValue(einfo.Result)
	l.PushGValue(accData)
	if l.PCall(3, 1, 0) != 0 {
		err := fmt.Errorf("run(%s) fail %s", linfo.FuncName, l.ToString(-1))
		l.Pop(1)
		return err
	}
	if l.IsBoolean(-1) {
		r := l.ToBoolean(-1)
		l.Pop(1)
		if r {
			logger.Debug(tag, "'%s' [%s] execute(%s) done", this.name, logstr, linfo)
			GLuaContext.End(ctx, nil)
		}
	} else if l.IsString(-1) {
		r := l.ToString(-1)
		l.Pop(1)
		logger.Debug(tag, "'%s' [%s] execute(%s) fail -> %s", this.name, logstr, linfo, r)
		GLuaContext.End(ctx, errors.New(r))
	}
	return nil
}

func (this *GLua) StartTask(taskName string, ctx context.Context, req map[string]interface{}, cb TaskCallback) error {
	pl, ok0 := this.plugins[taskName]
	if !ok0 {
		return fmt.Errorf("task '%s' not exists", taskName)
	}
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	if einfo != nil {
		einfo.Step = einfo.Step + 1
	}

	task := new(PluginTask)
	task.Context = ctx
	task.Request = req
	task.GLua = this
	task.Service = this.service
	task.GLuaName = this.name
	task.cb = cb
	s := ""
	if logger.EnableDebug(tag) {
		s = GLuaContext.String(ctx)
	}
	logger.Debug(tag, "'%s' [%s] task[%s] start", this.name, s, taskName)
	err := pl.Execute(task)
	if err != nil {
		logger.Debug(tag, "'%s' [%s] task[%s] fail - %s", this.name, s, taskName, err)
		return err
	}
	return nil
}

func (this *GLua) luaCallback4Task(taskName string, f string, ctx context.Context, cu ContextUpdater, taskErr error) {
	err := this.goo.DoNow(func() {
		if cu != nil {
			cu(ctx)
		}
		if taskErr != nil {
			GLuaContext.End(ctx, taskErr)
			return
		}
		if f != "" {
			linfo := new(ContextLuaInfo)
			linfo.FuncName = f
			err0 := this.processExecute(ctx, linfo)
			if err0 != nil {
				GLuaContext.End(ctx, err0)
			}
		}
	})
	if err != nil {
		GLuaContext.End(ctx, err)
	}
}

func (this *GLua) ReloadScript(n string) error {
	return this.goo.DoSync(func() error {
		l := this.l
		err1 := l.Eval(fmt.Sprintf("package.loaded[\"%s\"] = nil", n))
		if err1 != nil {
			return err1
		}
		err2 := l.Eval(fmt.Sprintf("require(\"%s\")", n))
		if err2 != nil {
			return err2
		}
		logger.Info(tag, "'%s' reload script '%s' done", this.name, n)
		return nil
	})
}

func (this *GLua) LoadScript(n string) error {
	return this.goo.DoSync(func() error {
		l := this.l
		err2 := l.Eval(fmt.Sprintf("require(\"%s\")", n))
		if err2 != nil {
			return err2
		}
		logger.Debug(tag, "'%s' load script '%s' done", this.name, n)
		return nil
	})
}
