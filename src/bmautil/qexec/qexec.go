package qexec

import (
	"bmautil/syncutil"
	"bmautil/valutil"
	"errors"
	"logger"
	"runtime/debug"
	"time"
)

type Request struct {
	name     string
	data     interface{}
	callback func(err error)
}

func SyncCallback(event chan error) func(err error) {
	return func(err error) {
		if event != nil {
			defer func() {
				recover()
			}()
			event <- err
		}
	}
}

func NewRequest(name string, data interface{}, cb func(err error)) Request {
	return Request{name, data, cb}
}

type RequestHandler func(req interface{}) (bool, error)
type ErrorHandler func(req interface{}, err error) bool
type StopHandler func()

type QueueExecutor struct {
	Tag        string
	EDebug     bool
	requests   chan Request
	oldQueue   chan Request
	closeState *syncutil.CloseState
	started    bool

	// handler
	RequestHandler RequestHandler
	ErrorHandler   ErrorHandler
	StopHandler    StopHandler
}

func safe(c func()) {
	defer func() {
		recover()
	}()
	c()
}

func NewQueueExecutor(tag string, bufsize int,
	rhandler RequestHandler) *QueueExecutor {
	r := new(QueueExecutor)
	r.InitQueueExecutor(tag, bufsize, rhandler)
	return r
}

func (this *QueueExecutor) InitQueueExecutor(tag string, bufsize int,
	rhandler RequestHandler) {
	this.Tag = tag
	if bufsize > 0 {
		this.InitRequests(bufsize)
	}
	this.closeState = syncutil.NewCloseState()
	this.RequestHandler = rhandler
}

func (this *QueueExecutor) InitRequests(sz int) {
	this.requests = make(chan Request, sz)
}

func (this *QueueExecutor) Execute(req *Request) (running bool, err error) {
	running = true
	defer func() {
		ex := recover()
		if ex != nil {
			if _, ok := ex.(error); ok {
				err = ex.(error)
			} else {
				err = errors.New(valutil.ToString(ex, "unknow error"))
			}
		}
		if err != nil {
			logger.Error(this.Tag, "execte fail - %s\n%s", err.Error(), string(debug.Stack()))
			if this.ErrorHandler != nil {
				safe(func() {
					running = this.ErrorHandler(req.data, err)
				})
			}
		}
		if req.callback != nil {
			req.callback(err)
		}
		if this.EDebug {
			logger.Debug(this.Tag, "request done - %s", req.name)
		}
	}()
	return this.RequestHandler(req.data)
}

func (this *QueueExecutor) run() {
	defer func() {
		if this.EDebug {
			logger.Debug(this.Tag, "stop")
		}
		if this.StopHandler != nil {
			safe(this.StopHandler)
		}
		close(this.requests)
		if this.oldQueue != nil {
			close(this.oldQueue)
			this.oldQueue = nil
		}
		this.closeState.DoneClose()
		this.started = false
	}()
	if this.EDebug {
		logger.Debug(this.Tag, "run queue executor")
	}
	for {
		var req Request
		if this.oldQueue != nil {
			select {
			case req = <-this.oldQueue:
				if req.name == "__CLOSE__" {
					logger.Info(this.Tag, "close old queue")
					close(this.oldQueue)
					this.oldQueue = nil
					continue
				}
			}
		} else {
			req = <-this.requests
		}
		if req.data == nil {
			if this.EDebug {
				logger.Debug(this.Tag, "receive nil request, stop")
			}
			return
		}
		if this.EDebug {
			logger.Debug(this.Tag, "popup new request - %s", req.name)
		}
		if r, _ := this.Execute(&req); !r {
			if this.EDebug {
				logger.Debug(this.Tag, "receive stop request")
			}
			return
		}
	}
}

func (this *QueueExecutor) IsRun() bool {
	return this.started
}

func (this *QueueExecutor) Run() bool {
	if this.requests == nil {
		logger.Error(this.Tag, "requests queue is nil")
		return false
	}
	this.started = true
	go this.run()
	return true
}

func (this *QueueExecutor) Do(name string, req interface{}, cb func(err error)) error {
	if this.closeState.IsClosing() {
		err := errors.New(this.Tag + " closed")
		if cb != nil {
			cb(err)
		}
		return err
	}
	if this.requests == nil {
		err := errors.New(this.Tag + " requests queue is nil")
		if cb != nil {
			cb(err)
		}
		return err
	}
	this.requests <- Request{name, req, cb}
	if this.EDebug {
		logger.Debug(this.Tag, "push new request - %s", name)
	}
	return nil
}

func (this *QueueExecutor) DoNow(name string, req interface{}) error {
	return this.Do(name, req, nil)
}

func (this *QueueExecutor) DoSync(name string, req interface{}) error {
	ev := make(chan error, 1)
	defer close(ev)
	this.Do(name, req, SyncCallback(ev))
	err := <-ev
	return err
}

func (this *QueueExecutor) DoTimeout(name string, req interface{}, timeout time.Duration) error {
	ev := make(chan error, 1)
	defer close(ev)
	tm := time.NewTimer(timeout)
	defer tm.Stop()
	this.Do(name, req, SyncCallback(ev))
	select {
	case err := <-ev:
		return err
	case <-tm.C:
		return errors.New("timeout")
	}
}

func (this *QueueExecutor) IsClosing() bool {
	return this.closeState.IsClosing()
}

func (this *QueueExecutor) Stop() bool {
	if this.closeState.AskClose() && this.requests != nil {
		this.requests <- Request{"STOP", nil, nil}
	}
	return true
}

func (this *QueueExecutor) WaitStop() bool {
	if this.closeState.IsClosing() && this.started {
		this.closeState.WaitClosed()
	}
	return true
}

func (this *QueueExecutor) Resize(newbufsize int) bool {
	if this.oldQueue != nil {
		return false
	}
	logger.Info(this.Tag, "resize to %d", newbufsize)
	if this.requests != nil {
		this.oldQueue = this.requests
	}
	this.InitRequests(newbufsize)
	q := this.oldQueue
	go func() {
		q <- Request{"__CLOSE__", nil, nil}
	}()
	return true
}
