package glua

import "sync"

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
