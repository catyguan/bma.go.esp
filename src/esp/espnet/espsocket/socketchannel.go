package espsocket

import (
	"bmautil/byteutil"
	"bmautil/socket"
	"bmautil/valutil"
	"esp/espnet/esnp"
	"fmt"
	"logger"
	"net"
	"sync/atomic"
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
	socket        *socket.Socket
	coder         SocketChannelCoder
	receiver      esnp.MessageListener
	closeListener func()

	closed       uint32
	socketReturn func(s *socket.Socket)
}

func NewSocketChannel(sock *socket.Socket, coderName string) *SocketChannel {
	if coderName == "" {
		coderName = SOCKET_CHANNEL_CODER_ESPNET
	}
	var c SocketChannelCoder
	fac, ok := globalSocketChannelCoders[coderName]
	if !ok {
		panic("unknow socket channel coder '" + coderName + "'")
	}
	c = fac()
	return NewSocketChannelC(sock, c)
}

func NewSocketChannelC(sock *socket.Socket, c SocketChannelCoder) *SocketChannel {
	this := new(SocketChannel)
	this.socket = sock
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
		msg, ok := ev.(*esnp.Message)
		if !ok {
			logger.Debug(tag, "%s not messsage[%t] pass in", this, ev)
			return nil
		}
		return rec(msg)
	}
	logger.Debug(tag, "%s no receiver", this)
	return nil
}

func (this *SocketChannel) onSocketClose(sock *socket.Socket) {
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		if this.closeListener != nil {
			this.closeListener()
			this.closeListener = nil
		}
	}
}

func (this *SocketChannel) String() string {
	s := this.socket
	if s == nil {
		return "closedSocketChannel[]"
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
	return false
}

func (this *SocketChannel) PostEvent(ev interface{}, cb SendCallback) error {
	s := this.socket
	if s == nil {
		return fmt.Errorf("closed")
	}

	if msg, ok := ev.(*esnp.Message); ok {
		p := msg.ToPackage()
		ctrl := esnp.FrameCoders.Trace
		if ctrl.Has(p) {
			info := fmt.Sprintf("%s -> %s", s.Conn.LocalAddr(), s.Conn.RemoteAddr())
			rmsg := ctrl.CreateReply(msg, info)
			go this.onReceiveEvent(rmsg)
		}
	}

	var f4send socket.SocketWriteListener
	if cb != nil {
		f4send = func(s *socket.Socket, err error) bool {
			go cb(err)
			return true
		}
	}
	callf := func(ev interface{}) error {
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
		req = socket.NewWriteReqB(v, f4send)
	case [][]byte:
		data := byteutil.NewBytesBufferA(v)
		req = socket.NewWriteReq(data, f4send)
	case *byteutil.BytesBuffer:
		req = socket.NewWriteReq(v, f4send)
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

func (this *SocketChannel) doClose(force bool) {
	s := this.socket
	if s == nil {
		return
	}

	if this.socketReturn != nil {
		if !force {
			s.RemoveCloseListener(this.closeId())
			this.onSocketClose(nil)
			this.socketReturn(s)
		}
		return
	}
	if !s.IsClosing() {
		s.Close()
	}
}

func (this *SocketChannel) AskClose() {
	this.doClose(false)
}

func (this *SocketChannel) Shutdown() {
	this.doClose(true)
}

func (this *SocketChannel) Bind(rec esnp.MessageListener, closeLis func()) {
	this.receiver = rec
	this.closeListener = closeLis
}

// Channel
func (this *SocketChannel) ToChannel() Channel {
	return Channel(this)
}

func (this *SocketChannel) SendMessage(msg *esnp.Message, cb SendCallback) error {
	return this.PostEvent(msg, cb)
}

func (this *SocketChannel) IsClosing() bool {
	s := this.socket
	return s == nil || s.IsClosing()
}
