package esnp

import (
	"io"
)

type BytesDecodeReader struct {
	data []byte
	pos  int
}

func NewBytesDecodeReader(b []byte) *BytesDecodeReader {
	r := new(BytesDecodeReader)
	r.data = b
	return r
}

func (this *BytesDecodeReader) ReadAll() ([]byte, error) {
	if this.pos < len(this.data) {
		r := this.data[this.pos:]
		this.pos = len(this.data)
		return r, nil
	}
	return []byte{}, nil
}

func (this *BytesDecodeReader) Read(p []byte) (n int, err error) {
	l := len(p)
	if this.pos+l-1 < len(this.data) {
		copy(p, this.data[this.pos:this.pos+l])
		this.pos = this.pos + l
		return l, nil
	}
	return 0, io.EOF
}

func (this *BytesDecodeReader) ReadByte() (byte, error) {
	if this.pos < len(this.data) {
		r := this.data[this.pos]
		this.pos = this.pos + 1
		return r, nil
	}
	return 0, io.EOF
}

func (this *BytesDecodeReader) Remain() int {
	return len(this.data) - this.pos
}

// BytesEncodeWriter
type BytesEncodeWriter struct {
	data    []byte
	pos     int
	linepos int
}

func NewBytesEncodeWriter(b []byte) *BytesEncodeWriter {
	r := new(BytesEncodeWriter)
	r.data = b
	return r
}

func (this *BytesEncodeWriter) grow(sz int) {
	if this.data == nil {
		this.data = make([]byte, 64)
	}
	if this.pos+sz > cap(this.data) {
		gl := ((sz+this.pos)/64 + 1) * 64
		buf := make([]byte, gl)
		copy(buf, this.data[:this.pos])
		this.data = buf
	}
}

func (this *BytesEncodeWriter) Write(b []byte) (int, error) {
	l := len(b)
	this.grow(l)
	copy(this.data[this.pos:], b)
	this.pos = this.pos + l
	return l, nil
}
func (this *BytesEncodeWriter) WriteByte(b byte) error {
	this.grow(1)
	this.data[this.pos] = b
	this.pos = this.pos + 1
	return nil
}

func (this *BytesEncodeWriter) WriteLine(mt byte, data []byte) error {
	l := 0
	if data != nil {
		l = len(data)
	}
	MessageLineHeaderWrite(this, mt, l)
	if l > 0 {
		this.Write(data)
	}
	this.linepos = 0
	return nil
}
func (this *BytesEncodeWriter) NewLine() error {
	this.linepos = this.pos
	this.grow(size_FHEADER)
	this.pos = this.pos + size_FHEADER
	return nil
}
func (this *BytesEncodeWriter) EndLine(mt byte) error {
	p := this.linepos
	old := this.pos
	this.pos = p
	sz := old - p - size_FHEADER
	MessageLineHeaderWrite(this, mt, sz)
	this.pos = old
	return nil
}

func (this *BytesEncodeWriter) ToBytes() []byte {
	return this.data[:this.pos]
}
