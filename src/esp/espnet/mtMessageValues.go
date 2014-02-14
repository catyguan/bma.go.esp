package espnet

import (
	"bmautil/byteutil"
	Coders "bmautil/coder"
	"errors"
	"esp/espnet/protpack"
)

type mvCoder byte

type struct_message_value struct {
	name    string
	value   interface{}
	reader  *byteutil.BytesBufferReader
	encoder protpack.Encoder
}

func (this *struct_message_value) Value(dec protpack.Decoder) (interface{}, error) {
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

func (O mvCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if mv, ok := v.(*struct_message_value); ok {
		Coders.LenString.DoEncode(w, mv.name)
		c := mv.encoder
		if c == nil {
			c = Coders.Varinat
		}
		err := c.Encode(w, mv.value)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("not struct_message_value")
}

func (O mvCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	s1, err := Coders.LenString.DoDecode(r, 0)
	if err != nil {
		return nil, err
	}
	return &struct_message_value{s1, nil, r, nil}, nil
}

func (O mvCoder) MT() byte {
	return byte(O)
}

func (O mvCoder) is(e *protpack.Frame, key string) (*struct_message_value, error) {
	if e.MessageType() == O.MT() {
		v, err := e.Value(O)
		if err != nil {
			return nil, err
		}
		if mv, ok := v.(*struct_message_value); ok {
			if mv.name == key {
				return mv, nil
			}
		} else {
			return nil, notValueErr
		}
	}
	return nil, nil
}

func (O mvCoder) Set(p *protpack.Package, key string, value interface{}, enc protpack.Encoder) {
	p.RemoveFrame(func(e *protpack.Frame) (bool, bool) {
		if mv, _ := O.is(e, key); mv != nil {
			return true, false
		}
		return false, false
	})
	f := protpack.NewFrameV(O.MT(), &struct_message_value{key, value, nil, enc}, O)
	if byte(O) == MT_DATA {
		p.PushBack(f)
	} else {
		p.PushFront(f)
	}
}

func (O mvCoder) Get(p *protpack.Package, key string, dec protpack.Decoder) (interface{}, error) {
	for e := p.Front(); e != nil; e = e.Next() {
		mv, err := O.is(e, key)
		if err != nil {
			continue
		}
		if mv != nil {
			return mv.Value(dec)
		}
	}
	return nil, nil
}

func (O mvCoder) Pop(p *protpack.Package, key string, dec protpack.Decoder) (interface{}, error) {
	for e := p.Front(); e != nil; e = e.Next() {
		mv, err := O.is(e, key)
		if err != nil {
			continue
		}
		if mv != nil {
			p.Remove(e)
			return mv.Value(dec)
		}
	}
	return nil, nil
}

func (O mvCoder) Remove(p *protpack.Package, key string) {
	p.RemoveFrame(func(e *protpack.Frame) (bool, bool) {
		if mv, _ := O.is(e, key); mv != nil {
			return true, false
		}
		return false, false
	})
}

func (O mvCoder) List(p *protpack.Package) []string {
	r := make([]string, 0)
	mt := O.MT()
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == mt {
			v, err := e.Value(O)
			if err != nil {
				continue
			}
			if mv, ok := v.(*struct_message_value); ok {
				r = append(r, mv.name)
			}
		}
	}
	return r
}
