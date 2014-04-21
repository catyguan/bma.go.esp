package espservice

import (
	"bmautil/socket"
	"bytes"
	"errors"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"runtime/debug"
	"sync/atomic"
)

type GoService struct {
	name    string
	handler ServiceHandler

	closed uint32
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

	}
	return true
}

func (this *GoService) PostRequest(sock *espsocket.Socket, msg *esnp.Message) error {
	if atomic.LoadUint32(&this.closed) > 0 {
		return errors.New("closed")
	}
	ctrl := esnp.FrameCoders.Trace
	p := msg.ToPackage()
	if ctrl.Has(p) {
		info := fmt.Sprintf("%s handled", this)
		rmsg := ctrl.CreateReply(msg, info)
		go sock.SendMessage(rmsg, nil)
	}
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logger.Warn(tag, "execute panic - %s\n%s", err, string(debug.Stack()))
			}
		}()
		err := this.handler(sock, msg)
		if err != nil {
			logger.Warn(tag, "execute fail - %s\n%s", err)
		}
	}()
	return nil
}

func (this *GoService) AcceptESP(sock *socket.Socket) error {
	ch := espsocket.NewSocketChannel(sock, "")
	s := espsocket.NewSocket(ch)
	return this.AcceptSocket(s)
}

func (this *GoService) AcceptSocket(sock *espsocket.Socket) error {
	sock.SetMessageListner(func(msg *esnp.Message) error {
		return this.PostRequest(sock, msg)
	})
	return nil
}
