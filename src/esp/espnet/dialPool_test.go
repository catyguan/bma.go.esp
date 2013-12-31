package espnet

import (
	"bmautil/socket"
	"logger"
	"testing"
	"time"
)

func TestDialPool(t *testing.T) {

	cfg := new(DialPoolConfig)
	cfg.Dial.Address = "127.0.0.1:1080"
	cfg.MaxSize = 3
	cfg.InitSize = 2
	pool := NewDialPool("Pool", cfg, socketInitor)

	pool.Start()
	pool.Run()

	time.Sleep(time.Duration(500) * time.Millisecond)

	for i := 0; i < 4; i++ {
		go func() {
			s, _ := pool.GetSocket(1*time.Second, true)
			if s != nil {
				logger.Info("TEST", "%s -> %p", pool, s)
				time.AfterFunc(1*time.Second, func() {
					pool.ReturnSocket(s)
				})
			}
		}()
	}

	time.Sleep(time.Duration(5) * time.Second)

	logger.Info("TEST", "before close - %s", pool)
	pool.Close()
	time.Sleep(time.Duration(1) * time.Millisecond)
	logger.Info("TEST", "after  close - %s", pool)
}

func TestDialRetry(t *testing.T) {

	cfg := new(DialPoolConfig)
	cfg.Dial.Address = "127.0.0.1:1080"
	cfg.MaxSize = 3
	cfg.InitSize = 2
	pool := NewDialPool("Pool", cfg, socketInitor)

	pool.Start()
	pool.Run()

	time.Sleep(time.Duration(5) * time.Second)

	logger.Info("TEST", "before close - %s", pool)
	pool.Close()
	time.Sleep(time.Duration(1) * time.Millisecond)
	logger.Info("TEST", "after  close - %s", pool)
}

func socketInitor(s *socket.Socket) error {
	s.Trace = 128
	return nil
}
