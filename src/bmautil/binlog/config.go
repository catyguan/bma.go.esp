package binlog

import (
	"bmautil/valutil"
	"fmt"
	"uprop"
)

type BinlogConfig struct {
	FileName  string
	Readonly  bool
	SyncOpNum int // 0-auto
}

func (this *BinlogConfig) Valid() error {
	if this.FileName == "" {
		return fmt.Errorf("binlog file name empty")
	}
	return nil
}

func (this *BinlogConfig) GetProperties() []*uprop.UProperty {
	b := new(uprop.UPropertyBuilder)
	b.NewProp("file", "binlog file name").Optional(false).BeValue(this.FileName, func(v string) error {
		this.FileName = v
		return nil
	})
	b.NewProp("readonly", "read only").BeValue(this.Readonly, func(v string) error {
		this.Readonly = valutil.ToBool(v, this.Readonly)
		return nil
	})
	b.NewProp("syncnum", "how many op to flush file write").BeValue(this.SyncOpNum, func(v string) error {
		this.SyncOpNum = valutil.ToInt(v, this.SyncOpNum)
		return nil
	})
	return b.AsList()
}
