package protpack

import "bmautil/byteutil"

type Encoder interface {
	Encode(w *byteutil.BytesBufferWriter, v interface{}) error
}

type Decoder interface {
	Decode(r *byteutil.BytesBufferReader) (interface{}, error)
}

type encoderFunc func(w *byteutil.BytesBufferWriter, v interface{}) error

func (this encoderFunc) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	return this(w, v)
}

func NewEncoderFuc(f func(w *byteutil.BytesBufferWriter, v interface{}) error) Encoder {
	return encoderFunc(f)
}
