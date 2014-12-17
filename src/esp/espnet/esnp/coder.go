package esnp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

// lenBytes
type LenBytesCoder int

func (this LenBytesCoder) DoEncode(w EncodeWriter, bs []byte) error {
	l := len(bs)
	Coders.Int32.DoEncode(w, int32(l))
	if l > 0 {
		_, err := w.Write(bs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this LenBytesCoder) Encode(w EncodeWriter, v interface{}) error {
	return this.DoEncode(w, v.([]byte))
}

func (this LenBytesCoder) DoDecode(r DecodeReader, maxlen int) ([]byte, error) {
	l, err := Coders.Int.DoDecode(r)
	if err != nil {
		return nil, err
	}
	if maxlen > 0 && l > maxlen {
		return nil, fmt.Errorf("too large bytes block - %d/%d", l, maxlen)
	}
	p := make([]byte, l)
	if l > 0 {
		_, err = r.Read(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (this LenBytesCoder) Decode(r DecodeReader) (interface{}, error) {
	s, err := this.DoDecode(r, int(this))
	if err != nil {
		return nil, err
	}
	return s, nil
}

// lenString
type LenStringCoder int

func (this LenStringCoder) DoEncode(w EncodeWriter, v string) error {
	bs := []byte(v)
	err := Coders.Int32.DoEncode(w, int32(len(bs)))
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (this LenStringCoder) Encode(w EncodeWriter, v interface{}) error {
	return this.DoEncode(w, v.(string))
}

func (this LenStringCoder) DoDecode(r DecodeReader, maxlen int) (string, error) {
	l, err := Coders.Int.DoDecode(r)
	if err != nil {
		return "", err
	}
	if maxlen > 0 && l > maxlen {
		return "", fmt.Errorf("too large string block - %d/%d", l, maxlen)
	}
	p := make([]byte, l)
	_, err = r.Read(p)
	if err != nil {
		return "", err
	}
	return string(p), nil
}

func (this LenStringCoder) Decode(r DecodeReader) (interface{}, error) {
	s, err := this.DoDecode(r, int(this))
	if err != nil {
		return nil, err
	}
	return s, nil
}

// string
type stringCoder int

func (this stringCoder) DoEncode(w EncodeWriter, v string) error {
	_, err := w.Write([]byte(v))
	return err
}

func (this stringCoder) Encode(w EncodeWriter, v interface{}) error {
	return this.DoEncode(w, v.(string))
}

func (this stringCoder) DoDecode(r DecodeReader) string {
	return string(r.ReadAll())
}

func (this stringCoder) Decode(r DecodeReader) (interface{}, error) {
	s := this.DoDecode(r)
	return s, nil
}

// bool
type boolCoder bool

func (this boolCoder) DoEncode(w EncodeWriter, v bool) error {
	b := byte(0)
	if v {
		b = 1
	}
	return w.WriteByte(b)
}

func (this boolCoder) Encode(w EncodeWriter, v interface{}) error {
	return this.DoEncode(w, v.(bool))
}

func (this boolCoder) DoDecode(r DecodeReader) (bool, error) {
	b, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	return b != 0, err
}

func (this boolCoder) Decode(r DecodeReader) (interface{}, error) {
	v, err := this.DoDecode(r)
	return v, err
}

// intx
type intCoder int
type int8Coder int
type int16Coder int
type int32Coder int
type int64Coder int
type uintCoder int
type uint8Coder int
type uint16Coder int
type uint32Coder int
type uint64Coder int

func (O intCoder) DoEncode(w EncodeWriter, v int) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(int32(v)))
	_, err := w.Write(b[:l])
	return err
}
func (O intCoder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(int))
}
func (O int8Coder) DoEncode(w EncodeWriter, v int8) error {
	return w.WriteByte(uint8(v))
}
func (O int8Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(int8))
}
func (O int16Coder) DoEncode(w EncodeWriter, v int16) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O int16Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(int16))
}
func (O int32Coder) DoEncode(w EncodeWriter, v int32) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O int32Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(int32))
}
func (O int64Coder) DoEncode(w EncodeWriter, v int64) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O int64Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(int64))
}
func (O uintCoder) DoEncode(w EncodeWriter, v uint) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O uintCoder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint))
}
func (O uint8Coder) DoEncode(w EncodeWriter, v uint8) error {
	return w.WriteByte(v)
}
func (O uint8Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint8))
}
func (O uint16Coder) DoEncode(w EncodeWriter, v uint16) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O uint16Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint16))
}
func (O uint32Coder) DoEncode(w EncodeWriter, v uint32) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O uint32Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint32))
}
func (O uint64Coder) DoEncode(w EncodeWriter, v uint64) error {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	_, err := w.Write(b[:l])
	return err
}
func (O uint64Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint64))
}

