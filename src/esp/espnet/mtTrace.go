package espnet

import (
	"bmautil/byteutil"
	"esp/espnet/protpack"
)

type mt_trace byte

func (O mt_trace) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	return nil
}

func (O mt_trace) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return true, nil
}

func (O mt_trace) Has(p *protpack.Package) bool {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			return true
		}
	}
	return false
}

func (O mt_trace) Remove(p *protpack.Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			p.Remove(e)
			break
		}
	}
}

func (O mt_trace) Set(p *protpack.Package) {
	O.Remove(p)
	f := protpack.NewFrameV(byte(O), true, O)
	p.PushFront(f)
}

func (O mt_trace) CreateReply(msg *Message, info string) *Message {
	r := msg.ReplyMessage()
	f := protpack.NewFrameV(MT_TRACE_RESP, true, O)
	r.ToPackage().PushFront(f)
	r.SetPayload([]byte(info))
	return r
}

func (O mt_trace) GetReplyInfo(msg *Message) string {
	bs := msg.GetPayloadB()
	if bs != nil {
		return string(bs)
	}
	return ""
}
