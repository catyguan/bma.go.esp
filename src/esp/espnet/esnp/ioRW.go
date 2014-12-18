package esnp

import (
	"bytes"
	"io"
)

type IODecodeReader struct {
	o  io.Reader
	bo io.ByteReader
}

func NewIODecodeReader(o io.Reader, bo io.ByteReader) *IODecodeReader {
	r := new(IODecodeReader)
	r.o = o
	r.bo = bo
	return r
}

func (this *IODecodeReader) Read(p []byte) (n int, err error) {
	n1, err1 := this.o.Read(p)
	if n1 > 0 {
		return n1, nil
	}
	return n1, err1
}

func (this *IODecodeReader) ReadByte() (byte, error) {
	if this.bo != nil {
		return this.bo.ReadByte()
	}
	p := []byte{0}
	n, err := this.o.Read(p)
	if n > 0 {
		return p[0], nil
	}
	return p[0], err
}

func (this *IODecodeReader) Remain() int {
	return -1
}

// IOEncodeWriter
type IOEncodeWriter struct {
	o     io.Writer
	bo    io.ByteWriter
	eline bool
	buf   *bytes.Buffer
}

func NewIOEncodeWriter(o io.Writer, bo io.ByteWriter) *IOEncodeWriter {
	r := new(IOEncodeWriter)
	r.o = o
	r.bo = bo
	return r
}

func (this *IOEncodeWriter) Write(b []byte) (int, error) {
	if this.eline {
		if this.buf == nil {
			this.buf = bytes.NewBuffer([]byte{})
		}
		return this.buf.Write(b)
	}
	return this.o.Write(b)
}
func (this *IOEncodeWriter) WriteByte(b byte) error {
	if this.eline {
		if this.buf == nil {
			this.buf = bytes.NewBuffer([]byte{})
		}
		return this.buf.WriteByte(b)
	}
	if this.bo != nil {
		return this.bo.WriteByte(b)
	} else {
		p := []byte{b}
		_, err := this.o.Write(p)
		return err
	}
}
func (this *IOEncodeWriter) WriteLine(mt byte, data []byte) error {
	this.eline = false
	l := 0
	if data != nil {
		l = len(data)
	}
	err1 := MessageLineHeaderWrite(this, mt, l)
	if err1 != nil {
		return err1
	}
	if l > 0 {
		_, err2 := this.Write(data)
		if err2 != nil {
			return err2
		}
	}
	return nil
}
func (this *IOEncodeWriter) NewLine() error {
	this.eline = true
	if this.buf != nil {
		this.buf.Reset()
	}
	return nil
}
func (this *IOEncodeWriter) EndLine(mt byte) error {
	this.eline = false
	sz := 0
	if this.buf != nil {
		sz = this.buf.Len()
	}
	MessageLineHeaderWrite(this, mt, sz)
	_, err := this.Write(this.buf.Bytes())
	return err
}
