package memserv

import (
	"bmautil/goo"
	"bmautil/memblock"
	"fmt"
	"time"
)

const (
	tag = "memserv"
)

type MBFuc func(mb *memblock.MemBlock) error

type MemGoConfig struct {
	QSize         int
	Max           int
	ClearStep     int
	ClearDuration time.Duration
}

func (this *MemGoConfig) Valid() error {
	if this.QSize <= 0 {
		this.QSize = 8
	}
	if this.ClearStep <= 0 {
		this.ClearStep = 100
	}
	if this.ClearDuration == 0 {
		this.ClearDuration = 5 * time.Millisecond
	}
	return nil
}

type MemGo struct {
	mem   *memblock.MemBlock
	goo   goo.Goo
	timer *time.Ticker
	tstop chan bool
	cfg   *MemGoConfig
}

func NewMemGo(cfg *MemGoConfig) *MemGo {
	cfg.Valid()

	r := new(MemGo)
	r.cfg = cfg
	r.mem = memblock.New()
	r.mem.MaxCount = cfg.Max
	r.goo.InitGoo(tag, cfg.QSize, r.doExit)
	return r
}

func (this *MemGo) Start() error {
	if !this.goo.Run() {
		return fmt.Errorf("goo run fail")
	}
	this.timer = time.NewTicker(this.cfg.ClearDuration)
	this.tstop = make(chan bool)
	go func() {
		for {
			this.goo.GetState()
			select {
			case <-this.timer.C:
			case <-this.tstop:
				return
			}
			err := this.goo.DoSync(func() {
				this.mem.Clear(this.cfg.ClearStep)
			})
			if err != nil {
				this.timer.Stop()
				return
			}
		}
	}()
	return nil
}

func (this *MemGo) Stop() {
	this.goo.Stop()
}

func (this *MemGo) doExit() {
	if this.timer != nil {
		this.timer.Stop()
		this.tstop <- true
	}
	this.mem.ClearAll(true)
}

func (this *MemGo) DoSync(f MBFuc) error {
	return this.goo.DoSync(func() error {
		return f(this.mem)
	})
}

func (this *MemGo) DoNow(f MBFuc) error {
	return this.goo.DoNow(func() error {
		return f(this.mem)
	})
}

func (this *MemGo) Size() (int, int32) {
	var c int
	var c2 int32
	this.goo.DoNow(func() {
		c, c2 = this.mem.Size()
	})
	return c, c2
}
