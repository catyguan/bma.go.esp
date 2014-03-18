package esptunnel

import (
	"bmautil/socket"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"fmt"
	"logger"
	"os"
	"testing"
	"time"
)

func TestTunnel1(t *testing.T) {
	time.AfterFunc(10*time.Second, func() {
		fmt.Println("i die!!!")
		os.Exit(1)
	})

	logger.Info("TEST", "new tunnel")
	tch := NewTunnel("tuch")
	defer tch.Stop()
	tch.SetCloseListener("", func() {
		fmt.Println("tunnel closed")
	})
	tch.SetMessageListner(func(msg *esnp.Message) error {
		fmt.Println("message income", msg.Dump())
		return nil
	})

	if true {
		cfg := new(socket.DialConfig)
		cfg.Address = "127.0.0.1:1080"
		sock, err := socket.Dial("ch1", cfg, nil)
		if err != nil {
			t.Error(err)
			return
		}
		ch := espchannel.NewSocketChannel(sock, espchannel.SOCKET_CHANNEL_CODER_ESPNET)
		tch.Add(ch)
	}

	if true {
		cfg := new(socket.DialConfig)
		cfg.Address = "127.0.0.1:1081"
		sock, err := socket.Dial("ch2", cfg, nil)
		if err != nil {
			t.Error(err)
			return
		}
		ch := espchannel.NewSocketChannel(sock, espchannel.SOCKET_CHANNEL_CODER_ESPNET)
		tch.Add(ch)
	}

	for i := 0; i < 5; i++ {
		logger.Info("TEST", "send message %d", i)
		msg := esnp.NewMessage()
		err := tch.PostMessage(msg)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
	logger.Info("TEST", "exit")
}
