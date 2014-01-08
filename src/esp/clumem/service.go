package clumem

import (
	"bmautil/qexec"
	"config"
	"esp/sqlite"
	"logger"
)

type serviceItem struct {
	group  *localMemGroup
	config *MemGroupConfig
}

type Service struct {
	name      string
	database  *sqlite.SqliteServer
	memgroups map[string]*serviceItem
	config    configInfo

	executor *qexec.QueueExecutor
}

func NewService(name string, db *sqlite.SqliteServer) *Service {
	this := new(Service)
	this.name = name
	this.database = db
	this.memgroups = make(map[string]*serviceItem)
	this.executor = qexec.NewQueueExecutor(tag, 128, this.requestHandler)
	this.executor.StopHandler = this.stopHandler

	this.initDatabase()

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
	if !this.executor.Run() {
		return false
	}
	err := this.executor.DoSync("startRun", func() error {
		return this.doRun()
	})
	if err != nil {
		logger.Error(tag, "%s start run fail - %s", this.name, err)
		return false
	}
	return true
}

func (this *Service) Stop() bool {
	this.executor.Stop()
	return true
}

func (this *Service) Cleanup() bool {
	this.executor.WaitStop()
	return true
}

func (this *Service) Save() error {
	return this.executor.DoNow("save", func() error {
		return this.doSave()
	})
}

func (this *Service) CreateMemGroup(cfg *MemGroupConfig) error {
	return this.executor.DoSync("create", func() error {
		return this.doCreateMemGroup(cfg)
	})
}

func (this *Service) Get() error {

}
