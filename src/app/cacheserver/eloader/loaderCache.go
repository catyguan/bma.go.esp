package eloader

import (
	"app/cacheserver"
	"bmautil/qexec"
	"bmautil/valutil"
	"bytes"
	"errors"
	"esp/shell"
	"fmt"
	"logger"
	"uprop"
)

const (
	tag        = "loaderCache"
	CACHE_TYPE = "loader"
)

func init() {
	cacheserver.RegCacheFactory(CACHE_TYPE, new(cacheFactory))
}

// Stats
type LoaderCacheStats struct {
	cacheserver.CacheStats
	PeerLoads  int64
	PeerHits   int64
	PeerErrors int64
}

func (this *LoaderCacheStats) BuildString(buf *bytes.Buffer) {

	var per float64

	buf.WriteString(fmt.Sprintf("PeerLoads=%d,", this.PeerLoads))

	per = 0
	if this.PeerLoads > 0 {
		per = float64(this.PeerHits*100) / float64(this.PeerLoads)
	}
	buf.WriteString(fmt.Sprintf("PeerHits=%d(%.2f", this.PeerHits, per))
	buf.WriteString("%),")

	per = 0
	if this.PeerLoads > 0 {
		per = float64(this.PeerErrors*100) / float64(this.PeerLoads)
	}
	buf.WriteString(fmt.Sprintf("PeerErrors=%d(%.2f", this.PeerErrors, per))
	buf.WriteString("%),")
}

type cacheConfigBase struct {
	cacheserver.LruCacheConfig
	QueueSize     int
	InvalidHolder bool  // create invalid holder
	UpdateSeconds int32 // update item which hold long, <=0 mean no updater
	UpdateStep    int32 // updater scan step, default 10
}

type cacheConfig struct {
	cacheConfigBase
	Loaders []*LoaderConfig
}

func (this *cacheConfig) Valid() error {
	s := this.LruCacheConfig.Valid()
	if s != "" {
		return errors.New(s)
	}
	if this.Loaders == nil || len(this.Loaders) == 0 {
		return errors.New("loaders empty")
	}
	names := make(map[string]bool)
	for _, lcfg := range this.Loaders {
		err := lcfg.Valid()
		if err != nil {
			return err
		}
		if names[lcfg.Name] {
			return errors.New("loader '" + lcfg.Name + "' conflict")
		}
		names[lcfg.Name] = true
	}
	if this.QueueSize <= 0 {
		this.QueueSize = 32
	}
	return nil
}

func (this *cacheConfig) GetProperties() []*uprop.UProperty {
	b := new(uprop.UPropertyBuilder)
	b.NewProp("maxsize", "cache max size").Optional(false).BeValue(this.Maxsize, func(v string) error {
		this.Maxsize = valutil.ToInt32(v, 0)
		return nil
	})
	b.NewProp("valid", "item valid second after put in cache").BeValue(this.ValidSeconds, func(v string) error {
		this.ValidSeconds = valutil.ToInt32(v, 0)
		return nil
	})
	b.NewProp("queuesize", "executor queue size").BeValue(this.QueueSize, func(v string) error {
		this.QueueSize = valutil.ToInt(v, 0)
		return nil
	})
	b.NewProp("iholder", "create invalid holder").BeValue(this.InvalidHolder, func(v string) error {
		this.InvalidHolder = valutil.ToBool(v, this.InvalidHolder)
		return nil
	})
	b.NewProp("update", "updater run duration(sec.), <=0 mean stop").BeValue(this.UpdateSeconds, func(v string) error {
		this.UpdateSeconds = valutil.ToInt32(v, 0)
		return nil
	})
	b.NewProp("step", "updater scan step, default 10").BeValue(this.UpdateStep, func(v string) error {
		this.UpdateStep = valutil.ToInt32(v, 0)
		return nil
	})
	loaderList := b.NewProp("loader", "loader").Optional(false).BeList(this.addLoader, this.removeLoader)
	if this.Loaders != nil {
		for _, lcfg := range this.Loaders {
			loaderList.AddFold(lcfg.Desc(), lcfg.GetProperties)
		}
	}
	return b.AsList()
}

