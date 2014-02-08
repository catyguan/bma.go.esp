package interfaces

import (
	"bmautil/byteutil"
	"fmt"
)

type Account interface {
	Get() (int64, error)

	Modify(v int64) (int64, error)
}

// Account::Get
type AccountGetReq struct {
}

func (this AccountGetReq) CodeType() string {
	return "AccountGetReq"
}
func (this AccountGetReq) Encode(w *byteutil.BytesBufferWriter) error {
	return nil
}
func (this AccountGetReq) Decode(r *byteutil.BytesBufferReader) error {
	return nil
}

type Resp4AccountResult struct {
	Value int64
}

// Account:Modify
type Req4AccountModify struct {
}
type Resp4AccountModify struct {
	Value int64
}

// Coder
type OpCoder4Account int

func (this OpCoder4Account) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	Int.DoEncode(w, len(bs))
	if bs != nil {
		w.Write(bs)
	}
	return nil
}

func (this OpCoder4Account) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	l, err := Int.DoDecode(r)
	if err != nil {
		return nil, err
	}
	if maxlen > 0 && l > maxlen {
		return nil, fmt.Errorf("too large bytes block - %d/%d", l, maxlen)
	}
	p := make([]byte, l)
	if l > 0 {
		_, err = r.Read(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}
