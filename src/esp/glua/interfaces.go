package glua

import (
	"context"
	"lua51"
)

type GLuaInit func(l *GLua)

type TaskCallback func(taskName string, ctx context.Context, cu ContextUpdater, err error)

type GLuaCallback func(ctx context.Context, err error)

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

type ContextUpdater func(ctx context.Context)

type StatisInfo struct {
	Active int32
}
