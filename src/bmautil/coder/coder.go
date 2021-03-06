package coder

import (
	"bmautil/byteutil"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

// lenBytes
type LenBytesCoder int

func (this LenBytesCoder) DoEncode(w *byteutil.BytesBufferWriter, bs []byte) {
	Int32.DoEncode(w, int32(len(bs)))
	if bs != nil {
		w.Write(bs)
	}
}

func (this LenBytesCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	this.DoEncode(w, v.([]byte))
	return nil
}

func (this LenBytesCoder) DoDecode(r *byteutil.BytesBufferReader, maxlen int) ([]byte, error) {
	l, err := Int.DoDecode(r)
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

func (this LenBytesCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	s, err := this.DoDecode(r, int(this))
	if err != nil {
		return nil, err
	}
	return s, nil
}

// lenString
type LenStringCoder int

func (this LenStringCoder) DoEncode(w *byteutil.BytesBufferWriter, v string) {
	bs := []byte(v)
	Int32.DoEncode(w, int32(len(bs)))
	w.Write(bs)
}

func (this LenStringCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	this.DoEncode(w, v.(string))
	return nil
}

func (this LenStringCoder) DoDecode(r *byteutil.BytesBufferReader, maxlen int) (string, error) {
	l, err := Int.DoDecode(r)
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

func (this LenStringCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	s, err := this.DoDecode(r, int(this))
	if err != nil {
		return nil, err
	}
	return s, nil
}

// string
type stringCoder int

func (this stringCoder) DoEncode(w *byteutil.BytesBufferWriter, v string) {
	w.Write([]byte(v))
}

func (this stringCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	this.DoEncode(w, v.(string))
	return nil
}

func (this stringCoder) DoDecode(r *byteutil.BytesBufferReader) string {
	return string(r.ReadAll())
}

func (this stringCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	s := this.DoDecode(r)
	return s, nil
}

// bool
type boolCoder bool

func (this boolCoder) DoEncode(w *byteutil.BytesBufferWriter, v bool) {
	b := byte(0)
	if v {
		b = 1
	}
	w.WriteByte(b)
}

func (this boolCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	this.DoEncode(w, v.(bool))
	return nil
}

func (this boolCoder) DoDecode(r *byteutil.BytesBufferReader) (bool, error) {
	b, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	return b != 0, err
}

func (this boolCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
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

func (O intCoder) DoEncode(w *byteutil.BytesBufferWriter, v int) {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(int32(v)))
	w.Write(b[:l])
}
func (O intCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(int))
	return nil
}
func (O int8Coder) DoEncode(w *byteutil.BytesBufferWriter, v uint8) {
	w.WriteByte(v)
}
func (O int8Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(uint8))
	return nil
}
func (O int16Coder) DoEncode(w *byteutil.BytesBufferWriter, v int16) {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(v))
	w.Write(b[:l])
}
func (O int16Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(int16))
	return nil
}
func (O int32Coder) DoEncode(w *byteutil.BytesBufferWriter, v int32) {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(v))
	w.Write(b[:l])
}
func (O int32Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(int32))
	return nil
}
func (O int64Coder) DoEncode(w *byteutil.BytesBufferWriter, v int64) {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutVarint(b, int64(v))
	w.Write(b[:l])
}
func (O int64Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(int64))
	return nil
}
func (O uintCoder) DoEncode(w *byteutil.BytesBufferWriter, v uint) {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	w.Write(b[:l])
}
func (O uintCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(uint))
	return nil
}
func (O uint8Coder) DoEncode(w *byteutil.BytesBufferWriter, v uint8) {
	w.WriteByte(v)
}
func (O uint8Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(uint8))
	return nil
}
func (O uint16Coder) DoEncode(w *byteutil.BytesBufferWriter, v uint16) {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	w.Write(b[:l])
}
func (O uint16Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(uint16))
	return nil
}
func (O uint32Coder) DoEncode(w *byteutil.BytesBufferWriter, v uint32) {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	w.Write(b[:l])
}
func (O uint32Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(uint32))
	return nil
}
func (O uint64Coder) DoEncode(w *byteutil.BytesBufferWriter, v uint64) {
	bs := [10]byte{}
	b := bs[:]
	l := binary.PutUvarint(b, uint64(v))
	w.Write(b[:l])
}
func (O uint64Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(uint64))
	return nil
}

func (O intCoder) DoDecode(r io.ByteReader) (int, error) {
	rv, err := binary.ReadVarint(r)
	return int(rv), err
}
func (O intCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O int8Coder) DoDecode(r io.ByteReader) (uint8, error) {
	return r.ReadByte()
}
func (O int8Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O int16Coder) DoDecode(r io.ByteReader) (int16, error) {
	rv, err := binary.ReadVarint(r)
	return int16(rv), err
}
func (O int16Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O int32Coder) DoDecode(r io.ByteReader) (int32, error) {
	rv, err := binary.ReadVarint(r)
	return int32(rv), err
}
func (O int32Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O int64Coder) DoDecode(r io.ByteReader) (int64, error) {
	rv, err := binary.ReadVarint(r)
	return int64(rv), err
}
func (O int64Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uintCoder) DoDecode(r io.ByteReader) (uint, error) {
	rv, err := binary.ReadVarint(r)
	return uint(rv), err
}
func (O uintCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint8Coder) DoDecode(r io.ByteReader) (uint8, error) {
	return r.ReadByte()
}
func (O uint8Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint16Coder) DoDecode(r io.ByteReader) (uint16, error) {
	rv, err := binary.ReadUvarint(r)
	return uint16(rv), err
}
func (O uint16Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint32Coder) DoDecode(r io.ByteReader) (uint32, error) {
	rv, err := binary.ReadUvarint(r)
	return uint32(rv), err
}
func (O uint32Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}
func (O uint64Coder) DoDecode(r io.ByteReader) (uint64, error) {
	rv, err := binary.ReadUvarint(r)
	return uint64(rv), err
}
func (O uint64Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}

