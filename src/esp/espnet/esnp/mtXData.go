package esnp

import (
	"errors"
	"fmt"
)

type mt_xdata byte

type struct_xdata struct {
	xid     int
	value   interface{}
	remain  []byte
	encoder Encoder
}

func (this *struct_xdata) String() string {
	return fmt.Sprintf("%d=%v", this.xid, this.value)
}

func (this *struct_xdata) Value(dec Decoder) (interface{}, error) {
	if this.value != nil {
		return this.value, nil
	}
	if this.remain != nil {
		var err error
		if dec == nil {
			dec = Coders.Varinat
		}
		var bdr BytesDecodeReader
		bdr.data = this.remain
		this.value, err = dec.Decode(&bdr)
		if err != nil {
			return nil, err
		}
		this.remain = nil
	}
	return this.value, nil
}

func (O mt_xdata) Encode(w EncodeWriter, v interface{}) error {
	if mv, ok := v.(*struct_xdata); ok {
		Coders.Int.DoEncode(w, mv.xid)
		var c Encoder
		c = mv.encoder
		if c == nil {
			c = Coders.Varinat
		}
		err := c.Encode(w, mv.value)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("not struct_xdata")
}

func (O mt_xdata) Decode(r DecodeReader) (interface{}, error) {
	c := Coders.Int
	v1, err := c.Decode(r)
	if err != nil {
		return nil, err
	}
	xid, ok := v1.(int)
	if !ok {
		return nil, errors.New("not messageValue.xid:int")
	}
	return &struct_xdata{xid, nil, r.Remain(), nil}, nil
}

func (O mt_xdata) MT() byte {
	return MT_XDATA
}

func (O mt_xdata) is(e *Frame, xid int) (*struct_xdata, error) {
	if e.MessageType() == O.MT() {
		v, err := e.Value(O)
		if err != nil {
			return nil, err
		}
		if mv, ok := v.(*struct_xdata); ok {
			if mv.xid == xid {
				return mv, nil
			}
		} else {
			return nil, notValueErr
		}
	}
	return nil, nil
}

func (O mt_xdata) xid(e *Frame) (int, error) {
	if e.MessageType() == O.MT() {
		v, err := e.Value(O)
		if err != nil {
			return 0, err
		}
		if mv, ok := v.(*struct_xdata); ok {
			return mv.xid, nil
		}
	}
	return 0, notValueErr
}

func (O mt_xdata) value(e *Frame, dec Decoder) (interface{}, error) {
	if e.MessageType() == O.MT() {
		v, err := e.Value(O)
		if err != nil {
			return nil, err
		}
		if mv, ok := v.(*struct_xdata); ok {
			return mv.Value(dec)
		}
	}
	return nil, notValueErr
}

func (O mt_xdata) Add(p *Package, xid int, value interface{}, enc Encoder) {
	f := NewFrameV(O.MT(), &struct_xdata{xid, value, nil, enc}, O)
	p.PushBack(f)
}

func (O mt_xdata) Get(p *Package, xid int, dec Decoder) (interface{}, error) {
	for e := p.Front(); e != nil; e = e.Next() {
		mv, err := O.is(e, xid)
		if err != nil {
			continue
		}
		if mv != nil {
			return mv.Value(dec)
		}
	}
	return nil, nil
}

func (O mt_xdata) Remove(p *Package, xid int) {
	p.RemoveFrame(func(e *Frame) (bool, bool) {
		if mv, _ := O.is(e, xid); mv != nil {
			return true, false
		}
		return false, false
	})
}

func (O mt_xdata) Iterator(p *Package) *XDataIterator {
	r := new(XDataIterator)
	r.frame = p.Front()
	r.mt = O
	return r
}

type XDataIterator struct {
	frame *Frame
	mt    mt_xdata
}

func (this *XDataIterator) moveFirst() {
	for {
		if this.frame == nil {
			return
		}
		if this.frame.MessageType() == this.mt.MT() {
			return
		}
		this.frame = this.frame.Next()
	}
}

func (this *XDataIterator) IsEnd() bool {
	if this.frame == nil {
		return true
	}
	return false
}

func (this *XDataIterator) Next() {
	if this.IsEnd() {
		return
	}
	for {
		this.frame = this.frame.Next()
		if this.frame == nil {
			return
		}
		if this.frame.MessageType() == this.mt.MT() {
			return
		}
	}
}

func (this *XDataIterator) Xid() int {
	if this.frame == nil {
		return 0
	}
	xid, err := this.mt.xid(this.frame)
	if err != nil {
		return 0
	}
	return xid
}

func (this *XDataIterator) Value(dec Decoder) (interface{}, error) {
	if this.frame == nil {
		return nil, errors.New("end")
	}
	return this.mt.value(this.frame, dec)
}

// MessageXData
type MessageXData struct {
	msg   *Message
	coder mt_xdata
}

func (this *MessageXData) Add(xid int, value interface{}, enc Encoder) {
	this.coder.Add(this.msg.ToPackage(), xid, value, enc)
}

func (this *MessageXData) Iterator() *XDataIterator {
	return this.coder.Iterator(this.msg.ToPackage())
}
