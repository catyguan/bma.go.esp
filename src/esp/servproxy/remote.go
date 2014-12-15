package servproxy

import (
	"sync/atomic"
	"time"
)

type RemoteObj struct {
	s        *Service
	name     string
	handler  RemoteHandler
	cfg      *RemoteConfigInfo
	Data     interface{}
	failFlag uint32
	failTime time.Time
}

func NewRemoteObj(s *Service, n string, cfg *RemoteConfigInfo, h RemoteHandler) *RemoteObj {
	r := new(RemoteObj)
	r.s = s
	r.name = n
	r.handler = h
	r.cfg = cfg
	return r
}

func (this *RemoteObj) Start() error {
	return this.handler.Start(this)
}

func (this *RemoteObj) Stop() error {
	return this.handler.Stop(this)
}

func (this *RemoteObj) Ping() bool {
	ff := atomic.LoadUint32(&this.failFlag)
	if ff == 1 {
		fr := this.cfg.FailRetryMS
		if fr <= 0 {
			fr = 30 * 1000
		}
		du := time.Since(this.failTime)
		if du > time.Duration(fr)*time.Millisecond {
			if atomic.CompareAndSwapUint32(&this.failFlag, 1, 0) {
				ff = 0
			}
		}
		if ff == 1 {
			return false
		}
	}
	cp, ok := this.handler.Ping(this)
	if cp {
		return ok
	}
	return true
}

func (this *RemoteObj) Fail() {
	if atomic.CompareAndSwapUint32(&this.failFlag, 0, 1) {
		this.failTime = time.Now()
	}
}
