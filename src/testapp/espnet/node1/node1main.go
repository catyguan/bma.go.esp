package main

import (
	"bmautil/socket"
	"boot"
	"esp/espnet"
	"esp/espnet/espnetutil"
	"logger"
)

func main() {
	cfile := "../test/empty-config.json"

	// espnet
	boot.Define(boot.INIT, "my", myInit)

	boot.Go(cfile)
}

const (
	tag = "N1"
)

var service *espnet.GoService
var service2 *espnet.GoService
var service3 *espnet.GoService

var publisher *espnetutil.Publisher
var hub *espnetutil.Broker

func myInit() bool {

	service = espnet.NewGoService("S1", newHandleRequest("S1"))
	boot.QuickDefine(service, "", true)
	service2 = espnet.NewGoService("S2", newHandleRequest("S2"))
	boot.QuickDefine(service2, "", true)
	service3 = espnet.NewGoService("S3", newHandleRequest("S3"))
	boot.QuickDefine(service3, "", true)

	cfgSer := new(espnet.ListenConfig)
	cfgSer.Port = 1080
	pointSER := espnet.NewListenPoint("PSER", cfgSer, socketAcceptSer)
	boot.QuickDefine(pointSER, "", true)

	publisher = espnetutil.NewPublisher("P1", 16)
	boot.QuickDefine(publisher, "", true)

	cfgPub := new(espnet.ListenConfig)
	cfgPub.Port = 1081
	pointPUB := espnet.NewListenPoint("PPUB", cfgPub, socketAcceptPub)
	boot.QuickDefine(pointPUB, "", true)

	hub = espnetutil.NewBroker("HUB", 16)
	boot.QuickDefine(hub, "", true)

	cfgPos := new(espnet.ListenConfig)
	cfgPos.Port = 1082
	pointPOS := espnet.NewListenPoint("POS", cfgPos, socketAcceptOSer)
	boot.QuickDefine(pointPOS, "", true)

	// ch, _ := service.NewChannel()
	// hub.AddRight(ch, true)
	// hub.AddRight(service2.NewChannel(), true)
	// hub.AddRight(service3.NewChannel(), true)

	return true
}

func socketAcceptSer(sock *socket.Socket) error {
	ch := espnet.NewSocketChannel(sock, espnet.SOCKET_CHANNEL_CODER_ESPNET)
	// ch := espnet.NewSocketChannel(sock, espnet.SOCKET_CHANNEL_CODER_TELNET)
	ch.SetProperty(espnet.PROP_ESPNET_MAXFRAME, 1024)
	ch.SetProperty(espnet.PROP_SOCKET_TRACE, 64)
	// ch.SetProperty(espnet.PROP_SOCKET_TIMEOUT, time.Duration(5)*time.Second)

	// sch := service.NewChannel()
	// p := new(espnet.Pipeline)
	// p.LeftHandler = pubPipelineHandler
	// p.CloseOnBreak = true
	// p.Create(ch, sch)
	hub.AddLeft(ch, true)

	return nil
}

func pubPipelineHandler(in espnet.Channel, msg *espnet.Message, out espnet.Channel) error {
	logger.Debug(tag, "pipeline %s >-> %s", in, out)
	if err := publisher.SendMessage(msg.Clone()); err != nil {
		logger.Debug(tag, "publish send fail - %s", err)
	}
	return out.SendMessage(msg)
}

func newHandleRequest(n string) espnet.ServiceHandler {
	return func(msg *espnet.Message, rep espnet.ServiceResponser) error {
		logger.Info(tag, "%s received request: [%s]", n, msg.Dump())
		bs := msg.GetPayloadB()
		if bs != nil {
			logger.Info(tag, "Payload = %s", string(bs))
		}

		rmsg := espnet.NewReplyMessage(msg)
		rmsg.Datas().Set("name", n)
		if bs != nil {
			rmsg.SetPayload(bs)
		}
		// espnet.FrameCoders.CloseChannel.Set(rmsg.ToPackage())
		rep(rmsg)
		return nil
	}
}

func socketAcceptPub(sock *socket.Socket) error {
	ch := espnet.NewSocketChannel(sock, espnet.SOCKET_CHANNEL_CODER_ESPNET)
	ch.SetProperty(espnet.PROP_ESPNET_MAXFRAME, 1024)
	ch.SetProperty(espnet.PROP_SOCKET_TRACE, 64)

	publisher.Add(ch)
	return nil
}

func socketAcceptOSer(sock *socket.Socket) error {
	ch := espnet.NewSocketChannel(sock, espnet.SOCKET_CHANNEL_CODER_ESPNET)
	ch.SetProperty(espnet.PROP_ESPNET_MAXFRAME, 1024)
	ch.SetProperty(espnet.PROP_SOCKET_TRACE, 64)

	hub.AddRight(ch, true)
	return nil
}
