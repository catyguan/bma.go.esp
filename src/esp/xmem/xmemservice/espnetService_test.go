package xmemservice

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

	msg := espnet.NewMessage()
	msg.SetAddress(espnet.NewAddress("xmem"))
	req := new(xmemprot.SHRequestSet)
	req.InitSet("test", xmemprot.MemKey{"a", "b"}, 1234, 8)
	req.Write(msg)
	fmt.Println("call request")
	rmsg, err2 := cl.Call(msg, time.NewTimer(2*time.Second))
	if err2 != nil {
		t.Error(err2)
		return
	}

	o := new(xmemprot.SHResponseSet)
	o.Read(rmsg)
	fmt.Println("RETURN", o)
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

	msg := espnet.NewMessage()
	msg.SetAddress(espnet.NewAddress("xmem"))
	req := new(xmemprot.SHRequestGet)
	req.Init("test", xmemprot.MemKey{"a", "e"})
	req.Write(msg)
	fmt.Println("call request")
	rmsg, err2 := cl.Call(msg, time.NewTimer(2*time.Second))
	if err2 != nil {
		t.Error(err2)
		return
	}

	o := new(xmemprot.SHResponseGet)
	o.Read(rmsg)
	fmt.Println("RETURN", o)
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
