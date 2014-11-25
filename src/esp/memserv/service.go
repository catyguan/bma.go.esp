package memserv

import (
	"boot"
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

func (this *MemoryServ) _create(n string, cfg *MemGoConfig) (*MemGo, error) {
	m := NewMemGo(n, cfg)
	err := m.Start()
	if err == nil {
		this.mgs[n] = m
	}
	return m, err
}

// return MemGo,IsCrete,error
func (this *MemoryServ) GetOrCreate(n string, cfg *MemGoConfig) (*MemGo, bool, error) {
	mg := this.Get(n)
	if mg != nil {
		return mg, false, nil
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	m, ok := this.mgs[n]
	if ok {
		return m, false, nil
	}
	m2, err := this._create(n, cfg)
	if err != nil {
		return nil, false, err
	}
	return m2, true, nil
}

func (this *MemoryServ) Create(n string, cfg *MemGoConfig) (*MemGo, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, ok := this.mgs[n]
	if ok {
		return nil, nil
	}
	return this._create(n, cfg)
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

func (this *MemoryServ) CreateBootService(n string) *boot.BootWrap {
	r := boot.NewBootWrap(n)
	r.SetClose(func() bool {
		this.CloseAll(true)
		return true
	})
	return r
}
