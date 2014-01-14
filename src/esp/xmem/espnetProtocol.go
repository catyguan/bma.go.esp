package xmem

import (
	"bmautil/binlog"
	"bmautil/coder"
	"esp/espnet"
)

type SHAction int8

const (
	SHA_NONE = iota
	SHA_SLAVE_JOIN
	SHA_BINLOG_EVENT
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
