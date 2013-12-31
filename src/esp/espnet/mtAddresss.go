package espnet

import (
	"bmautil/byteutil"
	"errors"
	"esp/espnet/protpack"
)

// Address
type addrCoder byte

func (O addrCoder) DoEncode(w *byteutil.BytesBufferWriter, a Address) {
	Coders.Int.DoEncode(w, a.Size())
	for _, v := range a {
		Coders.Int.DoEncode(w, len(v))
		w.WriteString(v)
	}
}

func (O addrCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if mv, ok := v.(Address); ok {
		O.DoEncode(w, mv)
		return nil
	}
	return errors.New("not address")
}

func (O addrCoder) DoDecode(r *byteutil.BytesBufferReader) (Address, error) {
	sz, err1 := Coders.Int.DoDecode(r)
	if err1 != nil {
		return nil, err1
	}

	a := make(Address, sz)
	for i := 0; i < sz; i++ {
		l, err2 := Coders.Int.DoDecode(r)
		if err2 != nil {
			return nil, err2
		}
		bs := make([]byte, l)
		_, err3 := r.Read(bs)
		if err3 != nil {
			return nil, err3
		}
		a[i] = string(bs)
	}
	return a, nil
}

func (O addrCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v, err := O.DoDecode(r)
	return v, err
}

func (O addrCoder) Get(p *protpack.Package) Address {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(Address); ok {
					return rv
				}
			}
			break
		}
	}
	return nil
}

func (O addrCoder) Remove(p *protpack.Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			p.Remove(e)
			break
		}
	}
}

func (O addrCoder) Set(p *protpack.Package, val Address) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			p.Remove(e)
			break
		}
	}
	f := protpack.NewFrameV(byte(O), val, O)
	p.PushFront(f)
}
