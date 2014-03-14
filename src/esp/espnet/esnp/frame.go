package esnp

import (
	"bmautil/byteutil"
	"bytes"
	"errors"
	"fmt"
)

const (
	size_FHEADER = 4
)

// FHeader
type FHeader struct {
	MessageType byte
	Size        uint32
}

func (this *FHeader) Write(buf []byte, pos int) int {
	buf[pos+0] = byte(this.MessageType)
	buf[pos+1] = byte(this.Size >> 16)
	buf[pos+2] = byte(this.Size >> 8)
	buf[pos+3] = byte(this.Size)

	return pos + size_FHEADER
}

func (this *FHeader) Read(b []byte, pos int) int {
	this.MessageType = byte(b[pos+0])
	this.Size = uint32(b[pos+3]) | uint32(b[pos+2])<<8 | uint32(b[pos+1])<<16
	return pos + size_FHEADER
}

func (this *FHeader) ToBytes() []byte {
	buf := make([]byte, size_FHEADER)
	this.Write(buf, 0)
	return buf
}

func (this *FHeader) String() string {
	return fmt.Sprintf("%d[%d]", this.MessageType, this.Size)
}

// Frame
type Frame struct {
	mtype byte
	data  *byteutil.BytesBuffer
	value interface{}

	pack *Package
	next *Frame
	prev *Frame

	encoder Encoder
}

func NewFrame(mt byte, data *byteutil.BytesBuffer) *Frame {
	r := new(Frame)
	r.mtype = mt
	r.data = data
	return r
}

func NewFrameB(mt byte, data []byte) *Frame {
	b := byteutil.NewBytesBufferB(data)
	return NewFrame(mt, b)
}

func newFrameH(h FHeader) *Frame {
	r := new(Frame)
	r.mtype = h.MessageType
	r.data = byteutil.NewBytesBuffer()
	return r
}

func NewFrameV(mt byte, v interface{}, enc Encoder) *Frame {
	r := new(Frame)
	r.mtype = mt
	r.value = v
	if enc != nil {
		r.encoder = enc
	} else {
		if re, ok := v.(Encoder); ok {
			r.encoder = re
		}
	}
	return r
}

func (this *Frame) Clone(mt byte, cloneData bool) *Frame {
	r := new(Frame)
	if mt != 0 {
		r.mtype = mt
	} else {
		r.mtype = this.mtype
	}
	if cloneData {
		if this.data != nil && this.data.DataList != nil {
			r.data = byteutil.NewBytesBuffer()
			for _, b := range this.data.DataList {
				r.data.Add(b)
			}
		}
	} else {
		r.data = this.data
	}
	r.value = this.value
	r.encoder = this.encoder

	return r
}

func (this *Frame) MessageType() byte {
	return this.mtype
}

func (this *Frame) Next() *Frame {
	return this.next
}

func (this *Frame) Prev() *Frame {
	return this.prev
}

func (this *Frame) Data() (*byteutil.BytesBuffer, error) {
	if this.data == nil {
		if this.value != nil {
			buf := byteutil.NewBytesBuffer()
			w := buf.NewWriter()
			var err error
			enc := this.encoder
			if enc == nil {
				return nil, errors.New(fmt.Sprintf("unknow encoder %T", this.value))
			}
			err = enc.Encode(w, this.value)
			if err != nil {
				return nil, err
			}
			this.data = w.End()
		}
	}
	return this.data, nil
}

func (this *Frame) RawData() *byteutil.BytesBuffer {
	return this.data
}

func (this *Frame) Value(dec Decoder) (interface{}, error) {
	if this.value == nil {
		if this.data != nil && this.data.Len() > 0 {
			var err error
			this.value, err = dec.Decode(this.data.NewReader())
			if err != nil {
				return nil, err
			}
		}
	}
	return this.value, nil
}

func (this *Frame) RawValue() interface{} {
	return this.value
}

func (this *Frame) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("FRAME[%d", this.mtype))
	if this.value != nil {
		buf.WriteString(fmt.Sprintf(",%v", this.value))
	} else {
		sz := 0
		if this.data != nil {
			sz = this.data.DataSize()
		}
		buf.WriteString(fmt.Sprintf(",%d", sz))
		if this.data != nil {
			buf.WriteString(",[")
			buf.WriteString(this.data.TraceString(16))
			if sz > 16 {
				buf.WriteString("...")
			}
			buf.WriteString("]")
		}
	}
	buf.WriteString("]")
	return buf.String()
}
