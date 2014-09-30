package glua

import (
	"bmautil/valutil"
	"context"
	"fmt"
	"lua51"
)

type plAllTask struct {
	Name    string
	Request map[string]interface{}
}

type plAllReq struct {
	Tasks map[string]*plAllTask
}

type PluginAll struct {
}

func (tshi *PluginAll) Name() string {
	return "all"
}

func (this *PluginAll) OnInitLua(l *lua51.State) error {
	return nil
}

func (this *PluginAll) OnCloseLua(l *lua51.State) {
}

func (this *PluginAll) Execute(task *PluginTask) error {
	req := task.Request
	if req != nil {
		info := new(plAllReq)
		if valutil.ToBean(req, info) && len(info.Tasks) > 0 {
			ctx := task.Context
			ei, _ := GLuaContext.ExecuteInfo(ctx)
			task.Attach = int(len(info.Tasks))
			scall := func(key string, tobj *plAllTask) {
				cb := func(taskName string, lctx context.Context, cu ContextUpdater, err error) {
					if GLuaContext.IsEnd(lctx) {
						return
					}
					glua := task.GetGLua()
					if glua != nil {
						glua.goo.DoNow(func() {
							if cu != nil {
								cu(lctx)
							}
							if err != nil {
								go task.Callback(this.Name(), nil, err)
								return
							}
							einfo, _ := GLuaContext.ExecuteInfo(ctx)
							einfo2, _ := GLuaContext.ExecuteInfo(lctx)
							if einfo.Result != nil {
								einfo.Result[key] = einfo2.Result
							}
							c := task.Attach.(int)
							c = c - 1
							task.Attach = c
							if c <= 0 {
								go task.Callback(this.Name(), nil, nil)
							}
						})
					}
				}
				glua := task.GetGLua()
				if glua != nil {
					einfo := new(ContextExecuteInfo)
					einfo.Title = fmt.Sprintf("plAll:%s,%s", key, tobj.Name)
					einfo.Result = make(map[string]interface{})
					if ei != nil {
						einfo.Step = ei.Step
						einfo.Data = ei.Data
					}
					sctx := GLuaContext.Sub(ctx, einfo)
					err := glua.StartTask(tobj.Name, sctx, tobj.Request, cb)
					if err != nil {
						cb(tobj.Name, sctx, nil, err)
					}
				}
			}
			for k, tobj := range info.Tasks {
				scall(k, tobj)
			}
			return nil
		}
	}
	go task.Callback(this.Name(), nil, nil)
	return nil
}
