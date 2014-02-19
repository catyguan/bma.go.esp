package espchannel

import (
	"fmt"
	"sync"
)

type CloseListenerGroup struct {
	lock      sync.Mutex
	listeners map[string]func()
}

func (this *CloseListenerGroup) Set(name string, lis func()) {
	if name == "" {
		if lis != nil {
			name = fmt.Sprintf("%p", lis)
		}
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	if lis == nil {
		if this.listeners != nil {
			delete(this.listeners, name)
		}
	} else {
		if this.listeners == nil {
			this.listeners = make(map[string]func())
		}
		this.listeners[name] = lis
	}
}

func (this *CloseListenerGroup) OnClose() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.listeners != nil {
		for k, lis := range this.listeners {
			delete(this.listeners, k)
			go func() {
				defer func() {
					recover()
				}()
				lis()
			}()
		}
	}
}
