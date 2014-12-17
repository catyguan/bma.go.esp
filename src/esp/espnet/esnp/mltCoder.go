package esnp

// Coders
var (
	MessageLineCoders messageLineCoder
)

type messageLineCoder struct {
	SessionInfo     mlt_key_values
	MessageId       mlt_message_id
	SourceMessageId mlt_message_id
	Header          mlt_key_values
	Data            mlt_key_values
	Payload         mlt_payload
	Address         mlt_address
	SourceAddress   mlt_address
	Trace           mlt_trace
	TraceResp       mlt_trace
	SeqNO           mlt_seq_no
	XData           mlt_xdata
	Error           mlt_error
	Flag            mlt_flag
	Version         mlt_version
}

func init() {
	MessageLineCoders.SessionInfo = mlt_key_values(MLT_SESSION_INFO)
	MessageLineCoders.MessageId = mlt_message_id(MLT_MESSAGE_ID)
	MessageLineCoders.SourceMessageId = mlt_message_id(MLT_SOURCE_MESSAGE_ID)
	MessageLineCoders.Header = mlt_key_values(MLT_HEADER)
	MessageLineCoders.Data = mlt_key_values(MLT_DATA)
	MessageLineCoders.Address = mlt_address(MLT_ADDRESS)
	MessageLineCoders.SourceAddress = mlt_address(MLT_SOURCE_ADDRESS)
	MessageLineCoders.Trace = mlt_trace(MLT_TRACE)
	MessageLineCoders.TraceResp = mlt_trace(MLT_TRACE_RESP)
}
