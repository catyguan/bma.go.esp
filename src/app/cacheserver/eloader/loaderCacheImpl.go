package eloader

import (

	// "fmt"

	"app/cacheserver"
	"errors"
	"logger"
	"time"
)

func (this *LoaderCache) requestHandler(ev interface{}) (bool, error) {
	switch rv := ev.(type) {
	case func() error:
		return true, rv()
	}
	return true, nil
}

func (this *LoaderCache) stopHandler() {
	this.executor = nil
	if this.updater != nil {
		this.updater.stop()
		this.updater = nil
	}

	if this.cache != nil {
		this.cache.Clear()
		this.cache = nil
	}

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

func logEvit(cache *cacheserver.Cache, name string) {
	if logger.EnableDebug(tag) {
		cache.Listener = func(key string, val []byte) {
			logger.Debug(tag, "evit '%s' %s", name, key)
		}
	}
}

func (this *LoaderCache) doCreateCache() error {
	cache := this.config.NewCache()
	logEvit(cache, this.name)
	this.cache = cache
	return nil
}

func (this *LoaderCache) doDeployCache(cfg *cacheConfig) error {
	logger.Info(tag, "reploy cache '%s'", this.name)
	*this.config = *cfg
	if this.cache.MaxSize() != cfg.Maxsize {
		// clone cache first
		old := this.cache
		logger.Info(tag, "resize cache %d -> %d", old.MaxSize(), cfg.Maxsize)
		ncache := cfg.NewCache()
		ncache.Clone(old)
		logEvit(ncache, this.name)

		this.cache = ncache

		logger.Info(tag, "resize cache %d -> %d done", old.MaxSize(), cfg.Maxsize)

		old.Clear()
		old = nil
	} else {
		if cfg.ValidSeconds > 0 {
			this.cache.ValidTime = int64(cfg.ValidSeconds)
		} else {
			this.cache.ValidTime = 0
		}
	}

	if this.updater != nil {
		if cfg.UpdateSeconds > 0 {
			this.updater.start(this)
		} else {
			this.updater.stop()
		}
	}

	// rebuild loaders
	names := make([]string, 0)
	for _, linfo := range this.loaders {
		names = append(names, linfo.config.Name)
	}
	for _, n := range names {
		this.doRemoveLoader(n)
	}
	for _, lcfg := range cfg.Loaders {
		this.doAddLoader(lcfg)
	}
	return nil
}

func (this *LoaderCache) doGet(req *cacheserver.GetRequest, rep chan *cacheserver.GetResult) error {
	this.stats.Gets++
	val, ok := this.cache.Get(req.Key)

	if !ok {
		return this.doCacheMiss(req, rep)
	}
	if val == nil {
		logger.Debug(tag, "cache '%s' invalid item '%s'", this.name, req.Key)
		r := cacheserver.NewGetResult(this.name, req.Key, req.Trace)
		r.Done = false
		if req.Trace {
			r.TraceInfo = []string{"invalid holder"}
		}
		rep <- r
		return nil
	}

	this.stats.CacheHits++
	if rep != nil {
		r := cacheserver.NewGetResult(this.name, req.Key, req.Trace)
		r.Done = true
		r.Value = val
		rep <- r
	}
	return nil
}

func (this *LoaderCache) doAddLoader(cfg *LoaderConfig) error {
	p := GetLoaderProvider(cfg.Type)
	if p == nil {
		return logger.Warn(tag, "loader type '%s' not exists, can't AddLoader '%s'", cfg.Type, cfg.Name)
	}

	if this.loaders != nil {
		for _, linfo := range this.loaders {
			if linfo.config.Name == cfg.Name {
				return logger.Warn(tag, "cache '%s' loader name '%s' exists", this.name, cfg.Name)
			}
		}
	}

	logger.Info(tag, "create loader '%s:%s' - %s", this.name, cfg.Name, cfg.Type)
	loader, err := p.CreateLoader(cfg, cfg.prop)
	if err != nil {
		logger.Warn(tag, "create loader '%s:%s' fail - %s", this.name, cfg.Name, err)
		return err
	}

	linfo := new(loaderInfo)
	linfo.loader = loader
	linfo.config = *cfg

	if this.loaders == nil {
		this.loaders = []*loaderInfo{linfo}
	} else {
		this.loaders = append(this.loaders, linfo)
	}

	return nil
}

func (this *LoaderCache) doRemoveLoader(lname string) error {
	if this.loaders != nil {
		for i, l := range this.loaders {
			if l.config.Name == lname {
				logger.Info(tag, "remove loader '%s:%s'", this.name, lname)
				sz := len(this.loaders)
				this.loaders[i] = nil
				this.loaders[i], this.loaders[sz-1] = this.loaders[sz-1], this.loaders[i]
				this.loaders = this.loaders[:sz-1]
				break
			}
		}
	}
	return nil
}

func (this *LoaderCache) doCacheMiss(req *cacheserver.GetRequest, rep chan *cacheserver.GetResult) error {
	logger.Debug(tag, "'%s:%s' miss", this.name, req.Key)

	// check timeout
	var wt *waiting
	if rep != nil {
		if req.TimeoutMs <= 0 || req.NotLoad {
			r := cacheserver.NewGetResult(this.name, req.Key, req.Trace)
			r.End(false, nil, []string{"local cache miss"})
			rep <- r
		} else {
			wt = this.wait(req, rep)
			wt.result.Traces([]string{"local cache miss"})
			f := func() {
				logger.Debug(tag, "'%s:%s' waiting timeout", this.name, req.Key)
				err := errors.New("timeout")
				wt.response(false, nil, err, []string{"local timeout"})
				this.QueueRemoveWaiting(wt)
			}
			wt.timer = time.AfterFunc(time.Duration(req.TimeoutMs)*time.Millisecond, f)
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "'%s:%s' wait %.3fs", this.name, req.Key, float32(req.TimeoutMs)/1000)
			}
		}
	}

	// load it
	if req.NotLoad {
		return nil
	}

	if wt != nil && req.Trace {
		wt.result.Traces([]string{"loading"})
	}
	this.doStartLoad(req)
	return nil
}
