package semaphore

import (
	"fmt"
	"logger"
	"sync"
	"sync/atomic"
	"time"
)

const (
	tag = "semaphore"
)

type semObject struct {
	info      *semInfo
	ch        chan int
	waitCount int32
	execCount int32
}

func (this *semObject) Close() {
	if this.ch != nil {
		close(this.ch)
	}
}

func (this *semObject) String() string {
	return fmt.Sprintf("%s, %d/%d/%d", this.info.Name, this.execCount, this.waitCount, this.info.Limit)
}

type Service struct {
	name   string
	config *configInfo

	slock sync.RWMutex
	sems  map[string]*semObject
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	this.sems = make(map[string]*semObject)
	return this
}

func (this *Service) getSem(name string) *semObject {
	var r *semObject
	this.slock.RLock()
	r = this.sems[name]
	this.slock.RUnlock()
	if r != nil {
		return r
	}

	this.slock.Lock()
	defer this.slock.Unlock()
	r = this.sems[name]
	if r != nil {
		return r
	}

	// create new semObject
	r = new(semObject)
	cfg := this.config.Find(name)
	if cfg == nil {
		cfg = this.config.DefaultSem
	}
	r.info = cfg
	r.ch = make(chan int, cfg.Limit)
	for i := 0; i < cfg.Limit; i++ {
		r.ch <- 1
	}
	this.sems[name] = r
	return r
}

func (this *Service) acquire(so *semObject, waitTimeout time.Duration) bool {
	atomic.AddInt32(&so.waitCount, 1)
	defer atomic.AddInt32(&so.waitCount, -1)

	if logger.EnableDebug(tag) {
		if len(so.ch) == 0 {
			logger.Debug(tag, "'%s' wait semObject[%s]", this.name, so)
		}
	}
	select {
	case v := <-so.ch:
		if v == 0 {
			logger.Debug(tag, "'%s' semObject[%s] closed", this.name, so)
			return false
		}
		return true
	default:
	}

	tm := time.NewTimer(waitTimeout)
	select {
	case v := <-so.ch:
		tm.Stop()
		if v == 0 {
			logger.Debug(tag, "'%s' semObject[%s] closed", this.name, so)
			return false
		}
	case <-tm.C:
		// wait timeout
		logger.Debug(tag, "'%s' semObject[%s] timeout", this.name, so)
		return false
	}
	return true
}

func (this *Service) release(so *semObject) {
	defer func() {
		recover()
	}()
	so.ch <- 1
}

func (this *Service) Execute(name string, callFunc func(), waitTimeout time.Duration) bool {
	so := this.getSem(name)
	execMaxMS := 30 * 1000
	if so != nil {
		if !this.acquire(so, waitTimeout) {
			return false
		}
		execMaxMS = so.info.ExecuteTimeoutMS
		atomic.AddInt32(&so.execCount, 1)
	} else {
		logger.Warn(tag, "'%s' semObject[%s] is nil", this.name, name)
	}
	now := time.Now()
	defer func() {
		end := time.Now()
		if so != nil {
			atomic.AddInt32(&so.execCount, -1)
			this.release(so)
		}
		sp := end.Sub(now)
		if sp > time.Duration(execMaxMS)*time.Millisecond {
			logger.Warn(tag, "'%s' semObject[%s] execute too long [%f]s", this.name, name, sp.Seconds())
		}
	}()
	callFunc()
	return true
}
