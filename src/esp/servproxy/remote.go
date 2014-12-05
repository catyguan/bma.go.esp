package servproxy

type RemoteObj struct {
	s       *Service
	name    string
	handler ProxyHandler
	cfg     *RemoteConfigInfo
	Data    interface{}
}

func NewRemoteObj(s *Service, n string, cfg *RemoteConfigInfo, h ProxyHandler) *RemoteObj {
	r := new(RemoteObj)
	r.s = s
	r.name = n
	r.handler = h
	r.cfg = cfg
	return r
}

func (this *RemoteObj) Start() error {
	return this.handler.Start(this)
}

func (this *RemoteObj) Stop() error {
	return this.handler.Stop(this)
}

func (this *RemoteObj) Valid() bool {
	return true
}
