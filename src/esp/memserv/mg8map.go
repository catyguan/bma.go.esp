package memserv

import (
	"bmautil/memblock"
	"time"
)

func (this *MemGo) Get(key string, tm *time.Time) (interface{}, error) {
	var rv interface{}
	err := this.goo.DoSync(func() error {
		var v interface{}
		var b bool
		if tm == nil {
			v, b = this.mem.Get(key)
		} else {
			v, b = this.mem.GetWithTimeout(key, *tm)
		}
		if !b {
			return nil
		}
		rv = v
		return nil
	})
	return rv, err
}

func (this *MemGo) Set(key string, val interface{}, timeoutMS int) error {
	return this.goo.DoSync(func() {
		nv, size := MemGoData(val)
		this.mem.Put(key, nv, int32(memblock.ItemSize+size), timeoutMS)
	})
	return nil
}

func (this *MemGo) Remove(key string) error {
	return this.goo.DoSync(func() {
		this.mem.Remove(key, false)
	})
	return nil
}
