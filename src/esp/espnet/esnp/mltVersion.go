package esnp

import (
	"errors"
	"fmt"
)

type Version struct {
	Major  uint8
	Minor  uint8
	Branch uint8
	Mutate uint8
}

func (this *Version) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", this.Major, this.Minor, this.Branch, this.Mutate)
}

type mlt_version int

func (O mlt_version) Encode(w EncodeWriter, v interface{}) error {
	if o, ok := v.(*Version); ok {
		err := w.WriteByte(byte(o.Major))
		if err != nil {
			return err
		}
		err = w.WriteByte(byte(o.Minor))
		if err != nil {
			return err
		}
		err = w.WriteByte(byte(o.Branch))
		if err != nil {
			return err
		}
		err = w.WriteByte(byte(o.Mutate))
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("not mtVersion")
}

func (O mlt_version) Decode(r DecodeReader) (interface{}, error) {
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
	o := new(Version)
	o.Major = v1
	o.Minor = v2
	o.Branch = v3
	o.Mutate = v4
	return o, nil
}

func (O mlt_version) Get(p *Message) *Version {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_VERSION {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(*Version); ok {
					return rv
				}
			}
			break
		}
	}
	return nil
}

func (O mlt_version) Remove(p *Message) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_VERSION {
			p.Remove(e)
			break
		}
	}
}

func (O mlt_version) Set(p *Message, val *Version) {
	O.Remove(p)
	f := NewMessageLineV(MLT_VERSION, val, O)
	p.PushFront(f)
}
