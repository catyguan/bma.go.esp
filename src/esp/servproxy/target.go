package servproxy

import "fmt"

type TargetObj struct {
	s       *Service
	name    string
	cfg     *TargetConfigInfo
	handler RemoteHandler
	remotes []*RemoteObj
}

func NewTargetObj(s *Service, n string, cfg *TargetConfigInfo, h RemoteHandler) *TargetObj {
	r := new(TargetObj)
	r.s = s
	r.name = n
	r.handler = h
	r.cfg = cfg
	r.remotes = make([]*RemoteObj, len(cfg.Remotes))
	for i, rcfg := range cfg.Remotes {
		r.remotes[i] = NewRemoteObj(s, fmt.Sprintf("%s_%d", n, i), rcfg, h)
	}
	return r
}

func (this *TargetObj) Start() error {
	for _, o := range this.remotes {
		err0 := o.Start()
		if err0 != nil {
			return err0
		}
	}
	return nil
}

func (this *TargetObj) Stop() {
	for _, o := range this.remotes {
		o.Stop()
	}
}
