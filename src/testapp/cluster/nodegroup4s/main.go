package main

import (
	"bmautil/socket"
	"boot"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"logger"
)

const (
	tag = "esp4s"
)

func main() {
	mux := espservice.NewServiceMux(nil, nil)
	mux.AddHandler("test", "add", H4Add)
	mux.AddHandler("sys", "reload", H4Reload)

	goservice := espservice.NewGoService("service", mux.Serve)

	lisPoint := socket.NewListenPoint("servicePoint", nil, goservice.AcceptESP)
	boot.Add(lisPoint, "", true)

	boot.Go("")
}

func H4Add(msg *esnp.Message, rep espservice.ServiceResponser) error {
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
		logger.Error(tag, "%d + %d = %d", a, b, c)
		rmsg := msg.ReplyMessage()
		rmsg.Datas().Set("c", c)
		return rep.SendMessage(rmsg)
	}
	return nil
}

func H4Reload(msg *esnp.Message, rep espservice.ServiceResponser) error {
	go func() {
		boot.Restart()
	}()
	return nil
}
