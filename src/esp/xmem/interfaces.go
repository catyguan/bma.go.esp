package xmem

import "strings"

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

type MemVer uint64

const (
	VERSION_INVALID = MemVer(0xFFFFFFFFFFFFFFFF)
)

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

type XMemEvent struct {
	Action    Action
	GroupName string
	Key       MemKey
	Value     interface{}
	Version   MemVer
}

type XMemListener func(events []*XMemEvent)

type WalkStep int

const (
	WALK_STEP_NEXT = iota
	WALK_STEP_OVER
	WALK_STEP_OUT
	WALK_STEP_END
)

type XMemWalker func(key MemKey, val interface{}, ver MemVer) WalkStep

type XMemCoder interface {
	Encode(val interface{}) (string, []byte, error)
	Decode(flag string, data []byte) (interface{}, int, error)
}

type XMemSnapshot struct {
	Key     string
	Version MemVer
	Kind    string
	Data    []byte
}

type XMem interface {
}
