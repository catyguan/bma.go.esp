package servicecall

import (
	"bmautil/valutil"
	"esp/espnet/espnetss"
	"esp/espnet/espsocket"
	"fmt"
)

type esnpnetpProvider struct {
	ss        *espnetss.SocketSource
	timeoutMS int
}

func (this *esnpnetpProvider) GetSocket() (*espsocket.Socket, error) {
	sock, err := this.ss.Open(this.timeoutMS)
	if err != nil {
		return nil, err
	}
	return sock, nil
}

func (this *esnpnetpProvider) Finish(sock *espsocket.Socket) {
	this.ss.Return(sock)
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

	ss, err0 := this.S.Open(&co.Config)
	if err0 != nil {
		return nil, err0
	}

	prov := new(esnpnetpProvider)
	prov.ss = ss

	r := new(ESNPServiceCaller)
	r.name = n
	r.provider = prov
	r.timeoutMS = co.TimeoutMS
	return r, nil
}
