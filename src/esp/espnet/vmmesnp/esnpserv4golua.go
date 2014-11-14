package vmmesnp

import (
	"bmautil/socket"
	"esp/espnet/espsocket"
	"logger"
	"strings"
	"sync"
)

const (
	tag = "vmmesnp"
)

func createESNP() (interface{}, error) {
	r := new(esnpserv)
	r.Init()
	return r, nil
}

type sockInfo struct {
	host  string
	token string
	sock  *espsocket.Socket
}

type esnpserv struct {
	lock  sync.RWMutex
	socks map[string]*sockInfo
}

func key(host, token string) string {
	return strings.ToLower(host) + "_" + token
}

func (this *esnpserv) Init() {
	this.socks = make(map[string]*sockInfo)
}

func (this *esnpserv) Create(host string, tms int) (*espsocket.Socket, error) {
	cfg := new(socket.DialConfig)
	cfg.Address = host
	cfg.TimeoutMS = tms
	if err := cfg.Valid(); err != nil {
		return nil, err
	}
	sock, err1 := espsocket.Dial("vmmesnp", cfg, espsocket.SOCKET_CHANNEL_CODER_ESPNET)
	if err1 != nil {
		logger.Debug(tag, "socket create '%s' fail - %s", host, err1)
		return nil, err1
	}
	logger.Debug(tag, "socket create '%s' -> %s", host, sock)
	return sock, err1
}

func (this *esnpserv) Open(host string, token string, tms int) (*espsocket.Socket, error) {
	k := key(host, token)
	this.lock.RLock()
	si, ok := this.socks[k]
	this.lock.RUnlock()
	if ok && !si.sock.IsBreak() {
		logger.Debug(tag, "socket open '%s' done", host)
		return si.sock, nil
	}

	if si == nil {
		si = new(sockInfo)
		si.host = host
		si.token = token
	}
	sock, err2 := this.Create(host, tms)
	if err2 != nil {
		return nil, err2
	}

	this.lock.Lock()
	defer this.lock.Unlock()
	si2, ok2 := this.socks[k]
	if ok2 {
		sock.AskClose()
		si = si2
	} else {
		si.sock = sock
		this.socks[k] = si
	}
	return si.sock, nil
}

func (this *esnpserv) TryClose() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	l := len(this.socks)
	for k, si := range this.socks {
		delete(this.socks, k)
		si.sock.AskClose()
	}
	return l > 0
}
