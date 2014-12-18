package esnp

import (
	"bmautil/byteutil"
	"errors"
	"fmt"
	"io"
)

type mlt_payload int

func (O mlt_payload) Encode(w EncodeWriter, v interface{}) error {
	if r, ok := v.([]byte); ok {
		_, err := w.Write(r)
		return err
	}
	if r, ok := v.(*byteutil.BytesBuffer); ok {
		for _, p := range r.DataList {
			_, err := w.Write(p)
			if err != nil {
				return err
			}
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

func (O mlt_payload) Decode(r DecodeReader) (interface{}, error) {
	sz := r.Remain()
	if sz == -1 {
		return nil, fmt.Errorf("unknow stream form xdata")
	}
	b := make([]byte, sz)
	_, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (O mlt_payload) Get(p *Message) []byte {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_PAYLOAD {
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

func (O mlt_payload) Remove(p *Message) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_PAYLOAD {
			p.Remove(e)
			break
		}
	}
}

func (O mlt_payload) Set(p *Message, val interface{}) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_PAYLOAD {
			p.Remove(e)
			break
		}
	}
	f := NewMessageLineV(MLT_PAYLOAD, val, O)
	p.PushFront(f)
}
