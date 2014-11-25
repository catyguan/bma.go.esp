package servicecall

import (
	"bmautil/valutil"
	"esp/espnet/espnetss"
	"esp/espnet/espsocket"
	"fmt"
)

type esnpnetpProvider struct {
	s    *espnetss.Service
	cfg  *esnpnetpConfig
	sock *espsocket.Socket
}

func (this *esnpnetpProvider) GetSocket() (*espsocket.Socket, error) {
	if this.sock == nil {
		ss := this.s.Get(this.cfg.Host, this.cfg.User)
		if ss == nil {
			return nil, fmt.Errorf("invalid espnet(%s)", espnetss.Key(this.cfg.Host, this.cfg.User))
		}
		if this.cfg.LoginType != "" {
			ss.Add(this.cfg.Certificate, this.cfg.LoginType)
		}
		sock, err := ss.Open(this.cfg.TimeoutMS)
		if err != nil {
			return nil, err
		}
		this.sock = sock
	}
	return this.sock, nil
}

func (this *esnpnetpProvider) Close() {
	this.sock.AskClose()
	this.sock = nil
}

type esnpnetpConfig struct {
	espnetss.Config
	TimeoutMS int
}

type ESNPNetServiceCallerFactory struct {
	S *espnetss.Service
}

func (this *ESNPNetServiceCallerFactory) Valid(cfg map[string]interface{}) error {
	var co esnpnetpConfig
	if valutil.ToBean(cfg, &co) {
		return co.Valid()
	}
	return fmt.Errorf("invalid ESNPNetServiceCallerFactory config")
}

func (this *ESNPNetServiceCallerFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) (same bool) {
	var co, oo esnpnetpConfig
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if !co.Compare(&oo.Config) {
		return false
	}
	if co.TimeoutMS != oo.TimeoutMS {
		return false
	}
	return true
}

func (this *ESNPNetServiceCallerFactory) Create(n string, cfg map[string]interface{}) (ServiceCaller, error) {
	err := this.Valid(cfg)
	if err != nil {
		return nil, err
	}

	var co esnpnetpConfig
	valutil.ToBean(cfg, &co)

	prov := new(esnpnetpProvider)
	prov.s = this.S
	prov.cfg = &co

	r := new(ESNPServiceCaller)
	r.name = n
	r.provider = prov
	r.timeoutMS = co.TimeoutMS
	return r, nil
}
