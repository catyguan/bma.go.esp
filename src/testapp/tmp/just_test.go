package tmp

import (
	"fmt"
	"sync/atomic"
	"testing"
)

type A struct {
	a int
}

func (this *A) print(v interface{}) {
	fmt.Printf("%p, %p\n", this, v)
}

func (this *A) p2(v interface{}) {
	this.print(v)
}

type B struct {
	c int
	A
	b int
}

func Test1(t *testing.T) {
	o := new(B)
	o.p2(o)
}

func Test2(t *testing.T) {
	var v interface{}
	v = Test1
	t.Errorf("%T", v)
}

func Test3(t *testing.T) {
	v := make([]int, 3)
	t.Errorf("%d", len(v[:3-1]))
}

func Test4(t *testing.T) {
	var v uint32
	v = 0xFFFFFFFF
	t.Errorf("%v", v)
	atomic.AddUint32(&v, 1)
	t.Errorf("%v", v)
}
