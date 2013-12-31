package uprop

import (
	"bytes"
	"errors"
	"fmt"
)

type UProperty struct {
	Name       string
	Value      interface{}
	IsOptional bool
	Tips       string
	Setter     func(v string) error
}

func NewUProperty(name string, val interface{}, def bool, tips string, setter func(v string) error) *UProperty {
	r := new(UProperty)
	r.Name = name
	r.Value = val
	r.IsOptional = def
	r.Tips = tips
	r.Setter = setter
	return r
}

func (this *UProperty) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	if this.IsOptional {
		buf.WriteString("*")
	}
	buf.WriteString(fmt.Sprintf("%s=%v", this.Name, this.Value))
	if this.Tips != "" {
		buf.WriteString(" (")
		buf.WriteString(this.Tips)
		buf.WriteString(")")
	}
	return buf.String()
}

func (this *UProperty) Commit(val string) (rerr error) {
	if this.Setter == nil {
		return errors.New("can't commit")
	}
	defer func() {
		err := recover()
		if err != nil {
			if re, ok := err.(error); ok {
				rerr = re
			} else {
				rerr = errors.New(fmt.Sprintf("%s", err))
			}
		}
	}()
	err := this.Setter(val)
	if err != nil {
		return err
	}
	return nil
}
