package smmapi4pprof

import (
	"bmautil/valutil"
	"fmt"
	"logger"
	"smmapi"
)

type smmObject int

func s(v uint64) string {
	return valutil.SizeString(v, 1024, valutil.SizeM)
}

func (this smmObject) GetInfo() (*smmapi.SMInfo, error) {
	r := new(smmapi.SMInfo)
	r.Title = "PProf"
	r.Content = ""

	r.Actions = make([]*smmapi.SMAction, 0)
	if true {
		a := new(smmapi.SMAction)
		a.Id = "pprof.heap"
		a.Title = "Create Heap PProf"
		a.Type = smmapi.SMA_API
		r.Actions = append(r.Actions, a)
	}
	if true {
		a := new(smmapi.SMAction)
		a.Id = "pprof.thread"
		a.Title = "Create Thread PProf"
		a.Type = smmapi.SMA_API
		r.Actions = append(r.Actions, a)
	}
	if true {
		a := new(smmapi.SMAction)
		a.Id = "boot.block"
		a.Title = "Create Block PProf"
		a.Type = smmapi.SMA_API
		r.Actions = append(r.Actions, a)
	}
	if true {
		a := new(smmapi.SMAction)
		a.Id = "boot.cpu"
		a.Title = "Create CPU PProf"
		a.Type = smmapi.SMA_API
		r.Actions = append(r.Actions, a)
	}
	if true {
		a := new(smmapi.SMAction)
		a.Id = "boot.gor"
		a.Title = "Create GoRoutines PProf"
		a.Type = smmapi.SMA_API
		r.Actions = append(r.Actions, a)
	}

	return r, nil
}

func (this smmObject) ExecuteAction(aid string, param map[string]interface{}) (interface{}, error) {
	msg := "OK"
	switch aid {
	case "pprof.heap":
		err := doSave("heap")
		if err != nil {
			return "", err
		}
		return msg, nil
	case "pprof.thread":
		err := doSave("threadcreate")
		if err != nil {
			return "", err
		}
		return msg, nil
	case "pprof.block":
		err := doSave("block")
		if err != nil {
			return "", err
		}
		return msg, nil
	case "pprof.gor":
		err := doSave("goroutine")
		if err != nil {
			return "", err
		}
		return msg, nil
	case "pprof.cpu":
		logger.Info(tag, "pprof cpu from smmapi")
		sec := 30
		f, errf := ofile("cpu")
		if errf != nil {
			return "", errf
		}
		err0 := doCPU(f, int(sec))
		if err0 != nil {
			return "", err0
		}
		return msg, nil
	}
	return nil, fmt.Errorf("unknow action(%s)", aid)
}

func init() {
	smmapi.Add("go.pprof", smmObject(0))
}
