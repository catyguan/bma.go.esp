package qpushpull

import "fmt"

const (
	tag   = "qpushpull"
	debug = false
)

type Handler func(req interface{})

type QueuePushPull struct {
	c1      chan interface{}
	buff    []interface{}
	c2      chan interface{}
	handler Handler
	closing bool
	closed  chan bool
}

type dopull int

func NewQueuePushPull(qsz int, h Handler) *QueuePushPull {
	this := new(QueuePushPull)
	this.c1 = make(chan interface{}, qsz)
	this.c2 = make(chan interface{}, qsz+1)
	this.buff = make([]interface{}, 0)
	this.handler = h
	return this
}

func (this *QueuePushPull) Run() {
	this.closed = make(chan bool, 1)
	go this.run1()
	go this.run2()
}

func (this *QueuePushPull) run1() {
	for {
		req := <-this.c1
		if req == nil {
			for _, br := range this.buff {
				this.c2 <- br
			}
			close(this.c2)
			return
		}
		if _, ok := req.(dopull); ok {
			if debug {
				fmt.Println(tag, "c2 pull")
			}
			c := cap(this.c2) - len(this.c2)
			if c > len(this.buff) {
				c = len(this.buff)
			}
			for _, bv := range this.buff[:c] {
				this.c2 <- bv
			}
			copy(this.buff, this.buff[c:])
			this.buff = this.buff[:len(this.buff)-c]
			continue
		}
		if len(this.buff) > 0 {
			this.buff = append(this.buff, req)
			continue
		}
		if len(this.c2) >= cap(this.c2) {
			if debug {
				fmt.Println(tag, "c2 full, wait pull")
			}
			this.buff = append(this.buff, req)
			continue
		}
		this.c2 <- req
	}
}

func (this *QueuePushPull) run2() {
	for {
		if len(this.c2) == 0 {
			if debug {
				fmt.Println(tag, "c2 send pull")
			}
			func() {
				defer func() {
					recover()
				}()
				this.c1 <- dopull(0)
			}()
		}
		req := <-this.c2
		if req == nil {
			close(this.closed)
			return
		}
		this.handler(req)
	}
}

func (this *QueuePushPull) Close() {
	if this.closing {
		return
	}
	this.closing = true
	close(this.c1)
}

func (this *QueuePushPull) WaitClose() {
	if this.closed != nil {
		<-this.closed
	}
}

func (this *QueuePushPull) IsClosing() bool {
	return this.closing
}

func (this *QueuePushPull) Push(req interface{}) error {
	if this.closing {
		return fmt.Errorf("closed")
	}
	this.c1 <- req
	return nil
}
