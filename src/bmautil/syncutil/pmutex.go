package syncutil

import (
	"sync"
)

type PMutex struct {
	mutex *sync.Mutex
}

func (this *PMutex) Lock() {
	if this.mutex != nil {
		this.mutex.Lock()
	}
}

func (this *PMutex) Unlock() {
	if this.mutex != nil {
		this.mutex.Unlock()
	}
}

func (this *PMutex) Enable() {
	if this.mutex != nil {
		this.mutex = new(sync.Mutex)
	}
}

type PRWMutex struct {
	mutex *sync.RWMutex
}

func (this *PRWMutex) Lock() {
	if this.mutex != nil {
		this.mutex.Lock()
	}
}

func (this *PRWMutex) Unlock() {
	if this.mutex != nil {
		this.mutex.Unlock()
	}
}

func (this *PRWMutex) RLock() {
	if this.mutex != nil {
		this.mutex.RLock()
	}
}

func (this *PRWMutex) RUnlock() {
	if this.mutex != nil {
		this.mutex.RUnlock()
	}
}

func (this *PRWMutex) Enable() {
	if this.mutex != nil {
		this.mutex = new(sync.RWMutex)
	}
}
