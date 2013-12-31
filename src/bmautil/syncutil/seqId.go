package syncutil

import (
	"sync/atomic"
)

type SeqId struct {
	id         uint64
	MaxValue   uint64
	ResetValue uint64
}

func (this *SeqId) Next() uint64 {
	mv := this.MaxValue
	if mv == 0 {
		mv = 100000000
		this.MaxValue = mv
	}
	for {
		v := atomic.AddUint64(&this.id, 1)
		if v < mv {
			return v
		}
		atomic.CompareAndSwapUint64(&this.id, v, this.ResetValue)
	}
}
