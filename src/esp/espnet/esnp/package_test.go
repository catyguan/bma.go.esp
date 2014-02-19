package esnp

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestPackageBase(t *testing.T) {

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, 1000)
	fmt.Println(bs)

	p := NewPackage()
	p.PushBack(NewFrameB(0x01, []byte{1, 2, 3, 4, 5, 6, 7}))
	// p.SetId(2013)

	fmt.Println(p.String())
	b, _ := p.ToBytesBuffer()
	fmt.Println(b.ToBytes())

	h := new(FHeader)
	h.Read(b.ToBytes(), 0)
	fmt.Println(h)

}

func TestPackageReader(t *testing.T) {

	// p := NewPackage()
	// p.PushBack(NewFrameB(0x01020304, []byte{1, 2, 3, 4, 5, 6, 7}))
	// p.SetId(uint64(time.Now().UnixNano()))

	// b := p.ToBytesBuffer().ToBytes()
	// t.Error(p.ToBytesBuffer().TraceString(64))
	// pr := NewPackageReader()
	// pr.Append(b[:len(b)-1])
	// pr.Append(b[7:])
	// t.Error(pr.buffer.TraceString(128))
	// if true {
	// 	pout, err := pr.ReadPackage(1024)
	// 	t.Error(pout, err)
	// 	t.Error(pr.buffer.TraceString(128))
	// }
	// if true {
	// 	pout, err := pr.ReadPackage(1024)
	// 	t.Error(pout, err)
	// 	t.Error(pr.buffer.TraceString(128))
	// }

}
