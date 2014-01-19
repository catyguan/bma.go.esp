package xmemservice

import (
	"bmautil/binlog"
	"bmautil/valutil"
	"uprop"
)

type MemGroupConfig struct {
	NoSave        bool
	BLConfig      *binlog.BinlogConfig
	BLSlaveConfig *BLSlaveConfig
}

func (this *MemGroupConfig) Valid() error {
	if this.BLConfig != nil {
		err := this.BLConfig.Valid()
		if err != nil {
			return err
		}
	}
	if this.BLSlaveConfig != nil {
		err := this.BLSlaveConfig.Valid()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MemGroupConfig) IsEnableBinlog() bool {
	return this.BLConfig != nil
}

func (this *MemGroupConfig) IsBinlogWrite() bool {
	if this.BLConfig == nil {
		return false
	}
	if this.BLConfig.Readonly {
		return false
	}
	return true
}

func (this *MemGroupConfig) IsEnableBinlogSlave() bool {
	return this.BLSlaveConfig != nil
}

func (this *MemGroupConfig) GetProperties() []*uprop.UProperty {
	b := new(uprop.UPropertyBuilder)
	b.NewProp("disableSave", "disable memory save").BeValue(this.NoSave, func(v string) error {
		this.NoSave = valutil.ToBool(v, this.NoSave)
		return nil
	})
	b.NewProp("binlog", "enable binlog").BeValue(this.IsEnableBinlog(), func(v string) error {
		e := valutil.ToBool(v, this.IsEnableBinlog())
		if e {
			if this.BLConfig == nil {
				this.BLConfig = new(binlog.BinlogConfig)
			}
		} else {
			this.BLConfig = nil
		}
		return nil
	})
	if this.BLConfig != nil {
		props := this.BLConfig.GetProperties()
		b.MergeWithPrex(props, "blog.")
	}
	b.NewProp("mss", "enable binlog master/slave sync").BeValue(this.IsEnableBinlogSlave(), func(v string) error {
		e := valutil.ToBool(v, this.IsEnableBinlogSlave())
		if e {
			if this.BLSlaveConfig == nil {
				this.BLSlaveConfig = new(BLSlaveConfig)
			}
		} else {
			this.BLSlaveConfig = nil
		}
		return nil
	})
	if this.BLSlaveConfig != nil {
		props := this.BLSlaveConfig.GetProperties()
		b.MergeWithPrex(props, "mss.")
	}
	return b.AsList()
}

func (this *MemGroupConfig) ToMap() map[string]interface{} {
	return valutil.BeanToMap(this)
}

func (this *MemGroupConfig) FromMap(data map[string]interface{}) error {
	valutil.ToBean(data, this)
	return this.Valid()
}
