package main

import (
	"bmautil/socket"
	"fmt"
	"net"
	"sync"
	"time"
)

type remoteResult struct {
	key  string
	resp bool
	err  error
}

type remoteExecutor func(sock *socket.Socket, key string, ch chan *remoteResult) error

type remoteRequest struct {
	conn    net.Conn
	execute remoteExecutor
	ch      chan *remoteResult
}

type remoteReqSession struct {
	sock *socket.Socket
}

type remote struct {
	key     string
	pool    *socket.DialPool
	timeout time.Duration
	lock    sync.Mutex
	reqs    map[*remoteRequest]*remoteReqSession
}

func (this *remote) InitRemote(k string) {
	this.key = k
	this.reqs = make(map[*remoteRequest]*remoteReqSession)
}

func (this *remote) Close() {
	if this.pool != nil {
		this.pool.Close()
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	for req, sess := range this.reqs {
		delete(this.reqs, req)
		sess.sock.Close()
		res := new(remoteResult)
		res.key = this.key
		res.err = fmt.Errorf("remote shutdown")
		req.ch <- res
	}
}

func (this *remote) PostRequest(req *remoteRequest) {
	go func() {
		sock, err := this.pool.GetSocket(this.timeout, false)
		if err != nil {
			res := new(remoteResult)
			res.key = this.key
			res.err = err
			req.ch <- res
			return
		}
		// sock.Trace = 32
		this.lock.Lock()
		sess := new(remoteReqSession)
		sess.sock = sock
		this.reqs[req] = sess
		this.lock.Unlock()

		err2 := req.execute(sock, this.key, req.ch)
		this.lock.Lock()
		delete(this.reqs, req)
		this.lock.Unlock()

		if err2 != nil {
			sock.Close()
			res := new(remoteResult)
			res.key = this.key
			res.err = err2
			req.ch <- res
		} else {
			sock.Receiver = nil
			this.pool.ReturnSocket(sock)
		}
	}()
}

func (this *remote) CancelRequest(c net.Conn) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for req, sess := range this.reqs {
		if req.conn == c {
			delete(this.reqs, req)
			sess.sock.Close()
		}
	}
}
