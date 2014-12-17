package esnp

import (
	"errors"
)

type Flag int32

type mlt_flag int

func (O mlt_flag) Encode(w EncodeWriter, v interface{}) error {
	if o, ok := v.(Flag); ok {
		return Coders.Int32.DoEncode(w, int32(o))
	}
	return errors.New("not Flag")
}

func (O mlt_flag) Decode(r DecodeReader) (interface{}, error) {
	v, err := Coders.Int32.DoDecode(r)
	if err != nil {
		return nil, err
	}
	return Flag(v), nil
}

func (O mlt_flag) Has(p *Message, f Flag) bool {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_FLAG {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(Flag); ok {
					return f == rv
				}
			}
			break
		}
	}
	return false
}

func (O mlt_flag) Remove(p *Message, f Flag) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_FLAG {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(Flag); ok {
					if rv == f {
						p.Remove(e)
					}
				}
			}
		}
	}
}

func (O mlt_flag) Set(p *Message, v Flag) {
	if O.Has(p, v) {
		return
	}
	f := NewMessageLineV(MLT_FLAG, v, O)
	p.PushFront(f)
}
