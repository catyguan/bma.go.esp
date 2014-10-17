package goluaserv

import (
	"golua"
	"logger"
	"sync"
)

const (
	tag = "goluaserv"
)

type Service struct {
	name    string
	config  *serviceConfigInfo
	vmgInit golua.VMGInitor

	lock sync.RWMutex
	gl   map[string]*golua.GoLua
}

func NewService(n string, initor golua.VMGInitor) *Service {
	r := new(Service)
	r.name = n
	r.vmgInit = initor
	r.gl = make(map[string]*golua.GoLua)
	return r
}

func (this *Service) GetGoLua(n string) *golua.GoLua {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.gl[n]
}

func (this *Service) removeGoLua(k string) *golua.GoLua {
	this.lock.Lock()
	defer this.lock.Unlock()
	gl := this.gl[k]
	if gl == nil {
		return nil
	}
	delete(this.gl, k)
	return gl
}

func (this *Service) ResetGoLua(k string) bool {
	gl := this.removeGoLua(k)
	defer func() {
		if gl != nil {
			gl.Close()
		}
	}()
	logger.Debug(tag, "'%s' reset '%s'", this.name, k)

	glcfg := this.config.GoLua[k]
	if glcfg == nil {
		return false
	}

	this.lock.Lock()
	defer this.lock.Unlock()
	return this._create(k, glcfg)
}
