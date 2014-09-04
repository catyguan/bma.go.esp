package glua

import (
	"context"
	"logger"
	"lua51"
)

type PluginTask struct {
	GLua     *GLua
	Service  *Service
	GLuaName string
	Context  context.Context
	Request  map[string]interface{}
	Attach   interface{}
	cb       TaskCallback
}

func (this *PluginTask) GetGLua() *GLua {
	if this.Service != nil {
		return this.Service.GetGLua(this.GLuaName)
	}
	return this.GLua
}

func (this *PluginTask) GetState() *lua51.State {
	gl := this.GetGLua()
	if gl == nil {
		return nil
	}
	return gl.l
}

func (this *PluginTask) Callback(taskName string, cu ContextUpdater, err error) {
	if logger.EnableDebug(tag) {
		glua := this.GetGLua()
		n := "unknow"
		if glua != nil {
			n = glua.name
		}
		s := GLuaContext.String(this.Context)
		if err != nil {
			logger.Debug(tag, "'%s' [%s] task[%s] fail - %s", n, s, taskName, err)
		} else {
			logger.Debug(tag, "'%s' [%s] task[%s] end", n, s, taskName)
		}
	}
	this.cb(taskName, this.Context, cu, err)
}
