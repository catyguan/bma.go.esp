package espnet

import (
	"bytes"
	"errors"
	"fmt"
	"logger"
	"runtime/debug"
	"sync/atomic"
)

type GoService struct {
	name    string
	handler ServiceHandler

	closed   uint32
	channels VChannelGroup
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

func (this *GoService) PostRequest(msg *Message, rep ServiceResponser) error {
	if atomic.LoadUint32(&this.closed) > 0 {
		return errors.New("closed")
	}
	ctrl := FrameCoders.Trace
	p := msg.ToPackage()
	if rep != nil && ctrl.Has(p) {
		info := fmt.Sprintf("%s handle", this)
		rmsg := ctrl.CreateReply(msg, info)
		go rep(rmsg)
	}
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logger.Warn(tag, "execute panic - %s\n%s", err, string(debug.Stack()))
			}
		}()
		err := this.handler(msg, rep)
		if err != nil {
			logger.Warn(tag, "execute fail - %s\n%s", err)
		}
	}()
	return nil
}

func (this *GoService) NewChannel() (Channel, error) {
	r := new(VChannel)
	r.InitVChannel(this.name)
	r.RemoveChannel = this.channels.Remove
	r.Sender = func(msg *Message) error {
		return DoServiceHandle(this.PostRequest, msg, r.ServiceResponse)
	}
	this.channels.Add(r)
	return r, nil
}
