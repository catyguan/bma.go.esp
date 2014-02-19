package espservice

import (
	"esp/espnet/esnp"
	"logger"
	"sync"
)

// ServiceMuxMatch
type ServiceMuxMatch func(msg *esnp.Message) bool

// ServiceMux
type ServiceMux struct {
	wlock    sync.Locker
	rlock    sync.Locker
	matchers []muxMatcher
	handlers map[string]ServiceHandler
}

type muxMatcher struct {
	matcher ServiceMuxMatch
	handler ServiceHandler
}

func NewServiceMux(wlock sync.Locker, rlock sync.Locker) *ServiceMux {
	this := new(ServiceMux)
	this.wlock = wlock
	this.rlock = rlock
	this.matchers = make([]muxMatcher, 0)
	this.handlers = make(map[string]ServiceHandler)
	return this
}

func (this *ServiceMux) AddHandler(op string, h ServiceHandler) {
	if this.wlock != nil {
		this.wlock.Lock()
		defer this.wlock.Unlock()
	}
	this.handlers[op] = h
}

func (this *ServiceMux) AddMatcher(m ServiceMuxMatch, h ServiceHandler) {
	if this.wlock != nil {
		this.wlock.Lock()
		defer this.wlock.Unlock()
	}
	this.matchers = append(this.matchers, muxMatcher{m, h})
}

func (this *ServiceMux) Match(msg *esnp.Message) ServiceHandler {
	if this.rlock != nil {
		this.rlock.Lock()
		defer this.rlock.Unlock()
	}
	for _, m := range this.matchers {
		if m.matcher(msg) {
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "%s match :-> %p", msg.GetAddress(), m.handler)
			}
			return m.handler
		}
	}
	a := msg.GetAddress()
	if a != nil {
		s := a.Get(esnp.ADDRESS_OP)
		if h, ok := this.handlers[s]; ok {
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "%s hit :-> %p", s, h)
			}
			return h
		}
	}
	return nil
}

func (this *ServiceMux) DoServe(msg *esnp.Message, rep ServiceResponser) error {
	h := this.Match(msg)
	if h != nil {
		return h(msg, rep)
	}
	err := logger.Warn(tag, "%s not found ServiceHandler", msg.GetAddress())
	return err
}

func (this *ServiceMux) Serve(msg *esnp.Message, rep ServiceResponser) error {
	return DoServiceHandle(this.DoServe, msg, rep)
}
