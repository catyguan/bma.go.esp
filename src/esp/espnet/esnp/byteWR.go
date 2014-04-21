package esnp

import

// BytesDecodeReader
"io"

type BytesDecodeReader struct {
	data []byte
	pos  int
}

func (this *BytesDecodeReader) ReadAll() []byte {
	if this.pos < len(this.data) {
		r := this.data[this.pos:]
		this.pos = len(this.data)
		return r
	}
	return []byte{}
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

func (this *BytesDecodeReader) Remain() []byte {
	if this.pos < len(this.data) {
		r := this.data[this.pos:]
		this.pos = len(this.data)
		return r
	}
	return nil
}

// BytesEncodeWriter
type BytesEncodeWriter struct {
	data []byte
	pos  int
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

func (this *BytesEncodeWriter) WriteFrame(mt byte, data []byte) error {
	l := 0
	if data != nil {
		l = len(data)
	}
	headerWrite(this, mt, l)
	if l > 0 {
		this.Write(data)
	}
	return nil
}
func (this *BytesEncodeWriter) NewFrame() (int, error) {
	r := this.pos
	this.grow(size_FHEADER)
	this.pos = this.pos + size_FHEADER
	return r, nil
}
func (this *BytesEncodeWriter) EndFrame(p int, mt byte) error {
	old := this.pos
	this.pos = p
	sz := old - p - size_FHEADER
	headerWrite(this, mt, sz)
	this.pos = old
	return nil
}

func (this *BytesEncodeWriter) ToBytes() []byte {
	return this.data[:this.pos]
}
