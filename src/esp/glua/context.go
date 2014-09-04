package glua

import (
	"bytes"
	"context"
	"esp/acclog"
	"fmt"
	"sync/atomic"
	"time"
)

type key int

var key4id key = 0
var key4einfo key = 1
var key4lua key = 2
var key4acclog key = 3
var key4end = 4
var key4callback = 5

type ContextExecuteInfo struct {
	Step   uint32
	Title  string
	Data   map[string]interface{}
	Result map[string]interface{}

	lua      *ContextLuaInfo
	end      uint32
	callback GLuaCallback
}

type ContextLuaInfo struct {
	Script   string
	Reload   bool
	FuncName string
}

func NewLuaInfo(s string, f string, reload bool) *ContextLuaInfo {
	r := new(ContextLuaInfo)
	r.Script = s
	r.FuncName = f
	r.Reload = reload
	return r
}

func (this *ContextLuaInfo) String() string {
	r := bytes.NewBuffer(make([]byte, 0))
	if this.Script != "" {
		r.WriteString(this.Script)
		r.WriteString(":")
	}
	r.WriteString(this.FuncName)
	return r.String()
}

type ContextAcclogInfo struct {
	Acclog    *acclog.Service
	AccName   string
	StartTime time.Time
	Logdata   map[string]interface{}
}

type clsGLuaContext int

var GLuaContext clsGLuaContext = 0

func (clsGLuaContext) New(id uint32, einfo *ContextExecuteInfo, ainfo *ContextAcclogInfo) context.Context {
	r := context.Background()
	r = context.WithValue(r, key4id, id)
	r = context.WithValue(r, key4acclog, ainfo)
	r = context.WithValue(r, key4einfo, einfo)
	return r
}

func (clsGLuaContext) Sub(ctx context.Context, einfo *ContextExecuteInfo) context.Context {
	r := ctx
	r = context.WithValue(r, key4einfo, einfo)
	return r
}

func (clsGLuaContext) SetExecuteInfo(ctx context.Context, title string, s *ContextLuaInfo, data map[string]interface{}) {
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	if einfo != nil {
		if title != "" {
			einfo.Title = title
		}
		if data != nil {
			einfo.Data = data
		}
		if s != nil {
			einfo.lua = s
		}
	}
}

func (clsGLuaContext) GetResult(ctx context.Context) map[string]interface{} {
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	if einfo != nil {
		return einfo.Result
	}
	return nil
}

func (clsGLuaContext) SetLua(ctx context.Context, s *ContextLuaInfo) {
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	if einfo != nil {
		einfo.lua = s
	}
}

func (clsGLuaContext) Lua(ctx context.Context) *ContextLuaInfo {
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	if einfo != nil {
		return einfo.lua
	}
	return nil
}

func (clsGLuaContext) Id(ctx context.Context) (uint32, bool) {
	r, ok := ctx.Value(key4id).(uint32)
	return r, ok
}

func (clsGLuaContext) ExecuteInfo(ctx context.Context) (*ContextExecuteInfo, bool) {
	r, ok := ctx.Value(key4einfo).(*ContextExecuteInfo)
	return r, ok
}

func (clsGLuaContext) AcclogInfo(ctx context.Context) (*ContextAcclogInfo, bool) {
	r, ok := ctx.Value(key4acclog).(*ContextAcclogInfo)
	return r, ok
}

func (clsGLuaContext) String(ctx context.Context) string {
	id, _ := GLuaContext.Id(ctx)
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	title := ""
	var step uint32
	if einfo != nil {
		title = einfo.Title
		step = einfo.Step
	}
	return fmt.Sprintf("%d/%d:%s", id, step, title)
}

func (clsGLuaContext) End(ctx context.Context, err error) {
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	if einfo != nil {
		if atomic.CompareAndSwapUint32(&einfo.end, 0, 1) {
			var dt map[string]interface{}
			if err != nil {
				dt = make(map[string]interface{})
				dt["error"] = err.Error()
			}
			GLuaContext.DoAccessLog(ctx, "end", dt)

			cb := einfo.callback
			if cb != nil {
				cb(ctx, err)
			}
		}
	}
}

func (clsGLuaContext) IsEnd(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
	}
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	if einfo != nil {
		r := atomic.LoadUint32(&einfo.end)
		return r == 1
	}
	return false
}

type ctxAccLogInfo struct {
	data      map[string]interface{}
	time      time.Time
	id        uint32
	step      uint32
	title     string
	act       string
	timeUseMS int
}

func (this *ctxAccLogInfo) Message(cfg map[string]string) string {
	out := bytes.NewBuffer(make([]byte, 0))
	out.WriteString("t=")
	out.WriteString(this.time.Format("2006-01-02 15:04:05"))
	out.WriteString("`")
	out.WriteString("id=")
	out.WriteString(fmt.Sprintf("reqi=%d`reqs=%d`reqn=%s`reqa=%s`", this.id, this.step, this.title, this.act))
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

func (clsGLuaContext) AccessLog(ctx context.Context, act string, dt map[string]interface{}) acclog.AccLogInfo {
	info, _ := GLuaContext.AcclogInfo(ctx)
	if info == nil {
		return nil
	}
	id, _ := GLuaContext.Id(ctx)
	einfo, _ := GLuaContext.ExecuteInfo(ctx)
	title := "unknow"
	var step uint32
	if einfo != nil {
		title = einfo.Title
		step = einfo.Step
	}

	now := time.Now()
	tu := int(now.Sub(info.StartTime).Seconds() * 1000)

	r := new(ctxAccLogInfo)
	r.time = now
	r.timeUseMS = tu
	r.act = act
	r.title = title
	r.id = id
	r.step = step
	if len(info.Logdata)+len(dt) > 0 {
		r.data = make(map[string]interface{})
		for k, v := range info.Logdata {
			r.data[k] = v
		}
		for k, v := range dt {
			r.data[k] = v
		}
	}
	return r
}

func (clsGLuaContext) DoAccessLog(ctx context.Context, act string, dt map[string]interface{}) {
	info, _ := GLuaContext.AcclogInfo(ctx)
	if info == nil {
		return
	}
	if info.Acclog == nil {
		return
	}
	r := GLuaContext.AccessLog(ctx, act, dt)
	info.Acclog.Write(info.AccName, r)
}
