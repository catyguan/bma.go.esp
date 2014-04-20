package esnp

import "errors"

type mt_error int

func (O mt_error) Encode(w EncodeWriter, v interface{}) error {
	if r, ok := v.(string); ok {
		_, err := w.Write([]byte(r))
		return err
	}
	return errors.New("not string")
}

func (O mt_error) Decode(r DecodeReader) (interface{}, error) {
	b := r.ReadAll()
	return string(b), nil
}

func (O mt_error) Get(p *Package) (bool, string) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_ERROR {
			v, err := e.Value(O)
			if err == nil {
				return true, v.(string)
			}
			break
		}
	}
	return false, ""
}

func (O mt_error) Remove(p *Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_ERROR {
			p.Remove(e)
			break
		}
	}
}

func (O mt_error) Set(p *Package, val string) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_ERROR {
			p.Remove(e)
			break
		}
	}
	f := NewFrameV(MT_ERROR, val, O)
	p.PushFront(f)
}
