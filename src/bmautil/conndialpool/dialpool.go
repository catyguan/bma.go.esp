package conndialpool

import (
	"bmautil/connutil"
	"bmautil/retryst"
	"bmautil/syncutil"
	"errors"
	"fmt"
	"logger"
	"net"
	"sync/atomic"
	"time"
)

const (
	tag = "ConnDialPool"
)

// DialPoolConfig
type DialPoolConfig struct {
	Net             string
	Address         string
	TimeoutMS       int
	MaxSize         int
	InitSize        int
	CheckMS         int
	Retry           *retryst.RetryConfig
	RetryFailInfoMS int
}

func (this *DialPoolConfig) Valid() error {
	if this.Address == "" {
		return errors.New("address empty")
	}
	if this.Net == "" {
		this.Net = "tcp"
	}
	if this.MaxSize <= 0 {
		return fmt.Errorf("maxsize(%d) invalid", this.MaxSize)
	}
	if this.RetryFailInfoMS <= 0 {
		this.RetryFailInfoMS = 30 * 1000
	}
	if this.Retry == nil {
		this.Retry = DefaultRetryConfig()
	}
	if this.CheckMS <= 0 {
		this.CheckMS = 100
	}
	return nil
}

func DefaultRetryConfig() *retryst.RetryConfig {
	rcfg := new(retryst.RetryConfig)
	rcfg.DelayMin = 100
	rcfg.DelayIncrease = 200
	rcfg.DelayMax = 1000
	return rcfg
}

type dialPoolItem struct {
	conn net.Conn
}

// DialPool
type DialPool struct {
	name       string
	config     *DialPoolConfig
	closeState syncutil.CloseState

	timer    *time.Ticker
	markTime time.Time
	count    int32
	wait     chan *connutil.ConnExt
}

func NewDialPool(name string, cfg *DialPoolConfig) *DialPool {
	err := cfg.Valid()
	if err != nil {
		logger.Error(tag, "config error - %s", err)
		panic(err.Error())
	}

	this := new(DialPool)
	this.name = name
	this.config = cfg
	this.closeState.InitCloseState()

	this.wait = make(chan *connutil.ConnExt, this.config.MaxSize)

	return this
}

func (this *DialPool) Name() string {
	return this.name
}

func (this *DialPool) String() string {
	return fmt.Sprintf("ConnDialPool[%s, %d:%d/%d]", this.name, len(this.wait), this.count, this.config.MaxSize)
}

func (this *DialPool) GetConn(timeout time.Duration, log bool) (net.Conn, error) {
	if this.IsClosing() {
		return nil, errors.New("closed")
	}
	var s net.Conn
	select {
	case s = <-this.wait:
		if s == nil {
			return nil, nil
		}
		return s, nil
	default:
	}
	if this.IsClosing() {
		return nil, errors.New("closed")
	}
	c := atomic.LoadInt32(&this.count)
	if c < int32(this.config.MaxSize) {
		// create it
		atomic.AddInt32(&this.count, 1)
		s, err := this.doDial(timeout, log)
		if err != nil {
			atomic.AddInt32(&this.count, -1)
			return nil, err
		}
		return s, nil
	}
	// wait it
	if log {
		logger.Debug(tag, "%s max, wait returnConn", this)
	}
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		select {
		case s := <-this.wait:
			timer.Stop()
			return s, nil
		case <-timer.C:
		}
	}
	return nil, errors.New("timeout")
}

func (this *DialPool) doCloseConn(conn net.Conn) int32 {
	conn.Close()
	c := atomic.AddInt32(&this.count, -1)
	return c
}

func (this *DialPool) CloseConn(conn net.Conn) {
	c := this.doCloseConn(conn)
	if c < int32(this.config.InitSize) && !this.IsClosing() {
		// reconnect
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "%s retry %s", this, this.config.Address)
		}
		go this.startRetry()
	}
}

func (this *DialPool) ReturnConn(conn net.Conn) {
	if this.IsClosing() {
		this.doCloseConn(conn)
		return
	}
	c := atomic.LoadInt32(&this.count)
	if c > int32(this.config.MaxSize) {
		// don't return,close it
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "Pool[%s] max(%d,%d), close returnConn", this.name, c, this.config.MaxSize)
		}
		this.doCloseConn(conn)
		return
	}
	defer func() {
		if recover() != nil {
			conn.Close()
		}
	}()
	this.wait <- connutil.NewConnExt(conn, nil)
}

