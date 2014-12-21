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
	prop    map[string]interface{}
}

func NewConnExt(conn net.Conn) *ConnExt {
	if o, ok := conn.(*ConnExt); ok {
		return o
	}

	r := new(ConnExt)
	r.conn = conn
	return r
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

func (this *ConnExt) BaseConn() net.Conn {
	return this.conn
}

func (this *ConnExt) GetProperty(name string) (interface{}, bool) {
	if this.prop != nil {
		r, ok := this.prop[name]
		return r, ok
	}
	return nil, false
}

func (this *ConnExt) SetProperty(name string, val interface{}) {
	if val == nil {
		if this.prop != nil {
			delete(this.prop, name)
		}
	} else {
		if this.prop == nil {
			this.prop = make(map[string]interface{})
		}
		this.prop[name] = val
	}
}

func (this *ConnExt) ClearProperty() {
	if this.prop != nil {
		this.prop = nil
	}
}

func (this *ConnExt) ListProperty() []string {
	r := make([]string, 0, len(this.prop))
	for k, _ := range this.prop {
		r = append(r, k)
	}
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

func (this *ConnExt) ReadCheckTimeout(b []byte) (n int, timeout bool, err error) {
	n, err = this.conn.Read(b)
	if err != nil {
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			timeout = true
		}
	}
	return
}

func (this *ConnExt) WriteCheckTimeout(b []byte) (n int, timeout bool, err error) {
	n, err = this.conn.Write(b)
	if err != nil {
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			timeout = true
		}
	}
	return
}

func (this *ConnExt) ClearReadDeadline() {
	this.conn.SetReadDeadline(time.Time{})
}

func (this *ConnExt) ClearWriteDeadline() {
	this.conn.SetWriteDeadline(time.Time{})
}

func (this *ConnExt) ClearDeadline() {
	this.conn.SetDeadline(time.Time{})
}

func (this *ConnExt) WriteString(s string) (int, error) {
	return this.Write([]byte(s))
}