// fixIntxCoder
type fixUint16Coder int
type fixUint32Coder int
type fixUint64Coder int

func (O fixUint16Coder) DoEncode(w *byteutil.BytesBufferWriter, v uint16) {
	bs := [2]byte{}
	b := bs[:]
	binary.BigEndian.PutUint16(b, uint16(v))
	w.Write(b)
}
func (O fixUint16Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(uint16))
	return nil
}
func (O fixUint32Coder) DoEncode(w *byteutil.BytesBufferWriter, v uint32) {
	bs := [4]byte{}
	b := bs[:]
	binary.BigEndian.PutUint32(b, uint32(v))
	w.Write(b)
}
func (O fixUint32Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(uint32))
	return nil
}
func (O fixUint64Coder) DoEncode(w *byteutil.BytesBufferWriter, v uint64) {
	bs := [8]byte{}
	b := bs[:]
	binary.BigEndian.PutUint64(b, uint64(v))
	w.Write(b)
}
func (O fixUint64Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
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
func (O fixUint16Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
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
func (O fixUint32Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
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
func (O fixUint64Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}

// Float32 Float64 Coder
type float32Coder int
type float64Coder int

func (O float32Coder) DoEncode(w *byteutil.BytesBufferWriter, v float32) {
	iv := math.Float32bits(v)
	FixUint32.DoEncode(w, iv)
}
func (O float32Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(float32))
	return nil
}
func (O float64Coder) DoEncode(w *byteutil.BytesBufferWriter, v float64) {
	iv := math.Float64bits(v)
	FixUint64.DoEncode(w, iv)
}
func (O float64Coder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	O.DoEncode(w, v.(float64))
	return nil
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
func (O float32Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
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
func (O float64Coder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return O.DoDecode(r)
}

// varCoder
type varCoder int

func (this varCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if v == nil {
		w.WriteByte(byte(reflect.Invalid))
		return nil
	}
	var b [binary.MaxVarintLen64]byte
	bs := b[:]

	if rb, ok := v.([]byte); ok {
		w.WriteByte(byte(reflect.Array))
		LenBytes.DoEncode(w, rb)
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
	w.WriteByte(tvk)
	switch tv.Kind() {
	case reflect.Bool:
		rv := tv.Bool()
		Bool.DoEncode(w, rv)
		return nil
	case reflect.Int8, reflect.Uint8:
		rv := byte(tv.Uint() & 0xFF)
		w.WriteByte(rv)
		return nil
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		rv := tv.Int()
		l := binary.PutVarint(bs, rv)
		w.Write(b[:l])
		return nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv := tv.Uint()
		l := binary.PutUvarint(bs, rv)
		w.Write(b[:l])
		return nil
	case reflect.Float32:
		Float32.DoEncode(w, float32(tv.Float()))
		return nil
	case reflect.Float64:
		Float64.DoEncode(w, tv.Float())
		return nil
	case reflect.String:
		LenString.DoEncode(w, tv.String())
		return nil
	case reflect.Map:
		if tv.Type().Key().Kind() != reflect.String {
			return errors.New("onlye encode map[string]value")
		}
		sz := tv.Len()
		l := binary.PutUvarint(bs, uint64(sz))
		w.Write(bs[:l])
		mkeys := tv.MapKeys()
		for _, k := range mkeys {
			sval := tv.MapIndex(k)
			LenString.DoEncode(w, k.String())
			err := this.Encode(w, sval.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Ptr:
		return this.Encode(w, tv.Interface())
	case reflect.Slice:
		sz := tv.Len()
		l := binary.PutUvarint(bs, uint64(sz))
		w.Write(bs[:l])
		for i := 0; i < sz; i++ {
			sval := tv.Index(i)
			err := this.Encode(w, sval.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		vt := tv.Type()
		sz := vt.NumField()
		l := binary.PutUvarint(bs, uint64(sz))
		w.Write(bs[:l])
		for i := 0; i < sz; i++ {
			tfield := vt.Field(i)
			sval := tv.Field(i)
			LenString.DoEncode(w, tfield.Name)
			err := this.Encode(w, sval.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New(fmt.Sprintf("unknow type %T", v))
	}
}

func (this varCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
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
		return LenBytes.DoDecode(r, 0)
	case reflect.Bool:
		return Bool.DoDecode(r)
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
		return Float32.Decode(r)
	case reflect.Float64:
		return Float64.Decode(r)
	case reflect.String:
		return LenString.Decode(r)
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
			kv, err2 := LenString.DoDecode(r, 0)
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

var (
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
)
