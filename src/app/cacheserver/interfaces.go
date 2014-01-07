package cacheserver

import "uprop"

type GetRequest struct {
	Key       string // cache key
	Trace     bool   // trace how to GET
	NotLoad   bool   // don't load if not exists
	TimeoutMs int32  // Timeout in ms, <=0 mean no timeout
	Update    bool   // do update not put
}

func NewGetRequest(key string) *GetRequest {
	this := new(GetRequest)
	this.Key = key
	return this
}

type GetResult struct {
	Done      bool
	Err       error
	Group     string
	Key       string
	Value     []byte
	TraceInfo []string
}

func NewGetResult(g, k string, trace bool) *GetResult {
	this := new(GetResult)
	this.Group = g
	this.Key = k
	if trace {
		this.TraceInfo = make([]string, 0)
	}
	return this
}

func (this *GetResult) End(ok bool, val []byte, trace []string) {
	this.Done = ok
	this.Value = val
	this.Traces(trace)
}

func (this *GetResult) Fail(err error, trace []string) {
	this.Done = false
	this.Err = err
	this.Value = nil
	this.Traces(trace)
}

func (this *GetResult) Traces(s []string) {
	if s == nil {
		return
	}
	if this.TraceInfo == nil {
		return
	}
	for _, t := range s {
		this.TraceInfo = append(this.TraceInfo, t)
	}
}

type ICacheConfig interface {
	GetProperties() []*uprop.UProperty

	Valid() error

	ToMap() map[string]interface{}

	FromMap(data map[string]interface{}) error
}

type CacheFactory interface {
	CreateCache(cfg ICacheConfig) (ICache, error)
	CreateConfig() ICacheConfig
}

// ICache, all methods is sync
type ICache interface {
	InitCache(s *CacheService, group string)

	Type() string

	GetConfig() ICacheConfig
	UpdateConfig(cfg ICacheConfig) error

	Get(req *GetRequest, rep chan *GetResult) error

	Put(key string, val []byte, deadUnixtime int64) error

	Delete(key string) (bool, error)

	Load(key string) error

	QueryStats() (string, error)

	IsStart() bool

	Start() error

	Run() error

	Stop() error
}
