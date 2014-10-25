package vmmsql

import (
	"context"
	"fileloader"
	"fmt"
	"golua"
	"os"
	"runtime"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func safeCall() {
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func TestSQL(t *testing.T) {
	if true {
		runtime.GOMAXPROCS(5)
		safeCall()

		data := make(map[string]interface{})

		dirs := []string{"../samplecodes/"}
		sr := new(fileloader.FileFileLoader)
		sr.Dirs = dirs

		gl := golua.NewGoLua("test", sr, func(vmg *golua.VMG) {
			golua.CoreModule(vmg)
			golua.GoModule().Bind(vmg)
			golua.TypesModule().Bind(vmg)
			Module().Bind(vmg)
		}, nil)
		defer gl.Close()

		trace := false
		f := "test_vmmSQL.lua"

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
