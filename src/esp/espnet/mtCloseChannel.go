package espnet

import (
	"bmautil/byteutil"
	"esp/espnet/protpack"
)

type mt_close_channel byte

func (O mt_close_channel) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	return nil
}

func (O mt_close_channel) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return true, nil
}

func (O mt_close_channel) Has(p *protpack.Package) bool {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_CLOSE_CHANNEL {
			return true
		}
	}
	return false
}

func (O mt_close_channel) Remove(p *protpack.Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_CLOSE_CHANNEL {
			p.Remove(e)
			break
		}
	}
}

func (O mt_close_channel) Set(p *protpack.Package) {
	O.Remove(p)
	f := protpack.NewFrameV(MT_CLOSE_CHANNEL, true, O)
	p.PushFront(f)
}

func ForceClose(ch Channel) {
	if ch != nil {
		msg := NewMessage()
		FrameCoders.CloseChannel.Set(msg.ToPackage())
		ch.SendMessage(msg)
	}
}
