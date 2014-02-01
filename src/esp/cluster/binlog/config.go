package binlog

import (
	"bmautil/valutil"
	"fmt"
	"path/filepath"
	"uprop"
)

type BinlogConfig struct {
	LogDir        string
	FileFormatter string
	FileMaxSize   int64
	Readonly      bool
	SyncOpNum     int // 0-auto
}

func (this *BinlogConfig) Valid() error {
	if this.LogDir == "" {
		return fmt.Errorf("binlog log dir name empty")
	}
	if this.FileFormatter == "" {
		this.FileFormatter = "binlog.%d.blog"
	}
	return nil
}

func (this *BinlogConfig) GetFileName(n int) string {
	return filepath.Join(this.LogDir, fmt.Sprintf(this.FileFormatter, n))
}

func (this *BinlogConfig) GetProperties() []*uprop.UProperty {
	b := new(uprop.UPropertyBuilder)
	b.NewProp("logdir", "binlog dir name").Optional(false).BeValue(this.LogDir, func(v string) error {
		this.LogDir = v
		return nil
	})
	b.NewProp("namef", "binlog file name formatter").BeValue(this.FileFormatter, func(v string) error {
		this.FileFormatter = v
		return nil
	})
	b.NewProp("fmax", "binlog file max(MB)").BeValue(valutil.SizeString(uint64(this.FileMaxSize), 1024, valutil.SizeM), func(v string) error {
		iv, err := valutil.ToSize(v, 1024, valutil.SizeM)
		if err != nil {
			return err
		}
		this.FileMaxSize = int64(iv)
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
