package thrift

type ThriftServer struct {
	// init
	name      string
	processor TProcessor

	// config
	config *configInfo

	// runtime
	server *TSimpleServer
}

func NewThriftServer(name string, p TProcessor) *ThriftServer {
	r := new(ThriftServer)
	r.name = name
	r.processor = p
	return r
}