func (O intCoder) DoDecode(r io.ByteReader) (int, error) {
	rv, err := binary.ReadVarint(r)
	return int(rv), err
}
func (O intCoder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O int8Coder) DoDecode(r io.ByteReader) (int8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return int8(b), nil
}
func (O int8Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O int16Coder) DoDecode(r io.ByteReader) (int16, error) {
	rv, err := binary.ReadVarint(r)
	return int16(rv), err
}
func (O int16Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O int32Coder) DoDecode(r io.ByteReader) (int32, error) {
	rv, err := binary.ReadVarint(r)
	return int32(rv), err
}
func (O int32Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O int64Coder) DoDecode(r io.ByteReader) (int64, error) {
	rv, err := binary.ReadVarint(r)
	return int64(rv), err
}
func (O int64Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uintCoder) DoDecode(r io.ByteReader) (uint, error) {
	rv, err := binary.ReadVarint(r)
	return uint(rv), err
}
func (O uintCoder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint8Coder) DoDecode(r io.ByteReader) (uint8, error) {
	return r.ReadByte()
}
func (O uint8Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint16Coder) DoDecode(r io.ByteReader) (uint16, error) {
	rv, err := binary.ReadUvarint(r)
	return uint16(rv), err
}
func (O uint16Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint32Coder) DoDecode(r io.ByteReader) (uint32, error) {
	rv, err := binary.ReadUvarint(r)
	return uint32(rv), err
}
func (O uint32Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint64Coder) DoDecode(r io.ByteReader) (uint64, error) {
	rv, err := binary.ReadUvarint(r)
	return uint64(rv), err
}
func (O uint64Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}

// fixIntxCoder
type fixUint16Coder int
type fixUint32Coder int
type fixUint64Coder int

func (O fixUint16Coder) DoEncode(w EncodeWriter, v uint16) error {
	bs := [2]byte{}
	b := bs[:]
	binary.BigEndian.PutUint16(b, uint16(v))
	_, err := w.Write(b)
	return err
}
func (O fixUint16Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(uint16))
}
func (O fixUint32Coder) DoEncode(w EncodeWriter, v uint32) error {
	bs := [4]byte{}
	b := bs[:]
	binary.BigEndian.PutUint32(b, uint32(v))
	_, err := w.Write(b)
	return err
}
func (O fixUint32Coder) Encode(w EncodeWriter, v interface{}) error {
	O.DoEncode(w, v.(uint32))
	return nil
}
func (O fixUint64Coder) DoEncode(w EncodeWriter, v uint64) error {
	bs := [8]byte{}
	b := bs[:]
	binary.BigEndian.PutUint64(b, uint64(v))
	_, err := w.Write(b)
	return err
}
func (O fixUint64Coder) Encode(w EncodeWriter, v interface{}) error {
	O.DoEncode(w, v.(uint64))
	return nil
}

func (O fixUint16Coder) DoDecode(r io.Reader) (uint16, error) {
	bs := [2]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}
func (O fixUint16Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O fixUint32Coder) DoDecode(r io.Reader) (uint32, error) {
	bs := [4]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}
