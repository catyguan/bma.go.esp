package nodeid

import (
	"bmautil/byteutil"
	"bmautil/coder"
	"config"
	"esp/sqlite"
	"logger"
	"sync"
	"uuid"
)

const (
	tag   = "nodeid"
	rtKey = ".nodeid"
)

type NodeId uint64

const (
	INVALID = NodeId(0)
)

var (
	Coder = mycoder(0)
)

type mycoder int

func (O mycoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	coder.Uint64.DoEncode(w, uint64(v.(NodeId)))
	return nil
}

func (O mycoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v, err := coder.Uint64.DoDecode(r)
	if err != nil {
		return nil, err
	}
	return NodeId(v), nil
}

type Listener func(nodeId NodeId)

type Service struct {
	name     string
	database *sqlite.SqliteServer

	lock      sync.Mutex
	nodeId    NodeId
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
		this.nodeId = NodeId(cfg.NodeId)
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
		this.nodeId = NodeId(val)
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
	this.database.InitRuntmeConfigTable()
}

func (this *Service) loadRuntimeConfig() bool {
	var cfg configInfo
	err := this.database.LoadRuntimeConfig(this.name+rtKey, &cfg)
	if err != nil {
		logger.Error(tag, "loadRuntimeConfig fail - %s", err)
		return false
	}
	if cfg.NodeId > 0 {
		this.nodeId = NodeId(cfg.NodeId)
	}
	return true
}

func (this *Service) storeRuntimeConfig() error {
	cfg := new(configInfo)
	cfg.NodeId = uint64(this.nodeId)
	return this.database.StoreRuntimeConfig(this.name+rtKey, cfg)
}

func (this *Service) GetId() NodeId {
	return this.nodeId
}

func (this *Service) GetAndListen(id string, lis Listener) NodeId {
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

func (this *Service) SetId(nid NodeId) error {
	var old NodeId
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
