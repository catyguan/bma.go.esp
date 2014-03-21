package nodegroup

import "esp/cluster/n2n"

const (
	tag = "seedService"
)

type Service struct {
	name   string
	n2n    *n2n.Service
	config *configInfo
}

func NewService(name string, s *n2n.Service) *Service {
	this := new(Service)
	this.name = name
	return this
}
