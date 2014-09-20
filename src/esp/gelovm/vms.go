package gelovm

import (
	"sync"

	"code.google.com/p/gelo"
)

type VMGroup struct {
	lock sync.RWMutex
	vms  map[string]*gelo.VM
}

func NewVMGroup() *VMGroup {
	r := new(VMGroup)
	r.vms = make(map[string]*gelo.VM)
	return r
}
