package socket

import (
	"bmautil/byteutil"
	"bmautil/syncutil"
	"errors"
	"logger"
	"net"
	"sync"
	"time"
)

const (
	tag = "Socket"
)

type SocketInit func(sock *Socket) error
type SocketReceiver func(sock *Socket, data []byte) error
type SocketWriteListener func(socket *Socket, err error) bool
type SocketCloseListener func(socket *Socket)

func Func4CloseAfterSend(socket *Socket, err error) bool {
	return false
}

type WriteReq struct {
	data     *byteutil.BytesBuffer
	datafn   func() *byteutil.BytesBuffer
	callback SocketWriteListener
}

func NewWriteReqB(data []byte, cb SocketWriteListener) *WriteReq {
	buf := byteutil.NewBytesBufferB(data)
	return NewWriteReq(buf, cb)
}

func NewWriteReqF(f func() *byteutil.BytesBuffer, cb SocketWriteListener) *WriteReq {
	this := new(WriteReq)
	this.datafn = f
	this.callback = cb
	return this
}

func NewWriteReq(data *byteutil.BytesBuffer, cb SocketWriteListener) *WriteReq {
	this := new(WriteReq)
	this.data = data
	this.callback = cb
	return this
}

type Socket struct {
	Conn net.Conn

	readBuffer    int
	writeBuffer   int
	wbuffer       []byte
	Timeout       time.Duration
	Trace         int
	Receiver      SocketReceiver
	WriteListener SocketWriteListener

	clisLock      sync.Mutex
	closeListener map[string]SocketCloseListener

	// runtime
	closeState *syncutil.CloseState
	writeq     chan *WriteReq
}

func NewSocket(conn net.Conn, wbuf int, timeout time.Duration) *Socket {
	this := new(Socket)
	this.Conn = conn
	this.closeState = syncutil.NewCloseState()
	this.readBuffer = 4 * 1024
	this.writeBuffer = this.readBuffer
	this.Timeout = timeout
	this.writeq = make(chan *WriteReq, wbuf)
	this.Trace = 0
	this.closeListener = make(map[string]SocketCloseListener)

	return this
}

// public
func (this *Socket) String() string {
	return this.Conn.RemoteAddr().String()
}

func (this *Socket) SetReadBuffer(v int) {
	if v > 0 {
		this.readBuffer = v
		if c, ok := this.Conn.(*net.TCPConn); ok {
			c.SetReadBuffer(v)
		}
	}
}

func (this *Socket) SetWriteBuffer(v int) {
	if v > 0 {
		this.writeBuffer = v
		if c, ok := this.Conn.(*net.TCPConn); ok {
			c.SetWriteBuffer(v)
		}
	}
}

func (this *Socket) AddCloseListener(lis SocketCloseListener, id string) {
	if lis == nil {
		return
	}
	if id == "" {
		id = logger.Sprintf("%p", lis)
	}
	this.clisLock.Lock()
	defer this.clisLock.Unlock()
	this.closeListener[id] = lis
}

func (this *Socket) RemoveCloseListener(id string) {
	this.clisLock.Lock()
	defer this.clisLock.Unlock()
	delete(this.closeListener, id)
}

func (this *Socket) Write(req *WriteReq) (err error) {
	if this.closeState.IsClosing() {
		return errors.New("closed")
	}
	defer func() {
		ex := recover()
		if ex != nil {
			if err2, ok := ex.(error); ok {
				err = err2
			} else {
				err = errors.New(logger.Sprintf("%v", ex))
			}
		}
	}()
	this.writeq <- req
	return nil
}

func (this *Socket) WriteNW(req *WriteReq) {
	go func() {
		this.Write(req)
	}()
}

func (this *Socket) IsClosing() bool {
	return this.closeState.IsClosing()
}

func (this *Socket) Close() {
	if this.closeState.AskClose() {
		defer func() {
			recover()
		}()
		this.writeq <- nil
	}
}

func (this *Socket) doCloseWriteq() {
	defer func() {
		recover()
	}()
	for {
		select {
		case req := <-this.writeq:
			if req != nil {
				func() {
					defer func() {
						recover()
					}()
					req.callback(this, errors.New("closed"))
				}()
			}
		default:
			return
		}
	}
	close(this.writeq)
}

func (this *Socket) doClose() {
	this.Conn.Close()
	this.doCloseWriteq()

	func() {
		this.clisLock.Lock()
		defer this.clisLock.Unlock()
		for k, lis := range this.closeListener {
			delete(this.closeListener, k)
			func() {
				defer func() {
					recover()
				}()
				lis(this)
			}()
		}
	}()

	this.closeState.DoneClose()
	logger.Debug(tag, "Socket[%s] closed", this)
}

