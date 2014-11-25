package servicecall

import (
	"bmautil/valutil"
	"esp/espnet/espsocket"
	"esp/espnet/mempipeline"
	"fmt"
	"strings"
)

type esnpmempProvider struct {
	sock *espsocket.Socket
}

func (this *esnpmempProvider) GetSocket() (*espsocket.Socket, error) {
	return this.sock, nil
}

func (this *esnpmempProvider) Close() {

}

type esnpmempConfig struct {
	Name      string
	TimeoutMS int
}

type ESNPMemPipelineServiceCallerFactory struct {
	S *mempipeline.Service
}

func (this *ESNPMemPipelineServiceCallerFactory) Split(n string) (string, string) {
	sp := strings.SplitN(n, ":", 2)
	if len(sp) > 1 {
		return sp[0], sp[1]
	}
	return n, ""
}

func (this *ESNPMemPipelineServiceCallerFactory) Valid(cfg map[string]interface{}) error {
	var co esnpmempConfig
	if valutil.ToBean(cfg, &co) {
		if co.Name == "" {
			return fmt.Errorf("Name empty")
		}
		n, p := this.Split(co.Name)
		if n == "" {
			return fmt.Errorf("Invalid Mempipeline name")
		}
		if p == "" {
			return fmt.Errorf("Invalid Mempipeline part")
		}
		return nil
	}
	return fmt.Errorf("invalid ESNPMemPipelineServiceCallerFactory config")
}

func (this *ESNPMemPipelineServiceCallerFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) (same bool) {
	var co, oo esnpmempConfig
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if co.Name != oo.Name {
		return false
	}
	if co.TimeoutMS != oo.TimeoutMS {
		return false
	}
	return true
}

func (this *ESNPMemPipelineServiceCallerFactory) Create(n string, cfg map[string]interface{}) (ServiceCaller, error) {
	err := this.Valid(cfg)
	if err != nil {
		return nil, err
	}

	var co esnpmempConfig
	valutil.ToBean(cfg, &co)
	mn, mp := this.Split(co.Name)

	sock := this.S.Open(mn, mp)
	prov := new(esnpmempProvider)
	prov.sock = sock

	r := new(ESNPServiceCaller)
	r.name = n
	r.provider = prov
	r.timeoutMS = co.TimeoutMS
	return r, nil
}
