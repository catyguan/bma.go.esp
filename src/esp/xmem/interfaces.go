package xmem

import (
	"bmautil/byteutil"
	xcoder "bmautil/coder"
	"strings"
)

// MemKey
type MemKey []string

func (list MemKey) ToString() string {
	return strings.Join(list, "/")
}

func MemKeyFromString(s string) MemKey {
	if s == "" {
		return MemKey{}
	}
	return strings.Split(s, "/")
}

func (list MemKey) At(idx int) (string, bool) {
	if idx < 0 {
		idx = len(list) + idx
	}
	if idx >= 0 && idx < len(list) {
		return list[idx], true
	}
	return "", false
}

// MemVer
type MemVer uint64

const (
	VERSION_INVALID = MemVer(0xFFFFFFFFFFFFFFFF)
)

// Action
type Action int

func (O Action) String() string {
	switch O {
	case ACTION_NONE:
		return "NONE"
	case ACTION_NEW:
		return "NEW"
	case ACTION_UPDATE:
		return "UDPATE"
	case ACTION_DELETE:
		return "DELETE"
	case ACTION_CLEAR:
		return "CLEAR"
	default:
		return "UNKNOW"
	}
}

const (
	ACTION_NONE = iota
	ACTION_NEW
	ACTION_UPDATE
	ACTION_DELETE
	ACTION_CLEAR
)

// XMemEvent & Listener
type XMemEvent struct {
	Action    Action
	GroupName string
	Key       MemKey
	Value     interface{}
	Version   MemVer
}

type XMemListener func(events []*XMemEvent)

// Walk
type WalkStep int

const (
	WALK_STEP_NEXT = iota
	WALK_STEP_OVER
	WALK_STEP_OUT
	WALK_STEP_END
)

type XMemWalker func(key MemKey, val interface{}, ver MemVer) WalkStep

// Coder
type XMemCoder interface {
	Encode(val interface{}) (string, []byte, error)
	Decode(flag string, data []byte) (interface{}, int, error)
}

// Snapshot
type XMemSnapshot struct {
	Key     string
	Version MemVer
	Kind    string
	Data    []byte
}

func (this *XMemSnapshot) Encode(w *byteutil.BytesBufferWriter) error {
	xcoder.LenString.DoEncode(w, this.Key)
	xcoder.Uint64.DoEncode(w, uint64(this.Version))
	xcoder.LenString.DoEncode(w, this.Kind)
	xcoder.Int.DoEncode(w, len(this.Data))
	if len(this.Data) > 0 {
		w.Write(this.Data)
	}
	return nil
}

func DecodeSnapshot(r *byteutil.BytesBufferReader) (*XMemSnapshot, error) {
	o := new(XMemSnapshot)
	var err error
	o.Key, err = xcoder.LenString.DoDecode(r)
	if err != nil {
		return nil, err
	}
	var v2 uint64
	v2, err = xcoder.Uint64.DoDecode(r)
	if err != nil {
		return nil, err
	}
	o.Version = MemVer(v2)
	o.Kind, err = xcoder.LenString.DoDecode(r)
	if err != nil {
		return nil, err
	}
	var v3 int
	v3, err = xcoder.Int.DoDecode(r)
	if err != nil {
		return nil, err
	}
	bs := make([]byte, v3)
	_, err = r.Read(bs)
	if err != nil {
		return nil, err
	}
	o.Data = bs
	return o, nil
}

// API
type XMemOP int

const (
	OP_SNAPSHOT = iota
	OP_SET
	OP_DELETE
	OP_CLEAR
)

type XMem interface {
	Get(key MemKey) (interface{}, MemVer, bool, error)
	GetAndListen(key MemKey, id string, lis XMemListener) (interface{}, MemVer, bool, error)
	List(key MemKey) ([]string, bool, error)
	ListAndListen(key MemKey, id string, lis XMemListener) ([]string, bool, error)
	AddListener(key MemKey, id string, lis XMemListener) error
	RemoveListener(key MemKey, id string) error

	Set(key MemKey, val interface{}, sz int) (MemVer, error)
	CompareAndSet(key MemKey, val interface{}, sz int, ver MemVer) (MemVer, error)
	SetIfAbsent(key MemKey, val interface{}, sz int) (MemVer, error)

	Delete(key MemKey, ver MemVer) (bool, error)
}
