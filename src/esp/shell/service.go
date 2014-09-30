package shell

import "esp/gelovm"

const (
	tag = "shell"
)

type Service struct {
	name   string
	config *configInfo
	initor gelovm.VMInitor

	vms *gelovm.VMGroup
}

func NewService(n string, initor gelovm.VMInitor) *Service {
	r := new(Service)
	r.name = n
	r.initor = initor
	r.vms = gelovm.NewVMGroup()
	return r
}
