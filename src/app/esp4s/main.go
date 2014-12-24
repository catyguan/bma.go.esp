package main

import (
	"bmautil/connutil"
	"boot"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"esp/espnet/secure"
	"esp/servbox"
	"esp/services/servboot"
	"fmt"
	"logger"
	"net"
	"netserver"
	"time"
)

const (
	tag = "esp4s"
)

func main() {
	cfile := "config/esp4s-config.json"

	mux := espservice.NewServiceMux(nil, nil)
	mux.AddHandler("test", "add", H4Add)
	mux.AddHandler("serviceCall", "hello", H4SC)
	servboot.InitMux(mux)

	goservice := espservice.NewGoService("service", mux.Serve)

	var se espservice.ServiceEntry
	useSecure := true
	if useSecure {
		ba := secure.NewBaseAuthEntry(secure.SimpleAppKeyProvider("123456"), goservice.DoServe)
		se = ba.AuthEntry
	} else {
		se = goservice.Serve
	}

	lisPoint := netserver.NewService("servicePoint", func(conn net.Conn) {
		ct := connutil.NewConnExt(conn)
		ct.Debuger = connutil.SimpleDebuger(tag)
		sock := espsocket.NewConnSocket(ct, 1024*1024)
		se(sock)
	})
	if lisPoint != nil {
		// boot.AddService(lisPoint)
	}

	boxc := servbox.NewClient("servboxClient", fmt.Sprintf("%d", time.Now().Unix()), goservice.Serve)
	boxc.Add("test")
	boxc.Add("serviceCall")
	boot.AddService(boxc)

	boot.Go(cfile)
}

func H4Add(sock espsocket.Socket, msg *esnp.Message) error {
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
		return sock.WriteMessage(rmsg)
	}
	return nil
}

func H4SC(sock espsocket.Socket, msg *esnp.Message) error {
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
		return sock.WriteMessage(rmsg)
	}
	return nil
}
