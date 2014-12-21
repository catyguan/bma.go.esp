package mempipelineconn

import (
	"bmautil/netutil"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	NETWORK = "mempipeline"
)

type memPipelineAddr struct {
	name string
}

func (this *memPipelineAddr) Network() string {
	return NETWORK
}

func (this *memPipelineAddr) String() string {
	return this.name
}

type memPipelineListener struct {
	n    string
	id   uint32
	wait chan net.Conn
}

func newListener(n string) *memPipelineListener {
	r := new(memPipelineListener)
	r.n = n
	r.wait = make(chan net.Conn, 4)
	return r
}

func (this *memPipelineListener) Accept() (c net.Conn, err error) {
	c = <-this.wait
	return
}

func (this *memPipelineListener) Close() error {
	CloseListener(this.n)
	return nil
}

func (this *memPipelineListener) Addr() net.Addr {
	return &memPipelineAddr{this.n}
}

type memPipelineConn struct {
	lname         string
	rname         string
	send          chan []byte
	receive       chan []byte
	rbuf          []byte
	readDeadline  time.Time
	writeDeadline time.Time

	closed uint32
}

func (this *memPipelineConn) Read(b []byte) (n int, err error) {
	if this.IsClosing() {
		return 0, io.EOF
	}
	l := len(this.rbuf)
	if l == 0 {
		if this.readDeadline.IsZero() {
			this.rbuf = <-this.receive
			if this.rbuf == nil {
				return 0, io.EOF
			}
			fmt.Println(this.lname, " << ", this.rbuf)
		} else {
			now := time.Now()
			if now.After(this.readDeadline) {
				return 0, netutil.TimeoutError()
			}
			t := time.NewTimer(this.readDeadline.Sub(now))
			select {
			case bs := <-this.receive:
				t.Stop()
				if bs == nil {
					return 0, io.EOF
				}
				this.rbuf = bs
			case <-t.C:
				return 0, netutil.TimeoutError()
			}
		}
		l = len(this.rbuf)
	}
	if l <= len(b) {
		n = l
		copy(b, this.rbuf)
		this.rbuf = nil
		return
	} else {
		n = len(b)
		copy(b, this.rbuf[:n])
		this.rbuf = this.rbuf[n:]
		return
	}
}

func (this *memPipelineConn) Write(b []byte) (n int, err error) {
	if this.IsClosing() {
		return 0, io.ErrClosedPipe
	}
	defer func() {
		x := recover()
		if x != nil {
			n = 0
			err = io.ErrClosedPipe
		}
	}()
	tmp := make([]byte, len(b))
	copy(tmp, b)
	b = tmp
	if this.writeDeadline.IsZero() {
		this.send <- b
		return len(b), nil
	} else {
		now := time.Now()
		if now.After(this.writeDeadline) {
			return 0, netutil.TimeoutError()
		}
		t := time.NewTimer(this.writeDeadline.Sub(now))
		select {
		case this.send <- b:
			t.Stop()
			return len(b), nil
		case <-t.C:
			return 0, netutil.TimeoutError()
		}
	}
}

func (this *memPipelineConn) Close() error {
	defer func() {
		recover()
	}()
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		close(this.send)
		close(this.receive)
	}
	return nil
}

func (this *memPipelineConn) LocalAddr() net.Addr {
	return &memPipelineAddr{this.lname}
}

func (this *memPipelineConn) RemoteAddr() net.Addr {
	return &memPipelineAddr{this.rname}
}

func (this *memPipelineConn) SetDeadline(t time.Time) error {
	this.readDeadline = t
	this.writeDeadline = t
	return nil
}

func (this *memPipelineConn) SetReadDeadline(t time.Time) error {
	this.readDeadline = t
	return nil
}

func (this *memPipelineConn) SetWriteDeadline(t time.Time) error {
	this.writeDeadline = t
	return nil
}

func (this *memPipelineConn) String() string {
	return this.lname
}

func (this *memPipelineConn) IsClosing() bool {
	return atomic.LoadUint32(&this.closed) != 0
}

func createPipeline(n string, id uint32, sz int) (*memPipelineConn, *memPipelineConn) {
	l := make(chan []byte, sz)
	r := make(chan []byte, sz)
	ln := fmt.Sprintf("%s.l.%d", n, id)
	rn := fmt.Sprintf("%s.r.%d", n, id)

	lc := newConn(ln, rn, l, r)
	rc := newConn(rn, ln, r, l)

	return lc, rc
}

func newConn(ln, rn string, s, r chan []byte) *memPipelineConn {
	o := new(memPipelineConn)
	o.lname = ln
	o.rname = rn
	o.send = s
	o.receive = r
	return o
}

type MemPipelineDialService struct {
	lock      sync.RWMutex
	listeners map[string]*memPipelineListener
}

var (
	gs MemPipelineDialService
)

func init() {
	gs.listeners = make(map[string]*memPipelineListener)
	netutil.AddDialService(NETWORK, &gs)
}

func CloseListener(n string) {
	gs.lock.Lock()
	defer gs.lock.Unlock()
	if p, ok := gs.listeners[n]; ok {
		delete(gs.listeners, n)
		close(p.wait)
	}
}

func (this *MemPipelineDialService) Listen(net, laddr string) (net.Listener, error) {
	gs.lock.Lock()
	defer gs.lock.Unlock()
	if _, ok := gs.listeners[laddr]; ok {
		return nil, fmt.Errorf("pipeline(%s) exists", laddr)
	}
	o := newListener(laddr)
	gs.listeners[laddr] = o
	return o, nil
}

func (this *MemPipelineDialService) Dial(network, address string) (net.Conn, error) {
	return this.DialTimeout(network, address, time.Second)
}

func (this *MemPipelineDialService) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	now := time.Now()
	for {
		if time.Since(now) >= timeout {
			return nil, netutil.TimeoutError()
		}
		gs.lock.RLock()
		p, ok := this.listeners[address]
		gs.lock.RUnlock()
		if ok {
			r := func() (c net.Conn) {
				defer func() {
					x := recover()
					if x != nil {
						c = nil
					}
				}()
				id := atomic.AddUint32(&p.id, 1)
				l, r := createPipeline(address, id, 64)
				p.wait <- l
				return r
			}()
			if r != nil {
				return r, nil
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
}
