package xmemclient

import (
	"esp/espnet"
	"esp/xmem/xmemprot"
	"fmt"
	"testing"
	"time"
)

func TestSet4Client(t *testing.T) {
	cl := espnet.NewChannelClient()
	cfg := new(espnet.DialConfig)
	cfg.Address = "127.0.0.1:8080"
	err := cl.Dial("test", cfg, "espnet")
	if err != nil {
		t.Error(err)
		return
	}
	defer cl.Close()
	defer fmt.Println("end")

	xm := NewClient(cl, espnet.NewAddress("xmem"), "test")
	fmt.Println(xm.Set(xmemprot.MemKey{"a", "c"}, 234, 8))
}

func TestGet4Client(t *testing.T) {
	cl := espnet.NewChannelClient()
	cfg := new(espnet.DialConfig)
	cfg.Address = "127.0.0.1:8080"
	err := cl.Dial("test", cfg, "espnet")
	if err != nil {
		t.Error(err)
		return
	}
	defer cl.Close()
	defer fmt.Println("end")

	xm := NewClient(cl, espnet.NewAddress("xmem"), "test")
	fmt.Println(xm.Get(xmemprot.MemKey{"a", "e"}))
}

func TestList4Client(t *testing.T) {
	cl := espnet.NewChannelClient()
	cfg := new(espnet.DialConfig)
	cfg.Address = "127.0.0.1:8080"
	err := cl.Dial("test", cfg, "espnet")
	if err != nil {
		t.Error(err)
		return
	}
	defer cl.Close()
	defer fmt.Println("end")

	xm := NewClient(cl, espnet.NewAddress("xmem"), "test")
	fmt.Println(xm.List(xmemprot.MemKey{"a"}))
}

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
			o := new(xmemprot.SHEventBinlog)
			o.Read(msg)
			fmt.Println("MESSAGE", o)
		}
		return nil
	})

	msg := espnet.NewMessage()
	msg.SetAddress(espnet.NewAddress("xmem"))
	req := new(xmemprot.SHRequestSlaveJoin)
	req.Group = "test"
	req.Version = 0
	req.Write(msg)
	fmt.Println("send request")
	cl.SendMessage(msg)

	time.Sleep(2 * time.Second)
	fmt.Println("end")
}
