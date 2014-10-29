package aclserv

const tag = "aclserv"

type Service struct {
	name   string
	config configInfo
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	return this
}
