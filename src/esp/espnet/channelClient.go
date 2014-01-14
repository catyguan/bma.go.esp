package espnet

import (
	"bmautil/socket"
	"errors"
	"logger"
)

type ChannelClient struct {
	C        Channel
	listener MessageListener
	own      bool
}

func NewChannelClient() *ChannelClient {
	r := new(ChannelClient)
	return r
}

func (this *ChannelClient) ConnectSocket(sock *socket.Socket, coderName string, own bool) error {
	r := NewSocketChannel(sock, coderName)
	return this.Connect(r, own)
}

func (this *ChannelClient) Connect(ch Channel, own bool) error {
	this.C = ch
	this.own = own
	ch.SetMessageListner(this.OnMessageIn)
	return nil
}

func (this *ChannelClient) Dial(name string, cfg *DialConfig, coderName string) error {
	sock, err := Dial(name, cfg, nil)
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
}

func (this *ChannelClient) IsOpen() bool {
	return this.C != nil
}

func (this *ChannelClient) SetMessageListner(rec MessageListener) {
	this.listener = rec
}

func (this *ChannelClient) SendMessage(ev *Message) error {
	if this.C != nil {
		return this.C.SendMessage(ev)
	}
	return errors.New("not open")
}

func (this *ChannelClient) OnMessageIn(msg *Message) error {
	if this.listener != nil {
		return this.listener(msg)
	}
	logger.Debug(tag, "%s discard message", this.C)
	return nil
}
