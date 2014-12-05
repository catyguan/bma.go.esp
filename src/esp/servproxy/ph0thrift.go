package servproxy

import (
	"bmautil/socket"
	"bmautil/valutil"
	"bytes"
	"encoding/binary"
	"fmt"
	"golua"
	"io"
	"logger"
	"net"
	"time"
)

const (
	VERSION_MASK = 0xffff0000
	VERSION_1    = 0x80010000
)

type ThriftProxyHandler int

func init() {
	AddProxyHandler("thrift", ThriftProxyHandler(0))
}

func (this ThriftProxyHandler) readFrameInfo(r io.Reader) (int, error) {
	buf := []byte{0, 0, 0, 0}
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	size := int(binary.BigEndian.Uint32(buf))
	if size < 0 {
		return 0, fmt.Errorf("Read a negative frame size (%d)", size)
	}
	return size, nil
}

func (this ThriftProxyHandler) readByte(r io.Reader) (value byte, n int, err error) {
	buf := []byte{0}
	n, err = r.Read(buf)
	return buf[0], n, err
}

func (this ThriftProxyHandler) readI32(r io.Reader) (value int32, n int, err error) {
	buf := []byte{0, 0, 0, 0}
	n, err = io.ReadFull(r, buf)
	if err != nil {
		return 0, 0, err
	}
	value = int32(binary.BigEndian.Uint32(buf))
	return value, n, nil
}

func (this ThriftProxyHandler) readStringBody(r io.Reader, size int) (value string, n int, err error) {
	if size < 0 {
		return "", 0, nil
	}
	isize := int(size)
	buf := make([]byte, isize)
	n, e := io.ReadFull(r, buf)
	if e != nil {
		return "", 0, e
	}
	return string(buf), n, nil
}

func (this ThriftProxyHandler) readString(r io.Reader) (value string, n int, err error) {
	size, l, e := this.readI32(r)
	if e != nil {
		return "", 0, e
	}
	s, l2, e2 := this.readStringBody(r, int(size))
	if e2 != nil {
		return "", 0, e2
	}
	return s, l2 + l, nil
}

func (this ThriftProxyHandler) readMessageHeader(r io.Reader) (name string, typeId int32, seqId int32, n int, err error) {
	size, c, e := this.readI32(r)
	if e != nil {
		return "", 0, 0, 0, e
	}
	if size < 0 {
		typeId = int32(size & 0x0ff)
		version := int64(int64(size) & VERSION_MASK)
		if version != VERSION_1 {
			return "", 0, 0, 0, fmt.Errorf("Bad version(%d) in ReadMessageBegin", version)
		}
		l := 0
		name, l, e = this.readString(r)
		if e != nil {
			return "", 0, 0, 0, e
		}
		c += l
		seqId, l, e = this.readI32(r)
		if e != nil {
			return "", 0, 0, 0, e
		}
		c += l
		return name, typeId, seqId, c, nil
	}
	name, l2, e2 := this.readStringBody(r, int(size))
	if e2 != nil {
		return "", 0, 0, 0, e2
	}
	c += l2
	b, l3, e3 := this.readByte(r)
	if e3 != nil {
		return "", 0, 0, 0, e3
	}
	c += l3
	typeId = int32(b)
	seqId, l4, e4 := this.readI32(r)
	if e4 != nil {
		return "", 0, 0, 0, e4
	}
	c += l4
	return name, typeId, seqId, c, nil
}

