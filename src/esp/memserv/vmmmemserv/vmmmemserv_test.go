package vmmmemserv

import (
	"context"
	"esp/memserv"
	"fileloader"
	"fmt"
	"golua"
	"os"
	"runtime"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func TestMemServ(t *testing.T) {
	if true {
		runtime.GOMAXPROCS(5)
		safeCall()

		data := make(map[string]interface{})

		s := memserv.NewMemoryServ()

		dirs := []string{"../../../golua/samplecodes/"}
		sr := new(fileloader.FileFileLoader)
		sr.Dirs = dirs

		gl := golua.NewGoLua("test", 8, sr, func(gl *golua.GoLua) {
			golua.InitCoreLibs(gl)
			InitGoLua(gl, s)
		}, nil)
		defer gl.Close()

		trace := false
		f := "test_memserv.lua"

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
			return
		}
		fmt.Println("Call => ", rval)

		// golua.Execute(ctx)

		time.Sleep(100 * time.Millisecond)
	}
}
