package xmem

import (
	"bmautil/qexec"
	"bytes"
	"config"
	"esp/sqlite"
	"fmt"
	"logger"
)

type serviceItem struct {
	group   *localMemGroup
	profile *MemGroupProfile
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
	this.StoreAllMemGroup()
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

func (this *Service) UpdateMemGroupConfig(name string, cfg *MemGroupConfig) error {
	return this.executor.DoSync("config", func() error {
		return this.doUpdateMemGroupConfig(name, cfg)
	})
}

func (this *Service) EnableMemGroup(prof *MemGroupProfile) error {
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

func (this *Service) StoreAllMemGroup() error {
	return this.executor.DoSync("storeAll", func() error {
		return this.doStoreAllMemGroup()
	})
}

func (this *Service) SaveMemGroup(name string, fileName string) error {
	return this.executor.DoSync("saveMG", func() error {
		return this.doMemSave(name, fileName, nil)
	})
}

func (this *Service) LoadMemGroup(name string, fileName string) error {
	return this.executor.DoSync("loadMG", func() error {
		return this.doMemLoad(name, fileName, nil)
	})
}

func (this *Service) SaveBinlogSnapshot(name string, fileName string) error {
	return this.executor.DoSync("saveBL", func() error {
		return this.doSaveBinlogSnapshot(name, fileName)
	})
}

func (this *Service) CreateXMem(name string) (XMem, error) {
	var r XMem
	err := this.executor.DoSync("createXMem", func() error {
		_, err := this.doGetGroup(name)
		if err != nil {
			return err
		}
		obj := new(XMem4Service)
		obj.Init(this, name)
		r = obj
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (this *Service) Dump(g string, key MemKey, all bool) (string, error) {
	str := ""
	err := this.executor.DoSync("dump", func() error {
		item, err := this.doGetGroup(g)
		if err != nil {
			return err
		}
		it, ok := item.group.Get(key)
		if !ok {
			return fmt.Errorf("<%s> not exists", key)
		}
		buf := bytes.NewBuffer([]byte{})
		it.Dump(key.ToString(), buf, 0, all)
		str = item.group.String() + "\n" + buf.String()
		return nil
	})
	return str, err
}