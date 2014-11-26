package espnetss

import (
	"boot"
	"sync"
)

type Service struct {
	lock sync.RWMutex
	ss   map[string]*SocketSource
}

func NewService() *Service {
	r := new(Service)
	r.ss = make(map[string]*SocketSource)
	return r
}

func (this *Service) Add(ss *SocketSource) bool {
	k := ss.Key()
	this.lock.Lock()
	defer this.lock.Unlock()
	_, ok := this.ss[k]
	if !ok {
		this.ss[k] = ss
		return true
	}
	return false
}

func (this *Service) Get(cfg *Config) *SocketSource {
	k := cfg.Key()
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.ss[k]
}

func (this *Service) Open(cfg *Config) (*SocketSource, error) {
	k := cfg.Key()
	this.lock.RLock()
	ss, ok := this.ss[k]
	this.lock.RUnlock()
	if ok {
		return ss, nil
	}
	err := cfg.Valid()
	if err != nil {
		return nil, err
	}
	ss = NewSocketSource(cfg)
	ss.Start()
	this.lock.Lock()
	defer this.lock.Unlock()
	ss2, ok2 := this.ss[k]
	if ok2 {
		ss.Close()
		return ss2, nil
	}
	this.ss[k] = ss
	return ss, nil
}

func (this *Service) CloseAll() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, s := range this.ss {
		s.Close()
	}
}

func (this *Service) CreateBootService(n string) *boot.BootWrap {
	r := boot.NewBootWrap(n)
	r.SetClose(func() bool {
		this.CloseAll()
		return true
	})
	return r
}
