package clusterbase

import (
	"bmautil/byteutil"
	"encoding/binary"
	"time"
)

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
