package memserv

import (
	"bmautil/memblock"
	"bmautil/valutil"
	"time"
)

type MemGoI struct {
	mg *MemGo
}

func (this *MemGoI) MemBlock() *memblock.MemBlock {
	return this.mg.mem
}

func msize(size int) int32 {
	return int32(memblock.ItemSize + size)
}

func (this *MemGoI) Set(key string, val interface{}, timeoutMS int) error {
	nv, size := MemGoData(val)
	this.mg.mem.Put(key, nv, msize(size), timeoutMS)
	return nil
}

func (this *MemGoI) Get(key string, tm *time.Time) (bool, interface{}, error) {
	var v interface{}
	var b bool
	if tm == nil {
		v, b = this.mg.mem.Get(key)
	} else {
		v, b = this.mg.mem.GetWithTimeout(key, *tm)
	}
	if !b {
		return false, nil, nil
	}
	return true, v, nil
}

func (this *MemGoI) Remove(key string) error {
	this.mg.mem.Remove(key, false)
	return nil
}

func (this *MemGoI) Touch(key string, timeoutMS int) error {
	this.mg.mem.Touch(key, timeoutMS)
	return nil
}

func (this *MemGoI) Put(key string, val interface{}, timeoutMS int) (bool, error) {
	b, _, err := this.Get(key, nil)
	if err != nil {
		return false, err
	}
	if b {
		return false, nil
	}
	err2 := this.Set(key, val, timeoutMS)
	if err2 != nil {
		return false, err2
	}
	return true, nil
}

func (this *MemGoI) Incr(key string, num int64, timeoutMS int) (int64, error) {
	_, v, err := this.Get(key, nil)
	if err != nil {
		return 0, err
	}
	val := valutil.ToInt64(v, 0)
	val = val + num
	err2 := this.Set(key, val, timeoutMS)
	if err2 != nil {
		return 0, err2
	}
	return val, nil
}
