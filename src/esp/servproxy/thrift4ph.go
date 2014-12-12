package servproxy

import (
	"bmautil/connutil"
	"bytes"
	"encoding/binary"
	"fmt"
	"golua"
	"io"
	"logger"
	"net"
)

const (
	VERSION_MASK = 0xffff0000
	VERSION_1    = 0x80010000
)

type ThriftPortHandler int

func init() {
	AddPortHandler("thrift", ThriftPortHandler(0))
}

func (this ThriftPortHandler) readFrameInfo(r io.Reader) (int, error) {
	buf := []byte{0, 0, 0, 0}
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	size := int(binary.BigEndian.Uint32(buf))
	if size < 0 {
		return 0, fmt.Errorf("Read a negative frame size (%d)", size)
	}
	return size, nil
}

func (this ThriftPortHandler) readByte(r io.Reader) (value byte, n int, err error) {
	buf := []byte{0}
	n, err = r.Read(buf)
	return buf[0], n, err
}

func (this ThriftPortHandler) readI32(r io.Reader) (value int32, n int, err error) {
	buf := []byte{0, 0, 0, 0}
	n, err = io.ReadFull(r, buf)
	if err != nil {
		return 0, 0, err
	}
	value = int32(binary.BigEndian.Uint32(buf))
	return value, n, nil
}

func (this ThriftPortHandler) readStringBody(r io.Reader, size int) (value string, n int, err error) {
	if size < 0 {
		return "", 0, nil
	}
	isize := int(size)
	buf := make([]byte, isize)
	n, e := io.ReadFull(r, buf)
	if e != nil {
		return "", 0, e
	}
	return string(buf), n, nil
}

func (this ThriftPortHandler) readString(r io.Reader) (value string, n int, err error) {
	size, l, e := this.readI32(r)
	if e != nil {
		return "", 0, e
	}
	s, l2, e2 := this.readStringBody(r, int(size))
	if e2 != nil {
		return "", 0, e2
	}
	return s, l2 + l, nil
}

func (this ThriftPortHandler) readMessageHeader(r io.Reader) (name string, typeId int32, seqId int32, n int, err error) {
	size, c, e := this.readI32(r)
	if e != nil {
		return "", 0, 0, 0, e
	}
	if size < 0 {
		typeId = int32(size & 0x0ff)
		version := int64(int64(size) & VERSION_MASK)
		if version != VERSION_1 {
			return "", 0, 0, 0, fmt.Errorf("Bad version(%d) in ReadMessageBegin", version)
		}
		l := 0
		name, l, e = this.readString(r)
		if e != nil {
			return "", 0, 0, 0, e
		}
		c += l
		seqId, l, e = this.readI32(r)
		if e != nil {
			return "", 0, 0, 0, e
		}
		c += l
		return name, typeId, seqId, c, nil
	}
	name, l2, e2 := this.readStringBody(r, int(size))
	if e2 != nil {
		return "", 0, 0, 0, e2
	}
	c += l2
	b, l3, e3 := this.readByte(r)
	if e3 != nil {
		return "", 0, 0, 0, e3
	}
	c += l3
	typeId = int32(b)
	seqId, l4, e4 := this.readI32(r)
	if e4 != nil {
		return "", 0, 0, 0, e4
	}
	c += l4
	return name, typeId, seqId, c, nil
}

func (this ThriftPortHandler) writeFrameInfo(w io.Writer, sz int) error {
	buf := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(buf, uint32(sz))
	_, err := w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (this ThriftPortHandler) writeString(w io.Writer, value string) error {
	err := this.writeI32(w, int32(len(value)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(value))
	return err
}

func (this ThriftPortHandler) writeByte(w io.Writer, value byte) error {
	v := []byte{value}
	_, err := w.Write(v)
	return err
}

func (this ThriftPortHandler) writeI32(w io.Writer, value int32) error {
	v := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(v, uint32(value))
	_, err := w.Write(v)
	return err
}

func (this ThriftPortHandler) writeI16(w io.Writer, value int16) error {
	v := []byte{0, 0}
	binary.BigEndian.PutUint16(v, uint16(value))
	_, e := w.Write(v)
	return e
}

func (this ThriftPortHandler) writeMessageHeader(w io.Writer, name string, typeId int32, seqId int32) error {
	if true {
		version := uint32(VERSION_1) | uint32(typeId)
		e := this.writeI32(w, int32(version))
		if e != nil {
			return e
		}
		e = this.writeString(w, name)
		if e != nil {
			return e
		}
		e = this.writeI32(w, seqId)
		return e
	} else {
		e := this.writeString(w, name)
		if e != nil {
			return e
		}
		e = this.writeByte(w, byte(typeId))
		if e != nil {
			return e
		}
		e = this.writeI32(w, seqId)
		return e
	}
}

func (this ThriftPortHandler) writeError(w io.Writer, req *ThriftProxyReq, err error) {
	str := err.Error()
	// err = oprot.WriteFieldBegin("message", STRING, 1)
	this.writeByte(w, 11)
	this.writeI16(w, 1)
	this.writeString(w, str)

	// err = oprot.WriteFieldBegin("type", I32, 2)
	this.writeByte(w, 8)
	this.writeI16(w, 2)
	this.writeI32(w, 6)

	// err = oprot.WriteFieldStop()
	this.writeByte(w, 0)
}

func (this ThriftPortHandler) Handle(s *Service, port *PortObj, conn net.Conn) {
	defer conn.Close()
	var dbg connutil.ConnDebuger
	if true {
		dbg = ConnDebuger
	}
	rconn := connutil.NewConnExt(conn, dbg)
	for {
		sz, err := this.readFrameInfo(rconn)
		if err != nil {
			if err == io.EOF {
				logger.Debug(tag, "%s closed", rconn)
			} else {
				logger.Warn(tag, "%s readFrameInfo fail - %s", rconn, err)
			}
			return
		}
		logger.Debug(tag, "%s readFrameInfo - %d", rconn, sz)
		req := new(ThriftProxyReq)
		req.s = s
		req.conn = rconn

		req.size = sz
		req.readed = 0
		name, tid, sid, n, err1 := this.readMessageHeader(rconn)
		if err1 != nil {
			logger.Warn(tag, "%s readMessageHeader fail - %s", rconn, err1)
			return
		}
		logger.Debug(tag, "%s readMessageHeader - %s, %d, %d/%d", rconn, name, tid, sid, n)
		req.name = name
		req.typeId = tid
		req.seqId = sid
		req.hsize = n
		req.readed += n

		_, errE := s.Execute(port, golua.NewGOO(req, gooThriftProxyReq(0)), req)
		if errE != nil {
			logger.Warn(tag, "%s execute fail - %s", dc, errE)
			return
		}
	}
}

func (this ThriftPortHandler) AnswerError(port *PortObj, preq interface{}, err error) error {
	req, ok := preq.(*ThriftProxyReq)
	if ok && !req.responsed {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "answer %s error - %s", req.conn.RemoteAddr(), err)
		}
		req.responsed = true
		buf := bytes.NewBuffer([]byte{})
		this.writeMessageHeader(buf, "", 3, req.seqId)
		this.writeError(buf, req, err)
		buf2 := bytes.NewBuffer(make([]byte, 0, 4+buf.Len()))
		this.writeFrameInfo(buf2, buf.Len())
		buf2.Write(buf.Bytes())
		req.conn.Write(buf2.Bytes())
	}
	return err
}
