package golua

import (
	"fileloader"
	"fmt"
	"logger"
	"sync"
	"sync/atomic"
)

// GoLua
type GoLua struct {
	name string
	ss   fileloader.FileLoader
	cfg  *VMConfig

	codeMutex sync.RWMutex
	codes     map[string]*ChunkCode

	vmid   uint32
	closed uint32

	vmpool chan *VM

	globalMutex sync.RWMutex
	globals     map[string]interface{}
}

func NewGoLua(n string, poolSize int, ss fileloader.FileLoader, init GoLuaInitor, cfg *VMConfig) *GoLua {
	r := new(GoLua)
	r.name = n
	r.vmpool = make(chan *VM, poolSize)
	r.globals = make(map[string]interface{})
	r.ss = ss
	init(r)
	r.cfg = cfg
	r.codes = make(map[string]*ChunkCode)
	return r
}

func (this *GoLua) Close() {
	if !atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		return
	}
	func() {
		for {
			select {
			case vm := <-this.vmpool:
				vm.Destroy()
			default:
				return
			}
		}
	}()
	close(this.vmpool)
}

func (this *GoLua) String() string {
	return fmt.Sprintf("GoLua[%s]", this.name)
}

func (this *GoLua) newVM() (*VM, error) {
	err := this.CheckClose()
	if err != nil {
		return nil, err
	}
	atomic.AddUint32(&this.vmid, 1)
	vm := newVM(this, this.vmid)
	if this.cfg != nil {
		vm.Setup(this.cfg)
	}
	return vm, nil
}

func (this *GoLua) CreateVM() (*VM, error) {
	vm, err := this.newVM()
	if err != nil {
		return nil, err
	}
	vm.name = fmt.Sprintf("%s:%d", this.name, vm.id)
	st := newVMStack(nil)
	st.chunkName = fmt.Sprintf("VM<%s>", vm.name)
	vm.initStack(st)
	logger.Debug(tag, "createVM -> %s", vm)
	return vm, nil
}

func (this *GoLua) GetVM() (*VM, error) {
	err := this.CheckClose()
	if err != nil {
		return nil, err
	}
	select {
	case vm := <-this.vmpool:
		if vm != nil {
			logger.Debug(tag, "%s leave pool", vm)
			return vm, nil
		}
	default:
	}
	return this.CreateVM()
}

func (this *GoLua) ReturnVM(vm *VM) {
	defer func() {
		x := recover()
		if x != nil {
			vm.Destroy()
		}
	}()
	if !this.IsClose() {
		select {
		case this.vmpool <- vm:
			logger.Debug(tag, "%s return pool", vm)
			return
		default:
		}
	}
	vm.Destroy()
}

func (this *GoLua) CheckClose() error {
	if atomic.LoadUint32(&this.closed) == 1 {
		return fmt.Errorf("%s closed", this)
	}
	return nil
}

func (this *GoLua) IsClose() bool {
	return atomic.LoadUint32(&this.closed) == 1
}
