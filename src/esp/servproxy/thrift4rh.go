package servproxy

import (
	"bmautil/conndialpool"
	"bmautil/socket"
	"bmautil/valutil"
	"bytes"
	"fmt"
	"io"
	"logger"
	"time"
)

type thriftRemoteParam struct {
	PoolMax  int
	PoolInit int
}

type thriftRemoteData struct {
	params *thriftRemoteParam
	pool   *conndialpool.DialPool
}

func (this ThriftProxyHandler) Ping(remote *RemoteObj) (bool, bool) {
	if remote.Data == nil {
		return true, false
	}
	tdata, ok := remote.Data.(*thriftRemoteData)
	if !ok {
		return true, false
	}
	if tdata.pool == nil {
		return true, false
	}
	if tdata.pool.GetInitSize() > 0 {
		return true, tdata.pool.ActiveConn() > 0
	}
	return false, false
}

func (this ThriftProxyHandler) Forward(port *PortObj, preq interface{}, remote *RemoteObj) (rerr error) {
	req, ok := preq.(*ThriftProxyReq)
	if !ok {
		return fmt.Errorf("only forward Thrift Request, can't foward %T", preq)
	}
	if remote.Data == nil {
		return fmt.Errorf("remote(%s) pool not create", remote.name)
	}
	tdata, ok := remote.Data.(*thriftRemoteData)
	if !ok {
		return fmt.Errorf("remote(%s) pool invalid", remote.name)
	}
	tm := remote.cfg.TimeoutMS
	if tm <= 0 {
		tm = 5000
	}
	sock, errSock := tdata.pool.GetSocket(time.Duration(tm)*time.Millisecond, true)
	if errSock != nil {
		return this.AnswerError(port, preq, errSock)
	}
	defer func() {
		if rerr != nil {
			sock.Close()
		} else {
			tdata.pool.ReturnSocket(sock)
		}
	}()
	ch := make(chan []byte, 1)
	defer close(ch)

	sock.Receiver = func(s *socket.Socket, data []byte) (err error) {
		func() {
			x := recover()
			if x != nil {
				err = io.EOF
			}
		}()
		ch <- data
		return nil
	}
	sock.AddCloseListener(func(s *socket.Socket) {
		ch <- nil
	}, "ph0thrift")
	defer sock.RemoveCloseListener("ph0thrift")
	sock.Trace = 32

	if true {
		buf := bytes.NewBuffer(make([]byte, 0, req.hsize))
		this.writeMessageHeader(buf, req.name, req.typeId, req.seqId)
		sz := req.size
		if req.hsize != buf.Len() {
			sz = sz - req.hsize + buf.Len()
		}
		buf2 := bytes.NewBuffer(make([]byte, 0, 4+buf.Len()))
		err1 := this.writeFrameInfo(buf2, sz)
		if err1 != nil {
			return this.AnswerError(port, preq, err1)
		}
		buf2.Write(buf.Bytes())
		err2 := sock.Write(socket.NewWriteReqB(buf2.Bytes(), nil))
		if err2 != nil {
			return this.AnswerError(port, preq, err2)
		}
	}
	buf := make([]byte, 4*1024)
	for {
		sz := req.Remain()
		if sz == 0 {
			break
		}
		if sz > 4*1024 {
			sz = 4 * 1024
		}
		rbuf := buf[:sz]
		n, err := req.conn.Read(rbuf)
		if n > 0 {
			req.readed += n
			// write it
			err3 := sock.Write(socket.NewWriteReqB(rbuf[:n], nil))
			if err3 != nil {
				return this.AnswerError(port, preq, err3)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return this.AnswerError(port, preq, err)
		}
	}
	logger.Debug(tag, "thrift '%s' send end", req)
	if req.IsOneway() {
		req.responsed = true
		return nil
	}
	rsize := 0
	recbuf := bytes.NewBuffer([]byte{})
	var err error
	for {
		rdata := <-ch
		if rdata == nil {
			return io.ErrUnexpectedEOF
		}
		_, err = req.conn.Write(rdata)
		if err != nil {
			return err
		}
		recbuf.Write(rdata)
		if recbuf.Len() < 4 {
			continue
		}
		rsize, err = this.readFrameInfo(recbuf)
		if err != nil {
			return nil
		}
		rsize = rsize - recbuf.Len()
		break
	}
	for {
		if rsize <= 0 {
			break
		}
		logger.Debug(tag, "waiting response %d", rsize)
		rdata := <-ch
		if rdata == nil {
			return io.ErrUnexpectedEOF
		}
		_, err = req.conn.Write(rdata)
		if err != nil {
			return err
		}
		rsize = rsize - len(rdata)
	}
	logger.Debug(tag, "thrift '%s' forward end", req)
	return nil
}

func (this ThriftProxyHandler) Valid(cfg *RemoteConfigInfo) error {
	if cfg.Host == "" {
		return fmt.Errorf("Host invalid")
	}
	return nil
}

func (this ThriftProxyHandler) Compare(cfg *RemoteConfigInfo, old *RemoteConfigInfo) bool {
	p1 := new(thriftRemoteParam)
	valutil.ToBean(cfg.Params, p1)
	p2 := new(thriftRemoteParam)
	valutil.ToBean(cfg.Params, p2)

	if p1.PoolMax != p2.PoolMax {
		return false
	}
	if p2.PoolInit != p2.PoolInit {
		return false
	}
	return true
}

func (this ThriftProxyHandler) Start(o *RemoteObj) error {
	rcfg := o.cfg
	p := new(thriftRemoteParam)
	valutil.ToBean(rcfg.Params, p)

	data := new(thriftRemoteData)
	o.Data = data

	data.params = p

	cfg := new(conndialpool.DialPoolConfig)
	cfg.Address = rcfg.Host
	tm := rcfg.TimeoutMS
	if tm <= 0 {
		tm = 5000
	}
	cfg.TimeoutMS = tm
	if p.PoolInit < 0 {
		cfg.InitSize = 0
	} else {
		cfg.InitSize = p.PoolInit
	}
	if p.PoolMax <= 0 {
		cfg.MaxSize = 10
	} else {
		cfg.MaxSize = p.PoolMax
	}
	err := cfg.Valid()
	if err != nil {
		return err
	}
	pool := conndialpool.NewDialPool(fmt.Sprintf("%s_remote", o.name), cfg)
	data.pool = pool
	if !pool.StartAndRun() {
		return fmt.Errorf("start remote pool fail")
	}
	return nil
}

func (this ThriftProxyHandler) Stop(o *RemoteObj) error {
	if o.Data == nil {
		return nil
	}
	data, ok := o.Data.(*thriftRemoteData)
	if !ok {
		return nil
	}
	if data.pool != nil {
		data.pool.AskClose()
	}
	return nil
}
