package estable

import (
	"app/cacheserver"
	"bmautil/qexec"
	"bmautil/valutil"
	"errors"
	"esp/shell"
	"fmt"
	"logger"
	"time"
	"uprop"
)

const (
	tag        = "stableCache"
	CACHE_TYPE = "stable"
)

func init() {
	cacheserver.RegCacheFactory(CACHE_TYPE, new(stableCacheFactory))
}

// Config
type CacheConfig struct {
	MaxSize     int // Cache Max Size
	QueueSize   int // Executor Size
	RecoverTime int // 恢复无效响应的时间，秒
}

func (this *CacheConfig) Valid() error {
	if this.MaxSize <= 0 {
		this.MaxSize = 100
	}
	if this.QueueSize <= 0 {
		this.QueueSize = 32
	}
	if this.RecoverTime <= 0 {
		this.RecoverTime = 5
	}
	return nil
}

func (this *CacheConfig) GetProperties() []*uprop.UProperty {
	r := make([]*uprop.UProperty, 0)
	r = append(r, uprop.NewUProperty("maxsize", this.MaxSize, true, "cache max size", func(v string) error {
		this.MaxSize = valutil.ToInt(v, 0)
		return nil
	}))
	r = append(r, uprop.NewUProperty("queuesize", this.QueueSize, true, "executor queue size", func(v string) error {
		this.QueueSize = valutil.ToInt(v, 0)
		return nil
	}))
	r = append(r, uprop.NewUProperty("recover", this.RecoverTime, true, "recover time when invalidResponse sended, SEC", func(v string) error {
		this.RecoverTime = valutil.ToInt(v, 0)
		return nil
	}))
	return r
}

func (this *CacheConfig) ToMap() map[string]interface{} {
	return valutil.BeanToMap(this)
}

func (this *CacheConfig) FromMap(data map[string]interface{}) error {
	valutil.ToBean(data, this)
	return this.Valid()
}

// Factory
type stableCacheFactory struct {
}

func (this *stableCacheFactory) CreateConfig() cacheserver.ICacheConfig {
	return new(CacheConfig)
}

func (this *stableCacheFactory) CreateCache(cfg cacheserver.ICacheConfig) (cacheserver.ICache, error) {
	r := NewStableCache()
	r.config = cfg.(*CacheConfig)
	return r, nil
}

// Cache
type CacheItem struct {
	data     []byte
	deadTime int64

	// runtime
	valid               bool      // 是否有效
	invalidResponseTime time.Time // 无效响应的时间，秒
}

type StableCache struct {
	name    string
	service *cacheserver.CacheService

	config   *CacheConfig
	items    map[string]*CacheItem
	executor *qexec.QueueExecutor
}

func NewStableCache() *StableCache {
	this := new(StableCache)
	this.items = make(map[string]*CacheItem)
	return this
}

func (this *StableCache) Type() string {
	return CACHE_TYPE
}

func (this *StableCache) GetConfig() cacheserver.ICacheConfig {
	return this.config
}

func (this *StableCache) UpdateConfig(cfg cacheserver.ICacheConfig) error {
	if err := cfg.Valid(); err != nil {
		return err
	}
	c := cfg.(*CacheConfig)
	*this.config = *c
	return nil
}

func (this *StableCache) requestHandler(ev interface{}) (bool, error) {
	if f, ok := ev.(func() error); ok {
		return true, f()
	}
	return true, nil
}

func (this *StableCache) stopHandler() {
	for k, _ := range this.items {
		delete(this.items, k)
	}
}

func (this *StableCache) InitCache(s *cacheserver.CacheService, n string) {
	this.service = s
	this.name = n
}

