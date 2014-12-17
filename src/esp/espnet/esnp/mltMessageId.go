package esnp

import "errors"

type mlt_message_id byte

func (O mlt_message_id) Encode(w EncodeWriter, v interface{}) error {
	if o, ok := v.(uint64); ok {
		return Coders.FixUint64.DoEncode(w, o)
	}
	return errors.New("not messageId")
}

func (O mlt_message_id) Decode(r DecodeReader) (interface{}, error) {
	return Coders.FixUint64.DoDecode(r)
}

func (O mlt_message_id) Get(p *Message) uint64 {
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

func (O mlt_message_id) Remove(p *Message) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			p.Remove(e)
			break
		}
	}
}

func (O mlt_message_id) Set(p *Message, val uint64) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			p.Remove(e)
			break
		}
	}
	f := NewMessageLineV(byte(O), val, O)
	p.PushFront(f)
}

func (O mlt_message_id) Sure(p *Message) uint64 {
	mid := O.Get(p)
	if mid == 0 {
		mid = NextMessageId()
	}
	O.Set(p, mid)
	return mid
}
