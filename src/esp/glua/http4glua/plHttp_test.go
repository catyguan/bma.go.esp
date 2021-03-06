package http4glua

import (
	"context"
	"esp/glua"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPluginHttp(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		fmt.Sprintln("os exit!!!")
		os.Exit(-1)
	})

	pw, _ := os.Getwd()
	pw = pw

	cfg := new(glua.ConfigInfo)
	cfg.QueueSize = 16
	cfg.Paths = []string{pw}
	cfg.Preloads = []string{"test"}

	gl := glua.NewGLua("test", cfg)
	gl.Add(new(glua.PluginAll))
	gl.Add(new(PluginHttp))

	gl.Run()
	defer func() {
		gl.StopAndWait()
		time.Sleep(1 * time.Millisecond)
	}()

	if true {
		ctx := gl.NewContext("http", false)
		lua := glua.NewLuaInfo("", "http", false)
		glua.GLuaContext.SetExecuteInfo(ctx, "", lua, nil)

		ctx, _ = context.WithTimeout(ctx, 3*time.Second)

		err := gl.ExecuteSync(ctx)
		str := glua.GLuaContext.String(ctx)
		rs := glua.GLuaContext.GetResult(ctx)
		// hresp := ctx.Result["http"]
		// if hresp != nil {
		// 	delete(hresp.(map[string]interface{}), "Content")
		// }
		fmt.Println(str, rs, err)
	}
}
