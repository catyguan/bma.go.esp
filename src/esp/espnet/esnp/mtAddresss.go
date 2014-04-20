package esnp

import (
	"errors"
)

// Address
type addrCoder byte

type struct_address_value struct {
	annotation int
	value      string
}

func (O addrCoder) Encode(w EncodeWriter, v interface{}) error {
	if mv, ok := v.(*struct_address_value); ok {
		Coders.Int.DoEncode(w, mv.annotation)
		Coders.LenString.DoEncode(w, mv.value)
		return nil
	}
	return errors.New("not struct_address_value")
}

func (O addrCoder) Decode(r DecodeReader) (interface{}, error) {
	v1, err := Coders.Int.DoDecode(r)
	if err != nil {
		return nil, err
	}
	v2, err2 := Coders.LenString.DoDecode(r, 0)
	if err2 != nil {
		return nil, err2
	}
	return &struct_address_value{v1, v2}, nil
}

func (O addrCoder) MT() byte {
	return byte(O)
}

func (O addrCoder) is(e *Frame, ann int) (*struct_address_value, error) {
	if e.MessageType() == O.MT() {
		v, err := e.Value(O)
		if err != nil {
			return nil, err
		}
		if mv, ok := v.(*struct_address_value); ok {
			if mv.annotation == ann {
				return mv, nil
			}
		} else {
			return nil, notValueErr
		}
	}
	return nil, nil
}

func (O addrCoder) Set(p *Package, ann int, value string) {
	p.RemoveFrame(func(e *Frame) (bool, bool) {
		if mv, _ := O.is(e, ann); mv != nil {
			return true, false
		}
		return false, false
	})
	f := NewFrameV(O.MT(), &struct_address_value{ann, value}, O)
	p.PushFront(f)
}

func (O addrCoder) Get(p *Package, ann int) (string, error) {
	for e := p.Front(); e != nil; e = e.Next() {
		mv, err := O.is(e, ann)
		if err != nil {
			continue
		}
		if mv != nil {
			return mv.value, nil
		}
	}
	return "", nil
}

func (O addrCoder) Remove(p *Package, ann int) {
	p.RemoveFrame(func(e *Frame) (bool, bool) {
		if mv, _ := O.is(e, ann); mv != nil {
			return true, false
		}
		return false, false
	})
}

func (O addrCoder) Clear(p *Package) {
	p.RemoveFrame(func(e *Frame) (bool, bool) {
		if e.MessageType() == O.MT() {
			return true, false
		}
		return false, false
	})
}

func (O addrCoder) List(p *Package) []int {
	r := make([]int, 0)
	mt := O.MT()
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() != mt {
			continue
		}
		v, err := e.Value(O)
		if err != nil {
			continue
		}
		if mv, ok := v.(*struct_address_value); ok {
			r = append(r, mv.annotation)
		}
	}
	return r
}
