package smmapi4pprof

import (
	"boot"
	"fmt"
	"logger"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const (
	tag = "smmapi4pprof"
)

func ofile(name string) (*os.File, error) {
	now := time.Now()
	tf := now.Format("20060102_150405")
	fn := fmt.Sprintf("pprof/%s_%s.prof", name, tf)
	ffn, err := boot.TempFile(fn)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(ffn, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}

func doSave(name string) error {
	logger.Info(tag, "%s profile begin", name)

	db := 1
	switch name {
	case "heap":
		runtime.GC()
	case "goroutine":
		db = 2
	}
	f, err0 := ofile(name)
	if err0 != nil {
		return err0
	}
	defer f.Close()
	return pprof.Lookup(name).WriteTo(f, db)
}

func doCPU(f *os.File, sec int) error {
	go func() {
		defer f.Close()

		logger.Info(tag, "cpu profile begin")
		err0 := pprof.StartCPUProfile(f)
		if err0 != nil {
			logger.Warn(tag, "cpu profile error - %s", err0)
			return
		}
		time.Sleep(time.Duration(sec) * time.Second)
		pprof.StopCPUProfile()
		logger.Info(tag, "cpu profile end")
	}()
	return nil
}
