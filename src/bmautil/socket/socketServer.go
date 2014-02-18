package socket

type SocketAcceptor func(sock *Socket) error

type SocketServer interface {
	SetAcceptor(sa SocketAcceptor)
}
