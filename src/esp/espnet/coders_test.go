package espnet

import (
	"bmautil/byteutil"
	"testing"
)

func TestCode1(t *testing.T) {
	buf := byteutil.NewBytesBuffer()
	w := buf.NewWriter()
	Coders.Uint32.Encode(w, uint32(3))
	w.End()
	t.Error(buf.TraceString(123))

	r := buf.NewReader()
	v, _ := Coders.Uint32.Decode(r)
	t.Error(v)
}
