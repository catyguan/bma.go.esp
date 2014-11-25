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
	host       string
	user       string
	pool       chan *espsocket.Socket
	logins     []*loginInfo
	closed     uint32
	doFail     bool
	doFailTime time.Time
}

func NewSocketSource(host string, user string, preconns int) *SocketSource {
	r := new(SocketSource)
	r.host = host
	r.user = user
	r.pool = make(chan *espsocket.Socket, preconns)
	r.logins = make([]*loginInfo, 0)
	return r
}

func (this *SocketSource) Name() string {
	return fmt.Sprintf("%s@%s", this.user, this.host)
}

func (this *SocketSource) String() string {
	return fmt.Sprintf("%s@%s(%d)", this.user, this.host, len(this.logins))
}

func (this *SocketSource) IsClose() bool {
	return atomic.LoadUint32(&this.closed) == 1
}

func (this *SocketSource) Add(cert string, lt string) bool {
	lh := GetLoginHandler(lt)
	if lh == nil {
		return false
	}
	for _, li := range this.logins {
		if li.certificate == cert && li.loginType == lt {
			return false
		}
	}
	this.logins = append(this.logins, &loginInfo{lt, lh, cert})
	return true
}

func (this *SocketSource) _login(sock *espsocket.Socket) (bool, error) {
	lgs := this.logins
	if len(lgs) > 0 {
		done := false
		var lastErr error
		for _, li := range lgs {
			ok, err1 := li.loginh(sock, this.user, li.certificate)
			if err1 != nil {
				logger.Debug(tag, "dologin(%s, %v) fail - %s", this.user, li.loginType, err1)
				lastErr = err1
			}
			if ok {
				lastErr = nil
				done = true
				break
			}
		}
		return done, lastErr
	}
	return true, nil
}

func (this *SocketSource) _create(timeoutMS int, log bool) (*espsocket.Socket, error) {
	cfg := new(socket.DialConfig)
	cfg.Address = this.host
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

func (this *SocketSource) preConn() {
	if this.IsClose() {
		return
	}
	if len(this.pool) == cap(this.pool) {
		return
	}
	sock, err := this._create(30*1000, false)
	if err != nil {
		this.logFail("preConn(%s) fail - %s", err)
	} else {
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
		return fmt.Errorf("preConns is 0, can't keeplive")
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
			this.pool <- s
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
