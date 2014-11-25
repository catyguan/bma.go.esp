package servicecall

import (
	"bmautil/valutil"
	"fmt"
	"sync"
	"time"
)

type Service struct {
	name     string
	config   *configInfo
	parent   ServiceCallHub
	lock     sync.RWMutex
	services map[string]ServiceCaller
	factorys map[string]ServiceCallerFactory
}

func NewService(n string, p ServiceCallHub) *Service {
	r := new(Service)
	r.name = n
	r.parent = p
	r.services = make(map[string]ServiceCaller)
	r.factorys = make(map[string]ServiceCallerFactory)
	return r
}

func (this *Service) Parent() ServiceCallHub {
	return this.parent
}

func (this *Service) AddServiceCallerFactory(n string, fac ServiceCallerFactory) {
	this.factorys[n] = fac
}

func (this *Service) GetServiceCallerFactory(n string) ServiceCallerFactory {
	fac, ok := this.factorys[n]
	if ok {
		return fac
	}
	if this.parent != nil {
		return this.parent.GetServiceCallerFactory(n)
	}
	return nil
}

func (this *Service) GetServiceCallerFactoryByType(cfg map[string]interface{}) (ServiceCallerFactory, string, error) {
	xt, ok := cfg["Type"]
	if !ok {
		return nil, "", fmt.Errorf("miss Type")
	}
	vxt := valutil.ToString(xt, "")
	if vxt == "" {
		return nil, "", fmt.Errorf("invalid Type(%v)", xt)
	}
	fac := this.GetServiceCallerFactory(vxt)
	if fac == nil {
		return nil, "", fmt.Errorf("invalid ServiceCaller Type(%s)", xt)
	}
	return fac, vxt, nil
}

func (this *Service) DoValid(cfg map[string]interface{}) error {
	fac, _, err := this.GetServiceCallerFactoryByType(cfg)
	if err != nil {
		return err
	}
	return fac.Valid(cfg)
}

func (this *Service) DoCompare(cfg map[string]interface{}, old map[string]interface{}) bool {
	fac1, xt1, err1 := this.GetServiceCallerFactoryByType(cfg)
	if err1 != nil {
		return false
	}
	_, xt2, err2 := this.GetServiceCallerFactoryByType(old)
	if err2 != nil {
		return false
	}
	if xt1 != xt2 {
		return false
	}
	return fac1.Compare(cfg, old)
}

func (this *Service) DoCreate(n string, cfg map[string]interface{}) (ServiceCaller, error) {
	fac, _, err := this.GetServiceCallerFactoryByType(cfg)
	if err != nil {
		return nil, err
	}
	return fac.Create(n, cfg)
}

func (this *Service) _create(k string, cfg map[string]interface{}) (ServiceCaller, error) {
	sc, err := this.DoCreate(k, cfg)
	if err != nil {
		return nil, err
	}
	err = sc.Start()
	if err != nil {
		return nil, err
	}
	this.services[k] = sc
	return sc, nil
}

func (this *Service) SetServiceCall(k string, sc ServiceCaller, overwrite bool) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	ss, ok := this.services[k]
	if ok {
		if !overwrite {
			return false
		}
	}
	this.services[k] = sc
	if ss != nil {
		ss.Stop()
	}
	return true
}

func (this *Service) RemoveServiceCall(k string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	ss, ok := this.services[k]
	if !ok {
		return
	}
	delete(this.services, k)
	ss.Stop()
}

func (this *Service) Assert(serviceName string, timeout time.Duration) (ServiceCaller, error) {
	sc, err := this.Get(serviceName, timeout)
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, fmt.Errorf("can't find service '%s'", serviceName)
	}
	return sc, nil
}

func (this *Service) Get(serviceName string, timeout time.Duration) (ServiceCaller, error) {
	this.lock.RLock()
	sc, ok := this.services[serviceName]
	this.lock.RUnlock()
	if ok {
		if sc.Ping() {
			return sc, nil
		}
		this.lock.Lock()
		sc2, ok2 := this.services[serviceName]
		if ok2 && sc2 == sc {
			delete(this.services, serviceName)
			sc.Stop()
		}
		this.lock.Unlock()
	}
	p := this.parent
	for p != nil {
		sc := p.LocalQuery(serviceName)
		if sc != nil && sc.Ping() {
			return sc, nil
		}
		p = p.Parent()
	}
	if serviceName == LOOKUP_SERVICE_NAME {
		return nil, nil
	}
	sc, err := this.Get(LOOKUP_SERVICE_NAME, timeout)
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, nil
	}
	rv, err1 := sc.Call("do", []interface{}{serviceName}, timeout)
	if err1 != nil {
		return nil, err1
	}
	if rv == nil {
		return nil, nil
	}
	if cfg, ok := rv.(map[string]interface{}); ok {
		this.lock.Lock()
		defer this.lock.Unlock()
		if sc2, ok := this.services[serviceName]; ok {
			return sc2, nil
		}
		sc, err1 = this._create(serviceName, cfg)
		if err1 != nil {
			return nil, err1
		}
		return sc, nil
	}
	return nil, nil
}

func (this *Service) LocalQuery(serviceName string) ServiceCaller {
	this.lock.RLock()
	defer this.lock.RUnlock()
	sc, ok := this.services[serviceName]
	if ok {
		return sc
	}
	return nil
}

func (this *Service) TryClose() bool {
	c := len(this.services)
	this.Stop()
	return c > 0
}
