package msgagent

import (
	"errors"
	"esp/espnet"
	"logger"
	"sync"
	"time"
)

const (
	tag = "msgAgent"
)

var (
	DefaultAgent *Agent
)

func init() {
	lock := &sync.RWMutex{}
	DefaultAgent = NewAgent(lock, lock.RLocker())
	DefaultAgent.Timeout = time.Duration(5) * time.Second
}

type task struct {
	msg      *espnet.Message
	deadTime time.Time
	resp     chan *espnet.Message
}

type Agent struct {
	Timeout time.Duration

	wlock sync.Locker
	rlock sync.Locker
	tasks map[uint64]*task
}

func NewAgent(wlock, rlock sync.Locker) *Agent {
	this := new(Agent)
	this.wlock = wlock
	this.rlock = rlock
	this.tasks = make(map[uint64]*task)
	return this
}

func S(a *Agent) *Agent {
	if a == nil {
		// TODO
		return DefaultAgent
	}
	return a
}

func (this *Agent) removeTask(mid uint64) *task {
	if this.wlock != nil {
		this.wlock.Lock()
		defer this.wlock.Unlock()
	}
	r, ok := this.tasks[mid]
	if ok {
		return r
	}
	return nil
}

func (this *Agent) SendMessage(ch espnet.Channel, msg *espnet.Message, timeout time.Duration) (*espnet.Message, error) {

	to := timeout
	if to <= 0 {
		to = this.Timeout
	}

	mid := espnet.FrameCoders.MessageId.Sure(msg.ToPackage())
	task := new(task)
	task.msg = msg
	task.resp = make(chan *espnet.Message, 1)
	defer close(task.resp)
	if timeout > 0 {
		task.deadTime = time.Now().Add(to)
	}

	func() {
		if this.wlock != nil {
			this.wlock.Lock()
			defer this.wlock.Unlock()
		}
		this.tasks[mid] = task
	}()
	defer this.removeTask(mid)

	ch.SetMessageListner(this.receiveMessage)
	defer ch.SetMessageListner(nil)

	logger.Debug(tag, "send message(%d)", mid)
	err := ch.SendMessage(msg)
	if err != nil {
		return nil, err
	}

	tm := time.NewTimer(to)
	select {
	case r := <-task.resp:
		tm.Stop()
		if r != nil {
			return r, nil
		}
	case <-tm.C:
	}
	return nil, errors.New("timeout")
}

func (this *Agent) receiveMessage(msg *espnet.Message) error {
	mid := espnet.FrameCoders.SourceMessageId.Get(msg.ToPackage())
	var task *task
	if mid > 0 {
		logger.Debug(tag, "receive message(%d) response", mid)
		func() {
			if this.rlock != nil {
				this.rlock.Lock()
				defer this.rlock.Unlock()
				r, ok := this.tasks[mid]
				if ok {
					task = r
				}
			}
		}()
	}
	if task != nil {
		defer func() {
			recover()
		}()
		task.resp <- msg
	} else {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "discard unknow message - %s", msg.Dump())
		}
	}
	return nil
}
