package glua

import (
	"bmautil/valutil"
	"lua51"
)

type plAllTask struct {
	Name    string
	Request map[string]interface{}
}

type plAllReq struct {
	Tasks map[string]plAllTask
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
			task.Attach = int(len(info.Tasks))
			cb := func(n string, cu ContextUpdater, err error) {
				if task.Context.IsEnd() {
					return
				}
				task.Service.goo.DoNow(func() {
					if cu != nil {
						cu(task.Context)
					}
					if err != nil {
						go task.Callback(n, nil, err)
						return
					}
					c := task.Attach.(int)
					c = c - 1
					task.Attach = c
					if c <= 0 {
						go task.Callback(this.Name(), nil, nil)
					}
				})
			}
			for _, tobj := range info.Tasks {
				err := task.Service.StartTask(tobj.Name, task.Context, tobj.Request, "", cb)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}
	go task.Callback(this.Name(), nil, nil)
	return nil
}
