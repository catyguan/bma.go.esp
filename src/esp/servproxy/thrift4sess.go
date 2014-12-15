package servproxy

import (
	"bmautil/conndialpool"
	"bmautil/connutil"
	"io"
	"time"
)

type ThriftRemoteSession struct {
	remote *RemoteObj
	pool   *conndialpool.DialPool
	conn   *connutil.ConnExt

	readStatus int
	remain     int
	hbuf       []byte
	buf        []byte
	forceClose bool
	fail       bool
}

func (this *ThriftRemoteSession) String() string {
	return this.conn.String()
}

func (this *ThriftRemoteSession) BeginWrite() error {
	return nil
}

func (this *ThriftRemoteSession) Write(b []byte) error {
	_, err := this.conn.Write(b)
	return err
}

func (this *ThriftRemoteSession) EndWrite() {
	return
}

func (this *ThriftRemoteSession) BeginRead(deadline time.Time) error {
	if !deadline.IsZero() {
		this.conn.SetReadDeadline(deadline)
	}
	return nil
}

func (this *ThriftRemoteSession) Read() (bool, []byte, error) {
	if this.readStatus == 0 {
		buf := []byte{0, 0, 0, 0}
		sz, err := OThriftProtocol.readFrameInfoB(this.conn, buf)
		if err != nil {
			return false, nil, err
		}
		this.remain = sz
		this.readStatus = 1
		return true, buf, nil
	}

	sz := this.remain
	l := 4 * 1024
	if this.buf == nil {
		this.buf = make([]byte, l)
	}
	if sz <= 0 {
		return false, nil, nil
	}
	if sz > l {
		sz = l
	}
	rbuf := this.buf[:sz]
	n, err := this.conn.Read(rbuf)
	if n > 0 {
		this.remain -= n
		b := rbuf[:n]
		return true, b, nil
	}
	if err != nil {
		if err == io.EOF {
			return true, nil, nil
		}
		return false, nil, err
	}
	return true, rbuf[:0], nil
}

func (this *ThriftRemoteSession) EndRead() {
	this.buf = nil
	this.conn.SetReadDeadline(time.Time{})
}

func (this *ThriftRemoteSession) Fail() {
	this.fail = true
	return
}

func (this *ThriftRemoteSession) ForceClose() {
	this.forceClose = true
	return
}

func (this *ThriftRemoteSession) Finish() {
	if this.fail {
		this.forceClose = true
		// this.remote.Fail()
	}
	if this.forceClose {
		this.pool.CloseConn(this.conn)
	} else {
		this.conn.Debuger = nil
		this.pool.ReturnConn(this.conn)
	}
	return
}
