package esnp

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	size_FHEADER = 4
)

func headerWrite(w EncodeWriter, mt byte, sz int) error {
	err := w.WriteByte(mt)
	if err != nil {
		return err
	}
	err = w.WriteByte(byte(sz >> 16))
	err = w.WriteByte(byte(sz >> 8))
	err = w.WriteByte(byte(sz))

	return nil
}

func headerRead(b []byte, pos int) (byte, int) {
	mt := byte(b[pos+0])
	sz := int(b[pos+3]) | int(b[pos+2])<<8 | int(b[pos+1])<<16
	return mt, sz
}

// Frame
type Frame struct {
	mtype byte
	data  []byte
	value interface{}

	pack *Package
	next *Frame
	prev *Frame

	encoder Encoder
}

func NewFrame(mt byte, data []byte) *Frame {
	r := new(Frame)
	r.mtype = mt
	r.data = data
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

func (this *Frame) Clone(mt byte) *Frame {
	r := new(Frame)
	if mt != 0 {
		r.mtype = mt
	} else {
		r.mtype = this.mtype
	}
	r.data = this.data
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

func (this *Frame) Encode(w EncodeWriter) error {
	if this.data == nil {
		if this.value != nil {
			var err error
			enc := this.encoder
			if enc == nil {
				return errors.New(fmt.Sprintf("unknow encoder %T", this.value))
			}
			p, err2 := w.NewFrame()
			if err2 != nil {
				return err2
			}
			err = enc.Encode(w, this.value)
			if err != nil {
				return err
			}
			return w.EndFrame(p, this.mtype)
		}
		return nil
	} else {
		return w.WriteFrame(this.mtype, this.data)
	}
}

func (this *Frame) RawData() []byte {
	return this.data
}

func (this *Frame) Value(dec Decoder) (interface{}, error) {
	if this.value == nil {
		if this.data != nil && len(this.data) > 0 {
			var bdr BytesDecodeReader
			bdr.data = this.data
			var err error
			this.value, err = dec.Decode(&bdr)
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
			sz = len(this.data)
		}
		buf.WriteString(fmt.Sprintf(",%d", sz))
		if this.data != nil {
			buf.WriteString(",[")
			for i := 0; i < 16; i++ {
				if i >= sz {
					break
				}
				buf.WriteString(fmt.Sprintf("%X", this.data[i]))
			}
			if sz > 16 {
				buf.WriteString("...")
			}
			buf.WriteString("]")
		}
	}
	buf.WriteString("]")
	return buf.String()
}
