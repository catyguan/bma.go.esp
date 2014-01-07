package clumem

import (
	"bmautil/qexec"
	"config"
	"esp/sqlite"
)

type Service struct {
	name      string
	database  *sqlite.SqliteServer
	memgroups map[string]*localMemGroup
	config    configInfo

	executor *qexec.QueueExecutor
}

func NewService(name string, db *sqlite.SqliteServer) *Service {
	this := new(Service)
	this.name = name
	this.database = db
	this.memgroups = make(map[string]*localMemGroup)
	this.executor = qexec.NewQueueExecutor(tag, 128, this.requestHandler)
	this.executor.StopHandler = this.stopHandler

	this.initDatabase()

	return this
}

func (this *Service) toIService() IService {
	return this
}

func (this *Service) Name() string {
	return this.name
}

type configInfo struct {
	AdminWord string
	SafeMode  bool
}

func (this *Service) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		this.config = cfg
	}
	return true
}

func (this *Service) Start() bool {
	cfg, ok := this.loadRuntimeConfig()
	if !ok {
		if !this.config.SafeMode {
			return false
		}
	}
	if !this.setupByConfig(cfg) {
		return false
	}

	return true
}

func (this *Service) Run() bool {
	return true
}

func (this *Service) Stop() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}

func (this *Service) Save() error {
	return this.executor.DoNow("save", func() error {
		return this.doSave()
	})
}
