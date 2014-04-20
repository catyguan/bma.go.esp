package esnp

import "testing"

func TestCode1(t *testing.T) {
	var w BytesEncodeWriter
	Coders.Uint32.Encode(&w, uint32(3))
	t.Error(w.ToBytes())

	var r BytesDecodeReader
	r.data = w.ToBytes()
	v, _ := Coders.Uint32.Decode(&r)
	t.Error(v)
}
