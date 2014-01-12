package xmem

import (
	"bmautil/qexec"
	"config"
	"esp/sqlite"
	"logger"
)

type serviceItem struct {
	group   *localMemGroup
	profile *memGroupProfile
	config  *MemGroupConfig
}

type Service struct {
	name      string
	database  *sqlite.SqliteServer
	memgroups map[string]*serviceItem

	config configInfo

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
	AdminWord     string
	SafeMode      bool
	NoWaitOnStart bool
}

func (this *Service) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {

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
	if this.config.NoWaitOnStart {
		this.executor.DoNow("startRunNow", this.doRun)
	} else {
		err := this.executor.DoSync("startRun", this.doRun)
		if err != nil {
			logger.Error(tag, "%s start run fail - %s", this.name, err)
			return false
		}
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

func (this *Service) EnableMemGroup(prof *memGroupProfile) error {
	return this.executor.DoSync("create", func() error {
		return this.doEnableMemGroup(prof)
	})
}

func (this *Service) ListMemGroupName() ([]string, error) {
	var r []string
	return r, this.executor.DoSync("list", func() error {
		r = this.doListMemGroupName()
		return nil
	})
}

func (this *Service) Get(key MemKey) error {
	return nil
}
