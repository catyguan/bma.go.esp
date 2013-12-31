package main

import (
	"bmautil/socket"
	"boot"
	"esp/espnet"
	"logger"
	"time"
)

const (
	tag = "NODE2"
)

func main() {
	cfile := "../test/empty-config.json"

	// espnet
	boot.Define(boot.INIT, "myinit", myInit)
	boot.Define(boot.RUN, "stop", func() bool {
		// time.AfterFunc(time.Duration(3)*time.Second, func() { boot.Shutdown() })
		return true
	})

	// boot.Define(boot.CLEANUP, "testTimeout", func() bool {
	// 	time.Sleep(17 * time.Second)
	// 	return true
	// })

	boot.Go(cfile)
}

func myInit() bool {
	cfg := new(espnet.DialConfig)
	cfg.Address = "127.0.0.1:1080"
	dp := espnet.NewDialPoint("N2", cfg, socketAccept)
	dp.RetryFailInfoDruation = 30 * time.Second

	boot.QuickDefine(dp, "N2", true)

	return true
}

func socketAccept(sock *socket.Socket) error {
	ch := espnet.NewSocketChannel(sock, espnet.SOCKET_CHANNEL_CODER_ESPNET)
	ch.SetProperty(espnet.PROP_SOCKET_TRACE, 128)

	// ch.SetPipelineListner(rec)
	ch.SetMessageListner(onMessageReceive)

	msg := espnet.NewRequestMessage()
	msg.SetAddress(espnet.NewAddress("s1"))
	msg.SetSourceAddress(espnet.NewAddress("node2"))
	msg.SetId(1234)
	msg.SetPayload([]byte("hello world"))
	msg.Headers().Set("host", "live.yy.com")
	msg.Datas().Set("uid", 808012345)
	// espnet.FrameCoders.Trace.Set(msg.ToPackage())

	if err := ch.PostEvent(msg); err != nil {
		return nil
	}
	return nil
}

func onMessageReceive(msg *espnet.Message) error {
	logger.Info(tag, "onMessageReceive -- %s", msg.Dump())
	bs := msg.GetPayloadB()
	if bs != nil {
		logger.Info(tag, "########### = %s", string(bs))
	}
	return nil
}
