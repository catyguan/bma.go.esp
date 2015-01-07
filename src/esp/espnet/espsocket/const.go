package espsocket

import (
	"bmautil/connutil"
	"bmautil/valutil"
	"esp/espnet/esnp"
	"net"
	"time"
)

const (
	tag                     = "espsocket"
	DEFAULT_MESSAGE_MAXSIZE = 10 * 1024 * 1024
)

const (
	PROP_MESSAGE_MAXSIZE = "socket.maxsize"

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
)

func TryRelyError(sock Socket, msg *esnp.Message, err error) {
	if msg.IsRequest() {
		rmsg := msg.ReplyMessage()
		rmsg.BeError(err)
		sock.WriteMessage(rmsg)
	}
}

func CloseAfterSend(sock Socket, msg *esnp.Message) {
	sock.WriteMessage(msg)
	sock.AskClose()
}

// Call & Handler
func Call(sock Socket, msg *esnp.Message) (*esnp.Message, error) {
	return CallTimeout(sock, msg, time.Time{})
}
func CallTimeout(sock Socket, msg *esnp.Message, deadline time.Time) (*esnp.Message, error) {
	conn := sock.BaseConn()
	if conn != nil {
		if !deadline.IsZero() {
			conn.SetDeadline(deadline)
			defer conn.SetDeadline(time.Time{})
		}
	}
	err0 := sock.WriteMessage(msg)
	if err0 != nil {
		return nil, err0
	}
	rmsg, err1 := ReadResponse(sock, msg)
	if err1 != nil {
		return nil, err1
	}
	err1 = rmsg.ToError()
	return rmsg, err1
}

func SetDeadline(sock Socket, tm time.Time) bool {
	return SetProperty(sock, PROP_SOCKET_DEAD_LINE, tm)
}

func ClearDeadline(sock Socket) bool {
	return SetProperty(sock, PROP_SOCKET_DEAD_LINE, time.Time{})
}

func SetProperty(sock Socket, name string, val interface{}) bool {
	conn := sock.BaseConn()
	if conn != nil {
		switch name {
		case PROP_SOCKET_DEAD_LINE:
			if tm, ok := val.(time.Time); ok {
				conn.SetDeadline(tm)
				return true
			}
			return false
		case PROP_SOCKET_READ_DEAD_LINE:
			if tm, ok := val.(time.Time); ok {
				conn.SetReadDeadline(tm)
				return true
			}
			return false
		case PROP_SOCKET_WRITE_DEAD_LINE:
			if tm, ok := val.(time.Time); ok {
				conn.SetWriteDeadline(tm)
				return true
			}
			return false
		case PROP_SOCKET_LINGER:
			if ce, ok := conn.(*connutil.ConnExt); ok {
				cb := ce.BaseConn()
				if cn, ok := cb.(*net.TCPConn); ok {
					cn.SetLinger(valutil.ToInt(val, -1))
				}
			}
			return false
		case PROP_SOCKET_NO_DELAY:
			if ce, ok := conn.(*connutil.ConnExt); ok {
				cb := ce.BaseConn()
				if cn, ok := cb.(*net.TCPConn); ok {
					cn.SetNoDelay(valutil.ToBool(val, true))
				}
			}
			return false
		case PROP_SOCKET_READ_BUFFER:
			if ce, ok := conn.(*connutil.ConnExt); ok {
				cb := ce.BaseConn()
				if cn, ok := cb.(*net.TCPConn); ok {
					v := valutil.ToInt(val, 0)
					if v > 0 {
						cn.SetReadBuffer(v)
					}
				}
			}
			return false
		case PROP_SOCKET_WRITE_BUFFER:
			if ce, ok := conn.(*connutil.ConnExt); ok {
				cb := ce.BaseConn()
				if cn, ok := cb.(*net.TCPConn); ok {
					v := valutil.ToInt(val, 0)
					if v > 0 {
						cn.SetWriteBuffer(v)
					}
				}
			}
			return false
		}
	}
	return sock.SetProperty(name, val)
}

func GetProperty(sock Socket, name string) (interface{}, bool) {
	v, ok := sock.GetProperty(name)
	if ok {
		return v, true
	}
	conn := sock.BaseConn()
	if conn != nil {
		switch name {
		case PROP_SOCKET_REMOTE_ADDR:
			return conn.RemoteAddr().String(), true
		case PROP_SOCKET_LOCAL_ADDR:
			return conn.LocalAddr().String(), true
		}
		return conn.RemoteAddr().String(), true
	}
	return "", false
}

func ReadResponse(sock Socket, msg *esnp.Message) (*esnp.Message, error) {
	return sock.ReadMessage(true)
}
