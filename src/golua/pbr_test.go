package golua

import (
	"fmt"
	"golua/goyacc"
	"io/ioutil"
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

func TestParserBuildRun(t *testing.T) {
	if true {
		runtime.GOMAXPROCS(5)
		safeCall()

		trace := false
		// f := "test1.lua"
		// f := "test_go_syn.lua"
		// f := "test_vmmGo.lua"
		f := "test_vmmTypes.lua"
		bs, err0 := ioutil.ReadFile("samplecodes/" + f)
		if err0 != nil {
			t.Error(err0)
			return
		}
		content := string(bs)

		// s = "a = 1"
		// s = "obj:print(1 + 2, true, a.b)"
		// s = "a.b = 1 + 2 - 3"
		// s = "function a(b, c) end"

		chunkName := f

		p := goyacc.NewParser(chunkName, content)
		node, err := p.Parse()
		if err != nil {
			fmt.Println(content)
			t.Error(err)
			return
		}
		fmt.Println("------------NODE---------------")
		fmt.Println(goyacc.DumpNode("", node))

		fmt.Println("--------------- RUN ---------------")
		vmg := NewVMG("test")
		CoreModule(vmg)
		GoModule().Bind(vmg)
		TypesModule().Bind(vmg)
		defer vmg.Close()

		chunk := NewChunk(chunkName, node)

		vm, err3 := vmg.CreateVM()
		if err3 != nil {
			t.Error("create vm error", err3)
			return
		}
		vm.EnableTrace(trace)
		vm.ResetExecutionTime()
		vm.SetMaxExecutionTime(100)

		vm.API_push(chunk)
		_, err4 := vm.Call(0, 1)
		if err4 != nil {
			t.Error("vm call error:", err4)
			return
		}
		fmt.Println(vm.DumpStack())
		rval, err5 := vm.API_pop1X(-1, true)
		if err5 != nil {
			t.Error("pop error:", err5)
			return
		}
		rval, err5 = vm.API_value(rval)
		if err5 != nil {
			t.Error("value error:", err5)
			return
		}
		fmt.Println("Call => ", rval)
		fmt.Println(vmg.globals)
		time.Sleep(100 * time.Millisecond)
	}
}
