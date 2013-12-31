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
)

const (
	tag = "loaderCache"
)

func init() {
	cacheserver.RegCacheFactory("loader", NewLoaderCache)
}

type cacheConfig struct {
	cacheserver.LruCacheConfig
	QueueSize     int
	InvalidHolder bool  // create invalid holder
	UpdateSeconds int32 // update item which hold long, <=0 mean no updater
	UpdateStep    int32 // updater scan step, default 10
}

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

type loaderInfo struct {
	loader Loader
	config LoaderConfig
}

type LoaderCache struct {
	name    string
	service *cacheserver.CacheService

	executor   *qexec.QueueExecutor
	cache      *cacheserver.Cache
	config     *cacheConfig
	editConfig *cacheConfig
	stats      LoaderCacheStats

	// loader
	loaders []*loaderInfo

	// waiting
	keyspot map[string]*keySpot

	// updater
	updater *cacheUpdater
}

func NewLoaderCache() cacheserver.ICache {
	this := new(LoaderCache)
	return this
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

func (this *LoaderCache) AddLoader(cfg *LoaderConfig) error {
	exec := this.executor
	if exec == nil {
		return errors.New("not start")
	}
	logger.Debug(tag, "AddLoader(%s, %s, %s)", this.name, cfg.Name, cfg.Type)
	call := func() error {
		return this.doAddLoader(cfg)
	}
	return exec.DoSync("AddLoader", call)
}

func (this *LoaderCache) RemoveLoader(name string) error {
	exec := this.executor
	if exec == nil {
		return errors.New("not start")
	}
	logger.Debug(tag, "RemoveLoader(%s)", name)
	call := func() error {
		return this.doRemoveLoader(name)
	}
	return exec.DoSync("RemoveLoader", call)
}

func (this *LoaderCache) Deploy() error {
	exec := this.executor
	if exec == nil {
		return errors.New("not start")
	}
	return exec.DoSync("Deploy", func() error {
		return this.doDeployCache()
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

func (this *LoaderCache) FromConfig(cfg map[string]interface{}) error {
	if cfg != nil {
		ccfg := valutil.ToStringMap(cfg["cache"])
		if ccfg != nil {
			cobj := new(cacheConfig)
			valutil.ToBean(ccfg, cobj)
			err := cobj.Valid()
			if err != "" {
				return errors.New(err)
			}
			this.config = cobj
			ncobj := new(cacheConfig)
			*ncobj = *cobj
			this.editConfig = ncobj
		}
		lcfg := valutil.ToArray(cfg["loader"])
		if lcfg != nil {
			for _, lv := range lcfg {
				m := valutil.ToStringMap(lv)
				if m != nil {
					l := new(LoaderConfig)
					l.Name = valutil.ToString(m["name"], "")
					if l.Name == "" {
						logger.Warn(tag, "'%s' loader name empty", this.name)
						continue
					}
					l.Type = valutil.ToString(m["type"], "")
					p := GetLoaderProvider(l.Type)
					if p == nil {
						logger.Warn(tag, "'%s' loader[%s] type '%s' invalid", this.name, l.Name, l.Type)
						continue
					}
					l.prop = p.CreateProperty()
					err := l.prop.FromMap(valutil.ToStringMap(m["prop"]))
					if err != nil {
						logger.Warn(tag, "'%s' loader[%s] prop fail - %s", this.name, l.Name, err)
						continue
					}
					err2 := this.doAddLoader(l)
					if err2 != nil {
						logger.Warn(tag, "'%s' loader[%s] addLoader fail - %s", this.name, l.Name, err)
						continue
					}
				}
			}
		}
	}
	return nil
}

func (this *LoaderCache) ToConfig() map[string]interface{} {
	r := make(map[string]interface{})
	if true {
		ccfg := valutil.BeanToMap(this.config)
		r["cache"] = ccfg
	}
	if this.loaders != nil && len(this.loaders) > 0 {
		lcfg := make([]map[string]interface{}, 0)
		for _, l := range this.loaders {
			m := make(map[string]interface{})
			m["name"] = l.config.Name
			m["type"] = l.config.Type
			if l.config.prop != nil {
				m["prop"] = l.config.prop.ToMap()
			}
			lcfg = append(lcfg, m)
		}
		r["loader"] = lcfg
	}
	return r
}

func (this *LoaderCache) CreateShell() shell.ShellDir {
	r := shell.NewShellDirCommon(this.name)
	this.service.BuildCacheCommands(this.name, r)
	r.AddCommand(&cmdEdit{this})
	r.AddCommand(&cmdDeploy{this})
	r.AddCommand(&cmdLoaders{this})
	r.AddCommand(&cmdUnload{this})
	r.AddCommand(&cmdNewLoader{this})
	return r
}
