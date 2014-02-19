package esnp

import (
	"bmautil/byteutil"
	Coders "bmautil/coder"
	"esp/espnet/protpack"
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

func TestMTXData(t *testing.T) {
	p1 := protpack.NewPackage()
	FrameCoders.XData.Add(p1, 1, 1234, nil)
	FrameCoders.XData.Add(p1, 2, "abcdef", nil)
	b, _ := p1.ToBytesBuffer()
	t.Error(b.ToBytes())

	pr := protpack.NewPackageReader()
	pr.Append(b.ToBytes())
	p2, _ := pr.ReadPackage(1024)
	it := FrameCoders.XData.Iterator(p2)
	for ; !it.IsEnd(); it.Next() {
		v, _ := it.Value(nil)
		t.Error(it.Xid(), v)
	}
}
