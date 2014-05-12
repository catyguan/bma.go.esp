package main

import (
	"bmautil/socket"
	"boot"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"logger"
	"mcserver"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

type Service struct {
	name   string
	config *configInfo

	poolId    int
	remotes   map[string]*remote
	plock     sync.RWMutex
	connCount int32
}

type remoteRequestMR interface {
	HandleResponse(key string, res *mcserver.MemcacheResult) (more bool, done bool, err error)
	CheckEnd(okc, failc, errc, total int) (end bool, done bool, iserr bool)
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	this.remotes = make(map[string]*remote)
	return this
}

func (this *Service) HandleMemcacheCommand(c net.Conn, cmd *mcserver.MemcacheCommand) (mcserver.HandleCode, error) {
	switch cmd.Action {
	case "reload":
		boot.Restart()
		c.Write([]byte("DONE\r\n"))
		return mcserver.DONE, nil
	case "version":
		c.Write([]byte("VERSION " + this.config.Version + "\r\n"))
		return mcserver.DONE, nil
	case "get":
		var mr mrGet
		mr.Init()
		done, err := this.executeMR(c, "", cmd, &mr)
		if err != nil {
			return mcserver.DONE, err
		}
		if done {
			data := mergeResults(mr.results)
			c.Write(data)
		} else {
			c.Write([]byte("END\r\n"))
		}
		return mcserver.DONE, nil
	case "getall":
		cmd2 := new(mcserver.MemcacheCommand)
		cmd2.Action = "get"
		cmd2.Params = cmd.Params

		var mr mrGetAll
		mr.Init()
		_, err := this.executeMR(c, "getall", cmd2, &mr)
		if err != nil {
			return mcserver.DONE, err
		}
		for _, o := range mr.results {
			bb := bytes.NewBuffer([]byte{})
			bb.WriteString(o.key)
			if o.result != nil {
				bb.WriteString(" -> ")
				bb.WriteString(o.result.Response)
				if o.result.Params != nil {
					bb.WriteString(" ")
					bb.WriteString(strings.Join(o.result.Params, " "))
				}
				bb.WriteString("\r\n")
			}
			c.Write(bb.Bytes())
		}
		c.Write([]byte("END\r\n"))
		return mcserver.DONE, nil
	case "set", "add":
		var mr mrUpdate
		mr.oks = "STORED"
		ok, err := this.executeMR(c, "", cmd, &mr)
		if err != nil {
			return mcserver.DONE, err
		}
		msg := "NOT_STORED"
		if ok {
			msg = "STORED"
		}
		c.Write([]byte(msg + "\r\n"))
		return mcserver.DONE, nil
	case "delete":
		var mr mrUpdate
		_, err := this.executeMR(c, "", cmd, &mr)
		if err != nil {
			return mcserver.DONE, err
		}
		msg := "DELETED"
		c.Write([]byte(msg + "\r\n"))
		return mcserver.DONE, nil
	}
	return mcserver.UNKNOW_COMMAND, nil
}

func (this *Service) OnMemcacheConnOpen(c net.Conn) bool {
	count := atomic.AddInt32(&this.connCount, 1)
	logger.Info(tag, "'%s' connected [%d] %s", this.name, count, c.RemoteAddr())
	return true
}

func (this *Service) OnMemcacheConnClose(c net.Conn) {
	count := atomic.AddInt32(&this.connCount, -1)
	logger.Info(tag, "'%s' disconnect[%d] %s", this.name, count, c.RemoteAddr())

	this.plock.RLock()
	defer this.plock.RUnlock()
	for _, rmt := range this.remotes {
		go rmt.CancelRequest(c)
	}
}

func (this *Service) _createRemote(remoteKey string) {
	logger.Debug(tag, "create remote '%s'", remoteKey)

	rmt := new(remote)
	rmt.InitRemote(remoteKey)

	this.poolId = this.poolId + 1
	cfg := new(socket.DialPoolConfig)
	cfg.Dial.Address = remoteKey
	cfg.InitSize = 1
	cfg.MaxSize = this.config.PoolMax

	p := socket.NewDialPool(fmt.Sprintf("mcr%d", this.poolId), cfg, nil)
	p.Start()
	p.Run()
	rmt.pool = p

	this.remotes[remoteKey] = rmt
}

func (this *Service) _closeRemote(k string, p *remote) {
	logger.Debug(tag, "close remote '%s'", k)
	p.Close()
	delete(this.remotes, k)
}

func writeCmd(sock *socket.Socket, cmd *mcserver.MemcacheCommand) error {
	bb := bytes.NewBuffer([]byte{})
	bb.WriteString(cmd.Action)
	bb.WriteString(" ")
	bb.WriteString(strings.Join(cmd.Params, " "))
	bb.WriteString("\r\n")
	if cmd.Data != nil {
		bb.Write(cmd.Data)
		bb.WriteString("\r\n")
	}
	return sock.Write(socket.NewWriteReqB(bb.Bytes(), nil))
}
func mergeResults(results []*mcserver.MemcacheResult) []byte {
	data := bytes.NewBuffer([]byte{})
	for _, r := range results {
		data.WriteString(r.Response)
		if len(r.Params) > 0 {
			data.WriteString(" ")
			data.WriteString(strings.Join(r.Params, " "))
		}
		data.WriteString("\r\n")
		if r.Data != nil {
			data.Write(r.Data)
			data.WriteString("\r\n")
		}
	}
	return data.Bytes()
}
func reader(sock *socket.Socket) io.Reader {
	r, w := io.Pipe()
	sock.Receiver = func(sock *socket.Socket, data []byte) error {
		w.Write(data)
		return nil
	}
	sock.AddCloseListener(func(sock *socket.Socket) {
		w.Close()
	}, "service")
	return r
}
func clearReader(sock *socket.Socket) {
	sock.Receiver = nil
	sock.RemoveCloseListener("service")
}

func (this *Service) executeMR(conn net.Conn, act string, cmd *mcserver.MemcacheCommand, mr remoteRequestMR) (bool, error) {
	if act == "" {
		act = cmd.Action
	}
	logger.Info(tag, "%s begin executeMR(%s)", conn.RemoteAddr(), act)
	defer logger.Info(tag, "%s end executeMR(%s)", conn.RemoteAddr(), act)

	tmp := make(map[string]*remote)
	this.plock.RLock()
	if true {
		for k, v := range this.remotes {
			tmp[k] = v
		}
	}
	this.plock.RUnlock()

	exec := func(sock *socket.Socket, key string, ch chan *remoteResult) error {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "remote[%s] execute(%s %v)", key, act, cmd.Params)
		}
		rd := reader(sock)
		defer clearReader(sock)

		err0 := writeCmd(sock, cmd)
		if err0 != nil {
			return err0
		}

		coder := mcserver.NewMemcacheCoder()
		in := bufio.NewReader(rd)
		buf := make([]byte, 1024)
		for {
			n, err := in.Read(buf)
			if err != nil {
				return err
			}
			coder.Write(buf[:n])
			for {
				ok, r := coder.DecodeResult()
				if !ok {
					break
				}
				more, done, err1 := mr.HandleResponse(key, r)
				if more {
					continue
				}
				res := new(remoteResult)
				res.key = key
				res.resp = done
				res.err = err1
				ch <- res
				return nil
			}
		}
	}

	ch := make(chan *remoteResult, len(tmp))
	for _, rmt := range tmp {
		req := new(remoteRequest)
		req.conn = conn
		req.ch = ch
		req.execute = exec
		rmt.PostRequest(req)
	}

	total := len(tmp)
	okc := 0
	failc := 0
	errc := 0
	count := 0
	for count < total {
		select {
		case res := <-ch:
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "remote[%s] execute done -> %v, %v", res.key, res.resp, res.err)
			}
			if _, ok := tmp[res.key]; ok {
				delete(tmp, res.key)
				count = count + 1
			}
			logger.Debug(tag, "request status = %d/%d", count, total)

			if !res.resp {
				logger.Info(tag, "remote[%s] execute(%s) fail", res.key, act)
			}

			if res.resp {
				okc = okc + 1
			} else {
				failc = failc + 1
			}
			if res.err != nil {
				errc = errc + 1
			}

			end, done, iserr := mr.CheckEnd(okc, failc, errc, total)
			if end {
				if iserr {
					return false, res.err
				}
				return done, nil
			}
		}
	}

	logger.Info(tag, "execute(%s) all fail", act)
	return false, errors.New("all fail")
}
