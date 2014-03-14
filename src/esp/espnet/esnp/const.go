package esnp

const (
	tag = "esnp"
)

const (
	MT_SESSION_INFO      = 0x10
	MT_MESSAGE_ID        = 0x11
	MT_SOURCE_MESSAGE_ID = 0x12
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
	MT_FLAG              = 0x1E
	MT_VERSION           = 0x1F
)

const (
	FLAG_TRACE      = MTFlag(1)
	FLAG_TRACE_INFO = MTFlag(2)
	FLAG_RESP       = MTFlag(3)
	FLAG_REQUEST    = MTFlag(4)
	FLAG_INFO       = MTFlag(5)
	FLAG_EVENT      = MTFlag(6)
	FLAG_APP_DEFINE = MTFlag(128)
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
