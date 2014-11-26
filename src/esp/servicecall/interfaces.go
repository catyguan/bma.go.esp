package servicecall

import (
	"esp/espnet/espsocket"
	"time"
)

const (
	tag                 = "servicecall"
	LOOKUP_SERVICE_NAME = "lookup"
)

type ServiceCaller interface {
	Start() error
	Ping() bool
	Stop()
	Call(method string, params []interface{}, timeout time.Duration) (interface{}, error)
}

type ServiceCallerFactory interface {
	Valid(cfg map[string]interface{}) error
	Compare(cfg map[string]interface{}, old map[string]interface{}) (same bool)
	Create(name string, cfg map[string]interface{}) (ServiceCaller, error)
}

type ServiceCallHub interface {
	Parent() ServiceCallHub
	Get(serviceName string, timeout time.Duration) (ServiceCaller, error)
	LocalQuery(serviceName string) ServiceCaller
	GetServiceCallerFactory(typ string) ServiceCallerFactory
}

type SocketProvider interface {
	GetSocket() (*espsocket.Socket, error)
	Finish(sock *espsocket.Socket)
}
