package espnet

import (
	"bmautil/socket"
	"errors"
	"logger"
	"sync"
	"time"
)

type ResponseListener func(msg *Message, err error) error

type ChannelClient struct {
	C        Channel
	listener MessageListener
	own      bool

	lock    sync.RWMutex
	waiting map[uint64]ResponseListener
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
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.waiting != nil {
		for k, lis := range this.waiting {
			delete(this.waiting, k)
			lis(nil, errors.New("closed"))
		}
	}
}

func (this *ChannelClient) IsOpen() bool {
	if this.C == nil {
		return false
	}
	if cb, ok := this.C.(BreakSupport); ok {
		bv := cb.IsBreak()
		if bv != nil && *bv {
			return false
		}
	}
	return true
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

func (this *ChannelClient) Call(msg *Message, to *time.Timer) (*Message, error) {
	var rmsg *Message
	var rerr error
	c := make(chan bool, 1)

	mid := msg.SureId()
	this.lock.Lock()
	if this.waiting == nil {
		this.waiting = make(map[uint64]ResponseListener)
	}
	this.waiting[mid] = func(msg *Message, err error) error {
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
	this.lock.Unlock()
	err := this.SendMessage(msg)
	if err != nil {
		close(c)
		return nil, err
	}
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

func (this *ChannelClient) popListener(mid uint64) ResponseListener {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.waiting != nil {
		rlis, ok := this.waiting[mid]
		if ok {
			delete(this.waiting, mid)
		}
		return rlis
	}
	return nil
}

func (this *ChannelClient) OnMessageIn(msg *Message) error {
	mid := FrameCoders.SourceMessageId.Get(msg.ToPackage())
	if mid > 0 {
		rlis := this.popListener(mid)
		if rlis != nil {
			return rlis(msg, nil)
		}
	}
	if this.listener != nil {
		return this.listener(msg)
	}
	logger.Debug(tag, "%s discard message", this.C)
	return nil
}
