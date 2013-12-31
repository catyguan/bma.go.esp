package main

import (
	"bmautil/socket"
	"boot"
	"esp/espnet"
	"flag"
	"logger"
	"time"
)

const (
	tag = "NODE3"
)

var (
	hostId string
)

func main() {
	cfile := "../test/empty-config.json"

	flag.Parse()
	if flag.NArg() > 0 {
		hostId = flag.Arg(0)
	} else {
		hostId = time.Now().Format("150405")
	}

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
	cfg.Address = "127.0.0.1:1081"
	dp := espnet.NewDialPoint("N3", cfg, socketAccept)
	dp.RetryFailInfoDruation = 30 * time.Second

	boot.QuickDefine(dp, "N3", true)

	return true
}

func socketAccept(sock *socket.Socket) error {
	ch := espnet.NewSocketChannel(sock, espnet.SOCKET_CHANNEL_CODER_ESPNET)
	ch.SetProperty(espnet.PROP_SOCKET_TRACE, 128)

	ch.SetMessageListner(onMessageReceive)

	logger.Info(tag, "WAITING message from %s", sock)

	return nil
}

func onMessageReceive(msg *espnet.Message) error {
	logger.Info(tag, "--- %s --------- %s", hostId, msg.Dump())
	return nil
}
