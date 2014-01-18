package xmemprot

import (
	"bmautil/binlog"
	"bmautil/byteutil"
	"bmautil/coder"
	"esp/espnet"
	"esp/espnet/protpack"
	"fmt"
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

func (O MemVer) String() string {
	if O == VERSION_INVALID {
		return "NOVER"
	}
	return fmt.Sprintf("%d", O)
}

func (O MemVer) Valid() bool {
	return O != VERSION_INVALID
}

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

// API
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

	Delete(key MemKey) (bool, error)
	CompareAndDelete(key MemKey, ver MemVer) (bool, error)
}

type SHAction int8

const (
	SHA_NONE = iota
	SHA_SLAVE_JOIN
	SHA_BINLOG_EVENT
	SHA_GET
	SHA_GET_RESP
	SHA_LIST
	SHA_LIST_RESP
	SHA_SET
	SHA_SET_RESP
	SHA_DELETE
	SHA_DELETE_RESP
	SHA_LISTEN
	SHA_LISTEN_EVENT
)

type SHObject interface {
	Write(msg *espnet.Message) error
	Read(msg *espnet.Message) error
}

// SHRequestSlaveJoin
type SHRequestSlaveJoin struct {
	Group   string
	Version binlog.BinlogVer
}

func (this *SHRequestSlaveJoin) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_SLAVE_JOIN), coder.Int8)
	xd.Add(1, this.Group, coder.LenString)
	xd.Add(2, this.Version, binlog.BinlogVerCoder(0))
	return nil
}

func (this *SHRequestSlaveJoin) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 1:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Group = v.(string)
		case 2:
			v, err := it.Value(binlog.BinlogVerCoder(0))
			if err != nil {
				return err
			}
			this.Version = v.(binlog.BinlogVer)
		}
	}
	return nil
}

// SHRequestBinlog
type SHEventBinlog struct {
	Group   string
	Version binlog.BinlogVer
	Data    []byte
}

func (this *SHEventBinlog) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_BINLOG_EVENT), coder.Int8)
	xd.Add(1, this.Group, coder.LenString)
	xd.Add(2, this.Version, binlog.BinlogVerCoder(0))
	xd.Add(3, this.Data, coder.LenBytes)
	return nil
}

func (this *SHEventBinlog) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 1:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Group = v.(string)
		case 2:
			v, err := it.Value(binlog.BinlogVerCoder(0))
			if err != nil {
				return err
			}
			this.Version = v.(binlog.BinlogVer)
		case 3:
			v, err := it.Value(coder.LenBytes)
			if err != nil {
				return err
			}
			this.Data = v.([]byte)
		}
	}
	return nil
}

// SHRequestGet
type SHRequestGet struct {
	Group string
	Key   string
}

func (this *SHRequestGet) Init(g string, key MemKey) {
	this.Group = g
	this.Key = key.ToString()
}

func (this *SHRequestGet) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_GET), coder.Int8)
	xd.Add(1, this.Group, coder.LenString)
	xd.Add(2, this.Key, coder.LenString)
	return nil
}

func (this *SHRequestGet) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 1:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Group = v.(string)
		case 2:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Key = v.(string)
		}
	}
	return nil
}

// SHResponseGet
type SHResponseGet struct {
	Miss    bool
	Value   interface{}
	Version MemVer
}

func (this *SHResponseGet) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_GET_RESP), coder.Int8)
	xd.Add(3, this.Miss, coder.Bool)
	xd.Add(4, this.Value, coder.Varinat)
	xd.Add(5, uint64(this.Version), coder.Uint64)
	return nil
}

func (this *SHResponseGet) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 3:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Miss = v.(bool)
		case 4:
			v, err := it.Value(coder.Varinat)
			if err != nil {
				return err
			}
			this.Value = v
		case 5:
			v, err := it.Value(coder.Uint64)
			if err != nil {
				return err
			}
			this.Version = MemVer(v.(uint64))
		}
	}
	return nil
}

// SHRequestSet
type SHRequestSet struct {
	Group   string
	Key     string
	Value   interface{}
	Size    int
	Version MemVer
	Absent  bool
}

func (this *SHRequestSet) InitSet(g string, key MemKey, val interface{}, sz int) {
	this.Group = g
	this.Key = key.ToString()
	this.Value = val
	this.Size = sz
	this.Version = VERSION_INVALID
}

func (this *SHRequestSet) InitCompareAndSet(g string, key MemKey, val interface{}, sz int, ver MemVer) {
	this.Group = g
	this.Key = key.ToString()
	this.Value = val
	this.Size = sz
	this.Version = ver
}

func (this *SHRequestSet) InitSetIfAbsent(g string, key MemKey, val interface{}, sz int) {
	this.Group = g
	this.Key = key.ToString()
	this.Value = val
	this.Size = sz
	this.Version = VERSION_INVALID
	this.Absent = true
}

