package conndialpool

import (
	"fmt"
	"logger"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(30*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func T2estDialPool(t *testing.T) {

	safeCall()

	cfg := new(DialPoolConfig)
	cfg.Address = "127.0.0.1:80"
	cfg.MaxSize = 3
	cfg.InitSize = 2
	pool := NewDialPool("Pool", cfg)

	pool.StartAndRun()

	time.Sleep(time.Duration(200) * time.Millisecond)

	for i := 1; i <= 4; i++ {
		go func(i int) {
			s, err := pool.GetConn(2*time.Second, true)
			if err != nil {
				logger.Warn("TEST", "%d GetConn fail - %s", i, err)
				return
			}
			if s != nil {
				logger.Info("TEST", "%d %s -> %p", i, pool, s)
				time.AfterFunc(1*time.Second, func() {
					logger.Info("TEST", "%d %s return %p", i, pool, s)
					pool.ReturnConn(s)
				})
			}
		}(i)
	}

	time.Sleep(time.Duration(5) * time.Second)

	logger.Info("TEST", "before close - %s", pool)
	pool.Close()
	time.Sleep(time.Duration(1) * time.Millisecond)
	logger.Info("TEST", "after  close - %s", pool)
}

func TestDialRetry(t *testing.T) {
	safeCall()

	cfg := new(DialPoolConfig)
	cfg.Address = "127.0.0.1:1080"
	cfg.MaxSize = 3
	cfg.InitSize = 2

	pool := NewDialPool("Pool", cfg)
	pool.StartAndRun()

	time.Sleep(time.Duration(5) * time.Second)

	logger.Info("TEST", "before close - %s", pool)
	pool.Close()
	time.Sleep(time.Duration(1) * time.Millisecond)
	logger.Info("TEST", "after  close - %s", pool)
	time.Sleep(time.Duration(1) * time.Millisecond)
}
