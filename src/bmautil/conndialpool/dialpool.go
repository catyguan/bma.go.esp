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
	IdleMS          int
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
	// if this.IdleMS <= 0 {
	// 	this.IdleMS = 60 * 1000
	// }
	return nil
}

func DefaultRetryConfig() *retryst.RetryConfig {
	rcfg := new(retryst.RetryConfig)
	rcfg.DelayMin = 100
	rcfg.DelayIncrease = 200
	rcfg.DelayMax = 1000
	return rcfg
}

type connManager struct {
	pool *DialPool
	done bool
}

func (this *connManager) CloseConn(conn *connutil.ConnExt) {
	if !this.done {
		this.done = true
		this.pool.CloseConn(conn)
	}
}

func (this *connManager) FinishConn(conn *connutil.ConnExt) {
	if !this.done {
		this.done = true
		this.pool.ReturnConn(conn)
	}
}

func (this *connManager) ReleaseConn(conn *connutil.ConnExt) {
	if !this.done {
		this.done = true
		this.pool.ReleaseConn(conn)
	}
}

type dialPoolItem struct {
	conn        *connutil.ConnExt
	idleOutTime time.Time
}

// DialPool
type DialPool struct {
	name       string
	config     *DialPoolConfig
	closeState syncutil.CloseState

	timer    *time.Ticker
	markTime time.Time
	count    int32
	active   int32
	wait     chan *dialPoolItem
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

	this.wait = make(chan *dialPoolItem, this.config.MaxSize)

	return this
}

func (this *DialPool) Name() string {
	return this.name
}

func (this *DialPool) SetName(n string) {
	this.name = n
}

func (this *DialPool) String() string {
	return fmt.Sprintf("ConnDialPool[%s, %d:%d/%d/%d]", this.name, len(this.wait), this.active, this.count, this.config.MaxSize)
}

func (this *DialPool) GetConn(deadline time.Time, log bool) (*connutil.ConnExt, error) {
	conn, err := this._getConn(deadline, log)
	if conn != nil {
		conn.Manager = &connManager{pool: this}
	}
	return conn, err
}

func (this *DialPool) apool(item *dialPoolItem) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	ok = true
	this.wait <- item
	return
}

func (this *DialPool) _getConn(deadline time.Time, log bool) (*connutil.ConnExt, error) {
	if this.IsClosing() {
		return nil, errors.New("closed")
	}
	var item *dialPoolItem
	select {
	case item = <-this.wait:
		if item == nil {
			return nil, nil
		}
		return item.conn, nil
	default:
	}
	if this.IsClosing() {
		return nil, errors.New("closed")
	}
	c := atomic.LoadInt32(&this.count)
	if c < int32(this.config.MaxSize) {
		// create it
		atomic.AddInt32(&this.count, 1)
		conn, err := this.doDial(deadline, log)
		if err != nil {
			atomic.AddInt32(&this.count, -1)
			return nil, err
		}
		atomic.AddInt32(&this.active, 1)
		rconn := connutil.NewConnExt(conn)
		return rconn, nil
	}
	// wait it
	if log {
		logger.Debug(tag, "%s max, wait returnConn", this)
	}
	if !deadline.IsZero() {
		timer := time.NewTimer(deadline.Sub(time.Now()))
		select {
		case item := <-this.wait:
			timer.Stop()
			if item == nil {
				return nil, nil
			}
			return item.conn, nil
		case <-timer.C:
		}
	}
	return nil, errors.New("timeout")
}

func (this *DialPool) ReleaseConn(conn *connutil.ConnExt) {
	atomic.AddInt32(&this.count, -1)
	atomic.AddInt32(&this.active, -1)
}

func (this *DialPool) doCloseConn(conn net.Conn) int32 {
	conn.Close()
	c := atomic.AddInt32(&this.count, -1)
	atomic.AddInt32(&this.active, -1)
	return c
}

func (this *DialPool) CloseConn(conn *connutil.ConnExt) {
	c := this.doCloseConn(conn)
	if c < int32(this.config.InitSize) && !this.IsClosing() {
		// reconnect
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "%s retry %s", this, this.config.Address)
		}
		go this.startRetry()
	}
}

func (this *DialPool) ReturnConn(conn *connutil.ConnExt) {
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
	conn.Manager = nil
	item := new(dialPoolItem)
	item.idleOutTime = time.Now().Add(time.Duration(this.config.IdleMS) * time.Millisecond)
	item.conn = conn
	if !this.apool(item) {
		this.doCloseConn(conn)
	}
}

func (this *DialPool) Start() bool {
	if err := this.config.Valid(); err != nil {
		logger.Warn(tag, "%s config invalid - %s", this, err)
		return false
	}
	return true
}

func (this *DialPool) doDial(deadline time.Time, log bool) (net.Conn, error) {
	var conn net.Conn
	var err error
	cfg := this.config
	var timeout time.Duration
	if deadline.IsZero() {
		tm := cfg.TimeoutMS
		if tm <= 0 {
			tm = 5 * 1000
		}
		timeout = time.Duration(tm) * time.Millisecond
	} else {
		timeout = deadline.Sub(time.Now())
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
		conn, err := this.GetConn(time.Time{}, false)
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
	return atomic.LoadInt32(&this.active)
}

func (this *DialPool) GetInitSize() int {
	return this.config.InitSize
}

func (this *DialPool) Run() bool {
	c := atomic.LoadInt32(&this.count)
	for i := int(c); i < this.config.InitSize; i++ {
		go func() {
			conn, err := this.GetConn(time.Time{}, true)
			if err != nil {
				if this.needRetry() {
					this.startRetry()
				}
				return
			}
			this.ReturnConn(conn)
		}()
	}
	if this.config.InitSize > 0 || this.config.IdleMS > 0 {
		this.timer = time.NewTicker(time.Duration(this.config.CheckMS) * time.Millisecond)
		go func() {
			for {
				dt := <-this.timer.C
				if dt.IsZero() {
					return
				}
				if this.IsClosing() {
					return
				}
				l := len(this.wait)
				// fmt.Println("checking", l)
				for i := 0; i < l; i++ {
					select {
					case item := <-this.wait:
						if item != nil {
							if item.conn.CheckBreak() {
								// fmt.Println("checking", "fail")
								this.CloseConn(item.conn)
								continue
							}
							// fmt.Println("checking idle", len(this.wait)+1, this.config.InitSize, time.Now().After(item.idleOutTime))
							if this.config.IdleMS > 0 && len(this.wait)+1 > this.config.InitSize && time.Now().After(item.idleOutTime) {
								logger.Debug(tag, "%s idle close", item.conn)
								this.CloseConn(item.conn)
								continue
							}
							// fmt.Println("checking", "ok")
							if !this.apool(item) {
								this.doCloseConn(item.conn)
							}
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
			case item := <-this.wait:
				if item != nil {
					this.doCloseConn(item.conn)
				}
			default:
				done = true
			}
		}
		close(this.wait)
	}
	return true
}

func (this *DialPool) IsClosing() bool {
	return this.closeState.IsClosing()
}
