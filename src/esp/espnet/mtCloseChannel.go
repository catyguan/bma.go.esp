package espnet

import (
	"bmautil/byteutil"
	"esp/espnet/protpack"
)

const (
	CLOSE_CHANNEL_NONE       = byte(0)
	CLOSE_CHANNEL_NOW        = byte(1)
	CLOES_CHANNEL_AFTER_SEND = byte(2)
)

type mt_close_channel byte

func (O mt_close_channel) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	return nil
}

func (O mt_close_channel) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	return true, nil
}

func (O mt_close_channel) Has(p *protpack.Package) byte {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_CLOSE_CHANNEL {
			return e.RawValue().(byte)
		}
	}
	return CLOSE_CHANNEL_NONE
}

func (O mt_close_channel) Remove(p *protpack.Package) {
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == MT_CLOSE_CHANNEL {
			p.Remove(e)
			break
		}
	}
}

func (O mt_close_channel) Set(p *protpack.Package, v byte) {
	O.Remove(p)
	f := protpack.NewFrameV(MT_CLOSE_CHANNEL, v, O)
	p.PushFront(f)
}

func CloseForce(ch Channel) {
	if ch != nil {
		msg := NewMessage()
		FrameCoders.CloseChannel.Set(msg.ToPackage(), CLOSE_CHANNEL_NOW)
		ch.SendMessage(msg)
	}
}

func CloseAfterSend(msg *Message) {
	FrameCoders.CloseChannel.Set(msg.ToPackage(), CLOES_CHANNEL_AFTER_SEND)
}
