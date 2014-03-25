package goo

import (
	"bmautil/valutil"
	"errors"
	"logger"
	"time"
)

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

type ExitHandler func()

type requestInfo struct {
	action   interface{}
	callback func(err error)
}

type Goo struct {
	Tag         string
	EDebug      bool
	sm          StateMachine
	queue       chan interface{}
	exitHandler ExitHandler
}

func (this *Goo) InitGoo(tag string, queueSize int, exithandler ExitHandler) {
	this.queue = make(chan interface{}, queueSize)
	this.exitHandler = exithandler
	this.sm.InitStateMachine(STATE_INIT, gooStates)
	this.sm.SetCanEnterF(canEnter4goo)
	this.sm.SetAfterEnterF(afterEnter4goo)
}

func (this *Goo) execute(req interface{}) (err error) {
	var cb func(err error)
	defer func() {
		ex := recover()
		if ex != nil {
			if _, ok := ex.(error); ok {
				err = ex.(error)
			} else {
				err = errors.New(valutil.ToString(ex, "unknow error"))
			}
		}
		if cb != nil {
			cb(err)
		} else {
			if err != nil {
				logger.Error(this.Tag, "execte '%T' fail - %s", req, err.Error())
			}
		}
		if this.EDebug {
			logger.Debug(this.Tag, "request done - %T", req)
		}
	}()
	switch act := req.(type) {
	case func():
		act()
	case func() error:
		err = act()
	case *requestInfo:
		if act.action != nil {
			cb = act.callback
			err = this.execute(act.action)
		}
	}
	return
}

func (this *Goo) run() {
	defer func() {
		if this.exitHandler != nil {
			this.exitHandler()
		}
		this.sm.Enter(this, STATE_CLOSE)
	}()
	this.sm.TryEnter(this, STATE_RUN)
	for {
		switch this.sm.GetState() {
		case STATE_CLOSE:
			return
		}
		dreq := <-this.queue
		if dreq == nil {
			return
		}
		if this.EDebug {
			logger.Debug(this.Tag, "popup new request - %T", dreq)
		}
		this.execute(dreq)
	}
}

func (this *Goo) Run() bool {
	this.sm.TryEnter(this, STATE_START)
	return true
}

func (this *Goo) Stop() bool {
	if !this.sm.CompareAndEnter(this, STATE_INIT, STATE_CLOSE) {
		this.sm.TryEnter(this, STATE_STOP)
	}
	return true
}

func (this *Goo) StopAndWait() {
	for {
		this.Stop()
		if this.sm.IsState(STATE_CLOSE) {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func (this *Goo) GetState() uint32 {
	return this.sm.GetState()
}

func (this *Goo) CanInvoke() bool {
	v := this.sm.GetState()
	return v == STATE_START || v == STATE_RUN
}

func (this *Goo) Do(req interface{}, cb func(err error)) (err error) {
	if !this.CanInvoke() {
		err := errors.New(this.Tag + " closed")
		if cb != nil {
			cb(err)
			return nil
		} else {
			return err
		}
	}
	defer func() {
		ex := recover()
		if ex != nil {
			err = ex.(error)
		}
	}()
	if cb != nil {
		o := new(requestInfo)
		o.action = req
		o.callback = cb
		this.queue <- o
	} else {
		this.queue <- req
	}
	if this.EDebug {
		logger.Debug(this.Tag, "push new request - %T", req)
	}
	return nil
}

func (this *Goo) DoNow(req interface{}) error {
	return this.Do(req, nil)
}

func (this *Goo) DoSync(req interface{}) error {
	ev := make(chan error, 1)
	defer close(ev)
	this.Do(req, SyncCallback(ev))
	err := <-ev
	return err
}

func (this *Goo) DoTimeout(req interface{}, timeout time.Duration) error {
	ev := make(chan error, 1)
	defer close(ev)
	tm := time.NewTimer(timeout)
	defer tm.Stop()
	this.Do(req, SyncCallback(ev))
	select {
	case err := <-ev:
		return err
	case <-tm.C:
		return errors.New("timeout")
	}
}
