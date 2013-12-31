package socket

import (
	"fmt"
	"sync"
)

type SocketGroup struct {
	lock      sync.Mutex
	socks     map[*Socket]bool
	closeDone func()
}

func NewSocketGroup() *SocketGroup {
	r := new(SocketGroup)
	r.Init()
	return r
}

func (this *SocketGroup) Init() {
	this.socks = make(map[*Socket]bool)
}

func (this *SocketGroup) Add(s *Socket) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.socks[s]; ok {
		return
	}
	this.socks[s] = true
	id := fmt.Sprintf("SG_%p", this)
	s.AddCloseListener(func(so *Socket) { this.Remove(so) }, id)
}

func (this *SocketGroup) Remove(s *Socket) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	// fmt.Println(s)
	if _, ok := this.socks[s]; ok {
		delete(this.socks, s)
		if len(this.socks) == 0 && this.closeDone != nil {
			this.closeDone()
			this.closeDone = nil
		}
		return true
	}
	return false
}

func (this *SocketGroup) Walk(f func(sock *Socket) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for s, _ := range this.socks {
		if !f(s) {
			return
		}
	}
}

func (this *SocketGroup) CloseFunc(f func()) {
	this.closeDone = f
	this.lock.Lock()
	defer this.lock.Unlock()
	if len(this.socks) > 0 {
		for s, _ := range this.socks {
			go s.Close()
		}
	} else {
		f()
	}
}

func (this *SocketGroup) Close() {
	id := fmt.Sprintf("SG_%p", this)
	this.lock.Lock()
	defer this.lock.Unlock()
	for s, _ := range this.socks {
		delete(this.socks, s)
		s.RemoveCloseListener(id)
		s.Close()
	}
}
