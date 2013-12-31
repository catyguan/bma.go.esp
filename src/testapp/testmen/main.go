package main

import (
	"bmautil/syncutil"
	"runtime"
	// "sync/atomic"
	"time"
)

/*
var a, b int
var l int32

func f() {
	a = 1
	b = 2
	atomic.StoreInt32(&l, 9)
}

func g() {
	v := atomic.LoadInt32(&l)
	print(b)
	print(a)
	print(v)
}

func main() {
	runtime.GOMAXPROCS(3)

	// go g()
	go f()
	go g()

	time.Sleep(time.Duration(100) * time.Millisecond)
}
*/

var a, b syncutil.MemHolder

func f() {
	a.Set(1)
	b.Set(2)
}

func g() {
	print(b.Value())
	print(a.Value())
}

func main() {
	runtime.GOMAXPROCS(3)

	// go g()
	go f()
	go g()

	time.Sleep(time.Duration(100) * time.Millisecond)
}
