package golua

import (
	"context"
	"fileloader"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	if true {
		runtime.GOMAXPROCS(5)
		safeCall()

		data := make(map[string]interface{})

		dirs := []string{"samplecodes/"}
		sr := new(fileloader.FileFileLoader)
		sr.Dirs = dirs

		golua := NewGoLua("test", sr, func(vmg *VMG) {
			CoreModule(vmg)
			GoModule().Bind(vmg)
			TypesModule().Bind(vmg)
		}, nil)
		defer golua.Close()

		trace := false
		f := "/s_add.lua"
		data["a"] = 1
		data["b"] = 2

		req := new(RequestInfo)
		req.Script = f
		req.Data = data
		req.Trace = trace
		ctx := context.Background()
		ctx, _ = context.CreateExecId(ctx)
		ctx = CreateRequest(ctx, req)

		fmt.Println("--------------- RUN ---------------")
		rval, err4 := golua.Execute(ctx)
		if err4 != nil {
			t.Error("golua call error:", err4)
			return
		}
		fmt.Println("Call => ", rval)

		// golua.Execute(ctx)

		time.Sleep(100 * time.Millisecond)
	}
}
