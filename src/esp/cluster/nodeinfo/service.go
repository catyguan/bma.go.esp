package nodeinfo

import (
	"bmautil/byteutil"
	"bmautil/coder"
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
	NodeIdCoder = nodeIdCoder(0)
)

type nodeIdCoder int

func (O nodeIdCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	coder.Uint64.DoEncode(w, uint64(v.(NodeId)))
	return nil
}

func (O nodeIdCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v, err := coder.Uint64.DoDecode(r)
	if err != nil {
		return nil, err
	}
	return NodeId(v), nil
}

type Service struct {
	name     string
	nodeId   NodeId
	nodeName string
}

func NewService(name string) *Service {
	this := new(Service)
	this.name = name
	return this
}

func (this *Service) GetId() NodeId {
	return this.nodeId
}

func (this *Service) GetNodeName() string {
	return this.nodeName
}
