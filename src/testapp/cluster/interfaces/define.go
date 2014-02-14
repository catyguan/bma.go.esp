package interfaces

import (
	"bmautil/byteutil"
	"bmautil/coder"
	"fmt"

	"code.google.com/p/goprotobuf/proto"
)

type Account interface {
	Get() (int64, error)

	Modify(v int64) (int64, error)
}

// Coder
type OpCoder4Account int

func (this OpCoder4Account) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	n := ""
	switch v.(type) {
	case *Req4AccountGet:
		n = "rg"
	case *Req4AccountModify:
		n = "rm"
	case *Resp4AccountGet:
		n = "pg"
	case *Resp4AccountModify:
		n = "pm"
	default:
		return fmt.Errorf("unknow type '%T'", v)
	}
	b, err := proto.Marshal(v.(proto.Message))
	if err != nil {
		return err
	}
	coder.LenString.DoEncode(w, n)
	coder.LenBytes.DoEncode(w, b)
	return nil
}

func (this OpCoder4Account) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	n, err := coder.LenString.DoDecode(r, 1024)
	if err != nil {
		return nil, err
	}
	bs, err2 := coder.LenBytes.DoDecode(r, 0)
	if err2 != nil {
		return nil, err2
	}
	var pb proto.Message
	switch n {
	case "rg":
		pb = new(Req4AccountGet)
	case "pg":
		pb = new(Resp4AccountGet)
	case "rm":
		pb = new(Req4AccountModify)
	case "pm":
		pb = new(Resp4AccountModify)
	}
	if pb == nil {
		return nil, fmt.Errorf("unknow type '%s'", n)
	}
	err3 := proto.Unmarshal(bs, pb)
	if err3 != nil {
		return nil, err3
	}
	return pb, nil
}
