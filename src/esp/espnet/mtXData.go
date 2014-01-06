package espnet

import (
	"bmautil/byteutil"
	"errors"
	"esp/espnet/protpack"
)

type mt_xdata byte

type struct_xdata struct {
	xid     int
	value   interface{}
	reader  *byteutil.BytesBufferReader
	encoder protpack.Encoder
}

func (this *struct_xdata) Value(dec protpack.Decoder) (interface{}, error) {
	if this.value != nil {
		return this.value, nil
	}
	if this.reader != nil {
		var err error
		if dec == nil {
			dec = Coders.Varinat
		}
		this.value, err = dec.Decode(this.reader)
		if err != nil {
			return nil, err
		}
		this.reader = nil
	}
	return this.value, nil
}

func (O mt_xdata) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if mv, ok := v.(*struct_xdata); ok {
		Coders.Int.DoEncode(w, mv.xid)
		var c protpack.Encoder
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

func (O mt_xdata) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	c := Coders.Int
	v1, err := c.Decode(r)
	if err != nil {
		return nil, err
	}
	xid, ok := v1.(int)
	if !ok {
		return nil, errors.New("not messageValue.xid:int")
	}
	return &struct_xdata{xid, nil, r, nil}, nil
}

func (O mt_xdata) MT() byte {
	return MT_XDATA
}

func (O mt_xdata) is(e *protpack.Frame, xid int) (*struct_xdata, error) {
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

func (O mt_xdata) xid(e *protpack.Frame) (int, error) {
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

func (O mt_xdata) value(e *protpack.Frame, dec protpack.Decoder) (interface{}, error) {
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

func (O mt_xdata) Add(p *protpack.Package, xid int, value interface{}, enc protpack.Encoder) {
	f := protpack.NewFrameV(O.MT(), &struct_xdata{xid, value, nil, enc}, O)
	p.PushBack(f)
}

func (O mt_xdata) Get(p *protpack.Package, xid int, dec protpack.Decoder) (interface{}, error) {
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

func (O mt_xdata) Remove(p *protpack.Package, xid int) {
	p.RemoveFrame(func(e *protpack.Frame) (bool, bool) {
		if mv, _ := O.is(e, xid); mv != nil {
			return true, false
		}
		return false, false
	})
}

func (O mt_xdata) Iterator(p *protpack.Package) *mtXDataIterator {
	r := new(mtXDataIterator)
	r.frame = p.Front()
	r.mt = O
	return r
}

type mtXDataIterator struct {
	frame *protpack.Frame
	mt    mt_xdata
}

func (this *mtXDataIterator) moveFirst() {
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

func (this *mtXDataIterator) IsEnd() bool {
	if this.frame == nil {
		return true
	}
	return false
}

func (this *mtXDataIterator) Next() {
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

func (this *mtXDataIterator) Xid() int {
	if this.frame == nil {
		return 0
	}
	xid, err := this.mt.xid(this.frame)
	if err != nil {
		return 0
	}
	return xid
}

func (this *mtXDataIterator) Value(dec protpack.Decoder) (interface{}, error) {
	if this.frame == nil {
		return nil, errors.New("end")
	}
	return this.mt.value(this.frame, dec)
}
