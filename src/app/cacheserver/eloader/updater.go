package eloader

import (
	"app/cacheserver"
	"logger"
	"time"
)

type cacheUpdater struct {
	started bool
	timer   *time.Timer
}

func newCacheUpdater() *cacheUpdater {
	this := new(cacheUpdater)
	return this
}

func (this *cacheUpdater) start(s *LoaderCache) {
	if this.started {
		return
	}
	this.started = true
	s.StepUpdate(this)
}

func (this *cacheUpdater) doUpdate(s *LoaderCache) {
	if this.timer != nil {
		this.timer.Stop()
		this.timer = nil
	}

	if !this.started {
		logger.Debug(tag, "cache[%s] updater stoped", s.name)
		this.started = false
		return
	}

	if s.updater != this {
		logger.Debug(tag, "cache[%s] not my target, skip updater", s.name)
		this.started = false
		return
	}

	if s.config.UpdateSeconds <= 0 {
		logger.Debug(tag, "cache[%s] update disabled, skip updater", s.name)
		this.started = false
		return
	}

	step := s.config.UpdateStep
	if step <= 0 {
		step = 10
	}
	usec := int64(s.config.UpdateSeconds)

	logger.Debug(tag, "ScanUpdate(%s,%d,%d)", s.name, step, usec)

	tm := time.Now().Unix() - usec
	var lastTime int64 = 0
	up := func(key string, utime int64) (r bool) {
		defer func() {
			err := recover()
			if err != nil {
				logger.Warn(tag, "ScanUpdate '%s' fail - %v", s.name, err)
				r = false
			}
		}()

		if utime > tm {
			lastTime = utime
			return false
		}

		if logger.EnableDebug(tag) {
			logger.Debug(tag, "update load '%s/%s' %v", s.name, key, time.Unix(utime, 0))
		}
		req := cacheserver.NewGetRequest(key)
		req.Update = true
		s.doStartLoad(req)
		return true
	}

	var sp int64 = 0
	end, empty := s.cache.ScanUpdate(up, step)
	if end {
		// step finish, do now
		if empty {
			sp = usec
		}
	} else {
		// have more data
		if lastTime > 0 {
			sp = lastTime + usec - time.Now().Unix()
		} else {
			sp = usec
		}
	}

	if sp <= 0 {
		// scan now
		logger.Debug(tag, "stepUpdate(%s) now/scan(%v, %v)", s.name, end, empty)
		go s.StepUpdate(this)
	} else {
		du := time.Duration(sp) * time.Second
		this.timer = time.AfterFunc(du, func() {
			s.StepUpdate(this)
		})

		logger.Debug(tag, "stepUpdate(%s) after %s", s.name, du)
	}

}

func (this *cacheUpdater) stop() {
	this.started = false
	if this.timer != nil {
		this.timer.Stop()
		this.timer = nil
	}
}
