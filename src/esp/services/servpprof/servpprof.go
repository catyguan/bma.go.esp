package servpprof

import (
	"boot"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const (
	tag = "servpprof"

	NAME_SERVICE   = "pprof"
	NAME_OP_HEAP   = "heap"
	NAME_OP_THREAD = "thread"
	NAME_OP_BLOCK  = "block"
	NAME_OP_CPU    = "cpu"
	NAME_OP_GOR    = "gor"
)

func InitMux(mux *espservice.ServiceMux) {
	mux.AddHandler(NAME_SERVICE, NAME_OP_HEAP, ServOP_Heap)
	mux.AddHandler(NAME_SERVICE, NAME_OP_THREAD, ServOP_Thread)
	mux.AddHandler(NAME_SERVICE, NAME_OP_BLOCK, ServOP_Block)
	mux.AddHandler(NAME_SERVICE, NAME_OP_GOR, ServOP_GOR)
	mux.AddHandler(NAME_SERVICE, NAME_OP_CPU, ServOP_CPU)
}

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

func ServOP_Heap(sock espsocket.Socket, msg *esnp.Message) error {
	logger.Info(tag, "op heap from %s", sock)
	err := doSave("heap")
	if err != nil {
		return err
	}
	sock.WriteMessage(msg.ReplyMessage())
	return nil
}

func ServOP_GOR(sock espsocket.Socket, msg *esnp.Message) error {
	logger.Info(tag, "op gor from %s", sock)
	err := doSave("goroutine")
	if err != nil {
		return err
	}
	sock.WriteMessage(msg.ReplyMessage())
	return nil
}

func ServOP_Thread(sock espsocket.Socket, msg *esnp.Message) error {
	logger.Info(tag, "op thread from %s", sock)
	err := doSave("threadcreate")
	if err != nil {
		return err
	}
	sock.WriteMessage(msg.ReplyMessage())
	return nil
}

func ServOP_Block(sock espsocket.Socket, msg *esnp.Message) error {
	logger.Info(tag, "op block from %s", sock)
	err := doSave("block")
	if err != nil {
		return err
	}
	sock.WriteMessage(msg.ReplyMessage())
	return nil
}

func ServOP_CPU(sock espsocket.Socket, msg *esnp.Message) error {
	logger.Info(tag, "op cpu from %s", sock)
	sec, _ := msg.Datas().GetInt("sec", 0)
	if sec == 0 {
		sec = 30
	}

	f, errf := ofile("cpu")
	if errf != nil {
		return errf
	}
	err0 := doCPU(f, int(sec))
	if err0 != nil {
		return err0
	}
	sock.WriteMessage(msg.ReplyMessage())
	return nil
}
