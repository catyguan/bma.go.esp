package logger

import (
	"config"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

func T2estLoggerBase(t *testing.T) {
	runtime.GOMAXPROCS(1)

	cfg := Config()
	cfg.InitLogger()

	Debug("test", "hello %s", "world")
	cfg.SetLevel("test", LEVEL_INFO)
	Debug("test", "i can't see DEBUG")
	Warn("test", "i can see WARN")

	fmt.Println("test2 enabled?", EnableDebug("test2"))
}

func TestLoggerInit(t *testing.T) {
	config.InitGlobalConfig("../../test/esp-config.json")

	cfg := Config()
	cfg.InitLogger()

	Debug("test", "i can't see DEBUG")
	Info("test", "i can see INFO")
	Debug("test2", "i can see DEBUG because disabled")
	Debug("test3", "i can see DEBUG because disabled")
	Debug("test4", "i can see DEBUG")

	time.Sleep(100 * time.Millisecond)
}

func T2estRotateFile(t *testing.T) {
	var wd, _ = os.Getwd()
	fn := wd + "/../../test/rotate"

	w := NewRotateFile(func(tm time.Time, num int) string {
		f := fn + "_" + tm.Format("20060102")
		if num != 0 {
			f += "." + fmt.Sprintf("%d", num)
		}
		f += ".log"
		return f
	}, true)
	if w == nil {
		t.Error("open rotateFile fail")
		return
	}
	defer w.Close()
	if !w.Write(0, "it's a test\n") {
		t.Error("write file fail")
		return
	}
	if !w.Rotate() {
		t.Error("rotate file fail")
		return
	}
	if !w.Write(0, "it's a test\n") {
		t.Error("write file fail")
		return
	}
}
