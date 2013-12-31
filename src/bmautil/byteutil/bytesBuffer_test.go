package byteutil

import (
	"io"
	"testing"
)

func TestByteBufferBase(t *testing.T) {
	buf := NewBytesBuffer(1)
	buf.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9})
	buf.Add([]byte{1, 2, 3, 4})
	t.Error(buf.TraceString(12))

	b := make([]byte, 4)
	for {
		if buf.PopSize(b, 4) {
			t.Error("read >> ", b, buf.TraceString(20))
		} else {
			break
		}
	}
	t.Error(buf.TraceString(20))
}

func TestReader(t *testing.T) {
	buf := NewBytesBuffer(1)
	buf.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9})
	buf.Add([]byte{1, 2, 3, 4})
	t.Error(buf.TraceString(20))

	if false {
		r1 := io.ByteReader(buf.NewReader())
		for {
			if v, err := r1.ReadByte(); err == nil {
				t.Error("->", v)
			} else {
				break
			}
		}
	}

	if false {

		r2 := io.Reader(buf.NewReader())
		b := make([]byte, 4)
		for {
			if n, err := r2.Read(b); err == nil {
				t.Errorf("%d -> %v", n, b)
			} else {
				break
			}
		}
	}

	if true {
		b := make([]byte, 4)
		robj := buf.NewReader()
		r2 := io.Reader(robj)
		s1 := io.Seeker(robj)
		s1.Seek(6, 0)
		s1.Seek(4, 1)
		n, err := r2.Read(b)
		t.Errorf("%d -> %v, %v", n, b, err)
	}
}
