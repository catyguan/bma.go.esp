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

type MBFuc func(mgi *MemGoI) error

type MemGoConfig struct {
	QSize     int
	Max       int
	ClearStep int
	ClearMS   int
}

func (this *MemGoConfig) Valid() error {
	if this.QSize <= 0 {
		this.QSize = 8
	}
	if this.ClearStep <= 0 {
		this.ClearStep = 100
	}
	if this.ClearMS == 0 {
		this.ClearMS = 5
	}
	return nil
}

func (this *MemGoConfig) Compare(old *MemGoConfig) bool {
	if this.QSize != old.QSize {
		return false
	}
	if this.Max != old.Max {
		return false
	}
	if this.ClearStep != old.ClearStep {
		return false
	}
	if this.ClearMS != old.ClearMS {
		return false
	}
	return true
}

var (
	DEFAULT_CONFIG *MemGoConfig
)

func init() {
	DEFAULT_CONFIG = new(MemGoConfig)
	DEFAULT_CONFIG.QSize = 128
	DEFAULT_CONFIG.Valid()
}

type MemGo struct {
	name    string
	mem     *memblock.MemBlock
	goo     goo.Goo
	timer   *time.Ticker
	tstop   chan bool
	cfg     *MemGoConfig
	scaners map[string]*memblock.MapItem
	RelEnv  interface{}
	RelAttr map[string]interface{}
}

func NewMemGo(n string, cfg *MemGoConfig) *MemGo {
	cfg.Valid()

	r := new(MemGo)
	r.name = n
	r.cfg = cfg
	r.mem = memblock.New()
	r.mem.Listener = r.memListener
	r.mem.MaxCount = cfg.Max
	r.goo.InitGoo(tag, cfg.QSize, r.doExit)
	r.RelAttr = make(map[string]interface{})
	return r
}

func (this *MemGo) memListener(k string, item *memblock.MapItem, rt memblock.REMOVE_TYPE) {
	if rt == memblock.RT_CLOSE {
		return
	}
	for k, pos := range this.scaners {
		if pos == item {
			this.scaners[k] = item.Next()
		}
	}
	fmt.Println("remove", k, item.Data, rt)
}

func (this *MemGo) Start() error {
	if !this.goo.Run() {
		return fmt.Errorf("goo run fail")
	}
	this.timer = time.NewTicker(time.Duration(this.cfg.ClearMS) * time.Millisecond)
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
	for k, _ := range this.scaners {
		delete(this.scaners, k)
	}
	this.mem.CloseClear(false)
}

func (this *MemGo) DoSync(f MBFuc) error {
	return this.goo.DoSync(func() error {
		mgi := MemGoI{this}
		return f(&mgi)
	})
}

func (this *MemGo) DoNow(f MBFuc) error {
	return this.goo.DoNow(func() error {
		mgi := MemGoI{this}
		return f(&mgi)
	})
}

func (this *MemGo) Size() (int, int32) {
	var c int
	var c2 int32
	this.goo.DoSync(func() {
		c, c2 = this.mem.Size()
	})
	return c, c2
}

func (this *MemGo) BeginScan(scanName string) error {
	return this.goo.DoSync(func() {
		if this.scaners == nil {
			this.scaners = make(map[string]*memblock.MapItem)
		}
		this.scaners[scanName] = this.mem.Head()
	})
}

// return isEnd, error
func (this *MemGo) Scan(scanName string, count int, f func(k string, v interface{})) (bool, error) {
	rb := true
	err := this.goo.DoSync(func() {
		var pos *memblock.MapItem
		if this.scaners != nil {
			pos = this.scaners[scanName]
		}
		for i := 0; i < count; i++ {
			if pos == nil {
				break
			}
			f(pos.Key, pos.Data)
			pos = pos.Next()
		}
		if pos == nil {
			if this.scaners != nil {
				delete(this.scaners, scanName)
			}
			return
		}
		if this.scaners == nil {
			this.scaners = make(map[string]*memblock.MapItem)
		}
		this.scaners[scanName] = pos
		rb = false
	})
	return rb, err
}

func (this *MemGo) EndScan(scanName string) error {
	return this.goo.DoSync(func() {
		if this.scaners != nil {
			delete(this.scaners, scanName)
		}
	})
}