func (this *SHRequestSet) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_SET), coder.Int8)
	xd.Add(1, this.Group, coder.LenString)
	xd.Add(2, this.Key, coder.LenString)
	xd.Add(3, this.Value, coder.Varinat)
	xd.Add(4, this.Size, coder.Int)
	xd.Add(5, uint64(this.Version), coder.Uint64)
	xd.Add(6, this.Absent, coder.Bool)
	return nil
}

func (this *SHRequestSet) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 1:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Group = v.(string)
		case 2:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Key = v.(string)
		case 3:
			v, err := it.Value(coder.Varinat)
			if err != nil {
				return err
			}
			this.Value = v
		case 4:
			v, err := it.Value(coder.Int)
			if err != nil {
				return err
			}
			this.Size = v.(int)
		case 5:
			v, err := it.Value(coder.Uint64)
			if err != nil {
				return err
			}
			this.Version = MemVer(v.(uint64))
		case 6:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Absent = v.(bool)
		}
	}
	return nil
}

// SHResponseSet
type SHResponseSet struct {
	Version MemVer
}

func (this *SHResponseSet) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_SET_RESP), coder.Int8)
	xd.Add(5, uint64(this.Version), coder.Uint64)
	return nil
}

func (this *SHResponseSet) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 5:
			v, err := it.Value(coder.Uint64)
			if err != nil {
				return err
			}
			this.Version = MemVer(v.(uint64))
		}
	}
	return nil
}

// SHRequestDelete
type SHRequestDelete struct {
	Group   string
	Key     string
	Version MemVer
}

func (this *SHRequestDelete) Init(g string, key MemKey) {
	this.Group = g
	this.Key = key.ToString()
	this.Version = VERSION_INVALID
}

func (this *SHRequestDelete) InitCompareAndDelete(g string, key MemKey, ver MemVer) {
	this.Group = g
	this.Key = key.ToString()
	this.Version = ver
}

func (this *SHRequestDelete) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_DELETE), coder.Int8)
	xd.Add(1, this.Group, coder.LenString)
	xd.Add(2, this.Key, coder.LenString)
	xd.Add(5, this.Version, coder.Uint64)
	return nil
}

func (this *SHRequestDelete) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 1:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Group = v.(string)
		case 2:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Key = v.(string)
		case 5:
			v, err := it.Value(coder.Uint64)
			if err != nil {
				return err
			}
			this.Version = MemVer(v.(uint64))
		}
	}
	return nil
}

// SHResponseDelete
type SHResponseDelete struct {
	Done bool
}

func (this *SHResponseDelete) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_DELETE_RESP), coder.Int8)
	xd.Add(1, this.Done, coder.Bool)
	return nil
}

func (this *SHResponseDelete) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 1:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Done = v.(bool)
		}
	}
	return nil
}

// SHRequestList
type SHRequestList struct {
	Group string
	Key   string
}

func (this *SHRequestList) Init(g string, key MemKey) {
	this.Group = g
	this.Key = key.ToString()
}

func (this *SHRequestList) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_LIST), coder.Int8)
	xd.Add(1, this.Group, coder.LenString)
	xd.Add(2, this.Key, coder.LenString)
	return nil
}

func (this *SHRequestList) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 1:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Group = v.(string)
		case 2:
			v, err := it.Value(coder.LenString)
			if err != nil {
				return err
			}
			this.Key = v.(string)
		}
	}
	return nil
}

// SHResponseList
type SHResponseList struct {
	Names []string
	Miss  bool
}

func (this *SHResponseList) encodeNames(w *byteutil.BytesBufferWriter, v interface{}) error {
	nlist := v.([]string)
	coder.Int.Encode(w, len(nlist))
	for _, n := range nlist {
		coder.LenString.DoEncode(w, n)
	}
	return nil
}

func (this *SHResponseList) decodeNames(r *byteutil.BytesBufferReader) (interface{}, error) {
	l, err := coder.Int.DoDecode(r)
	if err != nil {
		return nil, err
	}
	if l < 0 || l > 1000*1000 {
		return nil, fmt.Errorf("invalid list len %d", l)
	}
	nlist := make([]string, l)
	for i := 0; i < l; i++ {
		nlist[i], err = coder.LenString.DoDecode(r, 1000*1000)
		if err != nil {
			return nil, err
		}
	}
	return nlist, nil
}

func (this *SHResponseList) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_LIST_RESP), coder.Int8)
	xd.Add(1, this.Names, protpack.NewEncoderFuc(this.encodeNames))
	xd.Add(2, this.Miss, coder.Bool)
	return nil
}

func (this *SHResponseList) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 0:
		case 1:
			v, err := it.Value(protpack.NewDecoderFuc(this.decodeNames))
			if err != nil {
				return err
			}
			this.Names = v.([]string)
		case 2:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Miss = v.(bool)
		}
	}
	return nil
}
