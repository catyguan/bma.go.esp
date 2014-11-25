package espnetss

import (
	"boot"
	"fmt"
	"sync"
)

type Config struct {
	Host        string
	User        string
	LoginType   string
	Certificate string
	PreConns    int
}

func (this *Config) Valid() error {
	if this.Host == "" {
		return fmt.Errorf("Host empty")
	}
	return nil
}

func (this *Config) Compare(o *Config) bool {
	if this.Host != o.Host {
		return false
	}
	if this.User != o.User {
		return false
	}
	if this.LoginType != o.LoginType {
		return false
	}
	if this.Certificate != o.Certificate {
		return false
	}
	return true
}

type Service struct {
	lock sync.RWMutex
	ss   map[string]*SocketSource
}

func NewService() *Service {
	r := new(Service)
	r.ss = make(map[string]*SocketSource)
	return r
}

func Key(host string, user string) string {
	return fmt.Sprintf("%s@%s", user, host)
}

func (this *Service) Add(ss *SocketSource) bool {
	k := ss.Name()
	this.lock.Lock()
	defer this.lock.Unlock()
	_, ok := this.ss[k]
	if !ok {
		this.ss[k] = ss
		return true
	}
	return false
}

func (this *Service) Register(cfg *Config) {
	k := Key(cfg.Host, cfg.User)
	this.lock.Lock()
	defer this.lock.Unlock()
	ss, ok := this.ss[k]
	if !ok {
		ss = NewSocketSource(cfg.Host, cfg.User, cfg.PreConns)
		this.ss[k] = ss
	}
	ss.Add(cfg.Certificate, cfg.LoginType)
}

func (this *Service) Get(host string, user string) *SocketSource {
	k := Key(host, user)
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.ss[k]
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
