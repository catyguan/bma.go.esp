package vmmmemserv

import (
	"esp/memserv"
	"fmt"
	"golua"
	"sync"
)

const (
	KEY_GLOBAL = "global"
	KEY_APP    = "app"
)

type ObjectMemServ struct {
	gl      *golua.GoLua
	s       *memserv.MemoryServ
	cfg     *memserv.MemGoConfig
	lock    sync.Mutex
	appkeys map[string]bool
}

func (this *ObjectMemServ) Get(vm *golua.VM, key string, cfg *memserv.MemGoConfig) (*memserv.MemGo, error) {
	typ, n := memserv.SplitTypeName(key)
	if n == "" {
		if typ == KEY_APP {
			key = "app@" + this.gl.GetName()
		}
	}
	r, iscreate, err := this.s.GetOrCreate(key, cfg)
	if err != nil {
		return nil, err
	}
	if iscreate && typ == KEY_APP {
		this.lock.Lock()
		this.appkeys[key] = true
		this.lock.Unlock()
	}
	return r, nil
}

func (this *ObjectMemServ) Create(vm *golua.VM, key string, cfg *memserv.MemGoConfig) (*memserv.MemGo, error) {
	typ, n := memserv.SplitTypeName(key)
	if n == "" {
		return nil, fmt.Errorf("can't create '%s'", key)
	}
	r, err := this.s.Create(key, cfg)
	if err != nil {
		return nil, err
	}
	if typ == KEY_APP {
		this.lock.Lock()
		this.appkeys[key] = true
		this.lock.Unlock()
	}
	return r, nil
}

func (this *ObjectMemServ) Close(vm *golua.VM, key string) (bool, error) {
	_, n := memserv.SplitTypeName(key)
	if n == "" {
		return false, fmt.Errorf("can't close '%s'", key)
	}
	if this.s.CloseIt(key) {
		this.lock.Lock()
		delete(this.appkeys, key)
		this.lock.Unlock()
		return true, nil
	}
	return false, nil
}

func (this *ObjectMemServ) TryClose() bool {
	r := false
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, _ := range this.appkeys {
		delete(this.appkeys, k)
		if this.s.CloseIt(k) {
			r = true
		}
	}
	return r
}