func (this *cacheConfig) ToMap() map[string]interface{} {
	r := valutil.BeanToMap(this.cacheConfigBase)
	la := make([]interface{}, 0)
	if this.Loaders != nil {
		for _, lcfg := range this.Loaders {
			lm := make(map[string]interface{})
			lm["name"] = lcfg.Name
			lm["type"] = lcfg.Type
			lm["prop"] = lcfg.prop.ToMap()
			la = append(la, lm)
		}
	}
	r["loaders"] = la
	return r
}

func (this *cacheConfig) FromMap(data map[string]interface{}) error {
	valutil.ToBean(data, &this.cacheConfigBase)
	lcfg := valutil.ToArray(data["loaders"])
	if lcfg != nil {
		this.Loaders = make([]*LoaderConfig, 0)
		for _, lv := range lcfg {
			m := valutil.ToStringMap(lv)
			if m != nil {
				l := new(LoaderConfig)
				l.Name = valutil.ToString(m["name"], "")
				l.Type = valutil.ToString(m["type"], "")
				p := GetLoaderProvider(l.Type)
				if p != nil {
					l.prop = p.CreateProperty()
					l.prop.FromMap(valutil.ToStringMap(m["prop"]))
				}
				this.Loaders = append(this.Loaders, l)
			}
		}
	}
	return this.Valid()
}

func (this *cacheConfig) addLoader(slist []string) error {
	if this.Loaders == nil {
		this.Loaders = make([]*LoaderConfig, 0)
	}
	if slist != nil && len(slist) > 0 {
		for _, n := range slist {
			lcfg := new(LoaderConfig)
			lcfg.Name = n
			this.Loaders = append(this.Loaders, lcfg)
		}
	} else {
		lcfg := new(LoaderConfig)
		this.Loaders = append(this.Loaders, lcfg)
	}
	return nil
}

func (this *cacheConfig) removeLoader(slist []string) error {
	if this.Loaders == nil {
		return nil
	}
	for _, n := range slist {
		c := len(this.Loaders)
		idx := uprop.ToIndex(n, c)
		if idx == 0 {
			for i, lcfg := range this.Loaders {
				if lcfg.Name == n {
					idx = i
					break
				}
			}
		}
		idx = idx - 1
		if idx < 0 || idx >= c {
			return fmt.Errorf("'%s' invalid", n)
		}
		l := make([]*LoaderConfig, c-1)
		copy(l[0:], this.Loaders[0:idx])
		copy(l[idx:], this.Loaders[idx+1:])
	}
	return nil
}

// Factory
type cacheFactory struct {
}

func (this *cacheFactory) CreateConfig() cacheserver.ICacheConfig {
	return new(cacheConfig)
}

