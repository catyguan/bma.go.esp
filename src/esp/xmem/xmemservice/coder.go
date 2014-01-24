package xmemservice

import (
	"bmautil/byteutil"
	"bmautil/coder"
)

type SimpleCoder int

func (O SimpleCoder) Encode(val interface{}) (string, []byte, error) {
	buf := byteutil.NewBytesBuffer()
	w := buf.NewWriter()
	err := coder.Varinat.Encode(w, val)
	if err != nil {
		return "", nil, err
	}
	bs := w.End().ToBytes()
	return "", bs, nil
}

func (O SimpleCoder) Decode(flag string, data []byte) (interface{}, int, error) {
	buf := byteutil.NewBytesBufferB(data)
	r := buf.NewReader()
	val, err := coder.Varinat.Decode(r)
	if err != nil {
		return nil, 0, err
	}
	return val, len(data), nil
}
