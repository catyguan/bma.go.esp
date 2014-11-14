package esnp

import (
	"errors"
	"fmt"
)

type mvCoder byte

type struct_message_value struct {
	name    string
	value   interface{}
	remain  []byte
	encoder Encoder
}

func (this *struct_message_value) String() string {
	return fmt.Sprintf("%s=%v", this.name, this.value)
}

func (this *struct_message_value) Value(dec Decoder) (interface{}, error) {
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

func (O mvCoder) Encode(w EncodeWriter, v interface{}) error {
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

func (O mvCoder) Decode(r DecodeReader) (interface{}, error) {
	s1, err := Coders.LenString.DoDecode(r, 0)
	if err != nil {
		return nil, err
	}
	return &struct_message_value{s1, nil, r.Remain(), nil}, nil
}

func (O mvCoder) MT() byte {
	return byte(O)
}

func (O mvCoder) is(e *Frame, key string) (*struct_message_value, error) {
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

func (O mvCoder) Set(p *Package, key string, value interface{}, enc Encoder) {
	p.RemoveFrame(func(e *Frame) (bool, bool) {
		if mv, _ := O.is(e, key); mv != nil {
			return true, false
		}
		return false, false
	})
	f := NewFrameV(O.MT(), &struct_message_value{key, value, nil, enc}, O)
	if byte(O) == MT_DATA {
		p.PushBack(f)
	} else {
		p.PushFront(f)
	}
}

func (O mvCoder) Get(p *Package, key string, dec Decoder) (interface{}, error) {
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

func (O mvCoder) Pop(p *Package, key string, dec Decoder) (interface{}, error) {
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

func (O mvCoder) Remove(p *Package, key string) {
	p.RemoveFrame(func(e *Frame) (bool, bool) {
		if mv, _ := O.is(e, key); mv != nil {
			return true, false
		}
		return false, false
	})
}

func (O mvCoder) List(p *Package) []string {
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

func (O mvCoder) Map(p *Package) (map[string]interface{}, error) {
	r := make(map[string]interface{}, 0)
	mt := O.MT()
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == mt {
			v, err := e.Value(O)
			if err != nil {
				continue
			}
			if mv, ok := v.(*struct_message_value); ok {
				var err error
				r[mv.name], err = mv.Value(nil)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return r, nil
}
