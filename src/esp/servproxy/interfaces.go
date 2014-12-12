package servproxy

import (
	"fmt"
	"logger"
	"net"
)

type ProxyRequest interface {
	Type() string
	BeginRead() error
	Read() (bool, []byte, error)
	Finish()
	HasResponse() bool
}

type PortHandler interface {
	Handle(s *Service, port *PortObj, conn net.Conn)
	Write(port *PortObj, req interface{}, b []byte) error
	AnswerError(port *PortObj, req interface{}, err error) error
}

type RemoteSession interface {
	Write(b []byte) error
	Read() (bool, []byte, error)
	Fail()
	Finish()
}

type RemoteHandler interface {
	Valid(cfg *RemoteConfigInfo) error
	Compare(cfg *RemoteConfigInfo, old *RemoteConfigInfo) bool
	Start(obj *RemoteObj) error
	Stop(obj *RemoteObj) error
	Ping(remote *RemoteObj) (canPing bool, ok bool)
	Begin(remote *RemoteObj) (RemoteSession, error)
}

var (
	gphlibs map[string]PortHandler   = make(map[string]PortHandler)
	grhlibs map[string]RemoteHandler = make(map[string]RemoteHandler)
)

func AddPortHandler(n string, h PortHandler) {
	gphlibs[n] = h
}

func GetPortHandler(n string) PortHandler {
	return gphlibs[n]
}

func AssertPortHandler(typ string) (PortHandler, error) {
	h := GetPortHandler(typ)
	if h == nil {
		return nil, fmt.Errorf("invalid PortHandler Type(%s)", typ)
	}
	return h, nil
}

func AddRemoteHandler(n string, h RemoteHandler) {
	grhlibs[n] = h
}

func GetRemoteHandler(n string) RemoteHandler {
	return grhlibs[n]
}

func AssertRemoteHandler(typ string) (RemoteHandler, error) {
	h := GetRemoteHandler(typ)
	if h == nil {
		return nil, fmt.Errorf("invalid RemoteHandler Type(%s)", typ)
	}
	return h, nil
}

func ConnDebuger(conn net.Conn, b []byte, read bool) {
	if logger.EnableDebug(tag) {
		if read {
			logger.Debug(tag, "%s -> %X", conn.RemoteAddr(), b)
		} else {
			logger.Debug(tag, "%s <- %X", conn.RemoteAddr(), b)
		}
	}
}
