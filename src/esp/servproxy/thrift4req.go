package servproxy

import (
	"bytes"
	"fmt"
	"io"
	"logger"
	"net"
	"time"
)

type ThriftProxyReq struct {
	s          *Service
	conn       net.Conn
	size       int
	readed     int
	hsize      int
	name       string
	typeId     int32
	seqId      int32
	responsed  bool
	oneway     bool
	write      bool
	deadline   time.Time
	readStatus int
	cache      [][]byte
}

func (this *ThriftProxyReq) Type() string {
	return "thrift"
}

func (this *ThriftProxyReq) BeginRead() error {
	this.readStatus = 0
	if !this.deadline.IsZero() {
		this.conn.SetReadDeadline(this.deadline)
	}
	return nil
}

func (this *ThriftProxyReq) Read() (bool, []byte, error) {
	if this.readStatus == 0 {
		buf := bytes.NewBuffer(make([]byte, 0, this.hsize))
		OThriftProtocol.writeMessageHeader(buf, this.name, this.typeId, this.seqId)
		sz := this.size
		l := buf.Len()
		if this.hsize != l {
			sz = sz - this.hsize + l
		}
		buf2 := bytes.NewBuffer(make([]byte, 0, 4+l))
		err1 := OThriftProtocol.writeFrameInfo(buf2, sz)
		if err1 != nil {
			return false, nil, err1
		}
		buf2.Write(buf.Bytes())
		this.readStatus = 1
		return true, buf2.Bytes(), nil
	}
	n := this.readStatus - 1
	if n < len(this.cache) {
		b := this.cache[n]
		this.readStatus++
		return true, b, nil
	}
	this.readStatus++
	l := 4 * 1024
	buf := make([]byte, l)
	sz := this.Remain()
	if sz == 0 {
		return false, nil, nil
	}
	if sz > l {
		sz = l
	}
	rbuf := buf[:sz]
	n, err := this.conn.Read(rbuf)
	if n > 0 {
		this.readed += n
		b := rbuf[:n]
		this.cache = append(this.cache, b)
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

func (this *ThriftProxyReq) EndRead() {
	this.cache = nil
	this.conn.SetReadDeadline(time.Time{})
}

func (this *ThriftProxyReq) Deadline() time.Time {
	return this.deadline
}

func (this *ThriftProxyReq) SetDeadline(tm time.Time) {
	this.deadline = tm
}

func (this *ThriftProxyReq) Finish() {
	logger.Debug(tag, "thrift '%s' request end", this)
}

func (this *ThriftProxyReq) CheckFlag(n string) bool {
	switch n {
	case PRF_WRITE:
		return this.write
	case PRF_NO_RESPONSE:
		return this.IsOneway()
	}
	return false
}

func (this *ThriftProxyReq) IsOneway() bool {
	if this.typeId == 4 {
		return true
	}
	return this.oneway
}

func (this *ThriftProxyReq) Remain() int {
	return this.size - this.readed
}

func (this *ThriftProxyReq) String() string {
	return fmt.Sprintf("[%s, %d, %d](%d:%d)", this.name, this.typeId, this.seqId, this.size, this.readed)
}
