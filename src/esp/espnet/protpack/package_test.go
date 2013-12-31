package protpack

import (
	"testing"
	"time"
)

func TestPackageBase(t *testing.T) {

	p := NewPackage()
	p.PushBack(NewFrameB(0x01020304, []byte{1, 2, 3, 4, 5, 6, 7}))
	p.SetId(2013)

	t.Error(p.String())
	b := p.ToBytesBuffer().ToBytes()
	t.Error(b)

	h := new(FHeader)
	h.Read(b, 0)
	t.Error(h)

}

func TestPackageReader(t *testing.T) {

	p := NewPackage()
	p.PushBack(NewFrameB(0x01020304, []byte{1, 2, 3, 4, 5, 6, 7}))
	p.SetId(uint64(time.Now().UnixNano()))

	b := p.ToBytesBuffer().ToBytes()
	t.Error(p.ToBytesBuffer().TraceString(64))
	pr := NewPackageReader()
	pr.Append(b[:len(b)-1])
	pr.Append(b[7:])
	t.Error(pr.buffer.TraceString(128))
	if true {
		pout, err := pr.ReadPackage(1024)
		t.Error(pout, err)
		t.Error(pr.buffer.TraceString(128))
	}
	if true {
		pout, err := pr.ReadPackage(1024)
		t.Error(pout, err)
		t.Error(pr.buffer.TraceString(128))
	}

}
