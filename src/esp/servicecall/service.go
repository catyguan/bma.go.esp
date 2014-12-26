package servicecall

import (
	"bmautil/gotask"
	"bmautil/valutil"
	"fmt"
	"logger"
	"sync"
	"time"
)

type scItem struct {
	sc        ServiceCaller
	runtime   bool
	lookup    bool
	cfg       map[string]interface{}
	checkTime time.Time
}

type Service struct {
	name     string
	config   *configInfo
	parent   ServiceCallHub
	lock     sync.RWMutex
	services map[string]*scItem
	factorys map[string]ServiceCallerFactory
	gtask    gotask.GoTask
}

func NewService(n string, p ServiceCallHub) *Service {
	r := new(Service)
	r.name = n
	r.parent = p
	r.services = make(map[string]*scItem)
	r.factorys = make(map[string]ServiceCallerFactory)
	r.gtask.Init()
	tm := time.NewTicker(1 * time.Second)
	go func() {
		select {
		case tm := <-tm.C:
			r.checkLookup(tm)
		case <-r.gtask.C:
			logger.Debug(tag, "checker exit")
			return
		}
	}()
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

func (this *Service) _create(k string, cfg map[string]interface{}) (ServiceCaller, *scItem, error) {
	sc, err := this.DoCreate(k, cfg)
	if err != nil {
		return nil, nil, err
	}
	err = sc.Start()
	if err != nil {
		return nil, nil, err
	}
	si := new(scItem)
	si.sc = sc
	this.services[k] = si
	return sc, si, nil
}

func (this *Service) SetServiceCall(k string, sc ServiceCaller, overwrite bool) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	old, ok := this.services[k]
	if ok {
		if !overwrite {
			return false
		}
	}
	si := new(scItem)
	si.runtime = true
	si.sc = sc
	this.services[k] = si
	if old != nil {
		old.sc.Stop()
	}
	return true
}

func (this *Service) RemoveServiceCall(k string, removeRuntime bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	si, ok := this.services[k]
	if !ok {
		return
	}
	if si.runtime && !removeRuntime {
		return
	}
	delete(this.services, k)
	si.sc.Stop()
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
	si, ok := this.services[serviceName]
	this.lock.RUnlock()
	if ok {
		sc := si.sc
		if sc.Ping() {
			return sc, nil
		}
		if si.lookup {
			// remove break sc
			this.lock.Lock()
			si2, ok2 := this.services[serviceName]
			if ok2 && si2 == si {
				delete(this.services, serviceName)
				sc.Stop()
			}
			this.lock.Unlock()
		}
	}
	p := this.parent
	for p != nil {
		sc := p.LocalQuery(serviceName)
		if sc != nil && sc.Ping() {
			return sc, nil
		}
		p = p.Parent()
	}
	if serviceName == NAME_LOOKUP_SERVICE {
		return nil, nil
	}
	sc, err := this.Get(NAME_LOOKUP_SERVICE, timeout)
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, nil
	}
	ps := make(map[string]interface{})
	ps["name"] = serviceName
	rv, err1 := sc.Call("findServiceCall", ps, timeout)
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
			return sc2.sc, nil
		}
		sc, si, err1 = this._create(serviceName, cfg)
		if err1 != nil {
			return nil, err1
		}
		si.runtime = true
		si.lookup = true
		si.cfg = cfg
		return sc, nil
	}
	return nil, nil
}

func (this *Service) LocalQuery(serviceName string) ServiceCaller {
	this.lock.RLock()
	defer this.lock.RUnlock()
	si, ok := this.services[serviceName]
	if ok {
		return si.sc
	}
	return nil
}

func (this *Service) TryClose() bool {
	c := len(this.services)
	this.Stop()
	return c > 0
}

func (this *Service) checkLookup(tm time.Time) {
	this.lock.RLock()
	this.lock.RUnlock()
	if this.config != nil {
		du := time.Duration(this.config.LookupCheckSec) * time.Second
		for n, si := range this.services {
			if tm.Sub(si.checkTime) >= du {
				go this.doCheckLookup(n, si)
			}
		}
	}
}

func (this *Service) doCheckLookup(n string, si *scItem) {

}
