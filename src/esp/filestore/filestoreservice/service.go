package filestoreservice

import (
	"bmautil/qexec"
	"config"
	"esp/filestore/filestoreprot"
	"fmt"
	"os"
)

const (
	tag = "filestore"
)

type ServiceConfig struct {
	AppId filestoreprot.AppId
	Key   string
	Root  string
	Temp  string
}

func (this *ServiceConfig) Valid() error {
	if this.AppId == "" {
		return fmt.Errorf("appId empty")
	}
	if this.Root == "" {
		return fmt.Errorf("root empty")
	}
	return nil
}

type Service struct {
	name   string
	config configInfo

	executor *qexec.QueueExecutor
	sessions map[string]*session
}

func NewService(name string) *Service {
	this := new(Service)
	this.name = name
	this.executor = qexec.NewQueueExecutor(tag, 128, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.sessions = make(map[string]*session)
	return this
}

func (this *Service) Name() string {
	return this.name
}

type configInfo struct {
	TimeoutSec int
	TempDir    string
	Apps       map[string]*ServiceConfig
}

func (this *Service) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		this.config = cfg
	}
	if this.config.TimeoutSec <= 0 {
		this.config.TimeoutSec = 10
	}
	if this.config.TempDir == "" {
		this.config.TempDir = os.TempDir()
	}
	if this.config.Apps == nil {
		this.config.Apps = make(map[string]*ServiceConfig)
	}
	return true
}

func (this *Service) Start() bool {
	if !this.executor.Run() {
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
