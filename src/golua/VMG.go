package golua

import "sync"

type VMG struct {
	name    string
	mux     sync.RWMutex
	vmid    uint32
	vms     map[uint32]*VM
	globals map[string]interface{}
	closed  uint32
	gl      *GoLua
}
