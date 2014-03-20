package espclient

import (
	"bmautil/socket"
	"bmautil/syncutil"
	"errors"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"esp/espnet/espterminal"
	"time"
)

type ChannelClient struct {
	C   espchannel.Channel
	own bool

	tm espterminal.Terminal
}

func NewChannelClient() *ChannelClient {
	r := new(ChannelClient)
	r.tm.InitTerminal("channelClient")
	return r
}

func (this *ChannelClient) ConnectSocket(sock *socket.Socket, coderName string, own bool) error {
	r := espchannel.NewSocketChannel(sock, coderName)
	return this.Connect(r, own)
}

func (this *ChannelClient) Connect(ch espchannel.Channel, own bool) error {
	this.C = ch
	this.own = own
	this.tm.SetName(ch.String())
	this.tm.Connect(ch)
	return nil
}

func (this *ChannelClient) Dial(name string, cfg *socket.DialConfig, coderName string) error {
	sock, err := socket.Dial(name, cfg, nil)
	if err != nil {
		return err
	}
	return this.ConnectSocket(sock, coderName, true)
}

func (this *ChannelClient) Close() {
	if this.C != nil {
		if this.own {
			this.C.AskClose()
		}
		this.C.SetMessageListner(nil)
		this.C = nil
	}
	this.tm.Close()
}

func (this *ChannelClient) IsOpen() bool {
	if this.C == nil {
		return false
	}
	if cb, ok := this.C.(espchannel.BreakSupport); ok {
		return cb.IsBreak()
	}
	return true
}

func (this *ChannelClient) SetMessageListner(rec esnp.MessageListener) {
	this.tm.SetMessageListner(rec)
}

func (this *ChannelClient) SendMessage(ev *esnp.Message) error {
	if this.C != nil {
		return this.C.PostMessage(ev)
	}
	return errors.New("not open")
}

func (this *ChannelClient) FutureCall(msg *esnp.Message) *syncutil.Future {
	return this.tm.FutureCall(this.C, msg)
}

func (this *ChannelClient) Call(msg *esnp.Message, to *time.Timer) (*esnp.Message, error) {
	return this.tm.Call(this.C, msg, to)
}
