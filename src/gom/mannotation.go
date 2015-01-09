package gom

import (
	"bytes"
	"fmt"
)

type MAnnotation struct {
	Name  string
	Value interface{}
}

func (this *MAnnotation) String() string {
	return fmt.Sprintf("%s=%v", this.Name, this.Value)
}

type MAnnotations struct {
	list []*MAnnotation
}

func (this *MAnnotations) Has(n string) bool {
	for _, a := range this.list {
		if a.Name == n {
			return true
		}
	}
	return false
}

func (this *MAnnotations) Get(n string) interface{} {
	for _, a := range this.list {
		if a.Name == n {
			return a.Value
		}
	}
	return nil
}

func (this *MAnnotations) All() []*MAnnotation {
	return this.list
}

func (this *MAnnotations) List(n string) []interface{} {
	r := make([]interface{}, 0)
	for _, a := range this.list {
		if a.Name == n {
			r = append(r, a.Value)
		}
	}
	return r
}

func (this MAnnotations) Dump(buf *bytes.Buffer, prex string) {
	for _, a := range this.list {
		buf.WriteString(prex)
		buf.WriteString(fmt.Sprintf("@%s\n", a))
	}
}
