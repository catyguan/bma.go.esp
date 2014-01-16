package xmemprot

import (
	"bmautil/binlog"
	"bmautil/coder"
	"esp/espnet"
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

type SHAction int8

const (
	SHA_NONE = iota
	SHA_SLAVE_JOIN
	SHA_BINLOG_EVENT
	SHA_GET
	SHA_GET_RESP
	SHA_LIST
	SHA_SET
	SHA_SET_RESP
	SHA_DELETE
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
	Group string
	Key   string
}

func (this *SHRequestDelete) Init(g string, key MemKey) {
	this.Group = g
	this.Key = key.ToString()
}

func (this *SHRequestDelete) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_DELETE), coder.Int8)
	xd.Add(1, this.Group, coder.LenString)
	xd.Add(2, this.Key, coder.LenString)
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
		}
	}
	return nil
}
