package socket

import (
	"bmautil/retryst"
	"bmautil/syncutil"
	"errors"
	"fmt"
	"logger"
	"net"
	"sync/atomic"
	"time"
)

// DialPoint
type DialConfig struct {
	Net       string
	Address   string
	TimeoutMS int
}

func (this *DialConfig) Valid() error {
	if this.Address == "" {
		return errors.New("address empty")
	}
	if this.Net == "" {
		this.Net = "tcp"
	}
	return nil
}

// DialPoolConfig
type DialPoolConfig struct {
	Dial                  DialConfig
	MaxSize               int
	InitSize              int
	Retry                 *retryst.RetryConfig
	RetryFailInfoDruation time.Duration
}

func (this *DialPoolConfig) Valid() error {
	if err := this.Dial.Valid(); err != nil {
		return err
	}
	if this.MaxSize <= 0 {
		return fmt.Errorf("maxsize(%d) invalid", this.MaxSize)
	}
	if this.RetryFailInfoDruation <= 0 {
		this.RetryFailInfoDruation = 30 * time.Second
	}
	if this.Retry == nil {
		this.Retry = this.DefaultRetryConfig()
	}
	return nil
}

func (this *DialPoolConfig) DefaultRetryConfig() *retryst.RetryConfig {
	rcfg := new(retryst.RetryConfig)
	rcfg.DelayMin = 100
	rcfg.DelayIncrease = 200
	rcfg.DelayMax = 1000
	return rcfg
}

// DialPool
type DialPool struct {
	name       string
	config     *DialPoolConfig
	socketInit SocketInit
	closeState syncutil.CloseState

	markTime time.Time
	count    int32
	wait     chan *Socket
}

func NewDialPool(name string, cfg *DialPoolConfig, sinit SocketInit) *DialPool {
	err := cfg.Valid()
	if err != nil {
		logger.Error(tag, "config error - %s", err)
		panic(err.Error())
	}

	this := new(DialPool)
	this.name = name
	this.config = cfg
	this.socketInit = sinit
	this.closeState.InitCloseState()

	this.wait = make(chan *Socket, this.config.MaxSize)

	return this
}

func (this *DialPool) Name() string {
	return this.name
}

func (this *DialPool) String() string {
	return fmt.Sprintf("DialPool[%s, %d:%d/%d]", this.name, len(this.wait), this.count, this.config.MaxSize)
}

func (this *DialPool) GetSocket(timeout time.Duration, log bool) (*Socket, error) {
	if this.IsClosing() {
		return nil, errors.New("closed")
	}
	s := func() *Socket {
		for {
			select {
			case s := <-this.wait:
				if s == nil {
					return nil
				}
				if s.IsClosing() {
					// atomic.AddInt32(&this.count, -1)
					continue
				}
				return s
			default:
				return nil
			}
		}
	}()
	if s != nil {
		return s, nil
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
	logger.Debug(tag, "%s max, wait returnSocket", this)
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

func (this *DialPool) ReturnSocket(sock *Socket) {
	if sock.IsClosing() {
		return
	}
	if this.IsClosing() {
		this.removeSocket(sock)
		return
	}
	c := atomic.LoadInt32(&this.count)
	if c > int32(this.config.MaxSize) {
		// don't return,close it
		this.removeSocket(sock)
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "Pool[%s] max(%d,%d), close returnSocket", this.name, c, this.config.MaxSize)
		}
		return
	}
	defer func() {
		if recover() != nil {
			sock.Close()
		}
	}()
	sock.Receiver = nil
	this.wait <- sock
}

func (this *DialPool) removeSocket(sock *Socket) {
	atomic.AddInt32(&this.count, -1)
	sock.RemoveCloseListener(this.closeId())
	sock.Close()
}

func (this *DialPool) closeId() string {
	return fmt.Sprintf("DP_%p", this)
}

func (this *DialPool) Start() bool {
	if err := this.config.Valid(); err != nil {
		logger.Warn(tag, "%s config invalid - %s", this, err)
		return false
	}
	return true
}

