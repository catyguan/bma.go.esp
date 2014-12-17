package esnp

import (
	"errors"
	"fmt"
)

type struct_seq_no struct {
	seqno  int
	seqmax int
}

func (this *struct_seq_no) String() string {
	return fmt.Sprintf("%d/%d", this.seqno, this.seqmax)
}

type mlt_seq_no byte

func (O mlt_seq_no) Encode(w EncodeWriter, v interface{}) error {
	if o, ok := v.(*struct_seq_no); ok {
		err := Coders.Int.DoEncode(w, o.seqno)
		if err != nil {
			return err
		}
		err = Coders.Int.DoEncode(w, o.seqmax)
		return err
	}
	return errors.New("not SeqNo")
}

func (O mlt_seq_no) Decode(r DecodeReader) (interface{}, error) {
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

func (O mlt_seq_no) Get(p *Message) (int, int) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_SEQ_NO {
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

func (O mlt_seq_no) Remove(p *Message) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MLT_SEQ_NO {
			p.Remove(e)
			break
		}
	}
}

func (O mlt_seq_no) Set(p *Message, seqno, seqmax int) {
	O.Remove(p)
	f := NewMessageLineV(MLT_SEQ_NO, &struct_seq_no{seqno, seqmax}, O)
	p.PushFront(f)
}

func (O mlt_seq_no) IsLastSeq(p *Message) bool {
	n, m := O.Get(p)
	return O.IsLast(n, m)
}

func (O mlt_seq_no) FirstSeqno() int {
	return 1
}

func (O mlt_seq_no) IsFirst(seqno, seqmax int) bool {
	return seqno == 1
}
func (O mlt_seq_no) IsLast(seqno, seqmax int) bool {
	if seqmax == 0 {
		return false
	}
	return seqno == seqmax
}
