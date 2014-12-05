package servproxy

import (
	"fmt"
	"logger"
	"net"
	"time"
)

type ProxyHandler interface {
	Handle(s *Service, port *PortObj, conn net.Conn)
	AnswerError(port *PortObj, req interface{}, err error) error

	Valid(cfg *RemoteConfigInfo) error
	Compare(cfg *RemoteConfigInfo, old *RemoteConfigInfo) bool
	Start(obj *RemoteObj) error
	Stop(obj *RemoteObj) error
	Forward(port *PortObj, req interface{}, remote *RemoteObj) error
}

var (
	ghlibs map[string]ProxyHandler = make(map[string]ProxyHandler)
)

func AddProxyHandler(n string, h ProxyHandler) {
	ghlibs[n] = h
}

func GetProxyHandler(n string) ProxyHandler {
	return ghlibs[n]
}

func AssertProxyHandler(typ string) (ProxyHandler, error) {
	h := GetProxyHandler(typ)
	if h == nil {
		return nil, fmt.Errorf("invalid ProxyHandler Type(%s)", typ)
	}
	return h, nil
}

type DebugConn struct {
	conn net.Conn
}

func (this *DebugConn) Read(b []byte) (n int, err error) {
	n, err = this.conn.Read(b)
	if n > 0 {
		logger.Debug(tag, "%X", b[:n])
	}
	return n, err
}

func (this *DebugConn) Write(b []byte) (n int, err error) {
	return this.conn.Write(b)
}

func (this *DebugConn) Close() error {
	return this.conn.Close()
}

func (this *DebugConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *DebugConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *DebugConn) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}

func (this *DebugConn) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}

func (this *DebugConn) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}

func (this *DebugConn) String() string {
	return this.conn.RemoteAddr().String()
}
