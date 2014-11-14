package vmmesnp

import (
	"context"
	"fileloader"
	"fmt"
	"golua"
	"os"
	"runtime"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(2*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func TestESNP(t *testing.T) {
	if true {
		runtime.GOMAXPROCS(5)
		safeCall()

		data := make(map[string]interface{})

		dirs := []string{"../../../golua/samplecodes/"}
		sr := new(fileloader.FileFileLoader)
		sr.Dirs = dirs

		gl := golua.NewGoLua("test", 8, sr, func(gl *golua.GoLua) {
			golua.InitCoreLibs(gl)
			InitGoLua(gl)
		}, nil)
		defer gl.Close()

		trace := false
		f := "test_esnp.lua"

		req := golua.NewRequestInfo()
		req.Script = f
		req.Data = data
		req.Trace = trace
		ctx := context.Background()
		ctx, _ = context.CreateExecId(ctx)
		ctx = golua.CreateRequest(ctx, req)

		fmt.Println("--------------- RUN ---------------")
		rval, err4 := gl.Execute(ctx)
		if err4 != nil {
			t.Error("golua call error:", err4)
			time.Sleep(100 * time.Millisecond)
			return
		}
		fmt.Println("Call => ", rval)

		// golua.Execute(ctx)

		time.Sleep(100 * time.Millisecond)
	}
}
