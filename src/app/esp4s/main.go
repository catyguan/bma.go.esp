package main

import (
	"bmautil/socket"
	"boot"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"logger"
)

const (
	tag = "esp4s"
)

func main() {
	cfile := "config/esp4s-config.json"

	mux := espservice.NewServiceMux(nil, nil)
	mux.AddHandler("test", "add", H4Add)
	mux.AddHandler("sys", "reload", H4Reload)
	mux.AddHandler("serviceCall", "hello", H4SC)

	goservice := espservice.NewGoService("service", mux.Serve)

	lisPoint := socket.NewListenPoint("servicePoint", nil, goservice.AcceptESP)
	boot.Add(lisPoint, "", true)

	boot.Go(cfile)
}

func H4Add(sock *espsocket.Socket, msg *esnp.Message) error {
	ds := msg.Datas()
	if true {
		a, err1 := ds.GetInt("a", 0)
		if err1 != nil {
			return err1
		}
		b, err2 := ds.GetInt("b", 0)
		if err2 != nil {
			return err2
		}
		c := int(a + b)
		logger.Info(tag, "%d + %d = %d", a, b, c)
		rmsg := msg.ReplyMessage()
		rmsg.Datas().Set("c", c)
		return sock.SendMessage(rmsg, nil)
	}
	return nil
}

func H4SC(sock *espsocket.Socket, msg *esnp.Message) error {
	ds := msg.Datas()
	if true {
		params, err1 := ds.Get("p")
		if err1 != nil {
			return err1
		}
		if params == nil {
			logger.Info(tag, "params empty")
		} else {
			logger.Info(tag, "hello %v", params)
		}

		rmsg := msg.ReplyMessage()
		dt := rmsg.Datas()
		dt.Set("s", 200)
		dt.Set("r", true)
		return sock.SendMessage(rmsg, nil)
	}
	return nil
}

func H4Reload(sock *espsocket.Socket, msg *esnp.Message) error {
	go func() {
		boot.Restart()
	}()
	return nil
}
