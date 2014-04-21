package esnp

import "io"

type EncodeWriter interface {
	io.Writer
	io.ByteWriter

	WriteFrame(mt byte, data []byte) error
	NewFrame() (int, error)
	EndFrame(p int, mt byte) error
}

type DecodeReader interface {
	io.ByteReader
	io.Reader

	ReadAll() []byte
	Remain() []byte
}

type Encoder interface {
	Encode(w EncodeWriter, v interface{}) error
}

type Decoder interface {
	Decode(r DecodeReader) (interface{}, error)
}

type encoderFunc func(w EncodeWriter, v interface{}) error

func (this encoderFunc) Encode(w EncodeWriter, v interface{}) error {
	return this(w, v)
}

func NewEncoderFuc(f func(w EncodeWriter, v interface{}) error) Encoder {
	return encoderFunc(f)
}

type decoderFunc func(r DecodeReader) (interface{}, error)

func (this decoderFunc) Decode(r DecodeReader) (interface{}, error) {
	return this(r)
}

func NewDecoderFuc(f func(r DecodeReader) (interface{}, error)) Decoder {
	return decoderFunc(f)
}

// Message
type MessageListener func(msg *Message) error
type MessageSender func(msg *Message) error
type ResponseListener func(msg *Message, err error) error
type MessageHandler func(msg *Message) (bool, error)
