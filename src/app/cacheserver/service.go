package cacheserver

import (
	"bmautil/sqlutil"
	"boot"
	"config"
	"database/sql"
	"encoding/json"
	"errors"
	"esp/sqlite"
	"logger"
	"sync"
)

const (
	tag      = "CacheService"
	database = "local"
)

type cacheInfo struct {
	typeName string
	cache    ICache
}

var (
	factorties map[string]CacheFactory = make(map[string]CacheFactory)
)

func RegCacheFactory(typeName string, fac CacheFactory) {
	factorties[typeName] = fac
}

type CacheService struct {
	name     string
	database *sqlite.SqliteServer
	caches   map[string]*cacheInfo
	wcaches  map[string]*cacheInfo
	safeMode bool

	lock     sync.Mutex
	stopWait *sync.WaitGroup
}

func NewCacheService(name string, db *sqlite.SqliteServer) *CacheService {
	this := new(CacheService)
	this.name = name
	this.database = db
	this.caches = make(map[string]*cacheInfo)
	this.wcaches = make(map[string]*cacheInfo)

	this.initDatabase()

	return this
}

func (this *CacheService) Name() string {
	return this.name
}

type configInfo struct {
	SafeMode bool
}

func (this *CacheService) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		this.safeMode = cfg.SafeMode
		logger.Debug(tag, "FUCK %v", this.safeMode)
	} else {
		logger.Debug(tag, "FUCK2 %v", this.name)
	}
	return true
}

func (this *CacheService) Start() bool {
	cfg, ok := this.loadRuntimeConfig()
	if !ok {
		if !this.safeMode {
			return false
		}
	}
	if !this.setupByConfig(cfg) {
		return false
	}

	return this.startAllCache()
}

func (this *CacheService) Run() bool {
	return this.runAllCache()
}

func (this *CacheService) Stop() bool {
	this.stopAllCache()
	return true
}

func (this *CacheService) Cleanup() bool {
	this.waitAllCacheStop()
	return true
}

func (this *CacheService) DefaultBoot() {
	boot.Define(boot.INIT, this.name, this.Init)
	boot.Define(boot.START, this.name, this.Start)
	boot.Define(boot.RUN, this.name, this.Run)
	boot.Define(boot.STOP, this.name, this.Stop)
	boot.Define(boot.CLEANUP, this.name, this.Cleanup)

	boot.Install(this.name, this)
}

func (this *CacheService) cache(name string, isStart bool) (ICache, error) {
	m := this.caches
	cinfo, ok := m[name]
	if !ok {
		return nil, logger.Warn(tag, "cache[%s] not exists", name)
	}
	cache := cinfo.cache
	if isStart && !cache.IsStart() {
		return nil, errors.New(logger.Sprintf("cache[%s] not start", name))
	}
	return cinfo.cache, nil
}

func (this *CacheService) GetCache(name string, Started bool) (ICache, error) {
	return this.cache(name, Started)
}

func (this *CacheService) Get(name string, req *GetRequest, rep chan *GetResult) error {
	logger.Debug(tag, "Get(%s, %v)", name, req)
	cache, err := this.cache(name, true)
	if err != nil {
		return err
	}
	err = cache.Get(req, rep)
	if err != nil && rep != nil {
		r := new(GetResult)
		r.Err = err
		rep <- r
	}
	return err
}

func (this *CacheService) Put(name string, key string, val []byte, deadUnixtime int64) error {
	logger.Debug(tag, "Put(%s, %s)", name, key)

	cache, err := this.cache(name, true)
	if err != nil {
		return err
	}
	return cache.Put(key, val, deadUnixtime)
}

func (this *CacheService) Delete(name string, key string) (bool, error) {
	logger.Debug(tag, "Delete(%s, %s)", name, key)

	cache, err := this.cache(name, true)
	if err != nil {
		return false, err
	}
	return cache.Delete(key)
}

func (this *CacheService) Load(name string, key string) error {
	logger.Debug(tag, "Load(%s, %s)", name, key)

	cache, err := this.cache(name, true)
	if err != nil {
		return err
	}
	return cache.Load(key)
}

