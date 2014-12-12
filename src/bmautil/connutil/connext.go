package connutil

import (
	"fmt"
	"net"
	"time"
)

type ConnDebuger func(conn net.Conn, b []byte, read bool)

type ConnExt struct {
	Debuger ConnDebuger
	conn    net.Conn
	buf     []byte
}

func NewConnExt(conn net.Conn, debuger ConnDebuger) *ConnExt {
	if o, ok := conn.(*ConnExt); ok {
		if debuger != nil {
			o.Debuger = debuger
		}
		return o
	}

	r := new(ConnExt)
	r.conn = conn
	r.Debuger = debuger
	return r
}

func (this *ConnExt) CheckBreak() bool {
	this.conn.SetReadDeadline(time.Now().Add(1))
	one := make([]byte, 1)
	n, err := this.conn.Read(one)
	if n > 0 {
		if this.buf == nil {
			this.buf = one
		} else {
			this.buf = append(this.buf, one[0])
		}
	}
	this.conn.SetReadDeadline(time.Time{})
	if err == nil {
		return false
	}
	// fmt.Println("error", err)
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		return false
	}
	return true
}

func (this *ConnExt) WriteString(s string) (int, error) {
	return this.Write([]byte(s))
}

func (this *ConnExt) Read(b []byte) (n int, err error) {
	if this.buf != nil && len(this.buf) > 0 {
		if len(b) >= len(this.buf) {
			n = len(this.buf)
			copy(b, this.buf)
			this.buf = nil
		} else {
			n = len(b)
			copy(b, this.buf[:n])
			this.buf = this.buf[n:]
		}
	} else {
		n, err = this.conn.Read(b)
	}
	if n > 0 && this.Debuger != nil {
		this.Debuger(this.conn, b[:n], true)
	}
	return n, err
}

func (this *ConnExt) Write(b []byte) (n int, err error) {
	n, err = this.conn.Write(b)
	if n > 0 && this.Debuger != nil {
		this.Debuger(this.conn, b[:n], false)
	}
	return n, err
}

func (this *ConnExt) Close() error {
	return this.conn.Close()
}

func (this *ConnExt) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *ConnExt) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *ConnExt) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}

func (this *ConnExt) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}

func (this *ConnExt) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}

func (this *ConnExt) String() string {
	ra := this.conn.RemoteAddr()
	if ra != nil {
		return ra.String()
	}
	return fmt.Sprintf("%v", this.conn)
}
