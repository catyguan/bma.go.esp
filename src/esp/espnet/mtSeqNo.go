package espnet

import (
	"bmautil/byteutil"
	"errors"
	"esp/espnet/protpack"
)

type struct_seq_no struct {
	seqno  int
	seqmax int
}

type mt_seq_no byte

func (O mt_seq_no) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if o, ok := v.(*struct_seq_no); ok {
		Coders.Int.DoEncode(w, o.seqno)
		Coders.Int.DoEncode(w, o.seqmax)
		return nil
	}
	return errors.New("not SeqNo")
}

func (O mt_seq_no) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v1, e1 := Coders.Int.DoDecode(r)
	if e1 != nil {
		return nil, e1
	}
	v2, e2 := Coders.Int.DoDecode(r)
	if e2 != nil {
		return nil, e2
	}
	return &struct_seq_no{v1, v2}, nil
}

func (O mt_seq_no) Get(p *protpack.Package) (int, int) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_SEQ_NO {
			v, err := e.Value(O)
			if err == nil {
				if rv, ok := v.(*struct_seq_no); ok {
					return rv.seqno, rv.seqmax
				}
			}
			break
		}
	}
	return -1, -1
}

func (O mt_seq_no) Remove(p *protpack.Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_SEQ_NO {
			p.Remove(e)
			break
		}
	}
}

func (O mt_seq_no) Set(p *protpack.Package, seqno, seqmax int) {
	O.Remove(p)
	f := protpack.NewFrameV(MT_SEQ_NO, &struct_seq_no{seqno, seqmax}, O)
	p.PushFront(f)
}

func (O mt_seq_no) IsLastSeq(p *protpack.Package) bool {
	n, m := O.Get(p)
	return O.IsLast(n, m)
}

func (O mt_seq_no) FirstSeqno() int {
	return 1
}

func (O mt_seq_no) IsFirst(seqno, seqmax int) bool {
	return seqno == 1
}
func (O mt_seq_no) IsLast(seqno, seqmax int) bool {
	if seqmax == 0 {
		return false
	}
	return seqno == seqmax
}
