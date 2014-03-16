package goo

import (
	"logger"
	"time"
)

type ExitHandler func()

type Goo struct {
	Tag         string
	EDebug      bool
	sm          StateMachine
	queue       chan interface{}
	requests    []interface{}
	exitHandler ExitHandler
}

func (this *Goo) InitGoo(tag string, queueSize int, exithandler ExitHandler) {
	this.queue = make(chan interface{}, queueSize)
	this.exitHandler = exithandler
	this.sm.InitStateMachine(STATE_INIT, gooStates)
	this.sm.SetCanEnterF(canEnter4goo)
	this.sm.SetAfterEnterF(afterEnter4goo)
}

// func (this *Goo) Execute(req *Request) (running bool, err error) {
// 	running = true
// 	defer func() {
// 		ex := recover()
// 		if ex != nil {
// 			if _, ok := ex.(error); ok {
// 				err = ex.(error)
// 			} else {
// 				err = errors.New(valutil.ToString(ex, "unknow error"))
// 			}
// 		}
// 		if err != nil {
// 			logger.Error(this.Tag, "execte '%s' fail - %s", req.name, err.Error())
// 			if this.ErrorHandler != nil {
// 				safe(func() {
// 					running = this.ErrorHandler(req.data, err)
// 				})
// 			}
// 		}
// 		if req.callback != nil {
// 			req.callback(err)
// 		}
// 		if this.EDebug {
// 			logger.Debug(this.Tag, "request done - %s", req.name)
// 		}
// 	}()
// 	return this.RequestHandler(req.data)
// }

func (this *Goo) run() {
	defer func() {
		if this.exitHandler != nil {
			this.exitHandler()
		}
		this.sm.Enter(o, STATE_CLOSE)
	}()
	this.sm.TryEnter(this, STATE_RUN)
	for {
		switch this.sm.GetState() {
		case STATE_STOP, STATE_CLOSE:
			break
		}
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
				this.requests = append(this.requests, v)
			}
		}
		var dreq interface{}
		if len(this.requests) > 0 {
			dreq = this.requests[0]
			this.requests = 
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

func (this *Goo) Run() bool {
	this.sm.TryEnter(this, STATE_START)
	return true
}

func (this *Goo) Stop() bool {
	if this.sm.IsState(STATE_INIT) {
		this.sm.TryEnter(this, STATE_CLOSE)
	} else {
		this.sm.TryEnter(this, STATE_STOP)
	}
	return true
}

func (this *Goo) WaitClosed() {
	for {
		if this.sm.IsState(STATE_INIT) {
			this.sm.TryEnter(this, STATE_CLOSE)
		}
		if this.sm.IsState(this, STATE_CLOSE) {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
}
