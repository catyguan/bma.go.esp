package router

import (
	"bmautil/socket"
	"errors"
	"esp/espnet"
	"logger"
	"testing"
	"time"
)

type testDispatcher struct {
	fac espnet.ChannelFactory
}

func (this *testDispatcher) Dispatch(n string) (espnet.ChannelFactory, error) {
	logger.Info("testDispatcher", "call dispatch(%s)", n)
	if this.fac != nil {
		return this.fac, nil
	}
	return nil, errors.New("test no return!!!")
}
func (this *testDispatcher) Close() {
	logger.Info("testDispatcher", "call close")
}

func TestRouter(t *testing.T) {

	// pool
	cfg := new(espnet.DialConfig)
	cfg.Address = "127.0.0.1:1080"
	pool := espnet.NewDialPool("Pool", cfg, 3, socketInitor)
	cpool := pool.NewChannelFactory(time.Second)

	// router
	router := NewRouter("router")
	router.InitFactory("s1", cpool)
	router.InitDispatcher(&testDispatcher{cpool})

	pool.Start()
	router.Start()

	pool.Run()

	if true {
		ch, err := router.GetChannel("s0", time.Second)
		logger.Info("FUCK", "GetChannel -> %v, %v", ch, err)
	}

	if true {
		ch, err := router.GetChannel("s0", time.Second)
		logger.Info("FUCK", "GetChannel -> %v, %v", ch, err)
	}

	if true {
		ch, err := router.GetChannel("s1", time.Second)
		logger.Info("FUCK", "GetChannel -> %v, %v", ch, err)
	}

	time.Sleep(time.Duration(1) * time.Second)

	router.Stop()

	logger.Info("TEST", "before close - %s", pool)
	pool.Close()
	time.Sleep(time.Duration(1) * time.Millisecond)
	logger.Info("TEST", "after  close - %s", pool)

	router.Cleanup()
}

func socketInitor(s *socket.Socket) error {
	s.Trace = 128
	return nil
}
