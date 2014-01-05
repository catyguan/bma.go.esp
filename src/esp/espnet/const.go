package espnet

const (
	tag = "espnet"
)

const (
	PROP_QEXEC_QUEUE_SIZE = "qexec.QueueSize"
	PROP_QEXEC_DEBUG      = "qexec.Debug"

	PROP_ESPNET_MAXFRAME = "espnet.maxframe"

	PROP_SOCKET_REMOTE_ADDR       = "socket.RemoteAddr"
	PROP_SOCKET_LOCAL_ADDR        = "socket.LocalAddr"
	PROP_SOCKET_DEAD_LINE         = "socket.Deadline"
	PROP_SOCKET_READ_DEAD_LINE    = "socket.ReadDeadline"
	PROP_SOCKET_WRITE_DEAD_LINE   = "socket.WriteDeadline"
	PROP_SOCKET_TRACE             = "socket.Trace"
	PROP_SOCKET_TIMEOUT           = "socket.Timeout"
	PROP_SOCKET_LINGER            = "socket.Linger"
	PROP_SOCKET_KEEP_ALIVE        = "socket.KeepAlive"
	PROP_SOCKET_KEEP_ALIVE_PERIOD = "socket.KeepAlivePeriod"
	PROP_SOCKET_NO_DELAY          = "socket.NoDelay"
	PROP_SOCKET_READ_BUFFER       = "socket.ReadBuffer"
	PROP_SOCKET_WRITE_BUFFER      = "socket.WriteBuffer"
	PROP_SOCKET_WRITE_CHAN_SIZE   = "socket.WriteChanSize"
)

const (
	MT_CLOSE_CHANNEL     = 0x09
	MT_SESSION_INFO      = 0x10
	MT_MESSAGE_ID        = 0x11
	MT_MESSAGE_KIND      = 0x12
	MT_HEADER            = 0x13
	MT_DATA              = 0x14
	MT_PAYLOAD           = 0x15
	MT_ADDRESS           = 0x16
	MT_SOURCE_ADDRESS    = 0x17
	MT_SOURCE_MESSAGE_ID = 0x18
	MT_TRACE             = 0x19
	MT_TRACE_RESP        = 0x1A
	MT_SEQ_NO            = 0x1B
)

const (
	MK_UNKNOW   = MessageKind(0)
	MK_REQUEST  = MessageKind(1)
	MK_RESPONSE = MessageKind(2)
	MK_INFO     = MessageKind(3)
	MK_EVENT    = MessageKind(4)
)

const (
	SOCKET_CHANNEL_CODER_ESPNET = "espnet"
)
