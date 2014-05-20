package glua

import (
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
		this.callback(this)
	}
}

func (this *Context) IsEnd() bool {
	return this.state == stateEnd
}

type ContextUpdater func(ctx *Context)

type ExecuteCallback func(ctx *Context)

type StatisInfo struct {
	Active uint32
}
