package lua51

import (
	"fmt"
	"testing"
	"unsafe"
)

type Userdata struct {
	a, b int
}

func TestUserdata(t *testing.T) {
	var L *State
	L = NewState()
	L.OpenLibs()

	rawptr := L.NewUserdata(uintptr(unsafe.Sizeof(Userdata{})))
	var ptr *Userdata
	ptr = (*Userdata)(rawptr)
	ptr.a = 2
	ptr.b = 3

	fmt.Println(ptr)

	rawptr2 := L.ToUserdata(-1)
	ptr2 := (*Userdata)(rawptr2)

	fmt.Println(ptr2)
}
