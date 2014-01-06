package espnet

// Coders
var (
	FrameCoders frameCoder
)

type frameCoder struct {
	SessionInfo     mvCoder
	MessageId       mt_message_id
	SourceMessageId mt_message_id
	MessageKind     mt_message_kind
	Header          mvCoder
	Data            mvCoder
	Payload         mt_payload
	Address         addrCoder
	SourceAddress   addrCoder
	Trace           mt_trace
	TraceResp       mt_trace
	CloseChannel    mt_close_channel
	SeqNO           mt_seq_no
	XData           mt_xdata
}

func init() {
	FrameCoders.SessionInfo = mvCoder(MT_SESSION_INFO)
	FrameCoders.MessageId = mt_message_id(MT_MESSAGE_ID)
	FrameCoders.SourceMessageId = mt_message_id(MT_SOURCE_MESSAGE_ID)
	FrameCoders.Header = mvCoder(MT_HEADER)
	FrameCoders.Data = mvCoder(MT_DATA)
	FrameCoders.Address = addrCoder(MT_ADDRESS)
	FrameCoders.SourceAddress = addrCoder(MT_SOURCE_ADDRESS)
	FrameCoders.Trace = mt_trace(MT_TRACE)
	FrameCoders.TraceResp = mt_trace(MT_TRACE_RESP)
}
