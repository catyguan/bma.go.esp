package nodebase

import (
	"esp/espnet/esnp"
	"fmt"
)

type NodeId uint64

func (this NodeId) Is(v uint64) bool {
	return this == NodeId(v)
}

const (
	INVALID = NodeId(0)
)

var (
	NodeIdCoder = nodeIdCoder(0)
)

type nodeIdCoder int

func (O nodeIdCoder) Encode(w esnp.EncodeWriter, v interface{}) error {
	esnp.Coders.Uint64.DoEncode(w, uint64(v.(NodeId)))
	return nil
}

func (O nodeIdCoder) Decode(r esnp.DecodeReader) (interface{}, error) {
	v, err := esnp.Coders.Uint64.DoDecode(r)
	if err != nil {
		return nil, err
	}
	return NodeId(v), nil
}

type NodeInfo struct {
	Id   uint64
	Name string
}

func (this *NodeInfo) Valid() error {
	if this.Id == 0 {
		return fmt.Errorf("Id invalid")
	}
	if this.Name == "" {
		return fmt.Errorf("Name invalid")
	}
	return nil
}

var (
	Id   NodeId
	Name string
)
