package eloader

import (
	"app/cacheserver"
	"container/list"
	"errors"
	"logger"
	"sync"
	"time"
)

type loadTaskInfo struct {
	name string
	task LoadTask
}

type keySpot struct {
	loadTasks []loadTaskInfo
	waitList  *list.List
}

func (this *keySpot) isEmpty() bool {
	if this.loadTasks != nil && len(this.loadTasks) > 0 {
		return false
	}
	if this.waitList != nil && this.waitList.Len() > 0 {
		return false
	}
	return true
}

func (this *keySpot) close() {
	if this.loadTasks != nil {
		for _, lt := range this.loadTasks {
			if lt.task != nil {
				lt.task.Cancel()
			}
		}
		this.loadTasks = nil
	}
}

type waiting struct {
	// info
	req *cacheserver.GetRequest

	// response
	result *cacheserver.GetResult
	rep    chan *cacheserver.GetResult
	// control
	timer   *time.Timer
	listpos *list.Element
	once    sync.Once
}

func (this *waiting) response(done bool, val []byte, err error, traces []string) {
	if this.timer != nil {
		this.timer.Stop()
	}
	this.once.Do(func() {
		r := this.result
		if err != nil {
			r.Fail(err, traces)
		} else {
			r.End(done, val, traces)
		}
		this.rep <- r
	})
}

func (this *LoaderCache) wait(req *cacheserver.GetRequest, rep chan *cacheserver.GetResult) *waiting {
	r := new(waiting)
	r.req = req
	r.result = new(cacheserver.GetResult)
	if req.Trace {
		r.result.TraceInfo = make([]string, 0)
	}

	r.rep = rep

	if this.keyspot == nil {
		this.keyspot = make(map[string]*keySpot)
	}
	key := req.Key
	ks, ok := this.keyspot[key]
	if !ok {
		ks = new(keySpot)
		this.keyspot[key] = ks
	}
	l := ks.waitList
	if l == nil {
		l = list.New()
		ks.waitList = l
	}
	e := l.PushBack(r)
	r.listpos = e
	return r
}

func (this *LoaderCache) removeWait(wt *waiting) {
	if this.keyspot == nil {
		return
	}
	key := wt.req.Key
	ks, ok := this.keyspot[key]
	if !ok {
		return
	}
	l := ks.waitList
	if l == nil {
		return
	}
	if wt.listpos != nil {
		logger.Debug(tag, "remove '%s:%s' waiting", this.name, key)
		l.Remove(wt.listpos)
		wt.listpos = nil
	}
	if ks.isEmpty() {
		logger.Debug(tag, "delete keySpot '%s'", key)
		ks.close()
		delete(this.keyspot, key)
	}
}

func (this *LoaderCache) close() {
	if this.updater != nil {
		this.updater.stop()
		this.updater = nil
	}

	this.cache.Clear()

	err := errors.New("cache close")
	for _, ks := range this.keyspot {
		l := ks.waitList
		if l != nil {
			for e := l.Front(); e != nil; e = e.Next() {
				e.Value.(*waiting).response(false, nil, err, []string{"local cache close"})
			}
			l.Init()
		}
	}
}

func (this *LoaderCache) callLoad(req *cacheserver.GetRequest) []loadTaskInfo {
	this.stats.PeerLoads++
	l := make([]loadTaskInfo, 0, len(this.loaders))
	for _, linfo := range this.loaders {
		logger.Debug(tag, "cache '%s' loading '%s' use '%s/%s'", this.name, req.Key, linfo.config.Name, linfo.config.Type)
		lt := linfo.loader.Load(this, req)
		var lti loadTaskInfo
		lti.name = linfo.config.Name
		lti.task = lt
		l = append(l, lti)
	}
	return l
}

func (this *LoaderCache) doStartLoad(req *cacheserver.GetRequest) {
	key := req.Key
	if this.loaders == nil || len(this.loaders) == 0 {
		if !req.Update {
			err := logger.Warn(tag, "cache[%s] no valid loader", this.name)
			this.doLoadEnd("", key, false, nil, err, []string{err.Error()})
		}
		return
	}

	if this.keyspot == nil {
		this.keyspot = make(map[string]*keySpot)
	}

	ks, ok := this.keyspot[key]
	if ok && ks.loadTasks != nil && len(ks.loadTasks) > 0 {
		// loading, skip
		logger.Debug(tag, "cache[%s] loading, skip", this.name)
		return
	}

	if req.Update {
		this.callLoad(req)
		this.cache.UpdateTime(key)
		return
	}

	if !ok {
		// new key spot
		ks = new(keySpot)
		this.keyspot[key] = ks
	}
	ks.loadTasks = this.callLoad(req)
	return
}

func (this *LoaderCache) doLoadEnd(loaderName string, key string, done bool, val []byte, err error, traces []string) {
	if done {
		logger.Debug(tag, "'%s:%s' load '%s' data done", this.name, loaderName, key)
		this.stats.PeerHits++
		this.cache.Put(key, val, -1)
	} else {
		logger.Debug(tag, "'%s:%s' load '%s' data fail", this.name, loaderName, key)
		this.stats.PeerErrors++
	}

	if this.keyspot == nil {
		return
	}
	ks, ok := this.keyspot[key]
	if !ok {
		return
	}

	if !done {
		// fail
		if ks.loadTasks != nil {
			l := ks.loadTasks
			for i, lti := range l {
				if lti.name == loaderName {
					sz := len(l)
					if sz == 1 {
						l = nil
					} else {
						l[i], l[sz-1] = l[sz-1], l[i]
						l = l[:sz-1]
					}
					break
				}
			}
			ks.loadTasks = l
		}

		if ks.loadTasks != nil && len(ks.loadTasks) > 0 {
			// more loader
			if logger.EnableDebug(tag) {
				s := ""
				for _, lti := range ks.loadTasks {
					if s != "" {
						s += ","
					}
					s += lti.name
				}
				logger.Debug(tag, "still wait for %s", s)
			}
			return
		}

		// all loader fail
		logger.Debug(tag, "cache '%s' all loader fail - '%s'", this.name, key)
		if this.config.InvalidHolder {
			// create holder
			logger.Debug(tag, "cache '%s' create holder - '%s'", this.name, key)
			this.cache.Put(key, nil, -1)
		}
		val = nil
		if err == nil {
			err = errors.New("all loader fail")
		}
		if traces == nil {
			traces = []string{"all loader fail"}
		} else {
			traces = append(traces, "all loader fail")
		}
	}

	if ks.waitList != nil {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "'%s' release waiting list %d", this.name, ks.waitList.Len())
		}
		l := ks.waitList
		for e := l.Front(); e != nil; e = e.Next() {
			e.Value.(*waiting).response(done, val, err, traces)
		}
		ks.waitList = nil
	}
	ks.close()
	delete(this.keyspot, key)
}
