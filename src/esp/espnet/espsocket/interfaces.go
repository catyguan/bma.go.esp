package espsocket

import (
	"esp/espnet/esnp"
	"net"
	"time"
)

type Socket interface {
	String() string

	BaseConn() net.Conn

	// 关闭
	IsBreak() bool
	AskFinish()
	AskClose()

	// 获取属性/设置属性
	GetProperty(name string) (interface{}, bool)
	SetProperty(name string, val interface{}) bool

	// 读写消息
	WriteMessage(ev *esnp.Message) error
	ReadMessage(decodeErr bool) (*esnp.Message, error)
}

type SocketProvider interface {
	GetSocket(timeout time.Duration) (Socket, error)
	Close()
}
