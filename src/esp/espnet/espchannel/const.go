package espchannel

import (
	"esp/espnet/esnp"
	"sync/atomic"
)

const (
	tag = "espchannel"
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
	MT_CLOSE_CHANNEL = 0x09
)

const (
	SOCKET_CHANNEL_CODER_ESPNET = "espnet"
)

var (
	globalChanneIdSeq uint32
)

func NextChanneId() uint32 {
	for {
		v := atomic.AddUint32(&globalChanneIdSeq, 1)
		if v != 0 {
			return v
		}
	}
}

func TryRelyError(ch Channel, this *esnp.Message, err error) {
	if this.IsRequest() {
		rmsg := this.ReplyMessage()
		rmsg.BeError(err)
		ch.SendMessage(rmsg)
	}
}
