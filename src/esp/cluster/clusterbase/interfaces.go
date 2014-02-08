package clusterbase

import (
	"bmautil/byteutil"
	"encoding/binary"
	"time"
)

// OpVer
type OpVer int64

func NewOpVer() OpVer {
	return OpVer(time.Now().UnixNano())
}

type OpVerCoder int

func (this OpVerCoder) DoEncode(w *byteutil.BytesBufferWriter, v OpVerCoder) {
	binary.Write(w, binary.BigEndian, int64(v))
}

func (this OpVerCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	this.DoEncode(w, v.(OpVerCoder))
	return nil
}

func (this OpVerCoder) DoDecode(r *byteutil.BytesBufferReader) (OpVerCoder, error) {
	var v OpVerCoder
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func (this OpVerCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v, err := this.DoDecode(r)
	return v, err
}

// Role
type RoleType int

const (
	ROLE_NONE     = 0
	ROLE_LEADER   = 1
	ROLE_FOLLOWER = 2
	ROLE_LEANER   = 3
	ROLE_OBSERVER = 4
)

// OpCoder
type OpCoder interface {
	Encode(w *byteutil.BytesBufferWriter, v interface{}) error
	Decode(r *byteutil.BytesBufferReader) (interface{}, error)
}

// OpHandler
type OpHandler interface {
	Execute(op interface{}) error
}
