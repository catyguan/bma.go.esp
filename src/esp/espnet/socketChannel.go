package espnet

import (
	"bmautil/byteutil"
	"bmautil/socket"
	"bmautil/valutil"
	"fmt"
	"logger"
	"net"
	"sync"
	"time"
)

// SocketChannelCoder
type SocketChannelCoder interface {
	SupportProp

	EncodeMessage(ch *SocketChannel, ev interface{}, se func(ev interface{}) error) error

	DecodeMessage(ch *SocketChannel, data []byte, rec func(ev interface{}) error) error
}

type SocketChannelCoderFactory func() SocketChannelCoder

var globalSocketChannelCoders map[string]SocketChannelCoderFactory = make(map[string]SocketChannelCoderFactory)

func RegSocketChannelCoder(n string, c SocketChannelCoderFactory) {
	globalSocketChannelCoders[n] = c
}

func ListSocketChannelCoder() []string {
	r := make([]string, 0)
	for k, _ := range globalSocketChannelCoders {
		r = append(r, k)
	}
	return r
}

func GetSocketChannelCoder(n string) SocketChannelCoderFactory {
	return globalSocketChannelCoders[n]
}

// SocketChannel
type SocketChannel struct {
	id       uint32
	socket   *socket.Socket
	coder    SocketChannelCoder
	receiver func(ev interface{}) error

	propLock sync.Mutex
	props    map[string]interface{}

	lisGroup     CloseListenerGroup
	socketReturn func(s *socket.Socket)
}

func NewSocketChannel(sock *socket.Socket, coderName string) *SocketChannel {
	var c SocketChannelCoder
	if coderName != "" {
		fac, ok := globalSocketChannelCoders[coderName]
		if !ok {
			panic("unknow socket channel coder '" + coderName + "'")
		}
		c = fac()
	}
	return NewSocketChannelC(sock, c)
}

func NewSocketChannelC(sock *socket.Socket, c SocketChannelCoder) *SocketChannel {
	this := new(SocketChannel)
	this.id = NextChanneId()
	this.socket = sock
	this.coder = nil
	if c != nil {
		this.coder = c
	}
	this.socket.Receiver = this.onSocketReceive
	this.socket.AddCloseListener(this.onSocketClose, this.closeId())
	return this
}

func (this *SocketChannel) closeId() string {
	return fmt.Sprintf("SC_%p", this)
}

func (this *SocketChannel) onSocketReceive(sock *socket.Socket, data []byte) error {
	var err error
	if this.coder != nil {
		err = this.coder.DecodeMessage(this, data, this.onReceiveEvent)
	} else {
		err = this.onReceiveEvent(data)
	}
	return err
}

func (this *SocketChannel) onReceiveEvent(ev interface{}) error {
	rec := this.receiver
	if rec != nil {
		return rec(ev)
	}
	logger.Debug(tag, "%s no receiver", this)
	return nil
}

func (this *SocketChannel) onSocketClose(sock *socket.Socket) {
	this.lisGroup.OnClose()

	this.propLock.Lock()
	defer this.propLock.Unlock()
	this.props = nil
}

func (this *SocketChannel) String() string {
	s := this.socket
	if s == nil {
		return "closedSocketChannel"
	}
	return s.String()
}

