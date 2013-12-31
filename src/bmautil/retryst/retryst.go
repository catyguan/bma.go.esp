package retryst

import (
	"fmt"
	"time"
)

type RetryConfig struct {
	DelayMin      int // ms
	DelayMax      int
	DelayIncrease int
	Max           int
}

type RetryState struct {
	Config *RetryConfig

	lastDelay int
	count     int
	beginTime time.Time
	lastTime  time.Time
	timer     *time.Timer
}

func (this *RetryState) GetBeginTime() time.Time {
	return this.beginTime
}

func (this *RetryState) GetLastTime() time.Time {
	return this.beginTime
}

func (this *RetryState) GetRetryCount() int {
	return this.count
}

func (this *RetryState) Begin(f func(rs *RetryState, lastTry bool) bool) {
	this.Reset()

	this.beginTime = time.Now()
	this.doRetry(f)
}

func (this *RetryState) doRetry(f func(rs *RetryState, lastTry bool) bool) {
	if this.Config.Max > 0 && this.count+1 > this.Config.Max {
		go func() {
			defer func() {
				recover()
			}()
			f(this, true)
		}()
		return
	}
	delay := this.lastDelay
	if delay == 0 {
		delay = this.Config.DelayMin
	} else {
		delay += this.Config.DelayIncrease
	}
	if this.Config.DelayMax > 0 && delay > this.Config.DelayMax {
		delay = this.Config.DelayMax
	}
	this.lastDelay = delay
	this.count++
	this.timer = time.AfterFunc(time.Duration(delay)*time.Millisecond, func() {
		defer func() {
			if recover() != nil {
				this.doRetry(f)
			}
		}()
		this.lastTime = time.Now()
		if !f(this, false) {
			this.doRetry(f)
		}
	})
}

func (this *RetryState) Reset() {
	this.Cancel()

	this.lastDelay = 0
	this.beginTime = time.Time{}
	this.lastTime = time.Time{}
	this.count = 0
}

func (this *RetryState) Cancel() {
	if this.timer != nil {
		this.timer.Stop()
		this.timer = nil
	}
}

func (this *RetryState) String() string {
	return fmt.Sprintf("begin:%s, last:%s, count:%d, delay:%s", this.beginTime, this.lastTime, this.count, time.Duration(this.lastDelay)*time.Millisecond)
}
