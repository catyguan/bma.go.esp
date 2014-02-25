package esnp

import (
	"bmautil/byteutil"
	Coders "bmautil/coder"
	"errors"
)

type mt_message_id byte

func (O mt_message_id) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if o, ok := v.(uint64); ok {
		Coders.FixUint64.DoEncode(w, o)
		return nil
	}
	return errors.New("not messageId")
}

func (O mt_message_id) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return Coders.FixUint64.DoDecode(r)
}

func (O mt_message_id) Get(p *Package) uint64 {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(uint64); ok {
					return rv
				}
			}
			break
		}
	}
	return 0
}

func (O mt_message_id) Remove(p *Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			p.Remove(e)
			break
		}
	}
}

func (O mt_message_id) Set(p *Package, val uint64) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			p.Remove(e)
			break
		}
	}
	f := NewFrameV(byte(O), val, O)
	p.PushFront(f)
}

func (O mt_message_id) Sure(p *Package) uint64 {
	mid := O.Get(p)
	if mid == 0 {
		mid = NextMessageId()
	}
	O.Set(p, mid)
	return mid
}