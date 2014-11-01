package golua

import (
	"fmt"
	"logger"
	"sync/atomic"
)

func (this *GoLua) GetGlobal(n string) (interface{}, bool) {
	this.globalMutex.RLock()
	defer this.globalMutex.RUnlock()
	v, ok := this.globals[n]
	return v, ok
}

func (this *GoLua) SetGlobal(n string, v interface{}) interface{} {
	this.globalMutex.Lock()
	defer this.globalMutex.Unlock()
	old := this.globals[n]
	this.globals[n] = v
	return old
}

func (this *GoLua) SetObjectFactory(n string, of GoObjectFactory) {
	this.ofMap[n] = of
}

func (this *GoLua) NewObject(n string) (interface{}, error) {
	if of, ok := this.ofMap[n]; ok {
		return of(n)
	}
	return nil, fmt.Errorf("invalid object type '%s'", n)
}

func (this *GoLua) GetConfig(n string) (interface{}, bool) {
	this.configMutex.RLock()
	defer this.configMutex.RUnlock()
	v, ok := this.configs[n]
	return v, ok
}

func (this *GoLua) SetConfig(n string, v interface{}) bool {
	this.configMutex.Lock()
	defer this.configMutex.Unlock()
	if _, ok := this.configs[n]; ok {
		return false
	}
	this.configs[n] = v
	return true
}

func (this *GoLua) GetService(n string) (interface{}, bool) {
	this.serviceMutex.RLock()
	defer this.serviceMutex.RUnlock()
	v, ok := this.services[n]
	return v, ok
}

func (this *GoLua) SetService(n string, v interface{}) bool {
	this.serviceMutex.Lock()
	defer this.serviceMutex.Unlock()
	if _, ok := this.services[n]; ok {
		return false
	}
	this.services[n] = v
	return true
}

func (this *GoLua) AddService(title string, v interface{}) string {
	this.serviceMutex.Lock()
	defer this.serviceMutex.Unlock()
	for {
		sid := atomic.AddUint32(&this.sid, 1)
		id := fmt.Sprintf("__%d_%s", sid, title)
		if _, ok := this.services[id]; !ok {
			this.services[id] = v
			return id
		}
	}
}

func (this *GoLua) RemoveService(id string) bool {
	this.serviceMutex.Lock()
	_, ok := this.services[id]
	if ok {
		delete(this.services, id)
	}
	this.serviceMutex.Unlock()
	return ok
}

func (this *GoLua) CloseService(id string) bool {
	this.serviceMutex.Lock()
	o, ok := this.services[id]
	if ok {
		delete(this.services, id)
	}
	this.serviceMutex.Unlock()
	if doClose(o) {
		logger.Debug(tag, "%s close service '%s'", this, id)
	}
	return ok
}

type GoService struct {
	GL        *GoLua
	SID       string
	CloseFunc func()
	Data      interface{}
}

func (this *GoService) Close() {
	if this.CloseFunc != nil {
		this.CloseFunc()
	}
	if this.GL != nil {
		this.GL.RemoveService(this.SID)
	}
}

func (this *GoLua) CreateGoService(title string, data interface{}, closeFunc func()) *GoService {
	o := new(GoService)
	o.GL = this
	o.SID = this.AddService(title, o)
	o.CloseFunc = closeFunc
	o.Data = data
	return o
}
