package golua

import (
	"fmt"
	"logger"
	"sync"
	"sync/atomic"
)

var (
	_max_id uint32
)

const (
	tag = "golua"
)

type VMVar interface {
	Get() (interface{}, error)
	Set(v interface{}) (bool, error)
}

type _vmSH struct {
	parent *_vmSH
	heap   map[string]interface{}
	mux    *sync.RWMutex
}

func (this *_vmSH) Lock() {
	if this.mux != nil {
		this.mux.Lock()
	}
}

func (this *_vmSH) RLock() {
	if this.mux != nil {
		this.mux.RLock()
	}
}

func (this *_vmSH) Unlock() {
	if this.mux != nil {
		this.mux.Unlock()
	}
}

func (this *_vmSH) RUnlock() {
	if this.mux != nil {
		this.mux.RUnlock()
	}
}

type VM struct {
	id          uint32
	running     bool
	kill_switch chan bool
	heritage    *_heritage
	sh          *_vmSH
	stack       []interface{}
	stackTop    int
}

type _heritage struct {
	children map[uint32]chan bool
	parent   *VM
}

func NewVM() *VM {
	vm := new(VM)
	vm.id = atomic.AddUint32(&_max_id, 1)
	vm.kill_switch = make(chan bool)
	vm.sh = new(_vmSH)
	vm.sh.heap = make(map[string]interface{})
	vm.stack = make([]interface{}, 0, 8)
	return vm
}

func (this *VM) Id() uint32 {
	return this.id
}

func (this *VM) String() string {
	return fmt.Sprintf("VM(%d)", this.id)
}

func (this *VM) Spawn() *VM {
	vm2 := NewVM()
	vm2.heritage = &_heritage{parent: this}
	//no parent
	if this.heritage == nil {
		this.heritage = &_heritage{children: make(map[uint32]chan bool)}
	}
	//parent, no children
	if this.heritage.children == nil { //has a parent but no children
		this.heritage.children = make(map[uint32]chan bool)
	}
	this.Lock()
	this.heritage.children[vm2.id] = vm2.kill_switch
	this.Unlock()
	vm2.heritage.parent = this
	vm2.sh.parent = this.sh
	logger.Debug(tag, "%s spawn a child %s", this, vm2)
	return vm2
}

func (this *VM) Lock() {
	this.sh.Lock()
}

func (this *VM) RLock() {
	this.sh.RLock()
}

func (this *VM) Unlock() {
	this.sh.Unlock()
}

func (this *VM) RUnlock() {
	this.sh.RUnlock()
}

func (this *VM) Destroy() {
	if this == nil {
		return
	}
	if this.running {
		//therefore we are in a different goroutine
		this.kill_switch <- true
		return
	}
	logger.Debug(tag, "%s destoryed", this)
	if this.heritage != nil {
		h := this.heritage
		if h.parent != nil && h.parent.heritage != nil {
			//if we grab a reference before the field is set to nil in the
			//parent, it doesn't matter whether we delete our entry
			if pm := h.parent.heritage.children; pm != nil {
				h.parent.Lock()
				delete(pm, this.id)
				h.parent.Unlock()
			}
		}
		//if we spawned any VMs, kill them
		if h.children != nil {
			for _, child := range h.children {
				child <- true
			}
			h.children = nil
		} else {
			//If there were children we cannot free the ns pointers until
			//they are dead so they don't explode before they have a chance
			//to shut down, so we have to wait for the host to discard its
			//pointer to this VM for them to be collected.
			//If there were no children, however, we can safely discard them now
			this.sh = nil
		}
	}
	this.kill_switch = nil
	this.heritage = nil
	this.sh = nil
}

func Kill(vm *VM) {
	//if vm isn't nil but kill_switch is the vm has been destroyed but
	//the host is still holding on to a pointer
	if vm != nil {
		//grab a copy in case vm is destroyed in another thread
		//between the test and the send. Sending a kill to a destroyed VM
		//is safe.
		if kill_switch := vm.kill_switch; kill_switch != nil {
			logger.Debug(tag, "%s sent kill signal", vm)
			kill_switch <- true
		}
	}
}

func (this *VM) IsDead() bool {
	return this == nil || this.sh == nil
}

func (this *VM) IsRunning() bool {
	if this == nil {
		return false
	}
	return this.running
}

func (this *VM) IsIdle() bool {
	return !this.IsDead() && !this.IsIdle()
}

func (this *VM) Call(nargs int, nresults int) error {
	old := this.running
	this.running = true
	defer func() {
		this.running = old
	}()
	// c, err := root.Exec(this)
	// if err != nil {
	// 	return err
	// }
	return nil
}
