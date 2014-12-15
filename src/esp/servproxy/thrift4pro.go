package servproxy

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ThriftProtocol int

var (
	OThriftProtocol = ThriftProtocol(0)
)

func (this ThriftProtocol) readFrameInfo(r io.Reader) (int, error) {
	buf := []byte{0, 0, 0, 0}
	return this.readFrameInfoB(r, buf)
}

func (this ThriftProtocol) readFrameInfoB(r io.Reader, buf []byte) (int, error) {
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	size := int(binary.BigEndian.Uint32(buf))
	if size < 0 {
		return 0, fmt.Errorf("Read a negative frame size (%d)", size)
	}
	return size, nil
}

func (this ThriftProtocol) readByte(r io.Reader) (value byte, n int, err error) {
	buf := []byte{0}
	n, err = r.Read(buf)
	return buf[0], n, err
}

func (this ThriftProtocol) readI32(r io.Reader) (value int32, n int, err error) {
	buf := []byte{0, 0, 0, 0}
	n, err = io.ReadFull(r, buf)
	if err != nil {
		return 0, 0, err
	}
	value = int32(binary.BigEndian.Uint32(buf))
	return value, n, nil
}

func (this ThriftProtocol) readStringBody(r io.Reader, size int) (value string, n int, err error) {
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

func (this ThriftProtocol) readString(r io.Reader) (value string, n int, err error) {
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

func (this ThriftProtocol) readMessageHeader(r io.Reader) (name string, typeId int32, seqId int32, n int, err error) {
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

func (this ThriftProtocol) writeFrameInfo(w io.Writer, sz int) error {
	buf := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(buf, uint32(sz))
	_, err := w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (this ThriftProtocol) writeString(w io.Writer, value string) error {
	err := this.writeI32(w, int32(len(value)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(value))
	return err
}

func (this ThriftProtocol) writeByte(w io.Writer, value byte) error {
	v := []byte{value}
	_, err := w.Write(v)
	return err
}

func (this ThriftProtocol) writeI32(w io.Writer, value int32) error {
	v := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(v, uint32(value))
	_, err := w.Write(v)
	return err
}

func (this ThriftProtocol) writeI16(w io.Writer, value int16) error {
	v := []byte{0, 0}
	binary.BigEndian.PutUint16(v, uint16(value))
	_, e := w.Write(v)
	return e
}

func (this ThriftProtocol) writeMessageHeader(w io.Writer, name string, typeId int32, seqId int32) error {
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

func (this ThriftProtocol) writeError(w io.Writer, req *ThriftProxyReq, err error) {
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
