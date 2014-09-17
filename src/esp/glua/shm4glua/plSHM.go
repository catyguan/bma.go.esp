package shm4glua

import (
	"bmautil/valutil"
	"context"
	"esp/glua"
	"fmt"
	"lua51"
	"time"
)

const (
	tag = "http4glua"
)

type Request struct {
	Set       bool
	Delete    bool
	Key       string
	Keys      []string
	Value     interface{}
	Size      int
	Timeout   int
	ResultKey string
}

func (this *Request) IsQuery() bool {
	if this.Set {
		return false
	}
	if this.Delete {
		return false
	}
	return true
}

func (this *Request) Valid() error {
	if this.Set {
		if this.Key == "" {
			return fmt.Errorf("set key empty")
		}
		if this.Value == nil {
			return fmt.Errorf("set value null")
		}
	}
	if this.Key == "" && len(this.Keys) == 0 {
		return fmt.Errorf("key empty")
	}
	return nil
}

type PluginSHM struct {
	mem *LinkMap
}

func NewSHM() *PluginSHM {
	r := new(PluginSHM)
	r.mem = newLinkMap()
	return r
}

func (tshi *PluginSHM) Name() string {
	return "shm"
}

func (this *PluginSHM) OnInitLua(l *lua51.State) error {
	return nil
}

func (this *PluginSHM) OnCloseLua(l *lua51.State) {
}

func (this *PluginSHM) Execute(task *glua.PluginTask) error {
	req := new(Request)
	if task.Request != nil {
		valutil.ToBean(task.Request, req)
	}
	err := req.Valid()
	if err != nil {
		return err
	}
	go func() {
		err := this.doExecute(task, req)
		if err != nil {
			task.Callback(this.Name(), nil, err)
		}
	}()
	return nil
}

func (this *PluginSHM) doExecute(task *glua.PluginTask, req *Request) error {
	ctx := task.Context
	rk := req.ResultKey
	if rk == "" {
		rk = "shm"
	}

	var cu glua.ContextUpdater
	if req.IsQuery() {
		if req.Key != "" {
			val, ok := this.mem.Get(req.Key, time.Now())
			cu = func(ctx context.Context) {
				if ok {
					rs := glua.GLuaContext.GetResult(ctx)
					rs[rk] = val
				}
			}
		} else {
			m := this.mem.MGet(req.Keys, time.Now())

			cu = func(ctx context.Context) {
				rs := glua.GLuaContext.GetResult(ctx)
				rs[rk] = m
			}
		}
	} else if req.Delete {
		if req.Key != "" {
			ok := this.mem.Remove(req.Key)
			cu = func(ctx context.Context) {
				rs := glua.GLuaContext.GetResult(ctx)
				rs[rk] = ok
			}
		} else {
			c := this.mem.MRemove(req.Keys)
			cu = func(ctx context.Context) {
				rs := glua.GLuaContext.GetResult(ctx)
				rs[rk] = c
			}
		}
	} else if req.Set {
		tm := req.Timeout
		if tm <= 0 {
			tm = 5 * 60 * 100
		}
		this.mem.Put(req.Key, req.Value, int32(req.Size), tm)
		cu = func(ctx context.Context) {
			rs := glua.GLuaContext.GetResult(ctx)
			rs[rk] = true
		}
	} else {
		cu = func(ctx context.Context) {

		}
	}
	glua.GLuaContext.DoAccessLog(ctx, "shm:execute", nil)
	task.Callback(this.Name(), cu, nil)
	return nil
}
