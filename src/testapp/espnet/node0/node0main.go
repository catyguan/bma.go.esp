package main

import (
	"bmautil/socket"
	"boot"
	"esp/espnet"
	"logger"
	"time"
)

func main() {
	cfile := "../test/empty-config.json"

	// espnet
	boot.Define(boot.INIT, "my", myInit)

	boot.Go(cfile)
}

const (
	tag = "N0"
)

var service *espnet.GoService

func myInit() bool {

	smux := espnet.NewServiceMux(nil, nil)
	smux.AddHandler(espnet.NewAddress("s1"), newHandleRequest("S1"))
	smux.AddHandler(espnet.NewAddress("s2"), newHandleRequest("S2"))
	smux.AddHandler(espnet.NewAddress("s3"), newHandleRequest("S3"))

	service = espnet.NewGoService("SERVICE", smux.Serve)
	boot.QuickDefine(service, "", true)

	cfg := new(espnet.DialConfig)
	cfg.Address = "127.0.0.1:1082"
	dp := espnet.NewDialPoint("N0", cfg, socketAcceptSer)
	dp.RetryFailInfoDruation = 30 * time.Second

	boot.QuickDefine(dp, "N0", true)

	return true
}

func socketAcceptSer(sock *socket.Socket) error {
	ch := espnet.NewSocketChannel(sock, espnet.SOCKET_CHANNEL_CODER_ESPNET)
	ch.SetProperty(espnet.PROP_ESPNET_MAXFRAME, 1024)
	ch.SetProperty(espnet.PROP_SOCKET_TRACE, 64)

	espnet.ConnectService(ch, service.PostRequest)

	return nil
}

func newHandleRequest(n string) espnet.ServiceHandler {
	return func(msg *espnet.Message, rep espnet.ServiceResponser) error {
		logger.Info(tag, "%s received request: [%s]", n, msg.Dump())
		bs := msg.GetPayloadB()
		if bs != nil {
			logger.Info(tag, "Payload = %s", string(bs))
		}

		rmsg := espnet.NewReplyMessage(msg)
		rmsg.Datas().Set("node0", true)
		if bs != nil {
			rmsg.SetPayload(bs)
		}
		rep(rmsg)
		return nil
	}
}