func (O fixUint32Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O fixUint64Coder) DoDecode(r io.Reader) (uint64, error) {
	bs := [8]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}
func (O fixUint64Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}

// Float32 Float64 Coder
type float32Coder int
type float64Coder int

func (O float32Coder) DoEncode(w EncodeWriter, v float32) error {
	iv := math.Float32bits(v)
	return Coders.FixUint32.DoEncode(w, iv)
}
func (O float32Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(float32))
}
func (O float64Coder) DoEncode(w EncodeWriter, v float64) error {
	iv := math.Float64bits(v)
	return Coders.FixUint64.DoEncode(w, iv)
}
func (O float64Coder) Encode(w EncodeWriter, v interface{}) error {
	return O.DoEncode(w, v.(float64))
}

func (O float32Coder) DoDecode(r io.Reader) (float32, error) {
	bs := [4]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	iv := binary.BigEndian.Uint32(b)
	return math.Float32frombits(iv), nil
}
func (O float32Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O float64Coder) DoDecode(r io.Reader) (float64, error) {
	bs := [8]byte{}
	b := bs[:]
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	iv := binary.BigEndian.Uint64(b)
	return math.Float64frombits(iv), nil
}
func (O float64Coder) Decode(r DecodeReader) (interface{}, error) {
	return O.DoDecode(r)
}

// varCoder
type varCoder int

