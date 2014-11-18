package vmmmemserv

import (
	"esp/memserv"
	"fmt"
	"golua"
	"strings"
)

const (
	KEY_GLOBAL     = "global"
	KEY_APP        = "app"
	KEY_PREFIX_APP = "app-"
)

type ObjectMemServ struct {
	gl  *golua.GoLua
	s   *memserv.MemoryServ
	cfg *memserv.MemGoConfig
}

func (this *ObjectMemServ) DefaultConfig() *memserv.MemGoConfig {
	if this.cfg == nil {
		c := new(memserv.MemGoConfig)
		c.QSize = 128
		c.Valid()
		this.cfg = c
	}
	return this.cfg
}

func (this *ObjectMemServ) Get(vm *golua.VM, n string, cfg *memserv.MemGoConfig) (*memserv.MemGo, error) {
	switch n {
	case KEY_GLOBAL:
		if cfg == nil {
			cfg = this.DefaultConfig()
		}
		r, err := this.s.GetOrCreate(n, cfg)
		if err != nil {
			return nil, err
		}
		return r, nil
	case KEY_APP:
		key := KEY_PREFIX_APP + this.gl.GetName()
		if cfg == nil {
			cfg = this.DefaultConfig()
		}
		r, err := this.s.GetOrCreate(key, cfg)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if strings.HasPrefix(n, KEY_PREFIX_APP) {
		return nil, fmt.Errorf("can't get '%s*'", KEY_PREFIX_APP)
	}
	if cfg != nil {
		r, err := this.s.GetOrCreate(n, cfg)
		if err != nil {
			return nil, err
		}
		return r, nil
	} else {
		r := this.s.Get(n)
		return r, nil
	}
}

func (this *ObjectMemServ) Create(vm *golua.VM, n string, cfg *memserv.MemGoConfig) (*memserv.MemGo, error) {
	switch n {
	case KEY_GLOBAL, KEY_APP:
		return nil, fmt.Errorf("can't create '%s'", n)
	}
	if strings.HasPrefix(n, KEY_PREFIX_APP) {
		return nil, fmt.Errorf("can't create '%s*'", KEY_PREFIX_APP)
	}
	r, err := this.s.Create(n, cfg)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (this *ObjectMemServ) TryClose() bool {
	r := false
	key := KEY_PREFIX_APP + this.gl.GetName()
	if this.s.Close(key) {
		r = true
	}
	return r
}
