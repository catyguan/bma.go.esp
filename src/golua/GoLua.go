package golua

import (
	"fileloader"
	"fmt"
	"logger"
	"smmapi"
	"strings"
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

	vmid    uint32
	closed  uint32
	DevMode bool

	vmpool chan *VM

	globalMutex  sync.RWMutex
	globals      map[string]interface{}
	configMutex  sync.RWMutex
	configs      map[string]interface{}
	serviceMutex sync.RWMutex
	services     map[string]interface{}
	dgMutex      sync.RWMutex
	dgs          map[uint32]*VMDebugger
	breakpoints  map[*Breakpoint]bool
	sid          uint32

	ofMap map[string]GoObjectFactory

	ExtSMMApi smmapi.SMMObject
}

func NewGoLua(n string, poolSize int, ss fileloader.FileLoader, init GoLuaInitor, cfg *VMConfig) *GoLua {
	r := new(GoLua)
	r.name = n
	r.vmpool = make(chan *VM, poolSize)
	r.configs = make(map[string]interface{})
	r.globals = make(map[string]interface{})
	r.services = make(map[string]interface{})
	r.dgs = make(map[uint32]*VMDebugger)
	r.breakpoints = make(map[*Breakpoint]bool)
	r.ofMap = make(map[string]GoObjectFactory)

	r.ss = ss
	init(r)
	r.cfg = cfg
	r.codes = make(map[string]*ChunkCode)
	return r
}

func (this *GoLua) GetName() string {
	return this.name
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

	this.serviceMutex.Lock()
	tmp := this.services
	this.services = make(map[string]interface{})
	this.serviceMutex.Unlock()

	for k, o := range tmp {
		if doClose(o) {
			if logger.EnableDebug(tag) {
				s := k
				idx := strings.Index(s, "!!")
				if idx != -1 {
					s = s[:idx] + "..."
				}
				logger.Debug(tag, "%s shutdown service '%s'", this, s)
			}
		}
	}
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
	vm, err1 := func() (*VM, error) {
		select {
		case vm := <-this.vmpool:
			if vm != nil {
				logger.Debug(tag, "%s leave pool", vm)
				vm.ResetExecutionTime()
				return vm, nil
			}
		default:
		}
		return this.CreateVM()
	}()
	if err1 != nil {
		return nil, err1
	}
	if this.DevMode {
		this.dgMutex.RLock()
		l := len(this.breakpoints)
		this.dgMutex.RUnlock()
		if l > 0 {
			vm.DebugSet(3)
		}
	}
	return vm, nil
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
