package esnp

import (
	"bmautil/byteutil"
	"errors"
)

type MTVersion struct {
	Major  uint8
	Minor  uint8
	Branch uint8
	Mutate uint8
}

type mt_version int

func (O mt_version) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if o, ok := v.(*MTVersion); ok {
		w.WriteByte(byte(o.Major))
		w.WriteByte(byte(o.Minor))
		w.WriteByte(byte(o.Branch))
		w.WriteByte(byte(o.Mutate))
		return nil
	}
	return errors.New("not mtVersion")
}

func (O mt_version) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v1, err1 := r.ReadByte()
	if err1 != nil {
		return nil, err1
	}
	v2, err2 := r.ReadByte()
	if err2 != nil {
		return nil, err2
	}
	v3, err3 := r.ReadByte()
	if err3 != nil {
		return nil, err3
	}
	v4, err4 := r.ReadByte()
	if err4 != nil {
		return nil, err4
	}
	o := new(MTVersion)
	o.Major = v1
	o.Minor = v2
	o.Branch = v3
	o.Mutate = v4
	return o, nil
}

func (O mt_version) Get(p *Package) *MTVersion {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_VERSION {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(*MTVersion); ok {
					return rv
				}
			}
			break
		}
	}
	return nil
}

func (O mt_version) Remove(p *Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_VERSION {
			p.Remove(e)
			break
		}
	}
}

func (O mt_version) Set(p *Package, val *MTVersion) {
	O.Remove(p)
	f := NewFrameV(MT_VERSION, val, O)
	p.PushFront(f)
}
