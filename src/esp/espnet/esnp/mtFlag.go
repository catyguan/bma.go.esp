package esnp

import (
	"errors"
)

type MTFlag int32

type mt_flag int

func (O mt_flag) Encode(w EncodeWriter, v interface{}) error {
	if o, ok := v.(MTFlag); ok {
		Coders.Int32.DoEncode(w, int32(o))
		return nil
	}
	return errors.New("not MTFlag")
}

func (O mt_flag) Decode(r DecodeReader) (interface{}, error) {
	v, err := Coders.Int32.DoDecode(r)
	if err != nil {
		return nil, err
	}
	return MTFlag(v), nil
}

func (O mt_flag) Has(p *Package, f MTFlag) bool {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_FLAG {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(MTFlag); ok {
					return f == rv
				}
			}
			break
		}
	}
	return false
}

func (O mt_flag) Remove(p *Package, f MTFlag) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_FLAG {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(MTFlag); ok {
					if rv == f {
						p.Remove(e)
					}
				}
			}
		}
	}
}

func (O mt_flag) Set(p *Package, v MTFlag) {
	if O.Has(p, v) {
		return
	}
	f := NewFrameV(MT_FLAG, v, O)
	p.PushFront(f)
}
