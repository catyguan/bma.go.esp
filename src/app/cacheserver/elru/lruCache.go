package elru

import (
	"app/cacheserver"
	"bmautil/qexec"
	"bmautil/valutil"
	"errors"
	"esp/shell"
	"logger"
)

const (
	tag = "elru"
)

func init() {
	cacheserver.RegCacheFactory("lru", NewLruCache)
}

type cacheConfig struct {
	QueueSize int
	cacheserver.LruCacheConfig
}

type LruCache struct {
	name    string
	service *cacheserver.CacheService

	executor   *qexec.QueueExecutor
	cache      *cacheserver.Cache
	config     *cacheConfig
	editConfig *cacheConfig
	stats      cacheserver.CacheStats
}

func NewLruCache() cacheserver.ICache {
	this := new(LruCache)
	return this
}

func (this *LruCache) InitCache(s *cacheserver.CacheService, n string) {
	this.service = s
	this.name = n
}

func (this *LruCache) Get(req *cacheserver.GetRequest, rep chan *cacheserver.GetResult) error {
	exec := this.executor
	if exec == nil {
		return errors.New("not start")
	}
	return exec.DoSync("Get", func() error {
		return this.doGet(req, rep)
	})
}

func (this *LruCache) doGet(req *cacheserver.GetRequest, rep chan *cacheserver.GetResult) error {
	if this.cache == nil {
		return errors.New("cache nil")
	}

	this.stats.Gets++
	val, ok := this.cache.Get(req.Key)

	if !ok {
		if rep != nil {
			r := cacheserver.NewGetResult(this.name, req.Key, req.Trace)
			r.End(false, nil, []string{"miss"})
			rep <- r
		}
		return nil
	}
	if val == nil {
		logger.Debug(tag, "cache '%s' invalid item '%s'", this.name, req.Key)
		if rep != nil {
			r := cacheserver.NewGetResult(this.name, req.Key, req.Trace)
			r.End(false, nil, []string{"invalid holder"})
			rep <- r
		}
		return nil
	}

	this.stats.CacheHits++
	if rep != nil {
		r := cacheserver.NewGetResult(this.name, req.Key, req.Trace)
		r.End(true, val, nil)
		rep <- r
	}
	return nil
}

func (this *LruCache) Put(key string, val []byte, deadUnixtime int64) error {
	exec := this.executor
	if exec == nil {
		return errors.New("not start")
	}
	return exec.DoSync("Put", func() error {
		return this.doPut(key, val, deadUnixtime)
	})
}

func (this *LruCache) doPut(key string, val []byte, deadUnixtime int64) error {
	if this.cache == nil {
		return errors.New("cache nil")
	}
	dt := deadUnixtime
	if dt == 0 {
		dt = -1
	}
	this.cache.Put(key, val, dt)
	return nil
}

func (this *LruCache) Delete(key string) (bool, error) {
	exec := this.executor
	if exec == nil {
		return false, errors.New("not start")
	}
	ok := false
	err := exec.DoSync("Delete", func() error {
		r, err := this.doDelete(key)
		ok = r
		return err
	})
	return ok, err
}

func (this *LruCache) doDelete(key string) (bool, error) {
	if this.cache == nil {
		return false, errors.New("cache nil")
	}
	_, ok := this.cache.Remove(key)
	return ok, nil
}

func (this *LruCache) Load(key string) error {
	return nil
}

func (this *LruCache) QueryStats() (string, error) {
	exec := this.executor
	if exec == nil {
		return "", errors.New("not start")
	}
	r := ""
	err := exec.DoSync("QueryStats", func() error {
		var err error
		r, err = this.doQueryStats()
		return err
	})
	return r, err
}

func (this *LruCache) doQueryStats() (string, error) {
	var st cacheserver.LruCacheStats
	st.CopyLruCacheState(&this.stats, this.cache)
	return st.String(), nil
}

func logEvit(cache *cacheserver.Cache, name string) {
	if logger.EnableDebug(tag) {
		cache.Listener = func(key string, val []byte) {
			logger.Debug(tag, "evit '%s' %s", name, key)
		}
	}
}

func (this *LruCache) doCreateCache() error {
	this.cache = this.config.NewCache()
	logEvit(this.cache, this.name)
	return nil
}

func (this *LruCache) Deploy() error {
	exec := this.executor
	if exec == nil {
		return errors.New("not start")
	}
	return exec.DoSync("Deploy", func() error {
		return this.doDeployCache()
	})
}

func (this *LruCache) doDeployCache() error {
	logger.Info(tag, "deploy cache '%s'", this.name)
	cfg := this.config
	if this.cache.MaxSize() != cfg.Maxsize {
		// clone cache first
		old := this.cache
		logger.Info(tag, "resize cache %d -> %d", old.MaxSize(), cfg.Maxsize)
		ncache := cfg.NewCache()
		ncache.Clone(old)
		logEvit(ncache, this.name)

		this.cache = ncache

		logger.Info(tag, "resize cache %d -> %d done", old.MaxSize(), ncache.MaxSize())

		old.Clear()
		old = nil
	} else {
		if cfg.ValidSeconds > 0 {
			this.cache.ValidTime = int64(cfg.ValidSeconds)
		} else {
			this.cache.ValidTime = 0
		}
	}
	return nil
}

func (this *LruCache) IsStart() bool {
	return this.executor != nil
}

func (this *LruCache) Start() error {
	if this.executor != nil {
		return errors.New("started")
	}
	if this.config == nil {
		return errors.New("not config")
	}
	qs := this.config.QueueSize
	if qs <= 0 {
		qs = 16
	}
	e := qexec.NewQueueExecutor(tag, qs, this.requestHandler)
	e.StopHandler = this.stopHandler
	e.Run()
	e.DoNow("init", func() error {
		this.doCreateCache()
		return nil
	})
	this.executor = e
	return nil
}

func (this *LruCache) requestHandler(ev interface{}) (bool, error) {
	switch rv := ev.(type) {
	case func() error:
		return true, rv()
	}
	return true, nil
}

func (this *LruCache) stopHandler() {
	if this.cache != nil {
		this.cache.Clear()
	}
	this.executor = nil
}

func (this *LruCache) Run() error {
	return nil
}

func (this *LruCache) Stop() error {
	e := this.executor
	if e != nil {
		e.Stop()
	}
	return nil
}

func (this *LruCache) FromConfig(cfg map[string]interface{}) error {
	if cfg != nil {
		cobj := new(cacheConfig)
		if !valutil.ToBean(cfg, cobj) {
			return errors.New("invalid LruCacheConfig")
		}
		if err := cobj.Valid(); err != "" {
			return errors.New(err)
		}

		this.config = cobj
		ncobj := new(cacheConfig)
		*ncobj = *cobj
		this.editConfig = ncobj
	}
	return nil
}

func (this *LruCache) ToConfig() map[string]interface{} {
	if this.config != nil {
		r := valutil.BeanToMap(this.config)
		logger.Info(tag, "%v", r)
		return r
	}
	return make(map[string]interface{})
}

func (this *LruCache) CreateShell() shell.ShellDir {
	r := shell.NewShellDirCommon(this.name)
	this.service.BuildCacheCommands(this.name, r)
	r.AddCommand(&cmdEdit{this})
	r.AddCommand(&cmdDeploy{this})
	return r
}
