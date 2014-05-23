package glua

import (
	"logger"
	"sync"
)

type Service struct {
	name     string
	config   *serviceConfigInfo
	gluaInit GLuaInit

	lock  sync.RWMutex
	gluas map[string]*GLua
}

func NewService(n string) *Service {
	r := new(Service)
	r.name = n
	r.gluas = make(map[string]*GLua)
	return r
}

func NewServiceI(n string, initor GLuaInit) *Service {
	r := NewService(n)
	r.gluaInit = initor
	return r
}

func (this *Service) GetGLua(n string) *GLua {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.gluas[n]
}

func (this *Service) ResetGLua(k string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	gl := this.gluas[k]
	if gl == nil {
		return false
	}
	logger.Debug(tag, "'%s' reset '%s'", this.name, k)
	delete(this.gluas, k)
	glcfg := this.config.GLua[k]
	if glcfg == nil {
		return false
	}
	return this._create(k, glcfg)
}
