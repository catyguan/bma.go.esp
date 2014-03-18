package espservice

import (
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"logger"
	"sync"
)

// ServiceMuxMatch
type ServiceMuxMatch func(msg *esnp.Message) bool

type opHandlers map[string]ServiceHandler

// ServiceMux
type ServiceMux struct {
	wlock    sync.Locker
	rlock    sync.Locker
	matchers []muxMatcher
	handlers map[string]opHandlers
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
	this.handlers = make(map[string]opHandlers)
	return this
}

func (this *ServiceMux) AddHandler(s string, op string, h ServiceHandler) {
	if this.wlock != nil {
		this.wlock.Lock()
		defer this.wlock.Unlock()
	}
	var sh opHandlers
	ok := false
	sh, ok = this.handlers[s]
	if !ok {
		sh = make(opHandlers)
		this.handlers[s] = sh
	}
	sh[op] = h
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
		s := a.Get(esnp.ADDRESS_SERVICE)
		op := a.Get(esnp.ADDRESS_OP)
		if s != "" && op != "" {
			if sh, ok := this.handlers[s]; ok {
				if h, ok2 := sh[op]; ok2 {
					if logger.EnableDebug(tag) {
						logger.Debug(tag, "%s.%s hit :-> %p", s, op, h)
					}
					return h
				}
			}
		}
	}
	return nil
}

func (this *ServiceMux) DoServe(ch espchannel.Channel, msg *esnp.Message) error {
	h := this.Match(msg)
	if h != nil {
		return h(ch, msg)
	}
	err := logger.Warn(tag, "%s not found ServiceHandler", msg.GetAddress())
	return err
}

func (this *ServiceMux) Serve(ch espchannel.Channel, msg *esnp.Message) error {
	return DoServiceHandle(this.DoServe, ch, msg)
}