func (this *cacheFactory) CreateCache(cfg cacheserver.ICacheConfig) (cacheserver.ICache, error) {
	r := NewLoaderCache()
	r.config = cfg.(*cacheConfig)

	for _, lcfg := range r.config.Loaders {
		err := r.doAddLoader(lcfg)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

// Cache
type loaderInfo struct {
	loader Loader
	config LoaderConfig
}

type LoaderCache struct {
	name    string
	service *cacheserver.CacheService

	executor *qexec.QueueExecutor
	cache    *cacheserver.Cache
	config   *cacheConfig
	stats    LoaderCacheStats

	// loader
	loaders []*loaderInfo

	// waiting
	keyspot map[string]*keySpot

	// updater
	updater *cacheUpdater
}

func NewLoaderCache() *LoaderCache {
	this := new(LoaderCache)
	return this
}

func (this *LoaderCache) Type() string {
	return CACHE_TYPE
}

func (this *LoaderCache) GetConfig() cacheserver.ICacheConfig {
	return this.config
}

func (this *LoaderCache) InitCache(s *cacheserver.CacheService, n string) {
	this.service = s
	this.name = n
}

func (this *LoaderCache) Get(req *cacheserver.GetRequest, rep chan *cacheserver.GetResult) error {
	exec := this.executor
	if exec == nil {
		return errors.New("not start")
	}
	return exec.DoSync("Get", func() error {
		return this.doGet(req, rep)
	})
}

func (this *LoaderCache) Put(key string, val []byte, deadUnixtime int64) error {
	return errors.New("not implements")
}

func (this *LoaderCache) Delete(key string) (bool, error) {
	err := this.TryLoad(key)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *LoaderCache) TryLoad(key string) error {
	exec := this.executor
	if exec == nil {
		return errors.New("not start")
	}
	logger.Debug(tag, "TryLoad(%s)", key)
	call := func() error {
		req := cacheserver.NewGetRequest(key)
		req.Update = true
		this.doStartLoad(req)
		return nil
	}
	return this.executor.DoSync("TryLoad", call)
}

func (this *LoaderCache) QueueRemoveWaiting(wt *waiting) error {
	exec := this.executor
	if exec == nil {
		return errors.New("not start")
	}
	logger.Debug(tag, "queueRemoveWaiting(%s:%s)", this.name, wt.req.Key)
	call := func() error {
		this.removeWait(wt)
		return nil
	}
	go exec.DoNow("queueRemoveWaiting", call)
	return nil
}

func (this *LoaderCache) LoadEnd(loaderName string, done bool, key string, val []byte, err error, traces []string) {
	exec := this.executor
	if exec == nil {
		logger.Warn(tag, "%s not start when loadEnd", this.name)
		return
	}

	logger.Debug(tag, "LoadEnd(%v, %v, %v, %v, %v, %v)", this.name, loaderName, done, key, err, traces)
	call := func() error {
		this.doLoadEnd(loaderName, key, done, val, err, traces)
		return nil
	}
	exec.DoNow("LoadEnd", call)
}

func (this *LoaderCache) StepUpdate(updater *cacheUpdater) {
	exec := this.executor
	if exec == nil {
		logger.Warn(tag, "%s not start when StepUpdate", this.name)
		return
	}

	logger.Debug(tag, "stepUpdate(%s)", this.name)
	call := func() error {
		updater.doUpdate(this)
		return nil
	}
	go exec.DoNow("stepUpdate", call)
}

func (this *LoaderCache) UpdateConfig(newcfg cacheserver.ICacheConfig) error {
	cfg := newcfg.(*cacheConfig)
	exec := this.executor
	if exec == nil {
		*this.config = *cfg
		return nil
	}
	return exec.DoSync("Deploy", func() error {
		return this.doDeployCache(cfg)
	})
}

func (this *LoaderCache) Load(key string) error {
	return this.TryLoad(key)
}

func (this *LoaderCache) QueryStats() (string, error) {
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

func (this *LoaderCache) doQueryStats() (string, error) {
	var st cacheserver.LruCacheStats
	st.CopyLruCacheState(&this.stats.CacheStats, this.cache)
	buf := bytes.NewBuffer(make([]byte, 0))
	st.BuildString(buf)
	buf.WriteByte(',')
	this.stats.BuildString(buf)
	return buf.String(), nil
}

func (this *LoaderCache) IsStart() bool {
	return this.executor != nil
}

func (this *LoaderCache) Start() error {
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

func (this *LoaderCache) Run() error {
	if this.updater == nil {
		this.updater = newCacheUpdater()
	}
	if this.config.UpdateSeconds > 0 {
		this.updater.start(this)
	}
	return nil
}

func (this *LoaderCache) Stop() error {
	e := this.executor
	if e != nil {
		e.Stop()
	}
	return nil
}

func (this *LoaderCache) CreateShell() shell.ShellDir {
	r := shell.NewShellDirCommon(this.name)
	this.service.BuildCacheCommands(this.name, r)
	return r
}
