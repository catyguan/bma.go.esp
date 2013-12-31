package esimple

import (
	"app/cacheserver"
	"esp/shell"
	"fmt"
	"sync"
	"time"
)

func init() {
	cacheserver.RegCacheFactory("simple", NewSimpleCache)
}

type CacheItem struct {
	data     []byte
	deadTime int64
}

type SimpleCache struct {
	name    string
	service *cacheserver.CacheService
	lock    sync.Mutex
	items   map[string]CacheItem
	started bool
}

func NewSimpleCache() cacheserver.ICache {
	this := new(SimpleCache)
	this.items = make(map[string]CacheItem)
	return this
}

func (this *SimpleCache) InitCache(s *cacheserver.CacheService, n string) {
	this.service = s
	this.name = n
}

func (this *SimpleCache) Get(req *cacheserver.GetRequest, rep chan *cacheserver.GetResult) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if rep != nil {
		ci, ok := this.items[req.Key]
		res := cacheserver.NewGetResult(this.name, req.Key, req.Trace)
		if ok {
			if ci.deadTime != 0 && ci.deadTime <= time.Now().Unix() {
				// timeout
				delete(this.items, req.Key)
				if req.Trace {
					res.Traces([]string{"timeout"})
				}
			} else {
				res.End(ok, ci.data, nil)
				rep <- res
				return nil
			}
		}
		res.End(ok, nil, nil)
		rep <- res
	}
	return nil
}

func (this *SimpleCache) Put(key string, val []byte, deadUnixtime int64) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	var ci CacheItem
	ci.data = val
	ci.deadTime = deadUnixtime
	this.items[key] = ci
	return nil
}

func (this *SimpleCache) Delete(key string) (bool, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, ok := this.items[key]
	delete(this.items, key)
	return ok, nil
}

func (this *SimpleCache) Load(key string) error {
	return nil
}

func (this *SimpleCache) QueryStats() (string, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return fmt.Sprintf("size=%d", len(this.items)), nil
}

func (this *SimpleCache) IsStart() bool {
	return this.started
}

func (this *SimpleCache) Start() error {
	this.started = true
	return nil
}

func (this *SimpleCache) Run() error {
	return nil
}

func (this *SimpleCache) Stop() error {
	this.started = false
	return nil
}

func (this *SimpleCache) FromConfig(cfg map[string]interface{}) error {
	return nil
}

func (this *SimpleCache) ToConfig() map[string]interface{} {
	return nil
}

func (this *SimpleCache) CreateShell() shell.ShellDir {
	r := shell.NewShellDirCommon(this.name)
	this.service.BuildCacheCommands(this.name, r)
	return r
}
