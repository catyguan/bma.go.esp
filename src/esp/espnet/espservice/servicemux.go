package espservice

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"logger"
	"sync"
)

// ServiceMuxMatch
type ServiceMuxMatch func(msg *esnp.Message) bool

type sHandler struct {
	handler    ServiceHandler
	opHandlers map[string]ServiceHandler
}

// ServiceMux
type ServiceMux struct {
	wlock    sync.Locker
	rlock    sync.Locker
	matchers []muxMatcher
	handlers map[string]*sHandler
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
	this.handlers = make(map[string]*sHandler)
	return this
}

func (this *ServiceMux) AddHandler(s string, op string, h ServiceHandler) {
	if this.wlock != nil {
		this.wlock.Lock()
		defer this.wlock.Unlock()
	}
	sh, ok := this.handlers[s]
	if !ok {
		sh = new(sHandler)
		this.handlers[s] = sh
	}
	if sh.opHandlers == nil {
		sh.opHandlers = make(map[string]ServiceHandler)
	}
	sh.opHandlers[op] = h
}

func (this *ServiceMux) AddServiceHandler(s string, h ServiceHandler) {
	if this.wlock != nil {
		this.wlock.Lock()
		defer this.wlock.Unlock()
	}
	sh, ok := this.handlers[s]
	if !ok {
		sh = new(sHandler)
		this.handlers[s] = sh
	}
	sh.handler = h
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
		s := a.GetService()
		op := a.GetOp()
		if s != "" && op != "" {
			if sh, ok := this.handlers[s]; ok {
				if sh.handler != nil {
					logger.Debug(tag, "%s hit :-> %p", s, sh.handler)
					return sh.handler
				}
				if sh.opHandlers != nil {
					if h, ok2 := sh.opHandlers[op]; ok2 {
						logger.Debug(tag, "%s.%s hit :-> %p", s, op, h)
						return h
					}
				}
			}
		}
	}
	return nil
}

func (this *ServiceMux) Serve(sock espsocket.Socket, msg *esnp.Message) error {
	h := this.Match(msg)
	if h != nil {
		return h(sock, msg)
	}
	return Miss(msg)
}

func Miss(msg *esnp.Message) error {
	return logger.Warn(tag, "%s not found ServiceHandler", msg.GetAddress())
}
