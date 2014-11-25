package espsocket

import (
	"bmautil/socket"
	"bmautil/syncutil"
	"errors"
	"esp/espnet/esnp"
	"logger"
	"net"
	"sync"
	"time"
)

// waitInfo
type waitInfo struct {
	lis esnp.ResponseListener
	tm  *time.Timer
}

// Socket
type Socket struct {
	id      uint32
	channel Channel

	respLock sync.Mutex
	waiting  map[uint64]*waitInfo
	listener esnp.MessageListener
	handlers []esnp.MessageHandler

	propLock sync.Mutex
	props    map[string]interface{}

	lisGroup CloseListenerGroup
}

func NewSocket(ch Channel) *Socket {
	this := new(Socket)
	this.id = NextSocketId()
	this.channel = ch
	ch.Bind(this.OnMessageIn, this.onChannelClose)
	return this
}

func Dial(name string, cfg *socket.DialConfig, coderName string, log bool) (*Socket, error) {
	sock, err := socket.Dial2(name, cfg, nil, log)
	if err != nil {
		return nil, err
	}
	ch := NewSocketChannel(sock, coderName)
	return NewSocket(ch), nil
}

func (this *Socket) onChannelClose() {
	go func() {
		this.respLock.Lock()
		tmp := this.waiting
		this.waiting = nil
		this.respLock.Unlock()
		if tmp != nil && len(tmp) > 0 {
			err := errors.New("closed")
			for _, wi := range tmp {
				wi.lis(nil, err)
			}
		}

		this.lisGroup.OnClose()

		this.propLock.Lock()
		defer this.propLock.Unlock()
		this.props = nil

		this.channel = nil
	}()
}

func (this *Socket) String() string {
	ch := this.channel
	if ch == nil {
		return "closedSocket[]"
	}
	return ch.String()
}

func (this *Socket) GetProperty(name string) (interface{}, bool) {
	ch := this.channel
	if ch != nil {
		v, ok := ch.GetProperty(name)
		if ok {
			return v, true
		}
	}
	this.propLock.Lock()
	defer this.propLock.Unlock()
	if this.props != nil {
		if rv, ok := this.props[name]; ok {
			return rv, true
		}
	}
	return nil, false
}

func (this *Socket) SetProperty(name string, val interface{}) bool {
	ch := this.channel
	if ch != nil {
		ok := ch.SetProperty(name, val)
		if ok {
			return true
		}
	}
	this.propLock.Lock()
	defer this.propLock.Unlock()
	if this.props == nil {
		this.props = make(map[string]interface{})
	}
	this.props[name] = val
	return true
}

func (this *Socket) doClose(closeChannel bool, shutdown bool) {
	if closeChannel {
		ch := this.channel
		if ch != nil {
			if shutdown {
				ch.Shutdown()
			} else {
				ch.AskClose()
			}
		}
	}
}

func (this *Socket) AskClose() {
	this.doClose(true, false)
}

func (this *Socket) Shutdown() {
	this.doClose(true, true)
}

func (this *Socket) Id() uint32 {
	return this.id
}

func (this *Socket) PostMessage(msg *esnp.Message) error {
	ch := this.channel
	if ch == nil {
		return errors.New("closed")
	}
	return ch.SendMessage(msg, nil)
}

func (this *Socket) SendMessage(msg *esnp.Message, cb SendCallback) error {
	ch := this.channel
	if ch == nil {
		return errors.New("closed")
	}
	return ch.SendMessage(msg, cb)
}

func (this *Socket) SetCloseListener(name string, lis func()) error {
	this.lisGroup.Set(name, lis)
	return nil
}

func (this *Socket) IsBreak() bool {
	ch := this.channel
	return ch == nil || ch.IsClosing()
}

// Call & Handler
func (this *Socket) AddMessageHandler(mh esnp.MessageHandler) {
	this.respLock.Lock()
	defer this.respLock.Unlock()
	if this.handlers == nil {
		this.handlers = make([]esnp.MessageHandler, 0)
	}
	this.handlers = append(this.handlers, mh)
}

func (this *Socket) SetMessageListner(rec esnp.MessageListener) esnp.MessageListener {
	r := this.listener
	this.listener = rec
	return r
}

func (this *Socket) Invoke(msg *esnp.Message, cb esnp.ResponseListener, timeout time.Duration) {
	mid := msg.SureId()
	msg.SureRequest()
	this.respLock.Lock()
	if this.waiting == nil {
		this.waiting = make(map[uint64]*waitInfo)
	}
	wi := new(waitInfo)
	wi.lis = cb
	this.waiting[mid] = wi
	this.respLock.Unlock()

	err := this.PostMessage(msg)
	if err != nil {
		go this.OnError(mid, err)
	} else {
		wi.tm = time.AfterFunc(timeout, func() {
			this.OnTimeout(mid)
		})
	}
}

func (this *Socket) FutureCall(msg *esnp.Message, timeout time.Duration) *syncutil.Future {
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
	this.Invoke(msg, cb, timeout)
	return f
}

func (this *Socket) Call(msg *esnp.Message, timeout time.Duration) (*esnp.Message, error) {
	f := this.FutureCall(msg, timeout)
	f.WaitDone()
	_, v, err := f.Get()
	if err != nil {
		return nil, err
	}
	return v.(*esnp.Message), nil
}

func (this *Socket) Cancel(mid uint64) esnp.ResponseListener {
	this.respLock.Lock()
	defer this.respLock.Unlock()
	if this.waiting != nil {
		wi, ok := this.waiting[mid]
		if ok {
			wi.tm.Stop()
			delete(this.waiting, mid)
			return wi.lis
		}
	}
	return nil
}

func (this *Socket) OnMessageIn(msg *esnp.Message) error {
	p := msg.ToPackage()
	if esnp.FrameCoders.Flag.Has(p, esnp.FLAG_RESP) {
		mid := esnp.FrameCoders.SourceMessageId.Get(msg.ToPackage())
		if mid > 0 {
			rlis := this.Cancel(mid)
			if rlis != nil {
				return rlis(msg, nil)
			}
		}
		logger.Debug(tag, "'%s' discard response[%d]", this, mid)
		return nil
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
	lis := this.listener
	if lis != nil {
		return lis(msg)
	}
	logger.Debug(tag, "'%s' discard message", this)
	return nil
}

func (this *Socket) OnError(mid uint64, err error) {
	rlis := this.Cancel(mid)
	if rlis != nil {
		err2 := rlis(nil, err)
		if err2 != nil {
			logger.Debug(tag, "'%s' response fail - %s", this, err2)
		}
	}
}

func (this *Socket) OnTimeout(mid uint64) {
	this.OnError(mid, errors.New("timeout"))
}

func (this *Socket) GetRemoteAddr() (string, bool) {
	v, ok := this.GetProperty(PROP_SOCKET_REMOTE_ADDR)
	if !ok {
		return "", false
	}
	if v == nil {
		return "", false
	}
	if str, ok2 := v.(string); ok2 {
		return str, true
	}
	if addr, ok2 := v.(net.Addr); ok2 {
		return addr.String(), true
	}
	logger.Warn(tag, "unknow RemoteAddr(%T)", v)
	return "", false
}
