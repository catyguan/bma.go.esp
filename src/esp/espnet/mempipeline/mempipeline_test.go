package mempipeline

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"logger"
	"testing"
	"time"
)

func TestMemPipeline(t *testing.T) {
	mp := NewMemPipeline("test", 10)
	ch1 := mp.ChannelA()
	ch2 := mp.ChannelB()

	s1 := espsocket.NewSocket(ch1)
	defer s1.Shutdown()
	s1.SetMessageListner(func(msg *esnp.Message) error {
		logger.Debug("test", "s1 receive")
		time.AfterFunc(1*time.Second, func() {
			logger.Debug("test", "s1 send")
			s1.SendMessage(msg, nil)
		})
		return nil
	})

	s2 := espsocket.NewSocket(ch2)
	defer s2.Shutdown()
	s2.SetMessageListner(func(msg *esnp.Message) error {
		logger.Debug("test", "s2 receive")
		time.AfterFunc(1*time.Second, func() {
			logger.Debug("test", "s2 send")
			s2.SendMessage(msg, nil)
		})
		return nil
	})

	msg := esnp.NewMessage()
	msg.SureId()
	s1.SendMessage(msg, nil)

	time.Sleep(5 * time.Second)
}
