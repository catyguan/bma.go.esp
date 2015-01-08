package goyacc

import "fmt"

type Annotation struct {
	Name  string
	Value Node
}

func (this *Annotation) String() string {
	return fmt.Sprintf("%s:%s", this.Name, DumpNode("", this.Value))
}

type Annotations []*Annotation

func (this Annotations) Has(n string) bool {
	for _, a := range this {
		if a.Name == n {
			return true
		}
	}
	return false
}

func (this Annotations) Get(n string) Node {
	for _, a := range this {
		if a.Name == n {
			return a.Value
		}
	}
	return nil
}

func (this Annotations) List(n string) []Node {
	r := make([]Node, 0)
	for _, a := range this {
		if a.Name == n {
			r = append(r, a.Value)
		}
	}
	return r
}
