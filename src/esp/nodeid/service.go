package nodeid

import (
	"config"
	"esp/sqlite"
	"logger"
	"sync"
	"uuid"
)

const (
	tag       = "nodeid"
	tableName = "tbl_nodeid"
)

type NodeId uint64

type Listener func(nodeId uint64)

type Service struct {
	name     string
	database *sqlite.SqliteServer

	lock      sync.Mutex
	nodeId    uint64
	listeners map[string]Listener
}

func NewService(name string, db *sqlite.SqliteServer) *Service {
	this := new(Service)
	this.name = name
	this.database = db
	this.initDatabase()
	return this
}

func (this *Service) Name() string {
	return this.name
}

type configInfo struct {
	NodeId uint64
}

func (this *Service) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		if cfg.NodeId == 0 {
			logger.Error(tag, "config.NodeId invalid")
			return false
		}
		this.nodeId = cfg.NodeId
	}
	return true
}

func (this *Service) Start() bool {
	ok := this.loadRuntimeConfig()
	if !ok {
		return false
	}
	for {
		if this.nodeId != 0 {
			break
		}
		uid, err := uuid.NewV4()
		if err != nil {
			logger.Error(tag, "create uuid fail - %s", err)
			return false
		}
		var val uint64 = 1
		str := uid.String()
		sz := len(str)
		for i := 0; i < sz; i++ {
			val += (val * 37) + uint64(str[i])
		}
		this.nodeId = val
		if this.nodeId != 0 {
			err = this.storeRuntimeConfig()
			if err != nil {
				logger.Error(tag, "store nodeId fail - %s", err)
				return false
			}
		}
	}

	return true
}

func (this *Service) Save() error {
	return this.storeRuntimeConfig()
}

func (this *Service) initDatabase() {
	this.database.InitRuntmeConfigTable(tableName, []int{1})
}

func (this *Service) loadRuntimeConfig() bool {
	var cfg configInfo
	err := this.database.LoadRuntimeConfig(tableName, 1, &cfg)
	if err != nil {
		logger.Error(tag, "loadRuntimeConfig fail - %s", err)
		return false
	}
	this.nodeId = cfg.NodeId
	return true
}

func (this *Service) storeRuntimeConfig() error {
	cfg := new(configInfo)
	cfg.NodeId = this.nodeId
	return this.database.StoreRuntimeConfig(tableName, 1, cfg)
}

func (this *Service) GetId() uint64 {
	return this.nodeId
}

func (this *Service) GetAndListen(id string, lis Listener) uint64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.listeners == nil {
		this.listeners = make(map[string]Listener)
	}
	this.listeners[id] = lis
	return this.nodeId
}

func (this *Service) RemoveListen(id string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.listeners != nil {
		delete(this.listeners, id)
	}
}

func (this *Service) SetId(nid uint64) error {
	var old uint64
	this.lock.Lock()
	old, this.nodeId = this.nodeId, nid
	err := this.storeRuntimeConfig()
	if err != nil {
		this.lock.Unlock()
		this.nodeId = old
		return err
	}
	this.lock.Unlock()
	for _, lis := range this.listeners {
		lis(nid)
	}
	return nil
}
