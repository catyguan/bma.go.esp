package servicecall

import (
	"objfac"
	"time"
)

const (
	tag                 = "servicecall"
	NAME_GATE_SERVICE   = "*"
	NAME_LOOKUP_SERVICE = "lookup"
	NAME_LOOKUP_METHOD  = "findServiceCall"

	KIND_SERVICE_CALL = "serviceCall"
)

type PingSupported interface {
	Ping() bool
}

type ServiceCaller interface {
	Start() error
	Stop()
	SetName(n string)
	Call(serviceName, method string, params map[string]interface{}, deadline time.Time) (interface{}, error)
}

type ServiceCallLookup func(serviceName string, deadline time.Time) (map[string]interface{}, error)

func InitBaseFactory() {
	objfac.SetObjectFactory(KIND_SERVICE_CALL, "http", HttpServiceCallerFactory(0))
	objfac.SetObjectFactory(KIND_SERVICE_CALL, "esnp", ESNPServiceCallerFactory(0))
	objfac.SetObjectFactory(KIND_SERVICE_CALL, "socket", SocketServiceCallerFactory(0))
}
