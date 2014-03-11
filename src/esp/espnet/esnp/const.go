package esnp

const (
	tag = "esnp"
)

const (
	MT_SESSION_INFO      = 0x10
	MT_MESSAGE_ID        = 0x11
	MT_SOURCE_MESSAGE_ID = 0x12
	MT_MESSAGE_KIND      = 0x13
	MT_HEADER            = 0x14
	MT_DATA              = 0x15
	MT_ADDRESS           = 0x17
	MT_PAYLOAD           = 0x16
	MT_SOURCE_ADDRESS    = 0x18
	MT_TRACE             = 0x19
	MT_TRACE_RESP        = 0x1A
	MT_SEQ_NO            = 0x1B
	MT_XDATA             = 0x1C
	MT_ERROR             = 0x1D
)

const (
	MK_UNKNOW   = MessageKind(0)
	MK_REQUEST  = MessageKind(1)
	MK_RESPONSE = MessageKind(2)
	MK_INFO     = MessageKind(3)
	MK_EVENT    = MessageKind(4)
)

var (
	globalMessageId MessageIdGenerator = MessageIdGenerator{0, true}
)

func init() {
	globalMessageId.InitMessageIdGenerator()
}

func NextMessageId() uint64 {
	return globalMessageId.Next()
}
