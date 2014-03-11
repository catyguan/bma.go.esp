package esnp

import (
	"bmautil/byteutil"
	"errors"
	"io"
)

type mt_payload int

func (O mt_payload) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if r, ok := v.([]byte); ok {
		w.Append(r)
		return nil
	}
	if r, ok := v.(*byteutil.BytesBuffer); ok {
		for _, p := range r.DataList {
			w.Append(p)
		}
		return nil
	}
	if r, ok := v.(io.Reader); ok {
		_, err := io.Copy(w, r)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("not []byte")
}

func (O mt_payload) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return r.ReadAll(), nil
}

func (O mt_payload) Get(p *Package) []byte {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_PAYLOAD {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.([]byte); ok {
					return rv
				}
			}
			break
		}
	}
	return nil
}

func (O mt_payload) Remove(p *Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_PAYLOAD {
			p.Remove(e)
			break
		}
	}
}

func (O mt_payload) Set(p *Package, val interface{}) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_PAYLOAD {
			p.Remove(e)
			break
		}
	}
	f := NewFrameV(MT_PAYLOAD, val, O)
	p.PushFront(f)
}
