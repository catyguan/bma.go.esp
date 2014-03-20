package espterminal

import (
	"bmautil/socket"
	"bmautil/syncutil"
	"errors"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"logger"
	"sync"
	"time"
)

const (
	tag = "espTerminal"
)

type Terminal struct {
	name     string
	listener esnp.MessageListener
	handlers []esnp.MessageHandler

	lock    sync.Mutex
	waiting map[uint64]esnp.ResponseListener
}

func (this *Terminal) InitTerminal(n string) {
	this.name = n
}

func (this *Terminal) SetName(n string) {
	this.name = n
}

func (this *Terminal) ConnectSocket(sock *socket.Socket, coderName string) error {
	ch := espchannel.NewSocketChannel(sock, coderName)
	return this.Connect(ch)
}

func (this *Terminal) Connect(ch espchannel.Channel) error {
	ch.SetMessageListner(this.OnMessageIn)
	return nil
}

func (this *Terminal) Disconnect(ch espchannel.Channel) {
	ch.SetMessageListner(nil)
}

func (this *Terminal) Close() {
	this.lock.Lock()
	tmp := this.waiting
	this.waiting = nil
	this.lock.Unlock()
	if tmp != nil {
		for _, lis := range this.waiting {
			lis(nil, errors.New("closed"))
		}
	}
}

func (this *Terminal) AddMessageHandler(mh esnp.MessageHandler) {
	if this.handlers == nil {
		this.handlers = make([]esnp.MessageHandler, 0)
	}
	this.handlers = append(this.handlers, mh)
}

func (this *Terminal) SetMessageListner(rec esnp.MessageListener) {
	this.listener = rec
}

func (this *Terminal) Invoke(ch espchannel.Channel, msg *esnp.Message, cb esnp.ResponseListener) {
	mid := msg.SureId()
	msg.SureRequest()
	this.lock.Lock()
	if this.waiting == nil {
		this.waiting = make(map[uint64]esnp.ResponseListener)
	}
	this.waiting[mid] = cb
	this.lock.Unlock()
	err := ch.PostMessage(msg)
	if err != nil {
		go cb(nil, err)
	}
}

func (this *Terminal) FutureCall(ch espchannel.Channel, msg *esnp.Message) *syncutil.Future {
	f, fe := syncutil.NewFuture()
	cb := func(msg *esnp.Message, err error) error {
		rmsg := msg
		rerr := err
		if msg != nil {
			merr := msg.ToError()
			if merr != nil {
				rerr = merr
			}
		}
		fe(rmsg, rerr)
		return nil
	}
	this.Invoke(ch, msg, cb)
	return f
}

func (this *Terminal) Call(ch espchannel.Channel, msg *esnp.Message, to *time.Timer) (*esnp.Message, error) {
	var rmsg *esnp.Message
	var rerr error
	c := make(chan bool, 1)
	cb := func(msg *esnp.Message, err error) error {
		rmsg = msg
		rerr = err
		if msg != nil {
			merr := msg.ToError()
			if merr != nil {
				rerr = merr
			}
		}
		close(c)
		return nil
	}
	this.Invoke(ch, msg, cb)
	if to != nil {
		select {
		case <-c:
		case <-to.C:
			return nil, errors.New("timeout")
		}
	} else {
		<-c
	}
	return rmsg, rerr
}

func (this *Terminal) popListener(mid uint64) esnp.ResponseListener {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.waiting != nil {
		rlis, ok := this.waiting[mid]
		if ok {
			delete(this.waiting, mid)
		}
		return rlis
	}
	return nil
}

func (this *Terminal) OnMessageIn(msg *esnp.Message) error {
	mid := esnp.FrameCoders.SourceMessageId.Get(msg.ToPackage())
	if mid > 0 {
		rlis := this.popListener(mid)
		if rlis != nil {
			return rlis(msg, nil)
		}
	}
	if this.handlers != nil {
		for _, h := range this.handlers {
			ok, err := h(msg)
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
		}
	}
	if this.listener != nil {
		return this.listener(msg)
	}
	logger.Debug(tag, "'%s' discard message", this.name)
	return nil
}
