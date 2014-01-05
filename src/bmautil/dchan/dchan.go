package dchan

import "time"

var (
	empty     chan interface{} = make(chan interface{}, 0)
	closeItem oldItem
)

type oldItem struct {
	c  chan interface{}
	tm *time.Timer
}

type Chan struct {
	old  []*oldItem
	main chan interface{}
}

func NewDChan(sz int) *Chan {
	this := new(Chan)
	this.main = make(chan interface{}, sz)
	return this
}

func (this *Chan) CloseDChan() {
	defer func() {
		recover()
	}()
	close(this.main)
}

func (this *Chan) Write(v interface{}) {
	this.Send(v)
}

func (this *Chan) Send(v interface{}) {
	this.main <- v
}

func (this *Chan) Read(wait chan interface{}) (interface{}, interface{}) {
	r1, r2 := this.doRead(wait)
	return r1, r2
}

func mulRead(c1, c2 chan interface{}, c3 <-chan time.Time) (interface{}, interface{}, bool) {
	if c2 == nil {
		c2 = empty
	}
	if c3 != nil {
		select {
		case r1 := <-c1:
			return r1, nil, false
		case r2 := <-c2:
			return nil, r2, false
		case <-c3:
			return nil, nil, true
		}
	} else {
		select {
		case r1 := <-c1:
			return r1, nil, false
		case r2 := <-c2:
			return nil, r2, false
		}
	}
}

func (this *Chan) doRead(wait chan interface{}) (interface{}, interface{}) {
	if this.old != nil && len(this.old) > 0 {
		for {
			if len(this.old) == 0 {
				break
			}
			item := this.old[0]
			r1, r2, r3 := mulRead(item.c, wait, item.tm.C)
			if !r3 {
				return r1, r2
			}
			this.old = this.old[1:]
			close(item.c)
		}
	}
	r1, r2, _ := mulRead(this.main, wait, nil)
	return r1, r2
}

func (this *Chan) DoResize(newbufsize int) {
	if this.old == nil {
		this.old = make([]*oldItem, 0)
	}
	// fmt.Println("remain", len(this.main))
	o := this.main
	this.main = make(chan interface{}, newbufsize)
	this.old = append(this.old, &oldItem{o, time.NewTimer(1 * time.Millisecond)})
}
