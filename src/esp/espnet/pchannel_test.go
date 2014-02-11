package espnet

import (
	"boot"
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
	defer pch.Close()

	if true {
		cfg := new(DialPoolConfig)
		cfg.Dial.Address = "127.0.0.1:1081"
		cfg.MaxSize = 1
		cfg.InitSize = 1
		pool := NewDialPool("Pool", cfg, socketInitor)
		boot.RuntimeStartRun(pool)
		defer boot.RuntimeStopCloseClean(pool, true)
		time.Sleep(100 * time.Millisecond)

		cf := pool.NewChannelFactory("espnet", 1*time.Second)
		pch.Add(cf)
	}

	if true {
		cfg := new(DialPoolConfig)
		cfg.Dial.Address = "127.0.0.1:1080"
		cfg.MaxSize = 1
		cfg.InitSize = 1
		pool := NewDialPool("Pool", cfg, socketInitor)
		boot.RuntimeStartRun(pool)
		defer boot.RuntimeStopCloseClean(pool, true)
		time.Sleep(100 * time.Millisecond)

		cf := pool.NewChannelFactory("espnet", 1*time.Second)
		pch.Add(cf)
	}

	pch.OnReady()

	for i := 0; i < 5; i++ {
		logger.Info("TEST", "send message %d", i)
		msg := NewMessage()
		err := pch.SendMessage(msg)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
	logger.Info("TEST", "exit")
}
