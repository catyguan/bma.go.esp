package esnp

import (
	"errors"
	"fmt"
)

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

type mlt_xdata byte

func (O mlt_xdata) Encode(w EncodeWriter, v interface{}) error {
	if mv, ok := v.(*struct_xdata); ok {
		err := Coders.Int.DoEncode(w, mv.xid)
		if err != nil {
			return err
		}
		var c Encoder
		c = mv.encoder
		if c == nil {
			c = Coders.Varinat
		}
		err = c.Encode(w, mv.value)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("not struct_xdata")
}

func (O mlt_xdata) Decode(r DecodeReader) (interface{}, error) {
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

func (O mlt_xdata) MT() byte {
	return MLT_XDATA
}

func (O mlt_xdata) is(e *MessageLine, xid int) (*struct_xdata, error) {
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

func (O mlt_xdata) xid(e *MessageLine) (int, error) {
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

func (O mlt_xdata) value(e *MessageLine, dec Decoder) (interface{}, error) {
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

func (O mlt_xdata) Add(p *Message, xid int, value interface{}, enc Encoder) {
	f := NewMessageLineV(O.MT(), &struct_xdata{xid, value, nil, enc}, O)
	p.PushBack(f)
}

func (O mlt_xdata) Get(p *Message, xid int, dec Decoder) (interface{}, error) {
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

func (O mlt_xdata) Remove(p *Message, xid int) {
	p.RemoveMessageLine(func(e *MessageLine) (bool, bool) {
		if mv, _ := O.is(e, xid); mv != nil {
			return true, false
		}
		return false, false
	})
}

func (O mlt_xdata) Iterator(p *Message) *XDataIterator {
	r := new(XDataIterator)
	r.line = p.Front()
	r.mt = O
	return r
}

type XDataIterator struct {
	line *MessageLine
	mt   mlt_xdata
}

func (this *XDataIterator) moveFirst() {
	for {
		if this.line == nil {
			return
		}
		if this.line.MessageType() == this.mt.MT() {
			return
		}
		this.line = this.line.Next()
	}
}

func (this *XDataIterator) IsEnd() bool {
	if this.line == nil {
		return true
	}
	return false
}

func (this *XDataIterator) Next() {
	if this.IsEnd() {
		return
	}
	for {
		this.line = this.line.Next()
		if this.line == nil {
			return
		}
		if this.line.MessageType() == this.mt.MT() {
			return
		}
	}
}

func (this *XDataIterator) Xid() int {
	if this.line == nil {
		return 0
	}
	xid, err := this.mt.xid(this.line)
	if err != nil {
		return 0
	}
	return xid
}

func (this *XDataIterator) Value(dec Decoder) (interface{}, error) {
	if this.line == nil {
		return nil, errors.New("end")
	}
	return this.mt.value(this.line, dec)
}

// MessageXData
type MessageXData struct {
	msg   *Message
	coder mlt_xdata
}

func (this *MessageXData) Add(xid int, value interface{}, enc Encoder) {
	this.coder.Add(this.msg, xid, value, enc)
}

func (this *MessageXData) Iterator() *XDataIterator {
	return this.coder.Iterator(this.msg)
}