func (this *StableCache) doGet(req *cacheserver.GetRequest, rep chan *cacheserver.GetResult) error {
	ci, ok := this.items[req.Key]
	if !ok {
		if rep != nil {
			res := cacheserver.NewGetResult(this.name, req.Key, req.Trace)
			res.End(false, nil, nil)
			rep <- res
		}
		return nil
	}

	respData := true
	var traces []string
	if req.Trace {
		traces = make([]string, 0)
	}
	if ci.valid {
		if ci.deadTime != 0 && ci.deadTime <= time.Now().Unix() {
			ci.valid = false
			logger.Debug(tag, "Cache[%s] key[%s] timeout", this.name, req.Key)
			if req.Trace {
				traces = append(traces, "timeout")
			}
		}
	}
	if !ci.valid {
		// check has reponse?
		if ci.invalidResponseTime.IsZero() {
			// never invalidResponse, do invalidResponse
			respData = false
			ci.invalidResponseTime = time.Now()
			logger.Debug(tag, "Cache[%s] key[%s] do invalid response", this.name, req.Key)
			if req.Trace {
				traces = append(traces, "do invalid response")
			}
		} else {
			// response timeout?
			if this.config.RecoverTime > 0 {
				now := time.Now()
				tm := ci.invalidResponseTime.Add(time.Duration(this.config.RecoverTime) * time.Second)
				if now.After(tm) {
					respData = false
					ci.invalidResponseTime = now
					logger.Debug(tag, "Cache[%s] key[%s] recover invalid response", this.name, req.Key)
					if req.Trace {
						traces = append(traces, "recover invalid response")
					}
				} else {
					// logger.Debug(tag, "fuck2")
				}
			} else {
				// logger.Debug(tag, "fuck1")
			}
		}
		if respData {
			logger.Debug(tag, "Cache[%s] key[%s] stable response", this.name, req.Key)
			if req.Trace {
				traces = append(traces, "stable response")
			}
		}
	}

	if rep != nil {
		res := cacheserver.NewGetResult(this.name, req.Key, req.Trace)
		if respData {
			res.End(ok, ci.data, traces)
			rep <- res
		} else {
			res.End(false, nil, traces)
			rep <- res
		}
	}
	return nil
}

func (this *StableCache) Get(req *cacheserver.GetRequest, rep chan *cacheserver.GetResult) error {
	if this.executor == nil {
		return errors.New("not start")
	}
	return this.executor.DoSync("Get", func() error {
		return this.doGet(req, rep)
	})
}

func (this *StableCache) doPut(key string, val []byte, deadUnixtime int64) error {
	var ci *CacheItem
	ok := false
	if ci, ok = this.items[key]; !ok {
		if len(this.items) >= this.config.MaxSize {
			return errors.New("cache full")
		}
		ci = new(CacheItem)
		this.items[key] = ci
	}
	ci.data = val
	ci.deadTime = deadUnixtime
	ci.valid = true
	ci.invalidResponseTime = time.Time{}
	return nil
}
func (this *StableCache) Put(key string, val []byte, deadUnixtime int64) error {
	if this.executor == nil {
		return errors.New("not start")
	}
	return this.executor.DoSync("Put", func() error {
		return this.doPut(key, val, deadUnixtime)
	})
}

func (this *StableCache) doDelete(key string) (bool, error) {
	ci, ok := this.items[key]
	if ok {
		ci.valid = false
	}
	return ok, nil
}
func (this *StableCache) Delete(key string) (bool, error) {
	if this.executor == nil {
		return false, errors.New("not start")
	}
	r := false
	err := this.executor.DoSync("Delete", func() error {
		var err error
		r, err = this.doDelete(key)
		return err
	})
	return r, err
}

func (this *StableCache) Load(key string) error {
	return nil
}

func (this *StableCache) QueryStats() (string, error) {
	if this.executor == nil {
		return "", errors.New("not start")
	}
	r := ""
	err := this.executor.DoSync("QueryStats", func() error {
		var err error
		r, err = this.doQueryStats()
		return err
	})
	return r, err
}

func (this *StableCache) doQueryStats() (string, error) {
	return fmt.Sprintf("size=%d", len(this.items)), nil
}

func (this *StableCache) IsStart() bool {
	return this.executor != nil
}

func (this *StableCache) Start() error {
	if this.config == nil {
		return errors.New("not config")
	}
	this.executor = qexec.NewQueueExecutor(this.name, this.config.QueueSize, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.executor.Run()
	return nil
}

func (this *StableCache) Run() error {
	return nil
}

func (this *StableCache) Stop() error {
	if this.executor != nil {
		this.executor.Stop()
		this.executor.WaitStop()
		this.executor = nil
	}
	return nil
}

func (this *StableCache) CreateShell() shell.ShellDir {
	r := shell.NewShellDirCommon(this.name)
	this.service.BuildCacheCommands(this.name, r)
	return r
}
