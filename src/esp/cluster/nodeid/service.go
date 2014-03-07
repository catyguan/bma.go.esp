package nodeid

import (
	"bmautil/byteutil"
	"bmautil/coder"
	"config"
	"logger"
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

type Service struct {
	name   string
	nodeId NodeId
}

func NewService(name string) *Service {
	this := new(Service)
	this.name = name
	return this
}

func (this *Service) Name() string {
	return this.name
}

type configInfo struct {
	NodeId uint64
}

func (this *Service) parseConfig() *configInfo {
	cfg := configInfo{}
	if !config.GetBeanConfig(this.name, &cfg) {
		logger.Error(tag, "invalid config bean '%s'", this.name)
		return nil
	}
	if cfg.NodeId == 0 {
		logger.Error(tag, "config.NodeId invalid")
		return nil
	}
	return &cfg
}

func (this *Service) CheckConfig() bool {
	return this.parseConfig() != nil
}

func (this *Service) Init() bool {
	cfg := this.parseConfig()
	if cfg == nil {
		return false
	}
	this.nodeId = NodeId(cfg.NodeId)
	return true
}

func (this *Service) GetId() NodeId {
	return this.nodeId
}
