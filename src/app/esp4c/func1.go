package main

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"logger"
	"time"
)

func doAdd(address string) {
	sock := createSocket(address)
	if sock == nil {
		return
	}
	defer sock.AskClose()

	msg := esnp.NewMessage()
	msg.GetAddress().SetCall("test", "add")
	ds := msg.Datas()
	ds.Set("a", 1)
	ds.Set("b", 2)
	rmsg, err := espsocket.Call(sock, msg, 3*time.Second)
	if err != nil {
		logger.Warn(tag, "call 'add' fail - %s", err)
		return
	}
	ds2 := rmsg.Datas()
	res, err2 := ds2.GetInt("c", 0)
	if err2 != nil {
		logger.Warn(tag, "result fail - %s", err2)
		return
	}
	logger.Info(tag, "result = %d", res)
}

func doMAdd(address string) {
	/*
		sock := createSocket(address)
		if sock == nil {
			return
		}
		defer sock.AskClose()

		var f1 *syncutil.Future
		var f2 *syncutil.Future
		if true {
			msg := esnp.NewMessage()
			msg.GetAddress().SetCall("test", "add")
			ds := msg.Datas()
			ds.Set("a", 1)
			ds.Set("b", 2)
			f1 = sock.FutureCall(msg, 3*time.Second)
		}
		if true {
			msg := esnp.NewMessage()
			msg.GetAddress().SetCall("test", "add")
			ds := msg.Datas()
			ds.Set("a", 3)
			ds.Set("b", 4)
			f2 = sock.FutureCall(msg, 3*time.Second)
		}

		fg := syncutil.NewFutureGroup()
		fg.Add(f1)
		fg.Add(f2)

		if !fg.WaitAll(3 * time.Second) {
			logger.Error(tag, "call 'add' fail timeout")
			return
		}

		var r1 int
		var r2 int
		if true {
			_, v, err := f1.Get()
			if err != nil {
				logger.Warn(tag, "call 'add' fail - %s", err)
				return
			}
			rmsg := v.(*esnp.Message)
			ds2 := rmsg.Datas()
			res, err2 := ds2.GetInt("c", 0)
			if err2 != nil {
				logger.Warn(tag, "result fail - %s", err2)
				return
			}
			r1 = int(res)
		}
		if true {
			_, v, err := f2.Get()
			if err != nil {
				logger.Warn(tag, "call 'add' fail - %s", err)
				return
			}
			rmsg := v.(*esnp.Message)
			ds2 := rmsg.Datas()
			res, err2 := ds2.GetInt("c", 0)
			if err2 != nil {
				logger.Warn(tag, "result fail - %s", err2)
				return
			}
			r2 = int(res)
		}
		logger.Info(tag, "result = %d", r1+r2)
	*/
}
