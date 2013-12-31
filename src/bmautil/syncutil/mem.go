package syncutil

import (
	"sync/atomic"
	"unsafe"
)

type memHolder struct {
	value interface{}
}
type MemHolder struct {
	pointer unsafe.Pointer
}

func (this *MemHolder) Set(val interface{}) {
	np := &memHolder{val}
	atomic.StorePointer(&this.pointer, unsafe.Pointer(np))
}

func (this *MemHolder) Get() interface{} {
	np := (*memHolder)(atomic.LoadPointer(&this.pointer))
	if np == nil {
		return nil
	}
	return np.value
}
