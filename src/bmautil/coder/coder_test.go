package coder

import (
	"bmautil/byteutil"
	"testing"
)

func TestCode1(t *testing.T) {
	buf := byteutil.NewBytesBuffer()
	w := buf.NewWriter()
	Uint32.Encode(w, uint32(3))
	w.End()
	t.Error(buf.TraceString(123))

	r := buf.NewReader()
	v, _ := Uint32.Decode(r)
	t.Error(v)
}
