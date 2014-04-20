package esnp

import (
	"os"
	"testing"
	"time"
)

func TestMessageId(t *testing.T) {
	p1 := NewPackage()
	FrameCoders.MessageId.Set(p1, 88888888)
	b, _ := p1.ToBytes()
	t.Error(b)
}

func TestMTXData(t *testing.T) {

	time.AfterFunc(5*time.Second, func() {
		os.Exit(-1)
	})

	p1 := NewPackage()
	FrameCoders.XData.Add(p1, 1, 1234, nil)
	FrameCoders.XData.Add(p1, 2, "abcdef", nil)
	b, _ := p1.ToBytes()
	t.Error(b)

	pr := NewPackageReader()
	pr.Append(b)
	pr.Append(b)
	pr.Append([]byte{1, 2, 3})
	for {
		t.Error(pr.buffer[:pr.wpos], pr.rpos, pr.wpos)
		p2, _ := pr.ReadPackage(1024)
		if p2 != nil {
			it := FrameCoders.XData.Iterator(p2)
			for ; !it.IsEnd(); it.Next() {
				v, _ := it.Value(nil)
				t.Error(it.Xid(), v)
			}
		} else {
			break
		}
	}
}
