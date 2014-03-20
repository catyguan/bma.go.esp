package esnp

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

type decoderFunc func(r *byteutil.BytesBufferReader) (interface{}, error)

func (this decoderFunc) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return this(r)
}

func NewDecoderFuc(f func(r *byteutil.BytesBufferReader) (interface{}, error)) Decoder {
	return decoderFunc(f)
}

// Message
type MessageListener func(msg *Message) error
type MessageSender func(msg *Message) error
type ResponseListener func(msg *Message, err error) error
type MessageHandler func(msg *Message) (bool, error)
