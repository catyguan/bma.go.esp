package esnp

import "errors"

type mlt_error int

func (O mlt_error) Encode(w EncodeWriter, v interface{}) error {
	if r, ok := v.(string); ok {
		_, err := w.Write([]byte(r))
		return err
	}
	return errors.New("not string")
}

func (O mlt_error) Decode(r DecodeReader) (interface{}, error) {
	return Coders.String.Decode(r)
}

func (O mlt_error) Get(p *Message) (bool, string) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_ERROR {
			v, err := e.Value(O)
			if err == nil {
				return true, v.(string)
			}
			break
		}
	}
	return false, ""
}

func (O mlt_error) Remove(p *Message) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_ERROR {
			p.Remove(e)
			break
		}
	}
}

func (O mlt_error) Set(p *Message, val string) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_ERROR {
			p.Remove(e)
			break
		}
	}
	f := NewMessageLineV(MLT_ERROR, val, O)
	p.PushFront(f)
}
