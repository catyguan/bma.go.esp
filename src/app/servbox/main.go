package main

import (
	"bmautil/connutil"
	"boot"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"esp/espnet/secure"
	"esp/servbox"
	"esp/services/servboot"
	"esp/services/servpprof"
	"net"
	"netserver"
)

const (
	tag = "servbox"
)

func main() {
	cfile := "config/servbox-config.json"

	boxs := servbox.NewService("servbox")
	boot.AddService(boxs)

	if true {
		mux := espservice.NewServiceMux(nil, nil)
		servboot.InitMux(mux)
		servpprof.InitMux(mux)
		mux.DefaultHandler = boxs.Handler

		goservice := espservice.NewGoService("mainService", mux.Serve)
		var se espservice.ServiceEntry
		useSecure := false
		if useSecure {
			ba := secure.NewBaseAuthEntry(secure.SimpleAppKeyProvider("123456"), goservice.DoServe)
			se = ba.AuthEntry
		} else {
			se = goservice.Serve
		}

		lisPoint := netserver.NewService("servicePoint", func(conn net.Conn) {
			ct := connutil.NewConnExt(conn)
			// ct.Debuger = connutil.SimpleDebuger(tag)
			sock := espsocket.NewConnSocket(ct, 10*1024*1024)
			se(sock)
		})
		boot.AddService(lisPoint)
	}

	if true {
		lisPoint := netserver.NewService("managePoint", boxs.AcceptManageConn(""))
		boot.AddService(lisPoint)
	}

	boot.Go(cfile)
}
