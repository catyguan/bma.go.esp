package main

import (
	"bmautil/syncutil"
	"esp/espnet/esnp"
	"logger"
	"time"
)

func doAdd(address string) {
	c := createClient(address)
	if c == nil {
		return
	}
	defer c.Close()

	msg := esnp.NewMessage()
	msg.GetAddress().Set(esnp.ADDRESS_OP, "add")
	ds := msg.Datas()
	ds.Set("a", 1)
	ds.Set("b", 2)
	rmsg, err := c.Call(msg, nil)
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
	c := createClient(address)
	if c == nil {
		return
	}
	defer c.Close()

	var f1 *syncutil.Future
	var f2 *syncutil.Future
	if true {
		msg := esnp.NewMessage()
		msg.GetAddress().Set(esnp.ADDRESS_OP, "add")
		ds := msg.Datas()
		ds.Set("a", 1)
		ds.Set("b", 2)
		f1 = c.FutureCall(msg)
	}
	if true {
		msg := esnp.NewMessage()
		msg.GetAddress().Set(esnp.ADDRESS_OP, "add")
		ds := msg.Datas()
		ds.Set("a", 3)
		ds.Set("b", 4)
		f2 = c.FutureCall(msg)
	}

	fg := syncutil.NewFutureGroup()
	fg.Add(f1)
	fg.Add(f2)

	if !fg.WaitAll(1 * time.Second) {
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
}
