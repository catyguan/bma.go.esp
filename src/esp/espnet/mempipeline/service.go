package mempipeline

import (
	"boot"
	"esp/espnet/espsocket"
	"sync"
)

type item struct {
	memp *MemPipeline
	sa   *espsocket.Socket
	sb   *espsocket.Socket
}

func (this *item) Select(p string) *espsocket.Socket {
	if p == "a" || p == "A" {
		return this.sa
	}
	return this.sb
}

type Service struct {
	lock sync.RWMutex
	mems map[string]*item
}

func NewService() *Service {
	r := new(Service)
	r.mems = make(map[string]*item)
	return r
}

func (this *Service) _create(n string, sz int) *item {
	o := new(item)
	o.memp = NewMemPipeline(n, sz)
	o.sa = espsocket.NewSocket(o.memp.ChannelA())
	o.sb = espsocket.NewSocket(o.memp.ChannelB())
	return o
}

func (this *Service) Create(n string, sz int) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, ok := this.mems[n]
	if !ok {
		this.mems[n] = this._create(n, sz)
		return true
	}
	return false
}

func (this *Service) Open(n, p string) *espsocket.Socket {
	this.lock.Lock()
	defer this.lock.Unlock()
	s, ok := this.mems[n]
	if !ok {
		s = this._create(n, 16)
		this.mems[n] = s
	}
	return s.Select(p)
}

func (this *Service) Get(n, p string) *espsocket.Socket {
	this.lock.RLock()
	defer this.lock.RUnlock()
	s, ok := this.mems[n]
	if !ok {
		return nil
	}
	return s.Select(p)
}

func (this *Service) Close() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, o := range this.mems {
		delete(this.mems, k)
		o.sa.AskClose()
		o.sb.AskClose()
	}
}

func (this *Service) CreateBootService(n string) *boot.BootWrap {
	r := boot.NewBootWrap(n)
	r.SetClose(func() bool {
		this.Close()
		return true
	})
	return r
}
