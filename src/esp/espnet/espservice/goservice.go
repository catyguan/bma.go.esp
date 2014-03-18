package espservice

import (
	"bmautil/socket"
	"bytes"
	"errors"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"fmt"
	"logger"
	"runtime/debug"
	"sync/atomic"
)

type GoService struct {
	name    string
	handler ServiceHandler

	closed   uint32
	channels espchannel.VChannelGroup
}

func NewGoService(name string, h ServiceHandler) *GoService {
	this := new(GoService)
	this.name = name
	this.handler = h
	return this
}

func (this *GoService) Name() string {
	return this.name
}

func (this *GoService) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(this.name)
	buf.WriteString("(gos)")
	return buf.String()
}

func (this *GoService) Stop() bool {
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		this.channels.OnClose()
	}
	return true
}

func (this *GoService) PostRequest(ch espchannel.Channel, msg *esnp.Message) error {
	if atomic.LoadUint32(&this.closed) > 0 {
		return errors.New("closed")
	}
	ctrl := esnp.FrameCoders.Trace
	p := msg.ToPackage()
	if ctrl.Has(p) {
		info := fmt.Sprintf("%s handle", this)
		rmsg := ctrl.CreateReply(msg, info)
		go ch.PostMessage(rmsg)
	}
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logger.Warn(tag, "execute panic - %s\n%s", err, string(debug.Stack()))
			}
		}()
		err := this.handler(ch, msg)
		if err != nil {
			logger.Warn(tag, "execute fail - %s\n%s", err)
		}
	}()
	return nil
}

func (this *GoService) AcceptESP(sock *socket.Socket) error {
	ch := espchannel.NewSocketChannel(sock, espchannel.SOCKET_CHANNEL_CODER_ESPNET)
	ConnectService(ch, this.PostRequest)
	return nil
}

func (this *GoService) NewChannel() (espchannel.Channel, error) {
	r := new(espchannel.VChannel)
	r.InitVChannel(this.name)
	r.RemoveChannel = this.channels.Remove

	r.Sender = func(msg *esnp.Message) error {
		return DoServiceHandle(this.PostRequest, r, msg)
	}
	this.channels.Add(r)
	return r, nil
}