func (this ThriftProxyHandler) writeFrameInfo(w io.Writer, sz int) error {
	buf := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(buf, uint32(sz))
	_, err := w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (this ThriftProxyHandler) writeString(w io.Writer, value string) error {
	err := this.writeI32(w, int32(len(value)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(value))
	return err
}

func (this ThriftProxyHandler) writeByte(w io.Writer, value byte) error {
	v := []byte{value}
	_, err := w.Write(v)
	return err
}

func (this ThriftProxyHandler) writeI32(w io.Writer, value int32) error {
	v := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(v, uint32(value))
	_, err := w.Write(v)
	return err
}

func (this ThriftProxyHandler) writeI16(w io.Writer, value int16) error {
	v := []byte{0, 0}
	binary.BigEndian.PutUint16(v, uint16(value))
	_, e := w.Write(v)
	return e
}

func (this ThriftProxyHandler) writeMessageHeader(w io.Writer, name string, typeId int32, seqId int32) error {
	if true {
		version := uint32(VERSION_1) | uint32(typeId)
		e := this.writeI32(w, int32(version))
		if e != nil {
			return e
		}
		e = this.writeString(w, name)
		if e != nil {
			return e
		}
		e = this.writeI32(w, seqId)
		return e
	} else {
		e := this.writeString(w, name)
		if e != nil {
			return e
		}
		e = this.writeByte(w, byte(typeId))
		if e != nil {
			return e
		}
		e = this.writeI32(w, seqId)
		return e
	}
}

func (this ThriftProxyHandler) writeError(w io.Writer, req *ThriftProxyReq, err error) {
	str := err.Error()
	// err = oprot.WriteFieldBegin("message", STRING, 1)
	this.writeByte(w, 11)
	this.writeI16(w, 1)
	this.writeString(w, str)

	// err = oprot.WriteFieldBegin("type", I32, 2)
	this.writeByte(w, 8)
	this.writeI16(w, 2)
	this.writeI32(w, 6)

	// err = oprot.WriteFieldStop()
	this.writeByte(w, 0)
}

func (this ThriftProxyHandler) Handle(s *Service, port *PortObj, conn net.Conn) {
	defer conn.Close()
	dc := &DebugConn{conn}
	rconn := dc
	for {
		sz, err := this.readFrameInfo(rconn)
		if err != nil {
			if err == io.EOF {
				logger.Debug(tag, "%s closed", dc)
			} else {
				logger.Warn(tag, "%s readFrameInfo fail - %s", dc, err)
			}
			return
		}
		logger.Debug(tag, "%s readFrameInfo - %d", dc, sz)
		req := new(ThriftProxyReq)
		req.s = s
		req.conn = rconn

		req.size = sz
		req.readed = 0
		name, tid, sid, n, err1 := this.readMessageHeader(rconn)
		if err1 != nil {
			logger.Warn(tag, "%s readMessageHeader fail - %s", dc, err1)
			return
		}
		logger.Debug(tag, "%s readMessageHeader - %s, %d, %d/%d", dc, name, tid, sid, n)
		req.name = name
		req.typeId = tid
		req.seqId = sid
		req.hsize = n
		req.readed += n

		_, errE := s.Execute(port, golua.NewGOO(req, gooThriftProxyReq(0)), req)
		if errE != nil {
			logger.Warn(tag, "%s execute fail - %s", dc, errE)
			return
		}
	}
}

func (this ThriftProxyHandler) AnswerError(port *PortObj, preq interface{}, err error) error {
	req, ok := preq.(*ThriftProxyReq)
	if ok && !req.responsed {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "answer %s error - %s", req.conn.RemoteAddr(), err)
		}
		req.responsed = true
		buf := bytes.NewBuffer([]byte{})
		this.writeMessageHeader(buf, "", 3, req.seqId)
		this.writeError(buf, req, err)
		buf2 := bytes.NewBuffer(make([]byte, 0, 4+buf.Len()))
		this.writeFrameInfo(buf2, buf.Len())
		buf2.Write(buf.Bytes())
		req.conn.Write(buf2.Bytes())
	}
	return err
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

type thriftRemoteParam struct {
	PoolMax  int
	PoolInit int
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

type thriftRemoteData struct {
	params *thriftRemoteParam
	pool   *socket.DialPool
}

func (this ThriftProxyHandler) Start(o *RemoteObj) error {
	rcfg := o.cfg
	p := new(thriftRemoteParam)
	valutil.ToBean(rcfg.Params, p)

	data := new(thriftRemoteData)
	o.Data = data

	data.params = p

	cfg := new(socket.DialPoolConfig)
	cfg.Dial.Address = rcfg.Host
	tm := rcfg.TimeoutMS
	if tm <= 0 {
		tm = 5000
	}
	cfg.Dial.TimeoutMS = tm
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
	pool := socket.NewDialPool(fmt.Sprintf("%s_remote", o.name), cfg, func(sock *socket.Socket) error {
		return nil
	})
	data.pool = pool
	if !pool.Start() {
		return fmt.Errorf("start remote pool fail")
	}
	if !pool.Run() {
		return fmt.Errorf("run remote pool fail")
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
