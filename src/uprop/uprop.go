package uprop

import (
	"bmautil/valutil"
	"bytes"
	"errors"
	"fmt"
)

type UPropertyValueKind int

const (
	UPK_VALUE = UPropertyValueKind(0)
	UPK_LIST  = UPropertyValueKind(1)
	UPK_MAP   = UPropertyValueKind(3)
	UPK_FOLD  = UPropertyValueKind(4)
)

type UPropertyValue struct {
	Kind     UPropertyValueKind
	Value    interface{}
	Setter   func(v string) error
	Expender func() []*UProperty
}

func ToIndex(n string, l int) int {
	if n == "lastidx" {
		return l
	} else {
		return valutil.ToInt(n, 0)
	}
}

func (this *UPropertyValue) Find(ns []string) (*UProperty, *UPropertyValue) {
	if len(ns) == 0 {
		return nil, this
	}
	switch this.Kind {
	case UPK_LIST:
		if this.Value != nil {
			vals := this.Value.([]*UPropertyValue)
			n := ns[0]
			idx := ToIndex(n, len(vals))
			if idx > 0 && idx-1 < len(vals) {
				val := vals[idx-1]
				return val.Find(ns[1:])
			}
			return nil, nil
		}
	case UPK_MAP:
		if this.Value != nil {
			props := this.Value.([]*UProperty)
			return Find(props, ns)
		}
	}
	return nil, nil
}

func (this *UPropertyValue) Commit(val string) (rerr error) {
	if this.Setter == nil {
		return errors.New("can't commit")
	}
	err := this.Setter(val)
	if err != nil {
		return err
	}
	return nil
}

// UProperty
type UProperty struct {
	Name       string
	Tips       string
	IsOptional bool
	Value      *UPropertyValue

	Adder   func(vlist []string) error
	Remover func(vlist []string) error
}

func NewUProperty2(name string, tips string) *UProperty {
	r := new(UProperty)
	r.Name = name
	r.IsOptional = true
	r.Tips = tips
	return r
}

func NewUProperty(name string, val interface{}, def bool, tips string, setter func(v string) error) *UProperty {
	return NewUProperty2(name, tips).Optional(def).BeValue(val, setter)
}

func Find(props []*UProperty, ns []string) (*UProperty, *UPropertyValue) {
	if len(ns) == 0 {
		return nil, nil
	}
	for _, p := range props {
		rp, rv := p.Find(ns)
		if rp != nil {
			if rv == nil {
				rv = rp.Value
			}
			return rp, rv
		}
	}
	return nil, nil
}

func (this *UProperty) Find(ns []string) (*UProperty, *UPropertyValue) {
	if len(ns) == 0 {
		return this, nil
	}
	n := ns[0]
	if n == this.Name {
		if this.Value != nil {
			r, v := this.Value.Find(ns[1:])
			if v != nil {
				if r == nil {
					r = this
				}
				return r, v
			}
		}
	}
	return nil, nil
}

func (this *UProperty) Kind() UPropertyValueKind {
	if this.Value != nil {
		return this.Value.Kind
	}
	return UPK_VALUE
}

func (this *UProperty) Optional(v bool) *UProperty {
	this.IsOptional = v
	return this
}

func (this *UProperty) BeValue(val interface{}, setter func(v string) error) *UProperty {
	vo := new(UPropertyValue)
	this.Value = vo

	vo.Kind = UPK_VALUE
	vo.Value = val
	vo.Setter = setter
	return this
}

func (this *UProperty) BeList(adder func(vlist []string) error, remover func(vlist []string) error) *UProperty {
	vo := new(UPropertyValue)
	this.Value = vo
	this.Adder = adder
	this.Remover = remover

	vo.Kind = UPK_LIST
	vo.Value = []*UPropertyValue{}
	return this
}

