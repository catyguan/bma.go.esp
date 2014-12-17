package esnp

type mlt_trace byte

func (O mlt_trace) Has(p *Message) bool {
	return MessageLineCoders.Flag.Has(p, FLAG_TRACE)
}

func (O mlt_trace) Remove(p *Message) {
	MessageLineCoders.Flag.Remove(p, FLAG_TRACE)
}

func (O mlt_trace) Set(p *Message) {
	MessageLineCoders.Flag.Set(p, FLAG_TRACE)
}

func (O mlt_trace) IsReplyInfo(p *Message) bool {
	return MessageLineCoders.Flag.Has(p, FLAG_TRACE_INFO)
}

func (O mlt_trace) CreateReply(msg *Message, info string) *Message {
	r := msg.ReplyMessage()
	MessageLineCoders.Flag.Set(r, FLAG_TRACE_INFO)
	r.SetPayload([]byte(info))
	return r
}

func (O mlt_trace) GetReplyInfo(msg *Message) string {
	bs := msg.GetPayload()
	if bs != nil {
		return string(bs)
	}
	return ""
}
