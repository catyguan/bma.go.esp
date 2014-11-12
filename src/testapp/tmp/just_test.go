package tmp

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
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

func Test5(t *testing.T) {
	var bs []byte
	for _, v := range bs {
		t.Error(v)
	}
	t.Error(len(bs))
}

func TestGenRandome(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 1000; i++ {
		if i != 0 {
			fmt.Print(",")
		}
		fmt.Printf("%d", r.Uint32()%0x7FFFFF)
	}
	fmt.Println()
}

func TestSprintf(t *testing.T) {
	arr := []interface{}{1, 2}
	s := fmt.Sprintf("%d, %d", arr...)
	fmt.Println("ask = ", s)
}

func TestTime(t *testing.T) {
	tm := time.Now()
	bs, _ := tm.MarshalText()
	fmt.Println("TimeJson", string(bs))
}

func TestSliceAppend(t *testing.T) {
	var a []string
	a = append(a, "test")
	a = append(a, "test2")
	fmt.Println(a)
}
