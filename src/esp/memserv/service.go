package memserv

import (
	"sync"
)

type MemoryServ struct {
	lock sync.RWMutex
	mgs  map[string]*MemGo
}

func NewMemoryServ() *MemoryServ {
	r := new(MemoryServ)
	r.mgs = make(map[string]*MemGo)
	return r
}

func (this *MemoryServ) Get(n string) *MemGo {
	this.lock.RLock()
	defer this.lock.RUnlock()
	m, ok := this.mgs[n]
	if ok {
		return m
	}
	return nil
}

func (this *MemoryServ) Create(n string, cfg *MemGoConfig) (*MemGo, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	m, ok := this.mgs[n]
	if ok {
		return nil, nil
	}
	m = NewMemGo(cfg)
	err := m.Start()
	return m, err
}

func (this *MemoryServ) GetOrCreate(n string, cfg *MemGoConfig) (*MemGo, error, bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	m, ok := this.mgs[n]
	if ok {
		return m, nil, true
	}
	m = NewMemGo(cfg)
	err := m.Start()
	return m, err, false
}

func (this *MemoryServ) Remove(n string) *MemGo {
	this.lock.Lock()
	defer this.lock.Unlock()
	m, ok := this.mgs[n]
	if ok {
		delete(this.mgs, n)
		return m
	}
	return nil
}

func (this *MemoryServ) Close(n string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	m, ok := this.mgs[n]
	if ok {
		delete(this.mgs, n)
		m.Stop()
		return true
	}
	return false
}

func (this *MemoryServ) CloseAll(wait bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, m := range this.mgs {
		delete(this.mgs, k)
		if wait {
			m.goo.StopAndWait()
		} else {
			m.Stop()
		}
	}
}
