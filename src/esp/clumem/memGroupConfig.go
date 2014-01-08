package clumem

import (
	"bmautil/valutil"
	"errors"
	"uprop"
)

type MemGroupConfig struct {
	Name string
}

func (this *MemGroupConfig) Valid() error {
	if this.Name == "" {
		return errors.New("memory group name empty")
	}
	return nil
}

func (this *MemGroupConfig) GetProperties() []*uprop.UProperty {
	r := make([]*uprop.UProperty, 0)
	r = append(r, uprop.NewUProperty("name", this.Name, false, "memory group name", func(v string) error {
		this.Name = v
		return nil
	}))
	return r
}

func (this *MemGroupConfig) ToMap() map[string]interface{} {
	return valutil.BeanToMap(this)
}

func (this *MemGroupConfig) FromMap(data map[string]interface{}) error {
	valutil.ToBean(data, this)
	return this.Valid()
}
