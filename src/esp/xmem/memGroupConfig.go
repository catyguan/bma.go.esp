package xmem

import (
	"bmautil/valutil"
	"uprop"
)

type MemGroupConfig struct {
	NoSave bool
}

func (this *MemGroupConfig) Valid() error {
	return nil
}

func (this *MemGroupConfig) GetProperties() []*uprop.UProperty {
	b := new(uprop.UPropertyBuilder)
	b.NewProp("disableSave", "disable memory save").BeValue(this.NoSave, func(v string) error {
		this.NoSave = valutil.ToBool(v, this.NoSave)
		return nil
	})
	return b.AsList()
}

func (this *MemGroupConfig) ToMap() map[string]interface{} {
	return valutil.BeanToMap(this)
}

func (this *MemGroupConfig) FromMap(data map[string]interface{}) error {
	valutil.ToBean(data, this)
	return this.Valid()
}
