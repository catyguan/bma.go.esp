package golua

import (
	"boot"
	"config"
	"context"
	"fileloader"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func T2estExecute(t *testing.T) {
	if true {
		runtime.GOMAXPROCS(5)
		safeCall()

		config.InitGlobalConfig("../../bin/config/glserver-config.json")

		data := make(map[string]interface{})

		dirs := []string{"samplecodes/"}
		sr := new(fileloader.FileFileLoader)
		sr.Dirs = dirs

		gl := NewGoLua("test", 10, sr, func(gl *GoLua) {
			InitCoreLibs(gl)
		}, nil)
		defer func() {
			gl.Close()
			time.Sleep(100 * time.Millisecond)
		}()

		trace := false
		// f := "/s_add.lua"
		f := "test_vmmGo.lua"
		// f := "test_vmmConfig.lua"
		// f := "test_vmmStrings.lua"
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
		rval, err4 := gl.Execute(ctx)
		if err4 != nil {
			t.Error("gl call error:", err4)
			return
		}
		fmt.Println("Call => ", rval)

		// gl.Execute(ctx)

		time.Sleep(100 * time.Millisecond)
	}
}

func TestDebugger(t *testing.T) {
	if true {
		boot.DevMode = true
		runtime.GOMAXPROCS(5)
		safeCall()

		config.InitGlobalConfig("../../bin/config/glserver-config.json")

		data := make(map[string]interface{})

		dirs := []string{"samplecodes/"}
		sr := new(fileloader.FileFileLoader)
		sr.Dirs = dirs

		gl := NewGoLua("test", 10, sr, func(gl *GoLua) {
			InitCoreLibs(gl)
		}, nil)
		defer func() {
			gl.Close()
			time.Sleep(100 * time.Millisecond)
		}()
		bp := new(Breakpoint)
		bp.chunkName = "test_debugger.lua"
		bp.line = 9
		gl.AddBreakpoint(bp)

		trace := true
		f := "test_debugger.lua"
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
		time.AfterFunc(100*time.Millisecond, func() {
			dg := gl.GetDebugger(1)
			if dg == nil {
				fmt.Println("FUCK, DEBUGGER is nil")
				return
			}
			// dg.DoRun()
			dg.DoStepOut()
			time.Sleep(100 * time.Millisecond)
			dg.DoStep()
		})
		rval, err4 := gl.Execute(ctx)
		if err4 != nil {
			t.Error("gl call error:", err4)
			return
		}
		fmt.Println("Call => ", rval)

		// gl.Execute(ctx)

		time.Sleep(500 * time.Millisecond)
	}
}