func (this *Socket) Start(sinit SocketInit) error {

	// run writer
	this.doStartWrite()

	// init
	if sinit != nil {
		err := sinit(this)
		if err != nil {
			logger.Error(tag, "Socket[%s] init fail -%s", this, err)
			this.Close()
			return err
		}
	}

	// run reader
	this.doStartRead()

	logger.Debug(tag, "Socket[%s] start", this)

	return nil
}

// private do
func (this *Socket) doSend(c net.Conn, buf *byteutil.BytesBuffer) error {
	if this.Trace > 0 && logger.EnableDebug(tag) {
		logger.Debug(tag, "Socket[%s] << %s", this, buf.TraceString(this.Trace))
	}
	if this.wbuffer == nil || len(this.wbuffer) != this.writeBuffer {
		this.wbuffer = make([]byte, this.writeBuffer)
	}
	ws := len(this.wbuffer)
	sz := 0
	dof := func(b []byte) error {
		defer func() {
			sz = 0
		}()
		for {
			if this.Timeout > 0 {
				this.Conn.SetWriteDeadline(time.Now().Add(this.Timeout))
			}
			l, err := c.Write(b)
			if err != nil {
				return err
			}
			if l < len(b) {
				b = b[l:]
			} else {
				return nil
			}
		}
	}
	i := 0
	for {
		if i >= buf.Len() {
			break
		}
		bs := buf.DataList[i]
		bl := len(bs)
		if sz+bl > ws {
			if sz > 0 {
				err := dof(this.wbuffer[:sz])
				if err != nil {
					return err
				}
			} else {
				err := dof(bs)
				if err != nil {
					return err
				}
				i++
			}
		} else {
			copy(this.wbuffer[sz:], bs)
			sz += bl
			i++
		}

	}
	if sz > 0 {
		err := dof(this.wbuffer[:sz])
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Socket) doStartWrite() {
	logger.Debug(tag, "Socket[%s] start write", this)
	go func() {
		defer func() {
			logger.Debug(tag, "Socket[%s] close write", this)
			this.doClose()
		}()
		for {
			req := <-this.writeq
			if req == nil {
				return
			}
			d := req.data
			if d == nil && req.datafn != nil {
				d = func() *byteutil.BytesBuffer {
					defer func() {
						recover()
					}()
					return req.datafn()
				}()
			}
			var err error
			if d != nil {
				err = this.doSend(this.Conn, d)
			}
			if err != nil {
				this.closeState.BeginClose()
				if req.callback != nil {
					func() {
						defer func() {
							recover()
						}()
						req.callback(this, err)
					}()
				}
				return
			}
			if req.callback != nil {
				cl := func() bool {
					defer func() {
						recover()
					}()
					return req.callback(this, nil)
				}()
				if !cl {
					this.closeState.BeginClose()
					logger.Debug(tag, "Socket[%s] close after write", this)
					return
				}
			}
		}
	}()
}

func (this *Socket) doStartRead() {
	logger.Debug(tag, "Socket[%s] start read", this)
	go func() {
		defer func() {
			logger.Debug(tag, "Socket[%s] close read", this)
		}()
		var buf []byte
		for {
			rbs := this.readBuffer
			if buf == nil || (rbs > 0 && len(buf) != rbs) {
				buf = make([]byte, rbs)
			}
			if this.Timeout > 0 {
				this.Conn.SetReadDeadline(time.Now().Add(this.Timeout))
			}
			l, err := this.Conn.Read(buf)
			if l == 0 && err != nil {
				if !this.closeState.IsClosing() {
					this.Close()
					logger.Debug(tag, "Socket[%s] read fail - %s", this, err)
				}
				return
			}
			data := make([]byte, l)
			copy(data, buf[:l])
			if this.Trace > 0 && logger.EnableDebug(tag) {
				bbuf := byteutil.NewBytesBufferB(data)
				logger.Debug(tag, "Socket[%s] >> %s", this, bbuf.TraceString(this.Trace))
			}
			if this.Receiver != nil {
				err = this.Receiver(this, data)
				if err != nil {
					this.Close()
					logger.Debug(tag, "Socket[%s] receive fail - %s", this, err)
					return
				}
			} else {
				logger.Warn(tag, "Socket[%s] no receiver", this)
			}
		}
	}()
}
