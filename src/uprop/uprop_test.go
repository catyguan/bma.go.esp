package uprop

import (
	"testing"
)

func TestFind(t *testing.T) {
	obj := new(UPropertyStruct)
	obj.NewProp("f1", "f1").BeValue(1, nil)
	obj.NewProp("f2", "f2").BeValue(2, nil)
	p3 := obj.NewProp("f3", "f3").BeList(nil, nil)
	p3.Add(100, nil)
	p3.Add(110, nil)

	child := new(UPropertyStruct)
	child.NewProp("c1", "c1").BeValue("X1", nil)
	obj.NewProp("f4", "f4").BeMap(child.AsList(), nil)

	if true {
		r1, r2 := Find(obj.AsList(), []string{"f3", "lastidx"})
		t.Error(r1, r2)
	}
	if true {
		r1, r2 := Find(obj.AsList(), []string{"f4", "c1"})
		t.Error(r1, r2)
	}
}
