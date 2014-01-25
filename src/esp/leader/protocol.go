package leader

import (
	"bmautil/binlog"
	"bmautil/coder"
	"esp/espnet"
)

type SHAction int8

const (
	SHA_NONE = iota
	SHA_PING
)

// SHPing
type SHPing struct {
	NodeId  uint64
	Version binlog.BinlogVer
}

func (this *SHPing) Write(msg *espnet.Message) error {
	xd := msg.XDatas()
	xd.Add(0, int8(SHA_SLAVE_JOIN), coder.Int8)
	xd.Add(1, this.Group, coder.LenString)
	xd.Add(2, this.Version, binlog.BinlogVerCoder(0))
	return nil
}

func (this *SHPing) Read(msg *espnet.Message) error {
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
