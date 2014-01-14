package xmem

import (
	"esp/espnet"
	"fmt"
	"testing"
	"time"
)

func TestSlaveJoin4Client(t *testing.T) {
	cl := espnet.NewChannelClient()
	cfg := new(espnet.DialConfig)
	cfg.Address = "127.0.0.1:8080"
	err := cl.Dial("test", cfg, "espnet")
	if err != nil {
		t.Error(err)
		return
	}
	defer cl.Close()

	cl.SetMessageListner(func(msg *espnet.Message) error {
		err := msg.ToError()
		if err != nil {
			t.Error(err)
		} else {
			o := new(SHEventBinlog)
			o.Read(msg)
			fmt.Println("MESSAGE", o)
		}
		return nil
	})

	msg := espnet.NewMessage()
	msg.SetAddress(espnet.NewAddress("xmem"))
	req := new(SHRequestSlaveJoin)
	req.Group = "test"
	req.Version = 0
	req.Write(msg)
	fmt.Println("send request")
	cl.SendMessage(msg)

	time.Sleep(2 * time.Second)
	fmt.Println("end")
}
