package xmemservice

import (
	"bmautil/binlog"
	"bmautil/byteutil"
	xcoder "bmautil/coder"
	"esp/xmem/xmemprot"
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
	Key       xmemprot.MemKey
	Value     interface{}
	Version   xmemprot.MemVer
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

type XMemWalker func(key xmemprot.MemKey, val interface{}, ver xmemprot.MemVer) WalkStep

// Coder
type XMemCoder interface {
	Encode(val interface{}) (string, []byte, error)
	Decode(flag string, data []byte) (interface{}, int, error)
}

// Snapshot
type XMemSnapshot struct {
	Key     string
	Version xmemprot.MemVer
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
	o.Key, err = xcoder.LenString.DoDecode(r, 1024*100)
	if err != nil {
		return nil, err
	}
	var v2 uint64
	v2, err = xcoder.Uint64.DoDecode(r)
	if err != nil {
		return nil, err
	}
	o.Version = xmemprot.MemVer(v2)
	o.Kind, err = xcoder.LenString.DoDecode(r, 1024)
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

type XMemGroupSnapshot struct {
	BLVer     binlog.BinlogVer
	Snapshots []*XMemSnapshot
}

func (this *XMemGroupSnapshot) Encode() ([]byte, error) {
	buf := byteutil.NewBytesBuffer()
	w := buf.NewWriter()
	binlog.BinlogVerCoder(0).DoEncode(w, this.BLVer)
	xcoder.Int.DoEncode(w, len(this.Snapshots))
	for _, s := range this.Snapshots {
		s.Encode(w)
	}
	w.End()
	return buf.ToBytes(), nil
}

func (this *XMemGroupSnapshot) Decode(data []byte) error {
	buf := byteutil.NewBytesBufferB(data)
	r := buf.NewReader()
	blver, err0 := binlog.BinlogVerCoder(0).DoDecode(r)
	if err0 != nil {
		return err0
	}
	this.BLVer = blver
	l, err1 := xcoder.Int.DoDecode(r)
	if err1 != nil {
		return err1
	}
	slist := []*XMemSnapshot{}
	for i := 0; i < l; i++ {
		ss, err2 := DecodeSnapshot(r)
		if err2 != nil {
			return err2
		}
		slist = append(slist, ss)
	}
	this.Snapshots = slist
	return nil
}

// API
type XMem interface {
	Get(key xmemprot.MemKey) (interface{}, xmemprot.MemVer, bool, error)
	GetAndListen(key xmemprot.MemKey, id string, lis XMemListener) (interface{}, xmemprot.MemVer, bool, error)
	List(key xmemprot.MemKey) ([]string, bool, error)
	ListAndListen(key xmemprot.MemKey, id string, lis XMemListener) ([]string, bool, error)
	AddListener(key xmemprot.MemKey, id string, lis XMemListener) error
	RemoveListener(key xmemprot.MemKey, id string) error

	Set(key xmemprot.MemKey, val interface{}, sz int) (xmemprot.MemVer, error)
	CompareAndSet(key xmemprot.MemKey, val interface{}, sz int, ver xmemprot.MemVer) (xmemprot.MemVer, error)
	SetIfAbsent(key xmemprot.MemKey, val interface{}, sz int) (xmemprot.MemVer, error)

	Delete(key xmemprot.MemKey) (bool, error)
	CompareAndDelete(key xmemprot.MemKey, ver xmemprot.MemVer) (bool, error)
}
