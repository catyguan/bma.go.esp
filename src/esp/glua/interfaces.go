package glua

import (
	"bytes"
	"esp/acclog"
	"fmt"
	"logger"
	"lua51"
	"time"
)

type GLuaInit func(l *GLua)

type TaskCallback func(n string, cu ContextUpdater, err error)

type PluginTask struct {
	Service *GLua
	State   *lua51.State
	Context *Context
	Request map[string]interface{}
	Next    string
	Attach  interface{}
	cb      TaskCallback
}

func (this *PluginTask) Callback(n string, cu ContextUpdater, err error) {
	if logger.EnableDebug(tag) {
		if err != nil {
			logger.Debug(tag, "'%s' [%s] task[%s] fail - %s", this.Service.name, this.Context, n, err)
		} else {
			logger.Debug(tag, "'%s' [%s] task[%s] end", this.Service.name, this.Context, n)
		}
	}
	if this.cb == nil {
		this.Service.TaskCallback(n, this.Next, this.Context, cu, err)
	} else {
		this.cb(n, cu, err)
	}
}

type GLuaPlugin interface {
	Name() string
	OnInitLua(l *lua51.State) error
	OnCloseLua(l *lua51.State)
	Execute(task *PluginTask) error
}

const (
	stateInit = iota
	stateActive
	stateEnd
)

type Context struct {
	Id       uint32
	Step     uint32
	FuncName string
	Data     map[string]interface{}
	Timeout  time.Duration
	Result   map[string]interface{}
	Error    error

	Acclog    *acclog.Service
	AccName   string
	StartTime time.Time
	Logdata   map[string]interface{}

	callback ExecuteCallback
	state    int
	timer    *time.Timer
}

func (this *Context) String() string {
	return fmt.Sprintf("%d:%s/%d", this.Id, this.FuncName, this.Step)
}

func (this *Context) End(err error) {
	if this.state != stateEnd {
		this.Error = err
		this.state = stateEnd
		t := this.timer
		this.timer = nil
		if t != nil {
			t.Stop()
		}
		if this.Acclog != nil {
			var dt map[string]interface{}
			if err != nil {
				dt = make(map[string]interface{})
				dt["err"] = err.Error()
			}
			this.DoAccessLog("end", dt)
		}
		this.callback(this)
	}
}

func (this *Context) IsEnd() bool {
	return this.state == stateEnd
}

type ctxAccLogInfo struct {
	data      map[string]interface{}
	time      time.Time
	id        uint32
	step      uint32
	funcName  string
	act       string
	timeUseMS int
}

func (this *ctxAccLogInfo) Message(cfg map[string]string) string {
	out := bytes.NewBuffer(make([]byte, 0))
	out.WriteString("t=")
	out.WriteString(this.time.Format("2006-01-02 15:04:05"))
	out.WriteString("`")
	out.WriteString("id=")
	out.WriteString(fmt.Sprintf("reqi=%d`reqs=%d`reqn=%s`reqa=%s`", this.id, this.step, this.funcName, this.act))
	for k, v := range this.data {
		if v != nil {
			out.WriteString(k)
			out.WriteString("=")
			out.WriteString(fmt.Sprintf("%v", v))
			out.WriteString("`")
		}
	}
	out.WriteString("tu=")
	out.WriteString(fmt.Sprintf("%d", this.timeUseMS))
	out.WriteByte('\n')
	return out.String()
}

func (this *ctxAccLogInfo) TimeDay() int {
	return this.time.Day()
}

func (this *Context) DoAccessLog(act string, dt map[string]interface{}) {
	if this.Acclog == nil {
		return
	}
	now := time.Now()
	tu := int(now.Sub(this.StartTime).Seconds() * 1000)

	r := new(ctxAccLogInfo)
	r.time = now
	r.timeUseMS = tu
	r.act = act
	r.funcName = this.FuncName
	r.id = this.Id
	r.step = this.Step
	if len(this.Logdata)+len(dt) > 0 {
		r.data = make(map[string]interface{})
		for k, v := range this.Logdata {
			r.data[k] = v
		}
		for k, v := range dt {
			r.data[k] = v
		}
	}
	this.Acclog.Write(this.AccName, r)
}

type ContextUpdater func(ctx *Context)

type ExecuteCallback func(ctx *Context)

type StatisInfo struct {
	Active uint32
}
