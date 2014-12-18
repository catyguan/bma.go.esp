package esnp

import (
	"bmautil/valutil"
	"errors"
	"fmt"
)

var (
	notValueErr error = errors.New("not correct value")
)

type struct_key_value struct {
	name    string
	value   interface{}
	remain  []byte
	encoder Encoder
}

func (this *struct_key_value) String() string {
	return fmt.Sprintf("%s=%v", this.name, this.value)
}

func (this *struct_key_value) Value(dec Decoder) (interface{}, error) {
	if this.value != nil {
		return this.value, nil
	}
	if this.remain != nil {
		var err error
		if dec == nil {
			dec = Coders.Varinat
		}
		var bdr BytesDecodeReader
		bdr.data = this.remain
		this.value, err = dec.Decode(&bdr)
		if err != nil {
			return nil, err
		}
		this.remain = nil
	}
	return this.value, nil
}

type mlt_key_values byte

func (O mlt_key_values) Encode(w EncodeWriter, v interface{}) error {
	if mv, ok := v.(*struct_key_value); ok {
		err0 := Coders.LenString.DoEncode(w, mv.name)
		if err0 != nil {
			return err0
		}
		c := mv.encoder
		if c == nil {
			c = Coders.Varinat
		}
		err := c.Encode(w, mv.value)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("not struct_key_value")
}

func (O mlt_key_values) Decode(r DecodeReader) (interface{}, error) {
	s1, err := Coders.LenString.DoDecode(r, 0)
	if err != nil {
		return nil, err
	}
	sz := r.Remain()
	if sz == -1 {
		return nil, fmt.Errorf("unknow stream form xdata")
	}
	b := make([]byte, sz)
	_, err = r.Read(b)
	if err != nil {
		return nil, err
	}
	return &struct_key_value{s1, nil, b, nil}, nil
}

func (O mlt_key_values) MT() byte {
	return byte(O)
}

func (O mlt_key_values) is(e *MessageLine, key string) (*struct_key_value, error) {
	if e.MessageType() == O.MT() {
		v, err := e.Value(O)
		if err != nil {
			return nil, err
		}
		if mv, ok := v.(*struct_key_value); ok {
			if mv.name == key {
				return mv, nil
			}
		} else {
			return nil, notValueErr
		}
	}
	return nil, nil
}

func (O mlt_key_values) Set(p *Message, key string, value interface{}, enc Encoder) {
	p.RemoveMessageLine(func(e *MessageLine) (bool, bool) {
		if mv, _ := O.is(e, key); mv != nil {
			return true, false
		}
		return false, false
	})
	f := NewMessageLineV(O.MT(), &struct_key_value{key, value, nil, enc}, O)
	p.PushBack(f)
}

func (O mlt_key_values) Get(p *Message, key string, dec Decoder) (interface{}, error) {
	for e := p.Front(); e != nil; e = e.Next() {
		mv, err := O.is(e, key)
		if err != nil {
			continue
		}
		if mv != nil {
			return mv.Value(dec)
		}
	}
	return nil, nil
}

func (O mlt_key_values) Pop(p *Message, key string, dec Decoder) (interface{}, error) {
	for e := p.Front(); e != nil; e = e.Next() {
		mv, err := O.is(e, key)
		if err != nil {
			continue
		}
		if mv != nil {
			p.Remove(e)
			return mv.Value(dec)
		}
	}
	return nil, nil
}

func (O mlt_key_values) Remove(p *Message, key string) {
	p.RemoveMessageLine(func(e *MessageLine) (bool, bool) {
		if mv, _ := O.is(e, key); mv != nil {
			return true, false
		}
		return false, false
	})
}

func (O mlt_key_values) List(p *Message) []string {
	r := make([]string, 0)
	mt := O.MT()
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == mt {
			v, err := e.Value(O)
			if err != nil {
				continue
			}
			if mv, ok := v.(*struct_key_value); ok {
				r = append(r, mv.name)
			}
		}
	}
	return r
}

func (O mlt_key_values) Map(p *Message) (map[string]interface{}, error) {
	r := make(map[string]interface{}, 0)
	mt := O.MT()
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == mt {
			v, err := e.Value(O)
			if err != nil {
				continue
			}
			if mv, ok := v.(*struct_key_value); ok {
				var err error
				r[mv.name], err = mv.Value(nil)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return r, nil
}

type MessageValues struct {
	m     *Message
	coder mlt_key_values
}

func (this *MessageValues) Set(key string, value interface{}) {
	this.coder.Set(this.m, key, value, nil)
}

func (this *MessageValues) Get(key string) (interface{}, error) {
	return this.coder.Get(this.m, key, nil)
}

func (this *MessageValues) GetString(key string, defv string) (string, error) {
	v, err := this.Get(key)
	if err != nil {
		return "", err
	}
	return valutil.ToString(v, defv), nil
}

func (this *MessageValues) GetInt(key string, defv int64) (int64, error) {
	v, err := this.Get(key)
	if err != nil {
		return defv, err
	}
	return valutil.ToInt64(v, defv), nil
}

func (this *MessageValues) GetUint(key string, defv uint64) (uint64, error) {
	v, err := this.Get(key)
	if err != nil {
		return defv, err
	}
	return valutil.ToUint64(v, defv), nil
}

func (this *MessageValues) GetBool(key string) (bool, error) {
	v, err := this.Get(key)
	if err != nil {
		return false, err
	}
	r, ok := valutil.ToBoolNil(v)
	if ok {
		return r, nil
	}
	return false, errors.New("not bool")
}

func (this *MessageValues) Del(key string) {
	this.coder.Remove(this.m, key)
}

func (this *MessageValues) List() []string {
	return this.coder.List(this.m)
}

func (this *MessageValues) CopyFrom(m map[string]interface{}) {
	for k, v := range m {
		this.Set(k, v)
	}
}

func (this *MessageValues) ToMap() (map[string]interface{}, error) {
	return this.coder.Map(this.m)
}

func (this *MessageValues) ToBean(beanPtr interface{}) (bool, error) {
	m, err := this.ToMap()
	if err != nil {
		return false, err
	}
	return valutil.ToBean(m, beanPtr), nil
}
