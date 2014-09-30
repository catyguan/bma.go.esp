package acclog

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	runtime.GOMAXPROCS(1)

	fnc := NewFilenameCreator("test")
	cfg := make(map[string]string)
	f := NewFile(fnc, 16, true, cfg)
	defer f.Close()

	fmt.Println("start")
	f.Write(NewSimpleLog("hello world"))
	dt := make(map[string]interface{})
	dt["a"] = 1
	dt["b"] = "asdad"
	dt["c"] = time.Now()
	f.Write(NewCommonLog(dt, 100))
	fmt.Println("end")
}
