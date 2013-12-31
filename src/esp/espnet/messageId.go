package espnet

import (
	"sync/atomic"
	"time"
)

type MessageIdGenerator struct {
	seed uint64
	syn  bool
}

func (this *MessageIdGenerator) InitMessageIdGenerator() {
	this.seed = uint64(time.Now().UnixNano())
}

func (this *MessageIdGenerator) Next() uint64 {
	var r uint64
	if this.syn {
		r = atomic.AddUint64(&this.seed, 1)
	} else {
		this.seed++
		r = this.seed
	}
	return r
}