func (this *DialPool) Start() bool {
	if err := this.config.Valid(); err != nil {
		logger.Warn(tag, "%s config invalid - %s", this, err)
		return false
	}
	return true
}

func (this *DialPool) doDial(timeout time.Duration, log bool) (net.Conn, error) {
	var conn net.Conn
	var err error
	cfg := this.config
	if timeout == 0 {
		tm := cfg.TimeoutMS
		if tm <= 0 {
			tm = 5 * 1000
		}
		timeout = time.Duration(tm) * time.Millisecond
	}
	conn, err = net.DialTimeout(cfg.Net, cfg.Address, timeout)
	if err != nil {
		if log {
			logger.Debug(tag, "dial (%s %s) fail - %s", cfg.Net, cfg.Address, err.Error())
		}
		return nil, err
	}
	return conn, nil
}

func (this *DialPool) needRetry() bool {
	if this.IsClosing() {
		return false
	}
	c := atomic.LoadInt32(&this.count)
	return c < int32(this.config.InitSize)
}

func (this *DialPool) startRetry() {
	retry := new(retryst.RetryState)
	retry.Config = this.config.Retry
	retry.Begin(this.reconnectConn)
}

func (this *DialPool) reconnectConn(rs *retryst.RetryState, lastTry bool) bool {
	if this.needRetry() {
		conn, err := this.GetConn(0, false)
		if err != nil {
			if lastTry {
				logger.Warn(tag, "%s retry end for error - %s", this, err)
			} else {
				if this.config.RetryFailInfoMS > 0 {
					tm := this.markTime
					if tm.IsZero() {
						this.markTime = time.Now()
					} else {
						if time.Now().Sub(tm) >= time.Duration(this.config.RetryFailInfoMS)*time.Millisecond {
							this.markTime = time.Now()
							logger.Info(tag, "%s dial retry fail, begin %s (%d)", this, rs.GetBeginTime(), rs.GetRetryCount())
						}
					}
				}
			}
			return false
		}
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "%s reconnect %s done", this, this.config.Address)
		}
		this.ReturnConn(conn)
	}
	return true
}

func (this *DialPool) ActiveConn() int32 {
	return atomic.LoadInt32(&this.count)
}

func (this *DialPool) Run() bool {
	c := atomic.LoadInt32(&this.count)
	for i := int(c); i < this.config.InitSize; i++ {
		go func() {
			conn, err := this.GetConn(0, true)
			if err != nil {
				if this.needRetry() {
					this.startRetry()
				}
				return
			}
			this.ReturnConn(conn)
		}()
	}
	if this.config.InitSize > 0 {
		this.timer = time.NewTicker(time.Duration(this.config.CheckMS) * time.Millisecond)
		go func() {
			for {
				defer func() {
					recover()
				}()
				dt := <-this.timer.C
				if dt.IsZero() {
					return
				}
				l := len(this.wait)
				// fmt.Println("checking", l)
				for i := 0; i < l; i++ {
					select {
					case s := <-this.wait:
						if !s.CheckBreak() {
							// fmt.Println("checking", "ok")
							this.wait <- s
						} else {
							// fmt.Println("checking", "fail")
							this.CloseConn(s)
						}
					default:
						break
					}
				}
			}
		}()
	}
	return true
}

func (this *DialPool) StartAndRun() bool {
	if !this.Start() {
		return false
	}
	return this.Run()
}

func (this *DialPool) Close() bool {
	this.AskClose()
	return true
}

func (this *DialPool) AskClose() {
	if this.closeState.AskClose() {
		if this.timer != nil {
			this.timer.Stop()
		}
		done := false
		for {
			if done {
				break
			}
			select {
			case s := <-this.wait:
				this.doCloseConn(s)
			default:
				done = true
			}
		}
		close(this.wait)
	}
}

func (this *DialPool) IsClosing() bool {
	return this.closeState.IsClosing()
}
