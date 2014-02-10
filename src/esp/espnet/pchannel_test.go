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

	time.AfterFunc(5*time.Second, func() {
		fmt.Println("i die!!!")
		os.Exit(1)
	})

	cfg := new(DialPoolConfig)
	cfg.Dial.Address = "127.0.0.1:1080"
	cfg.MaxSize = 1
	cfg.InitSize = 1
	pool := NewDialPool("Pool", cfg, socketInitor)
	boot.RuntimeStartRun(pool)
	defer boot.RuntimeStopCloseClean(pool, true)
	time.Sleep(500 * time.Millisecond)

	cf := pool.NewChannelFactory("espnet", 1*time.Second)
	ch, err := cf.NewChannel()
	if err != nil {
		t.Error(err)
		return
	}
	// ch.AskClose()
	// ch.AskClose()

	logger.Info("TEST", "new pchannel %s", pool)
	pch := NewPChannel("tpch", cf, ch)
	defer pch.Close()

	time.Sleep(1 * time.Second)
	logger.Info("TEST", "do close channel %s", pool)
	ch.AskClose()
	ch.AskClose()
	logger.Info("TEST", "after close channel %s", pool)

	time.Sleep(1 * time.Second)
	logger.Info("TEST", "exit %s", pool)
}
