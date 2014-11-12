package espnode

import (
	"esp/cluster/n2n"
	"esp/espnet/auth"
	"esp/espnet/espservice"
)

const (
	tag = "espnode"
)

type Service struct {
	name   string
	config *configInfo
	n2n    *n2n.Service
	auth   *auth.Service
}

func NewService(name string) *Service {
	this := new(Service)
	this.name = name
	this.n2n = n2n.NewService(128)
	this.auth = auth.NewService()
	this.auth.Add(this.n2n.NodeAuth)
	return this
}

func (this *Service) InitServiceMux(mux *espservice.ServiceMux) {
	mux.AddServiceHandler(n2n.SN_N2N, this.n2n.Serve)
}

func (this *Service) BindAuth(h espservice.ServiceHandler) espservice.ServiceHandler {
	this.auth.Bind(h)
	return this.auth.DoServe
}
