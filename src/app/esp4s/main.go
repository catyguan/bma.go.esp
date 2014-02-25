package main

import (
	"bmautil/socket"
	"boot"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
)

const (
	tag = "esp4n"
)

func main() {
	cfile := "config/esp4s-config.json"

	mux := espservice.NewServiceMux(nil, nil)
	mux.AddHandler("add", H4Add)

	goservice := espservice.NewGoService("service", mux.Serve)

	lisPoint := socket.NewListenPoint("servicePoint", nil, goservice.AcceptESP)
	boot.QuickDefine(lisPoint, "", true)

	boot.Go(cfile)
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
		rmsg := msg.ReplyMessage()
		rmsg.Datas().Set("c", c)
		return rep.SendMessage(rmsg)
	}
	return nil
}
