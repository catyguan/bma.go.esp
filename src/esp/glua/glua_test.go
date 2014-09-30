package glua

import (
	"context"
	"fmt"
	"lua51"
	"os"
	"testing"
	"time"
)

type pl4test struct {
}

func (tshi *pl4test) Name() string {
	return "testpl"
}

func (this *pl4test) OnInitLua(l *lua51.State) error {
	return nil
}

func (this *pl4test) OnCloseLua(l *lua51.State) {
}

func (this *pl4test) Execute(task *PluginTask) error {
	go task.Callback(this.Name(), nil, nil)
	return nil
}

func TestGLuaBase(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		fmt.Sprintln("os exit!!!")
		os.Exit(-1)
	})

	pw, _ := os.Getwd()
	pw = pw

	cfg := new(ConfigInfo)
	cfg.Paths = []string{pw}
	cfg.Preloads = []string{"test"}

	gl := NewGLua("test", cfg)
	gl.Add(new(pl4test))
	gl.Add(new(PluginAll))

	gl.Run()
	defer func() {
		gl.StopAndWait()
		time.Sleep(1 * time.Millisecond)
	}()

	if true {
		ctx := gl.NewContext("hello", false)
		lua := NewLuaInfo("", "hello", false)
		GLuaContext.SetExecuteInfo(ctx, "helloT", lua, nil)

		ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)

		err := gl.ExecuteSync(ctx)
		str := GLuaContext.String(ctx)
		fmt.Println(str, err)
	}

	if true {
		ctx := gl.NewContext("add", false)
		lua := NewLuaInfo("", "add", false)

		dt := make(map[string]interface{})
		dt["a"] = 1
		dt["b"] = 2
		GLuaContext.SetExecuteInfo(ctx, "", lua, dt)

		err := gl.ExecuteSync(ctx)
		str := GLuaContext.String(ctx)
		rs := GLuaContext.GetResult(ctx)
		fmt.Println(str, dt, rs, err)
	}

	if true {
		ctx := gl.NewContext("async", false)
		lua := NewLuaInfo("", "async", false)
		GLuaContext.SetExecuteInfo(ctx, "", lua, nil)

		// ctx.Timeout = 1 * time.Second
		err := gl.ExecuteSync(ctx)
		str := GLuaContext.String(ctx)
		rs := GLuaContext.GetResult(ctx)
		fmt.Println(str, rs, err)
	}
	if true {
		ctx := gl.NewContext("all", false)
		lua := NewLuaInfo("", "all", false)
		GLuaContext.SetExecuteInfo(ctx, "", lua, nil)

		// ctx.Timeout = 1 * time.Second
		err := gl.ExecuteSync(ctx)
		str := GLuaContext.String(ctx)
		rs := GLuaContext.GetResult(ctx)
		fmt.Println(str, rs, err)
	}
}
