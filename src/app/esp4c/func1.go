package main

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"esp/espnet/secure"
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
	rmsg, err := espsocket.CallTimeout(sock, msg, 3*time.Second)
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

func doSAdd(address string) {
	sock := createSocket(address)
	if sock == nil {
		return
	}
	defer sock.AskClose()

	err0 := secure.DoBaseAuth(sock, "test", "123456", 5*time.Second)
	if err0 != nil {
		logger.Warn(tag, "BaseAuth fail - %s", err0)
		return
	}

	msg := esnp.NewMessage()
	msg.GetAddress().SetCall("test", "add")
	ds := msg.Datas()
	ds.Set("a", 1)
	ds.Set("b", 2)
	rmsg, err := espsocket.CallTimeout(sock, msg, 3*time.Second)
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
