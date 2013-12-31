package main

import (
	"bmautil/socket"
	"boot"
	"esp/espnet"
	"esp/protauth/authservice"
)

func main() {
	cfile := "config/seed-config.json"

	// espnet
	// service := seedservice.NewSeedService("seed")
	// service.DefaultBoot()

	serveMux := espnet.NewServiceMux(nil, nil)

	if true {
		executor := authservice.NewSimpleAuthExecutor("auth")
		boot.QuickDefine(executor, "", true)
		serveMux.AddHandler(espnet.NewAddress("auth"), authservice.NewAuthServiceHandler(executor).Serve)
	}

	service := espnet.NewGoService("service", serveMux.Serve)
	boot.QuickDefine(service, "", true)

	pointSER := espnet.NewListenPoint("servicePoint", nil, socketAcceptSer)
	boot.QuickDefine(pointSER, "", true)

	boot.Go(cfile)
}

func socketAcceptSer(sock *socket.Socket) error {
	ch := espnet.NewSocketChannel(sock, espnet.SOCKET_CHANNEL_CODER_ESPNET)
	// ch.SetProperty(espnet.PROP_ESPNET_MAXFRAME, 10*1024*10)
	// ch.SetProperty(espnet.PROP_SOCKET_TRACE, 64)
	// ch.SetProperty(espnet.PROP_SOCKET_TIMEOUT, time.Duration(5)*time.Second)
	obj := boot.ObjectFor("service")
	s := obj.(*espnet.GoService)
	espnet.ConnectService(ch, s.PostRequest)
	return nil
}
