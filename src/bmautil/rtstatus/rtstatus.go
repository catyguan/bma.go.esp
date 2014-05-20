package rtstatus

import (
	"bytes"
	"fmt"
)

type RTSProvider interface {
	WriteRuntimeStatus(stc *RTSCollections) bool
}

type RTStatus struct {
	value interface{}
}

func (this *RTStatus) Value() interface{} {
	return this.value
}

func (this *RTStatus) IsCollections() bool {
	if this.value == nil {
		return false
	}
	_, ok := this.value.(*RTSCollections)
	return ok
}

func (this *RTStatus) AsCollections() *RTSCollections {
	if this.value == nil {
		return nil
	}
	if r, ok := this.value.(*RTSCollections); ok {
		return r
	}
	return nil
}

type RTSCollections struct {
	statusList map[string]*RTStatus
}

func NewCollections() *RTSCollections {
	this := new(RTSCollections)
	this.statusList = make(map[string]*RTStatus)
	return this
}

func (this *RTSCollections) Set(n string, val interface{}) {
	st := new(RTStatus)
	st.value = val
	this.statusList[n] = st
}

func (this *RTSCollections) Remove(n string) {
	delete(this.statusList, n)
}

func (this *RTSCollections) SubCollections(n string) *RTSCollections {
	if ost, ok := this.statusList[n]; ok {
		col := ost.AsCollections()
		if col != nil {
			return col
		}
	}
	col := NewCollections()
	st := new(RTStatus)
	st.value = col
	this.statusList[n] = st
	return col
}

func (this *RTSCollections) Names() []string {
	r := make([]string, 0, len(this.statusList))
	for k, _ := range this.statusList {
		r = append(r, k)
	}
	return nil
}

func (this *RTSCollections) Get(n string) *RTStatus {
	if st, ok := this.statusList[n]; ok {
		return st
	}
	return nil
}

func (this *RTSCollections) Print(fm RTSFormatter) {
	for k, st := range this.statusList {
		col := st.AsCollections()
		if col == nil {
			fm.PrintValue(k, st.Value())
		} else {
			fm.NewCollections(k)
			col.Print(fm)
			fm.EndCollections()
		}
	}
}

type RTSFormatter interface {
	NewCollections(name string)
	PrintValue(name string, val interface{})
	EndCollections()
}

type RTSFormatter4String struct {
	Tab              string
	Newline          string
	Buffer           *bytes.Buffer
	FormatCollection func(n string) string
	FormatValue      func(n string, v interface{}) string
	ntab             int
}

func (this *RTSFormatter4String) Init() {
	if this.Tab == "" {
		this.Tab = "\t"
	}
	if this.Newline == "" {
		this.Newline = "\n"
	}
	if this.Buffer == nil {
		this.Buffer = bytes.NewBuffer([]byte{})
	}
	if this.FormatCollection == nil {
		this.FormatCollection = this.DefaultFormatCollection
	}
	if this.FormatValue == nil {
		this.FormatValue = this.DefaultFormatValue
	}
}

func (this *RTSFormatter4String) ptab() {
	for i := 0; i < this.ntab; i++ {
		this.Buffer.WriteString(this.Tab)
	}
}

func (this *RTSFormatter4String) pln() {
	this.Buffer.WriteString(this.Newline)
}

func (this *RTSFormatter4String) DefaultFormatCollection(n string) string {
	return fmt.Sprintf("%s:", n)
}

func (this *RTSFormatter4String) NewCollections(name string) {
	this.ptab()
	this.Buffer.WriteString(this.FormatCollection(name))
	this.pln()
	this.ntab = this.ntab + 1
}

func (this *RTSFormatter4String) DefaultFormatValue(name string, val interface{}) string {
	return fmt.Sprintf("%s: %v", name, val)
}

func (this *RTSFormatter4String) PrintValue(name string, val interface{}) {
	this.ptab()
	this.Buffer.WriteString(this.FormatValue(name, val))
	this.pln()
}

func (this *RTSFormatter4String) EndCollections() {
	this.ntab = this.ntab - 1
}

func (this *RTSFormatter4String) String() string {
	if this.Buffer == nil {
		return ""
	}
	return this.Buffer.String()
}
