package esppchannel

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

func TestPChannel(t *testing.T) {

	time.AfterFunc(10*time.Second, func() {
		fmt.Println("i die!!!")
		os.Exit(1)
	})

	logger.Info("TEST", "new pchannel")
	pch := NewPChannel("tpch")
	defer pch.Stop()
	pch.SetCloseListener("", func() {
		fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	})

	if true {
		cfg := new(socket.DialPoolConfig)
		cfg.Dial.Address = "127.0.0.1:1081"
		cfg.MaxSize = 1
		cfg.InitSize = 1
		pool := socket.NewDialPool("Pool", cfg, nil)
		pool.Start()
		pool.Run()
		defer func() {
			pool.AskClose()
		}()
		time.Sleep(100 * time.Millisecond)

		cf := espchannel.NewDialPoolChannelFactory(pool, espchannel.SOCKET_CHANNEL_CODER_ESPNET, 1*time.Second)
		pch.Add(cf)
	}

	if true {
		cfg := new(socket.DialPoolConfig)
		cfg.Dial.Address = "127.0.0.1:1080"
		cfg.MaxSize = 1
		cfg.InitSize = 1
		pool := socket.NewDialPool("Pool", cfg, nil)
		pool.Start()
		pool.Run()
		defer func() {
			pool.AskClose()
		}()
		time.Sleep(100 * time.Millisecond)

		cf := espchannel.NewDialPoolChannelFactory(pool, espchannel.SOCKET_CHANNEL_CODER_ESPNET, 1*time.Second)
		pch.Add(cf)
	}

	pch.OnReady()

	for i := 0; i < 6; i++ {
		logger.Info("TEST", "send message %d", i)
		msg := esnp.NewMessage()
		if i%2 == 1 {
			// CloseForce(pch)
			// CloseAfterSend(msg)
		}
		err := pch.SendMessage(msg)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
	logger.Info("TEST", "exit")
}
