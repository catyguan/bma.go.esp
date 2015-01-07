package servicecallcfg

import "esp/servicecall"

const (
	tag = "servicecallcfg"
)

type Service struct {
	name   string
	config *configInfo
	scids  map[string]uint32
}

func NewService(n string) *Service {
	r := new(Service)
	r.name = n
	r.scids = make(map[string]uint32)
	return r
}

func DefaultInit() {
	servicecall.DefaultInit()
}
