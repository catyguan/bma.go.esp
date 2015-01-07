package servicecall

import (
	"boot"
	"fmt"
	"logger"
	"objfac"
	"sync"
	"time"
)

type islookup bool

type scItem struct {
	scid uint32
	sc   ServiceCaller
	cfg  map[string]interface{}
	prop interface{}
}

func (this *scItem) isLookup() bool {
	if this.prop == nil {
		return false
	}
	_, okl := this.prop.(islookup)
	return okl
}

var (
	gid       uint32
	glock     sync.RWMutex
	gservices map[string]*scItem
	glookup   ServiceCallLookup
)

func init() {
	gservices = make(map[string]*scItem)
}

func SetLookup(l ServiceCallLookup) {
	glookup = l
}

func GetServiceCallerFactoryByType(cfg map[string]interface{}) (objfac.ObjectFactory, string, error) {
	fac, xt, err := objfac.QueryObjectFactory(KIND_SERVICE_CALL, cfg, nil)
	if err != nil {
		return nil, "", err
	}
	return fac, xt, nil
}

func DoValid(cfg map[string]interface{}) error {
	return objfac.DoValid(KIND_SERVICE_CALL, cfg, nil)
}

func DoCompare(cfg map[string]interface{}, old map[string]interface{}) bool {
	return objfac.DoCompare(KIND_SERVICE_CALL, cfg, old, nil)
}

func DoCreate(n string, cfg map[string]interface{}) (ServiceCaller, error) {
	fac, _, err := GetServiceCallerFactoryByType(cfg)
	if err != nil {
		return nil, err
	}
	o, err1 := fac.Create(cfg, nil)
	if err1 != nil {
		return nil, err1
	}
	sc, ok := o.(ServiceCaller)
	if !ok {
		return nil, fmt.Errorf("unknow ServiceCall(%T)", o)
	}
	sc.SetName(n)
	return sc, nil
}

func SetServiceCall(k string, cfg map[string]interface{}, prop interface{}, scid uint32) (bool, uint32, error) {
	glock.Lock()
	defer glock.Unlock()
	old, ok := gservices[k]
	if ok {
		if scid == 0 {
			logger.Debug(tag, "service(%s) exists - %T", k, old.sc)
			return false, old.scid, nil
		}
		if old.scid != scid {
			logger.Debug(tag, "replace service(%s) scid not same old:%d, rep:%d", k, old.scid, scid)
			return false, old.scid, nil
		}
		if DoCompare(cfg, old.cfg) {
			logger.Debug(tag, "replace service(%s) config same, skip", k)
			return true, old.scid, nil
		}
	}

	sc, err := DoCreate(k, cfg)
	if err != nil {
		return false, 0, err
	}

	err = sc.Start()
	if err != nil {
		return false, 0, err
	}

	gid++
	if gid == 0 {
		gid++
	}
	si := new(scItem)
	si.scid = gid
	si.sc = sc
	si.cfg = cfg
	si.prop = prop
	gservices[k] = si
	if old != nil {
		logger.Debug(tag, "service(%s) replace", k)
		old.sc.Stop()
	} else {
		logger.Debug(tag, "service(%s) set", k)
	}
	return true, si.scid, nil
}

func UpdateServiceCallProp(k string, scid uint32, prop interface{}) bool {
	glock.Lock()
	defer glock.Unlock()
	si, ok := gservices[k]
	if !ok {
		return false
	}
	if scid != 0 && si.scid != scid {
		return false
	}
	si.prop = prop
	return true
}

func RemoveServiceCall(k string, scid uint32) bool {
	glock.Lock()
	defer glock.Unlock()
	si, ok := gservices[k]
	if !ok {
		return true
	}
	if scid != 0 && si.scid != scid {
		return false
	}
	delete(gservices, k)
	si.sc.Stop()
	return true
}

func Assert(serviceName string, deadline time.Time) (ServiceCaller, error) {
	sc, err := Get(serviceName, deadline)
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, fmt.Errorf("can't find service '%s'", serviceName)
	}
	return sc, nil
}

func Query(serviceName string) (ServiceCaller, uint32) {
	glock.RLock()
	defer glock.RUnlock()
	si, ok := gservices[serviceName]
	if ok {
		return si.sc, si.scid
	}
	return nil, 0
}

func Get(serviceName string, deadline time.Time) (ServiceCaller, error) {
	glock.RLock()
	si, ok := gservices[serviceName]
	glock.RUnlock()

	if ok {
		sc := si.sc
		pok := true
		if psc, ok := sc.(PingSupported); ok {
			pok = psc.Ping()
		}
		if pok {
			return sc, nil
		}
		return nil, fmt.Errorf("service(%s) ping fail", serviceName)
	}
	// lookup
	if glookup != nil {
		cfg, err := glookup(serviceName, deadline)
		if err != nil {
			return nil, err
		}
		if cfg != nil {
			sok, _, err2 := SetServiceCall(serviceName, cfg, islookup(true), 0)
			if err2 != nil {
				return nil, err2
			}
			if sok {
				logger.Debug(tag, "service(%s) lookup ok", serviceName)
			}
			glock.RLock()
			si = gservices[serviceName]
			glock.RUnlock()
		} else {
			logger.Debug(tag, "service(%s) lookup miss", serviceName)
		}
	}
	if si != nil {
		return si.sc, nil
	}
	// gate
	if serviceName == NAME_GATE_SERVICE {
		return nil, nil
	}
	sc, _ := Query(NAME_GATE_SERVICE)
	return sc, nil
}

func CallTimeout(serviceName string, method string, params map[string]interface{}, to time.Duration) (interface{}, error) {
	return Call(serviceName, method, params, time.Now().Add(to))
}

func Call(serviceName string, method string, params map[string]interface{}, deadline time.Time) (interface{}, error) {
	sc, err := Assert(serviceName, deadline)
	if err != nil {
		return nil, err
	}
	return sc.Call(serviceName, method, params, deadline)
}

func RemoveAll() {
	glock.Lock()
	defer glock.Unlock()
	for k, s := range gservices {
		delete(gservices, k)
		s.sc.Stop()
	}
}

func BootWrap() *boot.BootWrap {
	r := boot.NewBootWrap("serviceCall")
	r.SetCleanup(func() bool {
		RemoveAll()
		return true
	})
	return r
}

func BindLookupService() {
	SetLookup(func(serviceName string, deadline time.Time) (map[string]interface{}, error) {
		if serviceName == NAME_LOOKUP_SERVICE {
			return nil, nil
		}
		sc, err := Get(NAME_LOOKUP_SERVICE, deadline)
		if err != nil {
			return nil, err
		}
		if sc != nil {
			ps := make(map[string]interface{})
			ps["name"] = serviceName
			v, err2 := sc.Call(NAME_LOOKUP_SERVICE, NAME_LOOKUP_METHOD, ps, deadline)
			if err2 != nil {
				return nil, err2
			}
			if v != nil {
				if cfg, ok := v.(map[string]interface{}); ok {
					return cfg, nil
				} else {
					logger.Debug(tag, "lookup service for(%s) invalid resultType(%T)", serviceName, v)
				}
			} else {
				logger.Debug(tag, "lookup service for(%s) miss", serviceName)
			}
			return nil, nil
		} else {
			logger.Debug(tag, "miss lookup service")
		}
		return nil, nil
	})
}

func DefaultInit() {
	InitBaseFactory()
	BindLookupService()
}
