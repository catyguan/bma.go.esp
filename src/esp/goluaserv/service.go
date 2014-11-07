package goluaserv

import (
	"golua"
	"logger"
	"sync"
	"time"
)

const (
	tag = "goluaserv"
)

type glInfo struct {
	status   int
	startErr error
	gl       *golua.GoLua
}

type Service struct {
	name   string
	config *serviceConfigInfo
	glInit golua.GoLuaInitor

	lock sync.RWMutex
	gli  map[string]*glInfo
}

func NewService(n string, initor golua.GoLuaInitor) *Service {
	r := new(Service)
	r.name = n
	r.glInit = initor
	r.gli = make(map[string]*glInfo)
	return r
}

func (this *Service) _getGoLua(n string) (*glInfo, bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	gli, ok := this.gli[n]
	if !ok {
		return nil, false
	}
	if gli.status == 0 {
		return gli, false
	}
	return gli, true
}

func (this *Service) GetGoLua(n string) (*golua.GoLua, error) {
	var tm time.Time
	for {
		gli, ok := this._getGoLua(n)
		if gli == nil {
			return nil, nil
		}
		if ok {
			return gli.gl, nil
		}
		if tm.IsZero() {
			tm = time.Now()
		}
		if int(time.Since(tm).Seconds()*1000) > this.config.GetTimeoutMS {
			return nil, logger.Warn(tag, "GetGoLua(%s) timeout - %v", n, gli.startErr)
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func (this *Service) removeGoLua(k string) *glInfo {
	this.lock.Lock()
	defer this.lock.Unlock()
	gli := this.gli[k]
	if gli == nil {
		return nil
	}
	delete(this.gli, k)
	gli.gl.CloseSMMApi()
	return gli
}

func (this *Service) ResetGoLua(k string) bool {
	gli := this.removeGoLua(k)
	defer func() {
		if gli != nil && gli.gl != nil {
			gli.gl.Close()
		}
	}()
	logger.Debug(tag, "'%s' reset '%s'", this.name, k)

	glcfg := this.config.GoLua[k]
	if glcfg == nil {
		return false
	}

	this.lock.Lock()
	defer this.lock.Unlock()
	return this._create(k, glcfg)
}
