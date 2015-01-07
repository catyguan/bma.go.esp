package servicecall

import (
	"bmautil/valutil"
	"esp/espnet/espsocket"
	"fmt"
	"objfac"
	"time"
)

type SocketServiceCaller struct {
	name      string
	sockp     espsocket.SocketProvider
	timeoutMS int
}

func (this *SocketServiceCaller) SetName(n string) {
	this.name = n
}

func (this *SocketServiceCaller) Ping() bool {
	if p, ok := this.sockp.(PingSupported); ok {
		return p.Ping()
	}
	return true
}

func (this *SocketServiceCaller) Start() error {
	return nil
}

func (this *SocketServiceCaller) Stop() {
	this.sockp.Close()
}

func (this *SocketServiceCaller) Call(serviceName, method string, params map[string]interface{}, deadline time.Time) (interface{}, error) {
	tm := this.timeoutMS
	if tm <= 0 {
		tm = 5000
	}
	tmd := time.Now().Add(time.Duration(tm) * time.Millisecond)
	if !deadline.IsZero() && tmd.After(deadline) {
		tmd = deadline
	}
	sock, err0 := this.sockp.GetSocket(tmd)
	if err0 != nil {
		return nil, err0
	}
	defer sock.AskFinish()
	sn := serviceName
	if sn == "" {
		sn = this.name
	}
	return ESNPCall(sock, sn, method, params, tmd)
}

type SocketServiceCallerFactory int

func (o SocketServiceCallerFactory) _valid(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) (map[string]interface{}, error) {
	if vo, ok := cfg["SP"]; ok {
		if m, ok2 := vo.(map[string]interface{}); ok2 {
			return m, objfac.DoValid(espsocket.KIND_SOCKET_PROVIDER, m, ofp)
		}
		return nil, fmt.Errorf("invalid SocketServiceCaller 'SP'")
	}
	return nil, fmt.Errorf("invalid SocketServiceCaller config 'SP'")
}

func (o SocketServiceCallerFactory) Valid(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) error {
	_, err := o._valid(cfg, ofp)
	return err
}

func (o SocketServiceCallerFactory) Compare(cfg map[string]interface{}, old map[string]interface{}, ofp objfac.ObjectFactoryProvider) (same bool) {
	tm1 := valutil.ToInt(cfg["TimeoutMS"], 5000)
	tm2 := valutil.ToInt(old["TimeoutMS"], 5000)
	if tm1 != tm2 {
		return false
	}
	var m1, m2 map[string]interface{}
	if vo, ok := cfg["SP"]; ok {
		m1, ok = vo.(map[string]interface{})
		if !ok {
			return false
		}
	}
	if vo, ok := old["SP"]; ok {
		m2, ok = vo.(map[string]interface{})
		if !ok {
			return false
		}
	}
	return objfac.DoCompare(espsocket.KIND_SOCKET_PROVIDER, m1, m2, ofp)
}

func (o SocketServiceCallerFactory) Create(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) (interface{}, error) {
	co, err := o._valid(cfg, ofp)
	if err != nil {
		return nil, err
	}
	v, err1 := objfac.DoCreate(espsocket.KIND_SOCKET_PROVIDER, co, ofp)
	if err1 != nil {
		return nil, err
	}
	sp, ok := v.(espsocket.SocketProvider)
	if !ok {
		return nil, fmt.Errorf("invalid SocketProvider(%T)", v)
	}

	v2 := cfg["TimeoutMS"]

	r := new(SocketServiceCaller)
	r.sockp = sp
	r.timeoutMS = valutil.ToInt(v2, 5*1000)
	return ServiceCaller(r), nil
}
