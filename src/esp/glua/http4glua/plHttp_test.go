package http4glua

import (
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
	cfg.Paths = []string{pw}
	cfg.Preloads = []string{"test"}

	gl := glua.NewGLua("test", 16, cfg)
	gl.Add(new(glua.PluginAll))
	gl.Add(new(PluginHttp))

	gl.Run()
	defer func() {
		gl.StopAndWait()
		time.Sleep(1 * time.Millisecond)
	}()

	if true {
		ctx := gl.NewContext("http")
		ctx.Timeout = 3 * time.Second
		gl.ExecuteSync(ctx)
		// hresp := ctx.Result["http"]
		// if hresp != nil {
		// 	delete(hresp.(map[string]interface{}), "Content")
		// }
		fmt.Println(ctx, ctx.Data, ctx.Result, ctx.Error)
	}
}