func (this varCoder) Encode(w EncodeWriter, v interface{}) error {
	if v == nil {
		w.WriteByte(byte(reflect.Invalid))
		return nil
	}
	var err error
	var b [binary.MaxVarintLen64]byte
	bs := b[:]

	if rb, ok := v.([]byte); ok {
		err = w.WriteByte(byte(reflect.Array))
		if err != nil {
			return err
		}
		err = Coders.LenBytes.DoEncode(w, rb)
		if err != nil {
			return err
		}
		return nil
	}

	tv := reflect.ValueOf(v)
	tvk := byte(tv.Kind())
	switch tv.Kind() {
	case reflect.Struct:
		tvk = byte(reflect.Map)
	case reflect.Int:
		rv := tv.Int()
		if rv <= 2147483647 && rv >= -2147483648 {
			tvk = byte(reflect.Int32)
		} else {
			tvk = byte(reflect.Int64)
		}
	case reflect.Uint:
		rv := tv.Uint()
		if rv <= 0xFFFFFFFF {
			tvk = byte(reflect.Uint32)
		} else {
			tvk = byte(reflect.Uint64)
		}
	}
	err = w.WriteByte(tvk)
	if err != nil {
		return err
	}
	switch tv.Kind() {
	case reflect.Bool:
		rv := tv.Bool()
		return Coders.Bool.DoEncode(w, rv)
	case reflect.Int8, reflect.Uint8:
		rv := byte(tv.Uint() & 0xFF)
		return w.WriteByte(rv)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		rv := tv.Int()
		l := binary.PutVarint(bs, rv)
		_, err = w.Write(b[:l])
		return err
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv := tv.Uint()
		l := binary.PutUvarint(bs, rv)
		_, err = w.Write(b[:l])
		return err
	case reflect.Float32:
		return Coders.Float32.DoEncode(w, float32(tv.Float()))
	case reflect.Float64:
		return Coders.Float64.DoEncode(w, tv.Float())
	case reflect.String:
		return Coders.LenString.DoEncode(w, tv.String())
	case reflect.Map:
		if tv.Type().Key().Kind() != reflect.String {
			return errors.New("onlye encode map[string]value")
		}
		sz := tv.Len()
		l := binary.PutUvarint(bs, uint64(sz))
		_, err = w.Write(bs[:l])
		if err != nil {
			return err
		}
		mkeys := tv.MapKeys()
		for _, k := range mkeys {
			sval := tv.MapIndex(k)
			err = Coders.LenString.DoEncode(w, k.String())
			if err != nil {
				return err
			}
			err = this.Encode(w, sval.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Ptr:
		return this.Encode(w, tv.Elem())
	case reflect.Slice:
		sz := tv.Len()
		l := binary.PutUvarint(bs, uint64(sz))
		_, err = w.Write(bs[:l])
		if err != nil {
			return err
		}
		for i := 0; i < sz; i++ {
			sval := tv.Index(i)
			err = this.Encode(w, sval.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		vt := tv.Type()
		sz := vt.NumField()
		l := binary.PutUvarint(bs, uint64(sz))
		_, err = w.Write(bs[:l])
		if err != nil {
			return err
		}
		for i := 0; i < sz; i++ {
			tfield := vt.Field(i)
			sval := tv.Field(i)
			err = Coders.LenString.DoEncode(w, tfield.Name)
			if err != nil {
				return err
			}
			err = this.Encode(w, sval.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New(fmt.Sprintf("unknow type %T", v))
	}
}

func (this varCoder) Decode(r DecodeReader) (interface{}, error) {
	var b [binary.MaxVarintLen64]byte
	bs := b[:]
	var err0 error
	bs[0], err0 = r.ReadByte()
	if err0 != nil {
		return nil, err0
	}
	k := reflect.Kind(bs[0])
	switch k {
	case reflect.Invalid:
		return nil, nil
	case reflect.Array:
		return Coders.LenBytes.DoDecode(r, 0)
	case reflect.Bool:
		return Coders.Bool.DoDecode(r)
	case reflect.Int8:
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		return int8(b), nil
	case reflect.Uint8:
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		return uint8(b), nil
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		rv, err := binary.ReadVarint(r)
		if err != nil {
			return nil, err
		}
		switch k {
		case reflect.Int:
			return int(rv), nil
		case reflect.Int16:
			return int16(rv), nil
		case reflect.Int32:
			return int32(rv), nil
		}
		return rv, nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, err
		}
		switch k {
		case reflect.Uint:
			return uint(rv), nil
		case reflect.Uint16:
			return uint16(rv), nil
		case reflect.Uint32:
			return uint32(rv), nil
		}
		return rv, nil
	case reflect.Float32:
		return Coders.Float32.Decode(r)
	case reflect.Float64:
		return Coders.Float64.Decode(r)
	case reflect.String:
		return Coders.LenString.Decode(r)
	case reflect.Map, reflect.Struct:
		l, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, err
		}
		if l == 0 {
			return nil, nil
		}

		rv := make(map[string]interface{})
		for i := 0; i < int(l); i++ {
			kv, err2 := Coders.LenString.DoDecode(r, 0)
			if err2 != nil {
				return nil, err2
			}
			fv, err3 := this.Decode(r)
			if err3 != nil {
				return nil, err3
			}
			rv[kv] = fv
		}
		return rv, nil
	case reflect.Slice:
		l, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, err
		}
		if l == 0 {
			return nil, nil
		}
		rv := make([]interface{}, l)
		for i := 0; i < int(l); i++ {
			fv, err2 := this.Decode(r)
			if err2 != nil {
				return nil, err2
			}
			rv[i] = fv
		}
		return rv, nil
	}
	return nil, nil
}

type NULL int

type allCoder struct {
	LenBytes  LenBytesCoder
	String    stringCoder
	LenString LenStringCoder
	Bool      boolCoder
	Int       intCoder
	Int8      int8Coder
	Int16     int16Coder
	Int32     int32Coder
	Int64     int64Coder
	Uint      uintCoder
	Uint8     uint8Coder
	Uint16    uint16Coder
	Uint32    uint32Coder
	Uint64    uint64Coder
	FixUint16 fixUint16Coder
	FixUint32 fixUint32Coder
	FixUint64 fixUint64Coder
	Float32   float32Coder
	Float64   float64Coder
	Varinat   varCoder
	NullValue NULL
}

var (
	Coders allCoder
)
