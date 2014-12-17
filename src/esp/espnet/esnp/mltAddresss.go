package esnp

import (
	"errors"
	"fmt"
)

type struct_address_value struct {
	annotation int
	value      string
}

func (this *struct_address_value) String() string {
	return fmt.Sprintf("%d:%s", this.annotation, this.value)
}

// Address
type mlt_address byte

func (O mlt_address) Encode(w EncodeWriter, v interface{}) error {
	if mv, ok := v.(*struct_address_value); ok {
		err0 := Coders.Int.DoEncode(w, mv.annotation)
		if err0 != nil {
			return err0
		}
		err1 := Coders.LenString.DoEncode(w, mv.value)
		if err1 != nil {
			return err1
		}
		return nil
	}
	return errors.New("not struct_address_value")
}

func (O mlt_address) Decode(r DecodeReader) (interface{}, error) {
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

func (O mlt_address) MT() byte {
	return byte(O)
}

func (O mlt_address) is(e *MessageLine, ann int) (*struct_address_value, error) {
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

func (O mlt_address) Set(p *Message, ann int, value string) {
	p.RemoveMessageLine(func(e *MessageLine) (bool, bool) {
		if mv, _ := O.is(e, ann); mv != nil {
			return true, false
		}
		return false, false
	})
	f := NewMessageLineV(O.MT(), &struct_address_value{ann, value}, O)
	p.PushFront(f)
}

func (O mlt_address) Get(p *Message, ann int) (string, error) {
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

func (O mlt_address) Remove(p *Message, ann int) {
	p.RemoveMessageLine(func(e *MessageLine) (bool, bool) {
		if mv, _ := O.is(e, ann); mv != nil {
			return true, false
		}
		return false, false
	})
}

func (O mlt_address) Clear(p *Message) {
	p.RemoveMessageLine(func(e *MessageLine) (bool, bool) {
		if e.MessageType() == O.MT() {
			return true, false
		}
		return false, false
	})
}

func (O mlt_address) List(p *Message) []int {
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