func (this *UProperty) Add(val interface{}, setter func(v string) error) *UProperty {
	if this.Value == nil || this.Value.Kind != UPK_LIST {
		panic("prop mus BeList")
	}

	vo := new(UPropertyValue)
	vo.Kind = UPK_VALUE
	vo.Value = val
	vo.Setter = setter

	tv := this.Value
	tv.Value = append(tv.Value.([]*UPropertyValue), vo)
	return this
}

func (this *UProperty) AddFold(desc string, ep func() []*UProperty) *UProperty {
	if this.Value == nil || this.Value.Kind != UPK_LIST {
		panic("prop mus BeList")
	}

	vo := new(UPropertyValue)
	vo.Kind = UPK_FOLD
	vo.Value = desc
	vo.Expender = ep

	tv := this.Value
	tv.Value = append(tv.Value.([]*UPropertyValue), vo)
	return this
}

func (this *UProperty) BeMap(m []*UProperty, adder func(vlist []string) error, remover func(vlist []string) error) *UProperty {
	vo := new(UPropertyValue)
	this.Value = vo
	this.Adder = adder
	this.Remover = remover

	vo.Kind = UPK_MAP
	if m == nil {
		m = []*UProperty{}
	}
	vo.Value = m
	return this
}

func (this *UProperty) BeFold(desc string, ep func() []*UProperty) *UProperty {
	vo := new(UPropertyValue)
	this.Value = vo

	vo.Kind = UPK_FOLD
	vo.Value = desc
	vo.Expender = ep
	return this
}

func (this *UProperty) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	if this.Value != nil {
		switch this.Value.Kind {
		case UPK_VALUE:
			buf.WriteString(fmt.Sprintf("%s=%v", this.Name, this.Value.Value))
		case UPK_LIST:
			buf.WriteString(fmt.Sprintf("%s=[...]", this.Name))
		case UPK_MAP:
			buf.WriteString(fmt.Sprintf("%s={...}", this.Name))
		case UPK_FOLD:
			buf.WriteString(fmt.Sprintf("%s : %s", this.Name, this.Value.Value))
		}
	} else {
		buf.WriteString("<nil>")
	}
	return buf.String()
}

func (this *UProperty) CallSet(v string) error {
	if this.Value == nil {
		return errors.New("can't set")
	}
	return this.Value.Commit(v)
}

func (this *UProperty) CallAdd(vlist []string) (rerr error) {
	if this.Adder == nil {
		return errors.New("can't add")
	}
	err := this.Adder(vlist)
	if err != nil {
		return err
	}
	return nil
}

func (this *UProperty) CallRemove(vlist []string) (rerr error) {
	if this.Remover == nil {
		return errors.New("can't remove")
	}
	err := this.Remover(vlist)
	if err != nil {
		return err
	}
	return nil
}

// UPropertyBuilder
type UPropertyBuilder struct {
	properties []*UProperty
}

func (this *UPropertyBuilder) NewProp(name string, tips string) *UProperty {
	r := NewUProperty2(name, tips)
	this.AddProp(r)
	return r
}

func (this *UPropertyBuilder) AddProp(p *UProperty) {
	this.properties = append(this.AsList(), p)
}

func (this *UPropertyBuilder) Merge(list []*UProperty) {
	if list != nil {
		for _, p := range list {
			this.AddProp(p)
		}
	}
}

func (this *UPropertyBuilder) MergeWithPrex(list []*UProperty, prex string) {
	if list != nil {
		for _, p := range list {
			p.Name = prex + p.Name
			this.AddProp(p)
		}
	}
}

func (this *UPropertyBuilder) ToMapProp(name, tips string, adder, remover func(vs []string) error) *UProperty {
	return this.NewProp(name, tips).BeMap(this.AsList(), adder, remover)
}

func (this *UPropertyBuilder) AsList() []*UProperty {
	if this.properties == nil {
		this.properties = make([]*UProperty, 0)
	}
	return this.properties
}
