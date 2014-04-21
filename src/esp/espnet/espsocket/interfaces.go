package espsocket

import "esp/espnet/esnp"

type SupportProp interface {
	GetProperty(name string) (interface{}, bool)
	SetProperty(name string, val interface{}) bool
}

type SendCallback func(err error)

type Channel interface {
	String() string

	// 关闭
	IsClosing() bool
	AskClose()
	Shutdown()

	// 获取属性/设置属性
	GetProperty(name string) (interface{}, bool)
	SetProperty(name string, val interface{}) bool

	// 绑定
	Bind(rec esnp.MessageListener, closeLis func())
	SendMessage(ev *esnp.Message, cb SendCallback) error
}

type BreakSupport interface {
	IsBreak() bool
}

// SocketFactory
type SocketFactory interface {
	NewSocket() (*Socket, error)
}

type SocketAcceptor func(sock *Socket) error
type SocketServer interface {
	SetAcceptor(acceptor SocketAcceptor)
}
