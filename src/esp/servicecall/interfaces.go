package servicecall

import "time"

const (
	tag                 = "servicecall"
	NAME_LOOKUP_SERVICE = "lookup"
)

type ServiceCaller interface {
	Start() error
	Ping() bool
	Stop()
	IsRuntime() bool
	Call(method string, params map[string]interface{}, timeout time.Duration) (interface{}, error)
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

func InitBaseFactory(s *Service) {
	s.AddServiceCallerFactory("http", HttpServiceCallerFactory(0))
	s.AddServiceCallerFactory("esnp", ESNPServiceCallerFactory(0))
}
