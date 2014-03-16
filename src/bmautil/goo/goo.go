package goo

import (
	"bmautil/valutil"
	"errors"
	"logger"
	"time"
)

type ExitHandler func()

type Goo struct {
	EDebug      bool
	TSM         StateMachine
	queue       chan interface{}
	requests    []interface{}
	exitHandler ExitHandler
}

func (this *Goo) InitGoo(queueSize int, exithandler ExitHandler) {
	this.queue = make(chan interface{}, queueSize)
	this.exitHandler = exithandler
}

func (this *Goo) Execute(req *Request) (running bool, err error) {
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
			logger.Error(this.Tag, "execte '%s' fail - %s", req.name, err.Error())
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

func (this *Goo) run() {
	defer func() {
		if this.EDebug {
			logger.Debug(this.Tag, "stop")
		}
		if this.StopHandler != nil {
			safe(this.StopHandler)
		}
		this.requests.CloseDChan()
		this.closeState.DoneClose()
		this.started = false
	}()
	if this.EDebug {
		logger.Debug(this.Tag, "run queue executor")
	}
	for {
		if len(this.requests) > 0 {
			select {
			case v := <-this.queue:
				if v != nil {
					this.requests = append(this.requests, v)
				}
			default:
			}
		} else {
			v := <-this.queue
			if v != nil {

			}
		}
		dreq, _ := this.requests.Read(nil)
		req := dreq.(*Request)

		if req.data == nil {
			if this.EDebug {
				logger.Debug(this.Tag, "receive nil request, stop")
			}
			return
		}
		if this.EDebug {
			logger.Debug(this.Tag, "popup new request - %s", req.name)
		}
		if r, _ := this.Execute(req); !r {
			if this.EDebug {
				logger.Debug(this.Tag, "receive stop request")
			}
			return
		}
	}
}

func (this *Goo) IsRun() bool {
	return this.started
}

func (this *Goo) Run() bool {
	if this.requests == nil {
		logger.Error(this.Tag, "requests queue is nil")
		return false
	}
	this.started = true
	go this.run()
	return true
}

func (this *Goo) Do(name string, req interface{}, cb func(err error)) error {
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
	this.requests.Write(&Request{name, req, cb})
	if this.EDebug {
		logger.Debug(this.Tag, "push new request - %s", name)
	}
	return nil
}

func (this *Goo) DoNow(name string, req interface{}) error {
	return this.Do(name, req, nil)
}

func (this *Goo) DoSync(name string, req interface{}) error {
	ev := make(chan error, 1)
	defer close(ev)
	this.Do(name, req, SyncCallback(ev))
	err := <-ev
	return err
}

func (this *Goo) DoTimeout(name string, req interface{}, timeout time.Duration) error {
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

func (this *Goo) IsClosing() bool {
	return this.closeState.IsClosing()
}

func (this *Goo) Stop() bool {
	if this.closeState.AskClose() && this.requests != nil {
		this.requests.Write(&Request{"STOP", nil, nil})
	}
	return true
}

func (this *Goo) WaitStop() bool {
	if this.closeState.IsClosing() && this.started {
		this.closeState.WaitClosed()
	}
	return true
}

func (this *Goo) Resize(newbufsize int) bool {
	logger.Info(this.Tag, "resize to %d", newbufsize)
	this.requests.DoResize(newbufsize)
	return true
}