func (this *CacheService) QueryStats(name string) (string, error) {
	logger.Debug(tag, "QueryStats(%s)", name)

	cache, err := this.cache(name, true)
	if err != nil {
		return "", err
	}
	return cache.QueryStats()
}

func (this *CacheService) StartCache(name string) error {
	logger.Debug(tag, "StartCache(%s)", name)

	cache, err := this.cache(name, false)
	if err != nil {
		return err
	}
	if !cache.IsStart() {
		err = cache.Start()
		if err != nil {
			return err
		}
	}
	err = cache.Run()
	if err != nil {
		return err
	}
	return nil
}

func (this *CacheService) StopCache(name string) error {
	logger.Debug(tag, "StopCache(%s)", name)

	cache, err := this.cache(name, false)
	if err != nil {
		return err
	}
	if cache.IsStart() {
		err = cache.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

// LOCKED
func (this *CacheService) copyOnWrite() {
	m := make(map[string]*cacheInfo)
	for k, ci := range this.wcaches {
		m[k] = ci
	}
	this.caches = m
}

func (this *CacheService) CreateCache(name, typeName string) (ICache, error) {
	logger.Debug(tag, "CreateCache(%s, %s)", name, typeName)

	fac, ok := factorties[typeName]
	if !ok {
		return nil, errors.New(logger.Sprintf("CacheType[%s] not exists", typeName))
	}
	cache := fac()
	cache.InitCache(this, name)

	done := func() bool {
		this.lock.Lock()
		defer this.lock.Unlock()
		_, ok := this.wcaches[name]
		if ok {
			return false
		}
		ci := new(cacheInfo)
		ci.cache = cache
		ci.typeName = typeName
		this.wcaches[name] = ci
		this.copyOnWrite()
		return true
	}()

	if !done {
		return nil, errors.New(logger.Sprintf("Cache[%s] exists", name))
	}
	return cache, nil
}

func (this *CacheService) DeleteCache(name string, stop bool) error {
	logger.Debug(tag, "DeleteCache(%s, %v)", name, stop)

	cache := func() ICache {
		this.lock.Lock()
		defer this.lock.Unlock()
		ci, ok := this.wcaches[name]
		if !ok {
			return nil
		}
		delete(this.wcaches, name)
		this.copyOnWrite()
		return ci.cache
	}()
	if cache != nil && stop {
		err := cache.Stop()
		if err != nil {
			logger.Debug(tag, "Cache[%s] stop fail - %s", name, err)
		}
		return err
	}
	return nil
}

func (this *CacheService) ListCacheName() []string {
	logger.Debug(tag, "ListCacheName()")

	m := this.caches
	r := make([]string, 0, len(m))
	for k, _ := range m {
		r = append(r, k)
	}
	return r
}

// impl
type runtimeConfig struct {
	Caches map[string]*cacheRuntime
}

type cacheRuntime struct {
	TypeName string
	Config   map[string]interface{}
}

func (this *CacheService) initDatabase() {
	sqlstr := []string{
		"create table tbl_cache_service (id integer not null primary key, content text)",
		"insert into tbl_cache_service values (1, '')",
	}
	this.database.AddInit(sqlite.InitTable(database, "tbl_cache_service", sqlstr))
}

func (this *CacheService) loadRuntimeConfig() (*runtimeConfig, bool) {
	content := ""
	rowScan := func(rows *sql.Rows) error {
		if rows.Next() {
			return rows.Scan(&content)
		}
		return nil
	}
	sqlstr := "SELECT content FROM tbl_cache_service WHERE id = ?"
	action := sqlutil.QueryAction(rowScan, sqlstr, 1)
	event := make(chan error)
	defer close(event)
	this.database.Do(database, action, event)
	if err := <-event; err != nil {
		logger.Error(tag, "load local data fail %s", err)
		return nil, false
	}
	logger.Debug(tag, "load runtime config = %s", content)
	var cfg runtimeConfig
	if content != "" {
		if err := json.Unmarshal([]byte(content), &cfg); err != nil {
			logger.Error(tag, "runtime config parse error => %s", err)
			return nil, false
		}
	}
	return &cfg, true
}

func (this *CacheService) setupByConfig(cfg *runtimeConfig) bool {
	if cfg.Caches != nil {
		for cname, cobj := range cfg.Caches {
			cache, err := this.CreateCache(cname, cobj.TypeName)
			if err != nil {
				logger.Error(tag, "Cache[%s] setup fail %s", cname, err)
				if this.safeMode {
					continue
				}
				return false
			}
			err = cache.FromConfig(cobj.Config)
			if err != nil {
				logger.Error(tag, "Cache[%s] load runtime config fail - %s", cname, err)
				if this.safeMode {
					continue
				}
				return false
			}
		}
	}
	return true
}

func (this *CacheService) storeRuntimeConfig(cfg *runtimeConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		logger.Error("ERROR: runtime config format error => %s", err)
		return err
	}
	f := func(r int64) {
		logger.Debug(tag, "store runtime config => %d", r)
	}
	content := string(data)
	logger.Debug(tag, "store runtime config = %s", content)
	sqlstr := "UPDATE tbl_cache_service SET content = ? WHERE id = ?"
	action := sqlutil.ExecuteAction(f, sqlstr, content, 1)
	logger.Info(tag, "do store runtime config")
	return this.database.Do(database, action, nil)
}

func (this *CacheService) buildRuntimeConfig() *runtimeConfig {
	r := new(runtimeConfig)
	this.lock.Lock()
	defer this.lock.Unlock()
	r.Caches = make(map[string]*cacheRuntime)
	for cname, cobj := range this.caches {
		cr := new(cacheRuntime)
		cr.TypeName = cobj.typeName
		cr.Config = cobj.cache.ToConfig()
		// logger.Info(tag, "%s config = %v", cobj.typeName, cr.Config)
		r.Caches[cname] = cr
	}
	return r
}

func (this *CacheService) save() error {
	cfg := this.buildRuntimeConfig()
	return this.storeRuntimeConfig(cfg)
}

func (this *CacheService) startAllCache() bool {
	m := this.caches
	for k, ci := range m {
		if ci.cache.IsStart() {
			logger.Debug(tag, "boot start Cache[%s] - skip", k)
			continue
		}
		logger.Debug(tag, "boot start Cache[%s]", k)
		err := ci.cache.Start()
		if err != nil {
			if this.safeMode {
				logger.Warn(tag, "start cache[%s] fail %s", k, err)
				continue
			}
			logger.Error(tag, "boot start Cache[%s] fail - %s", k, err)
			return false
		}
	}
	return true
}

func (this *CacheService) runAllCache() bool {
	m := this.caches
	for k, ci := range m {
		if !ci.cache.IsStart() {
			logger.Debug(tag, "boot run Cache[%s] - not start", k)
			continue
		}
		logger.Debug(tag, "boot run Cache[%s]", k)
		err := ci.cache.Run()
		if err != nil {
			if this.safeMode {
				logger.Warn(tag, "run cache[%s] fail %s", k, err)
				continue
			}
			logger.Error(tag, "boot run Cache[%s] fail - %s", k, err)
			return false
		}
	}
	return true
}

func (this *CacheService) stopAllCache() {
	this.stopWait = new(sync.WaitGroup)
	m := this.caches
	for k, ci := range m {
		this.stopWait.Add(1)
		cache := ci.cache
		go func(k string, c ICache) {
			defer this.stopWait.Done()
			if !c.IsStart() {
				logger.Debug(tag, "boot stop Cache[%s] - skip", k)
				return
			}
			logger.Debug(tag, "boot stop Cache[%s]", k)
			err := c.Stop()
			if err != nil {
				logger.Debug(tag, "boot stop Cache[%s] fail - %s", k, err)
			}
		}(k, cache)
	}
}

func (this *CacheService) waitAllCacheStop() {
	if this.stopWait != nil {
		this.stopWait.Wait()
	}
}