func (this *SocketChannel) GetProperty(name string) (interface{}, bool) {
	s := this.socket
	if s != nil {
		switch name {
		case PROP_SOCKET_REMOTE_ADDR:
			return s.Conn.RemoteAddr(), true
		case PROP_SOCKET_LOCAL_ADDR:
			return s.Conn.LocalAddr(), true
		}
	}
	if this.coder != nil {
		if rv, ok := this.coder.GetProperty(name); ok {
			return rv, true
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

func (this *SocketChannel) SetProperty(name string, val interface{}) bool {
	s := this.socket
	if s != nil {
		switch name {
		case PROP_SOCKET_DEAD_LINE:
			if v, ok := val.(time.Time); ok {
				s.Conn.SetDeadline(v)
				return true
			}
		case PROP_SOCKET_READ_DEAD_LINE:
			if v, ok := val.(time.Time); ok {
				s.Conn.SetReadDeadline(v)
				return true
			}
		case PROP_SOCKET_WRITE_DEAD_LINE:
			if v, ok := val.(time.Time); ok {
				s.Conn.SetWriteDeadline(v)
				return true
			}
		case PROP_SOCKET_TRACE:
			s.Trace = valutil.ToInt(val, 0)
			return true
		case PROP_SOCKET_TIMEOUT:
			if v, ok := val.(time.Duration); ok {
				s.Timeout = v
				return true
			}
		case PROP_SOCKET_READ_BUFFER:
			v := valutil.ToInt(val, -1)
			if v > 0 {
				s.SetReadBuffer(v)
				return true
			}
		case PROP_SOCKET_WRITE_BUFFER:
			v := valutil.ToInt(val, -1)
			if v > 0 {
				s.SetWriteBuffer(v)
				return true
			}
		case PROP_SOCKET_WRITE_CHAN_SIZE:
			v := valutil.ToInt(val, -1)
			if v > 0 {
				s.SetWriteChanSize(v)
				return true
			}
		}

		if c, ok := s.Conn.(*net.TCPConn); ok {
			switch name {
			case PROP_SOCKET_LINGER:
				v := valutil.ToInt(val, -1)
				if v >= 0 {
					c.SetLinger(v)
					return true
				}
			case PROP_SOCKET_KEEP_ALIVE:
				v, ok := valutil.ToBoolNil(val)
				if ok {
					c.SetKeepAlive(v)
					return true
				}
			case PROP_SOCKET_KEEP_ALIVE_PERIOD:
				if v, ok := val.(time.Duration); ok {
					c.SetKeepAlivePeriod(v)
					return true
				}
			case PROP_SOCKET_NO_DELAY:
				v, ok := valutil.ToBoolNil(val)
				if ok {
					c.SetNoDelay(v)
					return true
				}
			}
		}
	}
	if this.coder != nil {
		if this.coder.SetProperty(name, val) {
			return true
		}
	}
	this.propLock.Lock()
	defer this.propLock.Unlock()
	if this.props == nil {
		this.props = make(map[string]interface{})
		this.props[name] = val
	}
	return true
}

func (this *SocketChannel) PostEvent(ev interface{}) error {
	s := this.socket
	if s == nil {
		return fmt.Errorf("closed")
	}

	cctype := CLOSE_CHANNEL_NONE
	if msg, ok := ev.(*Message); ok {
		p := msg.ToPackage()
		ctrl := FrameCoders.Trace
		if ctrl.Has(p) {
			info := fmt.Sprintf("%s -> %s", s.Conn.LocalAddr(), s.Conn.RemoteAddr())
			rmsg := ctrl.CreateReply(msg, info)
			go this.onReceiveEvent(rmsg)
		}
		cctype = FrameCoders.CloseChannel.Has(p)
	}

	if cctype == CLOSE_CHANNEL_NOW {
		go func() {
			s := this.socket
			if s != nil && !s.IsClosing() {
				s.Close()
			}
		}()
		return nil
	}

	callf := func(ev interface{}) error {
		var f4send socket.SocketWriteListener
		if cctype == CLOES_CHANNEL_AFTER_SEND {
			f4send = socket.Func4CloseAfterSend
		}
		return this.doPostEvent(ev, f4send)
	}
	if this.coder != nil {
		return this.coder.EncodeMessage(this, ev, callf)
	}
	return callf(ev)
}

func (this *SocketChannel) doPostEvent(ev interface{}, f4send socket.SocketWriteListener) error {
	s := this.socket
	if s == nil {
		return fmt.Errorf("closed")
	}
	var req *socket.WriteReq
	switch v := ev.(type) {
	case []byte:
		req = socket.NewWriteReqB(v, nil)
	case [][]byte:
		data := byteutil.NewBytesBufferA(v)
		req = socket.NewWriteReq(data, nil)
	case *byteutil.BytesBuffer:
		req = socket.NewWriteReq(v, nil)
	case *socket.WriteReq:
		req = v
	default:
		logger.Debug(tag, "unknow event %T", ev)
	}
	if req != nil {
		return s.Write(req)
	}
	return nil
}

func (this *SocketChannel) AskClose() {
	this.Close()
}

func (this *SocketChannel) Close() {
	this.propLock.Lock()
	s := this.socket
	this.socket = nil
	this.propLock.Unlock()

	if s == nil {
		return
	}

	if this.socketReturn != nil {
		s.RemoveCloseListener(this.closeId())
		this.onSocketClose(nil)
		this.socketReturn(s)
	} else {
		if !s.IsClosing() {
			s.Close()
		}
	}
}

func (this *SocketChannel) SetPipelineListner(rec func(ev interface{}) error) {
	this.receiver = rec
}

func (this *SocketChannel) doRequestResponse(rmsg *Message) error {
	logger.Info(tag, "HERE")
	return this.PostEvent(rmsg)
}

// Channel
func (this *SocketChannel) ToChannel() Channel {
	return Channel(this)
}

func (this *SocketChannel) Id() uint32 {
	return this.id
}

func (this *SocketChannel) Name() string {
	return this.String()
}

func (this *SocketChannel) SendMessage(msg *Message) error {
	return this.PostEvent(msg)
}

func (this *SocketChannel) SetMessageListner(rec MessageListener) {
	this.SetPipelineListner(func(ev interface{}) error {
		if msg, ok := ev.(*Message); ok {
			return rec(msg)
		}
		return nil
	})
}

func (this *SocketChannel) SetCloseListener(name string, lis func()) error {
	this.lisGroup.Set(name, lis)
	return nil
}
