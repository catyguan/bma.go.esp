package esnp

type mt_trace byte

func (O mt_trace) Has(p *Package) bool {
	return FrameCoders.Flag.Has(p, FLAG_TRACE)
	for e := p.Front(); e != nil; e = e.Next() {
		if e.MessageType() == byte(O) {
			return true
		}
	}
	return false
}

func (O mt_trace) Remove(p *Package) {
	FrameCoders.Flag.Remove(p, FLAG_TRACE)
}

func (O mt_trace) Set(p *Package) {
	FrameCoders.Flag.Set(p, FLAG_TRACE)
}

func (O mt_trace) IsReplyInfo(p *Package) bool {
	return FrameCoders.Flag.Has(p, FLAG_TRACE_INFO)
}

func (O mt_trace) CreateReply(msg *Message, info string) *Message {
	r := msg.ReplyMessage()
	FrameCoders.Flag.Set(r.pack, FLAG_TRACE_INFO)
	r.SetPayload([]byte(info))
	return r
}

func (O mt_trace) GetReplyInfo(msg *Message) string {
	bs := msg.GetPayload()
	if bs != nil {
		return string(bs)
	}
	return ""
}
