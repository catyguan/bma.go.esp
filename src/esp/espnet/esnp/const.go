package esnp

const (
	tag = "esnp"
)

const (
	MLT_END               = byte(0x00)
	MLT_SESSION_INFO      = 0x10
	MLT_MESSAGE_ID        = 0x11
	MLT_SOURCE_MESSAGE_ID = 0x12
	MLT_HEADER            = 0x14
	MLT_DATA              = 0x15
	MLT_ADDRESS           = 0x17
	MLT_PAYLOAD           = 0x16
	MLT_SOURCE_ADDRESS    = 0x18
	MLT_TRACE             = 0x19
	MLT_TRACE_RESP        = 0x1A
	MLT_SEQ_NO            = 0x1B
	MLT_XDATA             = 0x1C
	MLT_ERROR             = 0x1D
	MLT_FLAG              = 0x1E
	MLT_VERSION           = 0x1F
)

const (
	FLAG_TRACE      = Flag(1)
	FLAG_TRACE_INFO = Flag(2)
	FLAG_RESP       = Flag(3)
	FLAG_REQUEST    = Flag(4)
	FLAG_INFO       = Flag(5)
	FLAG_EVENT      = Flag(6)
	FLAG_APP_DEFINE = Flag(128)
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
