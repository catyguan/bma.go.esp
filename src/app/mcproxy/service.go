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
)

type Service struct {
	name   string
	config *configInfo

	poolId  int
	remotes map[string]*remote
	plock   sync.RWMutex
}

type connInfo struct {
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	this.remotes = make(map[string]*remote)
	return this
}

func (this *Service) HandleMemcacheCommand(c net.Conn, cmd *mcserver.MemcacheCommand) (mcserver.HandleCode, error) {
	switch cmd.Action {
	case "restart":
		boot.Restart()
		c.Write([]byte("DONE\r\n"))
		return mcserver.DONE, nil
	case "version":
		c.Write([]byte("VERSION " + this.config.Version + "\r\n"))
		return mcserver.DONE, nil
	case "get":
		res, err := this.executeGet(c, cmd)
		if err != nil {
			return mcserver.DONE, err
		}
		if res.resp {
			c.Write(res.data.Bytes())
		} else {
			c.Write([]byte("END\r\n"))
		}
		return mcserver.DONE, nil
	case "set", "add":
		ok, err := this.executeUpdate(c, cmd, "STORED")
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
		_, err := this.executeUpdate(c, cmd, "")
		if err != nil {
			return mcserver.DONE, err
		}
		msg := "DELETED"
		c.Write([]byte(msg + "\r\n"))
		return mcserver.DONE, nil
	}
	return mcserver.UNKNOW_COMMAND, nil
}

func (this *Service) OnMemcacheConnClose(c net.Conn) {
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

func (this *Service) executeUpdate(conn net.Conn, cmd *mcserver.MemcacheCommand, oks string) (bool, error) {
	f := func(sock *socket.Socket, key string, ch chan *remoteResult) error {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "remote[%s] execute(%s %v)", key, cmd.Action, cmd.Params)
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
			ok, r := coder.DecodeResult()
			if ok {
				res := new(remoteResult)
				res.key = key
				done := false
				if oks == "" || r.Response == oks {
					done = true
				}
				res.resp = done
				ch <- res
				return nil
			}
		}
	}
	return this.executeAll(conn, cmd, f)
}

func (this *Service) executeAll(conn net.Conn, cmd *mcserver.MemcacheCommand, exec remoteExecutor) (bool, error) {
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "begin executeAll(%s %v)", cmd.Action, cmd.Params)
	}
	tmp := make(map[string]*remote)
	this.plock.RLock()
	if true {
		for k, v := range this.remotes {
			tmp[k] = v
		}
	}
	this.plock.RUnlock()

	ch := make(chan *remoteResult, len(tmp))
	for _, rmt := range tmp {
		req := new(remoteRequest)
		req.conn = conn
		req.ch = ch
		req.execute = exec
		rmt.PostRequest(req)
	}

	c := len(tmp)
	count := 0
	for count < c {
		select {
		case res := <-ch:
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "remote[%s] execute done -> %v, %v", res.key, res.resp, res.err)
			}
			if _, ok := tmp[res.key]; ok {
				delete(tmp, res.key)
				count = count + 1
			}
			logger.Debug(tag, "request status = %d/%d", count, c)

			if res.err != nil {
				return false, res.err
			}

			if !res.resp {
				logger.Info(tag, "remote[%s] execute(%s) fail", cmd.Action)
				return false, nil
			}
		}
	}

	return true, nil
}

func (this *Service) executeGet(conn net.Conn, cmd *mcserver.MemcacheCommand) (*remoteResult, error) {
	f := func(sock *socket.Socket, key string, ch chan *remoteResult) error {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "remote[%s] execute(%s %v)", key, cmd.Action, cmd.Params)
		}
		rd := reader(sock)
		defer clearReader(sock)

		err0 := writeCmd(sock, cmd)
		if err0 != nil {
			return err0
		}

		res := new(remoteResult)
		res.key = key
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
				isErr, errMsg := r.ToError()
				if isErr {
					res.err = errors.New(errMsg)
					ch <- res
					return nil
				}
				if res.data == nil {
					res.data = bytes.NewBuffer([]byte{})
				}
				res.data.WriteString(r.Response)
				if len(r.Params) > 0 {
					res.data.WriteString(" ")
					res.data.WriteString(strings.Join(r.Params, " "))
				}
				res.data.WriteString("\r\n")
				if r.Data != nil {
					res.data.Write(r.Data)
					res.data.WriteString("\r\n")
				}
				if r.Response == "END" {
					res.resp = true
					ch <- res
					return nil
				}
			}
		}
	}
	return this.executeOne(conn, cmd, f)
}

func (this *Service) executeOne(conn net.Conn, cmd *mcserver.MemcacheCommand, exec remoteExecutor) (*remoteResult, error) {
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "begin executeOne(%s %v)", cmd.Action, cmd.Params)
	}
	tmp := make(map[string]*remote)
	this.plock.RLock()
	if true {
		for k, v := range this.remotes {
			tmp[k] = v
		}
	}
	this.plock.RUnlock()

	ch := make(chan *remoteResult, len(tmp))
	for _, rmt := range tmp {
		req := new(remoteRequest)
		req.conn = conn
		req.ch = ch
		req.execute = exec
		rmt.PostRequest(req)
	}

	c := len(tmp)
	count := 0
	for count < c {
		select {
		case res := <-ch:
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "remote[%s] execute done -> %v, %v", res.key, res.resp, res.err)
			}
			if _, ok := tmp[res.key]; ok {
				delete(tmp, res.key)
				count = count + 1
			}

			if res.err != nil {
				logger.Debug(tag, "remote[%s] execute error - %s", res.key, res.err)
				continue
			}

			if !res.resp {
				logger.Debug(tag, "remote[%s] execute fail", res.key)
				continue
			}

			return res, nil
		}
	}
	return nil, fmt.Errorf("remotes all fail")
}
