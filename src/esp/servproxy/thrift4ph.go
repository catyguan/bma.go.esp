package servproxy

import (
	"bmautil/connutil"
	"bytes"
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

type ThriftPortHandler int

func init() {
	AddPortHandler("thrift", ThriftPortHandler(0))
}

func (this ThriftPortHandler) Handle(s *Service, port *PortObj, conn net.Conn) {
	defer conn.Close()
	var dbg connutil.ConnDebuger
	if debugTraffic {
		dbg = ConnDebuger
	}
	rconn := connutil.NewConnExt(conn, dbg)
	for {
		sz, err := OThriftProtocol.readFrameInfo(rconn)
		if err != nil {
			if err == io.EOF {
				logger.Debug(tag, "%s closed", rconn)
			} else {
				logger.Warn(tag, "%s readFrameInfo fail - %s", rconn, err)
			}
			return
		}
		logger.Debug(tag, "%s readFrameInfo - %d", rconn, sz)
		req := new(ThriftProxyReq)
		req.s = s
		req.conn = rconn

		req.size = sz
		req.readed = 0
		name, tid, sid, n, err1 := OThriftProtocol.readMessageHeader(rconn)
		if err1 != nil {
			logger.Warn(tag, "%s readMessageHeader fail - %s", rconn, err1)
			return
		}
		logger.Debug(tag, "%s readMessageHeader - %s, %d, %d/%d", rconn, name, tid, sid, n)
		req.name = name
		req.typeId = tid
		req.seqId = sid
		req.hsize = n
		req.readed += n

		_, errE := s.Execute(port, golua.NewGOO(req, gooThriftProxyReq(0)), req)
		if errE != nil {
			logger.Warn(tag, "%s execute fail - %s", rconn, errE)
			return
		}
	}
}

func (this ThriftPortHandler) AnswerError(port *PortObj, preq ProxyRequest, err error) error {
	req, ok := preq.(*ThriftProxyReq)
	if ok && !req.responsed {
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "answer %s error - %s", req.conn.RemoteAddr(), err)
		}
		req.responsed = true
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		OThriftProtocol.writeMessageHeader(buf, "", 3, req.seqId)
		OThriftProtocol.writeError(buf, req, err)
		l := buf.Len()
		buf2 := bytes.NewBuffer(make([]byte, 0, 4+l))
		OThriftProtocol.writeFrameInfo(buf2, l)
		buf2.Write(buf.Bytes())
		req.conn.Write(buf2.Bytes())
		time.Sleep(1 * time.Millisecond)
	}
	return err
}

func (this ThriftPortHandler) BeginWrite(port *PortObj, preq ProxyRequest) error {
	req, ok := preq.(*ThriftProxyReq)
	if !ok {
		return fmt.Errorf("unknow request(%T)", req)
	}
	if req.responsed {
		return fmt.Errorf("request already responsed")
	}
	return nil
}

func (this ThriftPortHandler) Write(port *PortObj, preq ProxyRequest, b []byte) error {
	req, ok := preq.(*ThriftProxyReq)
	if !ok {
		return nil
	}
	req.responsed = true
	_, err := req.conn.Write(b)
	return err
}

func (this ThriftPortHandler) EndWrite(port *PortObj, preq ProxyRequest) {
	return
}
