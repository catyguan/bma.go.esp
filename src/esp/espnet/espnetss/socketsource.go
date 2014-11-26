package espnetss

import (
	"bmautil/socket"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"sync/atomic"
	"time"
)

type loginInfo struct {
	loginType   string
	loginh      LoginHandler
	certificate string
}

type SocketSource struct {
	cfg        *Config
	pool       chan *espsocket.Socket
	closed     uint32
	doFail     bool
	doFailTime time.Time
}

func NewSocketSource(cfg *Config) *SocketSource {
	r := new(SocketSource)
	r.cfg = cfg
	if cfg.PoolSize > 0 {
		if cfg.PreConns > cfg.PoolSize {
			cfg.PreConns = cfg.PoolSize
		}
		r.pool = make(chan *espsocket.Socket, cfg.PoolSize)
	}
	return r
}

func (this *SocketSource) Name() string {
	return fmt.Sprintf("%s@%s", this.cfg.User, this.cfg.Host)
}

func (this *SocketSource) Key() string {
	return this.cfg.Key()
}

func (this *SocketSource) String() string {
	return this.Name()
}

func (this *SocketSource) IsClose() bool {
	return atomic.LoadUint32(&this.closed) == 1
}

func (this *SocketSource) _login(sock *espsocket.Socket) (bool, error) {
	var li LoginHandler
	if this.cfg.LoginType != "" {
		li = GetLoginHandler(this.cfg.LoginType)
		if li == nil {
			return false, fmt.Errorf("invalid login handler(%s)", this.cfg.LoginType)
		}
	}
	if li != nil {
		return li(sock, this.cfg.User, this.cfg.Certificate)
	}
	return true, nil
}

func (this *SocketSource) _create(timeoutMS int, log bool) (*espsocket.Socket, error) {
	cfg := new(socket.DialConfig)
	cfg.Address = this.cfg.Host
	cfg.TimeoutMS = timeoutMS
	sock, err := espsocket.Dial(this.Name(), cfg, espsocket.SOCKET_CHANNEL_CODER_ESPNET, log)
	if err != nil {
		return nil, err
	}
	done, err1 := this._login(sock)
	if err1 != nil {
		sock.Shutdown()
		return nil, err1
	}
	if !done {
		sock.Shutdown()
		return nil, fmt.Errorf("'%s' login fail", this.Name())
	}
	return sock, nil
}

func (this *SocketSource) logFail(msg string, err error) {
	if this.doFail {
		if this.doFailTime.Sub(time.Now()).Minutes() >= 5 {
			this.doFail = false
		}
	}
	if !this.doFail {
		logger.Info(tag, msg, this, err)
		this.doFail = true
		this.doFailTime = time.Now()
	}
}

func (this *SocketSource) isFull() bool {
	return len(this.pool) >= this.cfg.PreConns
}

func (this *SocketSource) preConn() {
	if this.IsClose() {
		return
	}
	if this.isFull() {
		return
	}
	sock, err := this._create(30*1000, false)
	if err != nil {
		this.logFail("preConn(%s) fail - %s", err)
	} else {
		if this.isFull() {
			sock.Shutdown()
			return
		}
		func() {
			defer func() {
				recover()
			}()
			this.doFail = false
			select {
			case this.pool <- sock:
				sock.SetCloseListener("SocketSource", func() {
					this.onSocketClose(sock)
				})
				logger.Debug(tag, "%s preconn -> %s", this, sock)
			default:
				sock.Shutdown()
			}
		}()
	}
	this.preConn()
}

func (this *SocketSource) Return(sock *espsocket.Socket) bool {
	if sock.IsBreak() {
		return false
	}
	if this.IsClose() {
		return false
	}
	defer func() {
		recover()
	}()
	select {
	case this.pool <- sock:
		sock.SetCloseListener("SocketSource", func() {
			this.onSocketClose(sock)
		})
		logger.Debug(tag, "%s return -> %s", this, sock)
		return true
	default:
		sock.Shutdown()
		return false
	}
}

func (this *SocketSource) Open(timeoutMS int) (*espsocket.Socket, error) {
	var sock *espsocket.Socket
	select {
	case sock = <-this.pool:
	default:
	}
	if sock != nil {
		sock.SetCloseListener("SocketSource", nil)
		go this.preConn()
		return sock, nil
	}
	return this._create(timeoutMS, true)
}

func (this *SocketSource) KeepLive(acceptor espsocket.SocketAcceptor) error {
	if cap(this.pool) == 0 {
		return fmt.Errorf("PoolSize is 0, can't keeplive")
	}
	if this.IsClose() {
		return nil
	}
	go func() {
		for {
			sock := <-this.pool
			if sock != nil {
				return
			}
			go this.preConn()
			err := acceptor(sock)
			if err != nil {
				this.logFail("keepLive(%s) fail - %s", err)
				sock.AskClose()
				time.Sleep(1 * time.Second)
				continue
			}
			sock.SetCloseListener("SocketSourceKeepLive", func() {
				this.KeepLive(acceptor)
			})
			return
		}
	}()
	return nil
}

func (this *SocketSource) Start() {
	go this.preConn()
}

func (this *SocketSource) onSocketClose(sock *espsocket.Socket) {
	sl := make([]*espsocket.Socket, 0, len(this.pool))
	for {
		select {
		case s := <-this.pool:
			if sock != s {
				sl = append(sl, s)
			}
		default:
			break
		}
	}
	func() {
		defer func() {
			recover()
		}()
		for _, s := range sl {
			select {
			case this.pool <- s:
			default:
				s.SetCloseListener("SocketSource", nil)
				s.Shutdown()
			}
		}
	}()
	go this.preConn()
}

func (this *SocketSource) Close() {
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		for {
			select {
			case sock := <-this.pool:
				sock.SetCloseListener("SocketSource", nil)
				sock.Shutdown()
			default:
				return
			}
		}
		close(this.pool)
	}
}
