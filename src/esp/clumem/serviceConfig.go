package clumem

import (
	"bmautil/valutil"
	"esp/espnet/cfprototype"
	"uprop"
)

type serviceConfig struct {
	Global  *cfprototype.DialPoolPrototype
	Remotes []*cfprototype.DialPoolPrototype
}

func (this *serviceConfig) Valid() error {
	if this.Global != nil {
		err := this.Global.Valid()
		if err != nil {
			return err
		}
	}
	for _, remote := range this.Remotes {
		err := remote.Valid()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *serviceConfig) GetProperties() []*uprop.UProperty {
	b := new(uprop.UPropertyBuilder)
	return b.AsList()
}

func (this *serviceConfig) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	if this.Global != nil {
		r["Global"] = this.Global.ToMap()
	}
	if len(this.Remotes) > 0 {
		list := make([]interface{}, 0)
		for _, rem := range this.Remotes {
			list = append(list, rem.ToMap())
		}
		r["Remotes"] = list
	}
	return r
}

func (this *serviceConfig) FromMap(data map[string]interface{}) error {
	gp := valutil.ToStringMap(data["Global"])
	if gp != nil {
		this.Global = new(cfprototype.DialPoolPrototype)
		this.Global.FromMap(gp)
	}
	list := valutil.ToArray(data["Remotes"])
	if list != nil {
		rs := make([]*cfprototype.DialPoolPrototype, 0)
		for _, ri := range this.Remotes {
			rim := valutil.ToStringMap(ri)
			if rim != nil {
				p := new(cfprototype.DialPoolPrototype)
				p.FromMap(rim)
				rs = append(rs, p)
			}
		}
		this.Remotes = rs
	}
	return this.Valid()
}