func (this *DialPool) doDial(timeout time.Duration, log bool) (*Socket, error) {
	var conn net.Conn
	var err error
	cfg := &this.config.Dial
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
	sock := NewSocket(conn, 32, 0)
	err = sock.Start(this.socketInit)
	if err != nil {
		logger.Debug(tag, "Socket[%s] start fail", sock)
		return nil, err
	}
	sock.AddCloseListener(this.onSocketClose, this.closeId())
	return sock, nil
}

func (this *DialPool) clearSocket(so *Socket) {
	tmp := make([]*Socket, 0)
	for {
		select {
		case s := <-this.wait:
			if s != so && !s.IsClosing() {
				tmp = append(tmp, s)
			}
		default:
			for _, s := range tmp {
				this.wait <- s
			}
			return
		}
	}
}

func (this *DialPool) onSocketClose(so *Socket) {
	c := atomic.AddInt32(&this.count, -1)
	this.clearSocket(so)
	if c < int32(this.config.InitSize) && !this.IsClosing() {
		// reconnect
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "%s retry %s", this, this.config.Dial.Address)
		}
		go this.startRetry()
	}
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
	retry.Begin(this.reconnectSocket)
}

func (this *DialPool) reconnectSocket(rs *retryst.RetryState, lastTry bool) bool {
	if this.needRetry() {
		sock, err := this.GetSocket(0, false)
		if err != nil {
			if lastTry {
				logger.Warn(tag, "%s retry end for error - %s", this, err)
			} else {
				if this.config.RetryFailInfoDruation > 0 {
					tm := this.markTime
					if tm.IsZero() {
						this.markTime = time.Now()
					} else {
						if time.Now().Sub(tm) >= this.config.RetryFailInfoDruation {
							this.markTime = time.Now()
							logger.Info(tag, "%s dial retry fail, begin %s (%d)", this, rs.GetBeginTime(), rs.GetRetryCount())
						}
					}
				}
			}
			return false
		}
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "%s reconnect %s done", this, this.config.Dial.Address)
		}
		this.ReturnSocket(sock)
	}
	return true
}

func (this *DialPool) IsBreak() *bool {
	r := false
	c := atomic.LoadInt32(&this.count)
	if c > 0 {
		return &r
	}
	if this.config.InitSize > 0 {
		r = true
		return &r
	}
	return nil
}

func (this *DialPool) Run() bool {
	c := atomic.LoadInt32(&this.count)
	for i := int(c); i < this.config.InitSize; i++ {
		go func() {
			sock, err := this.GetSocket(0, true)
			if err != nil {
				if this.needRetry() {
					this.startRetry()
				}
				return
			}
			this.ReturnSocket(sock)
		}()
	}
	return true
}

func (this *DialPool) Close() bool {
	this.AskClose()
	return true
}

func (this *DialPool) AskClose() {
	if this.closeState.AskClose() {
		done := false
		for {
			if done {
				break
			}
			select {
			case s := <-this.wait:
				s.Close()
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

// SocketServer
type SocketServer4Dial struct {
	pool      *DialPool
	timeoutMS int
	acceptor  SocketAcceptor
}

func NewSocketServer4Dial(pool *DialPool, timeoutMS int) *SocketServer4Dial {
	r := new(SocketServer4Dial)
	r.pool = pool
	r.timeoutMS = timeoutMS
	return r
}

func (this *SocketServer4Dial) SetAcceptor(sa SocketAcceptor) {
	this.acceptor = sa
	this.doConnect()
}

func (this *SocketServer4Dial) doConnect() {
	if this.pool.IsClosing() {
		return
	}
	go func() {
		s, err := this.pool.GetSocket(time.Duration(this.timeoutMS)*time.Millisecond, false)
		if this.pool.IsClosing() {
			return
		}
		if err == nil {
			err2 := this.acceptor(s)
			if err2 != nil {
				logger.Debug(tag, "accept(%s) fail - %s", s, err2)
			} else {
				s.AddCloseListener(func(sock *Socket) {
					this.doConnect()
				}, "")
				return
			}
		}
		time.AfterFunc(5*time.Millisecond, this.doConnect)
	}()
}
