package glua

import (
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

	gl := NewGLua("test", 16, cfg)
	gl.Add(new(pl4test))
	gl.Add(new(PluginAll))

	gl.Run()
	defer func() {
		gl.StopAndWait()
		time.Sleep(1 * time.Millisecond)
	}()

	if true {
		ctx := gl.NewContext("hello")
		ctx.Timeout = 100 * time.Millisecond
		gl.ExecuteSync(ctx)
		fmt.Println(ctx, ctx.Error)
	}

	if true {
		ctx := gl.NewContext("add")
		ctx.Timeout = 1 * time.Second

		dt := make(map[string]interface{})
		dt["a"] = 1
		dt["b"] = 2
		ctx.Data = dt

		gl.ExecuteSync(ctx)
		fmt.Println(ctx, ctx.Data, ctx.Result, ctx.Error)
	}

	if true {
		ctx := gl.NewContext("async")
		ctx.Timeout = 1 * time.Second
		gl.ExecuteSync(ctx)
		fmt.Println(ctx, ctx.Data, ctx.Result, ctx.Error)
	}
	if true {
		ctx := gl.NewContext("all")
		ctx.Timeout = 1 * time.Second
		gl.ExecuteSync(ctx)
		fmt.Println(ctx, ctx.Data, ctx.Result, ctx.Error)
	}
}
