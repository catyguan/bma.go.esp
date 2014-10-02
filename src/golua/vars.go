package golua

import (
	"fmt"
	"sync"
)

type globalVar struct {
	vmg  *VMG
	name string
}

func (this *globalVar) Get() (interface{}, error) {
	v, ok := this.vmg.GetGlobal(this.name)
	if ok {
		return v, nil
	}
	return nil, nil
}

func (this *globalVar) Set(v interface{}) (bool, error) {
	this.vmg.SetGlobal(this.name, v)
	return true, nil
}

func (this *globalVar) String() string {
	return fmt.Sprintf("%s:%s", this.vmg, this.name)
}

type localVar struct {
	value interface{}
	mux   *sync.RWMutex
}

func (this *localVar) Get() (interface{}, error) {
	if this.mux != nil {
		this.mux.RLock()
		defer this.mux.RUnlock()
	}
	return this.value, nil
}

func (this *localVar) Set(v interface{}) (bool, error) {
	if this.mux != nil {
		this.mux.Lock()
		defer this.mux.Unlock()
	}
	this.value = v
	return true, nil
}

func (this *localVar) String() string {
	return fmt.Sprintf("localVar(%v)", this.value)
}
