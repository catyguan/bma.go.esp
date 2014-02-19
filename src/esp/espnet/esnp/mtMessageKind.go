package esnp

import (
	"bmautil/byteutil"
	"errors"
)

type mt_message_kind int

func (O mt_message_kind) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if o, ok := v.(MessageKind); ok {
		w.WriteByte(byte(o))
		return nil
	}
	return errors.New("not messageType")
}

func (O mt_message_kind) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	return MessageKind(v), nil
}

func (O mt_message_kind) Get(p *Package) MessageKind {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_MESSAGE_KIND {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(MessageKind); ok {
					return rv
				}
			}
			break
		}
	}
	return MK_UNKNOW
}

func (O mt_message_kind) Remove(p *Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_MESSAGE_KIND {
			p.Remove(e)
			break
		}
	}
}

func (O mt_message_kind) Set(p *Package, val MessageKind) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_MESSAGE_KIND {
			p.Remove(e)
			break
		}
	}
	f := NewFrameV(MT_MESSAGE_KIND, val, O)
	p.PushFront(f)
}
