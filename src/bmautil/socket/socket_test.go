package socket

import (
	"logger"
	"net"
	"testing"
	"time"
)

func TestSocketId(t *testing.T) {
	t.Errorf("%p", TestSocketBase)
}

func TestSocketBase(t *testing.T) {

	ln, err := net.Listen("tcp", ":1080")
	if err != nil {
		t.Error(err)
		return
	}

	sg := NewSocketGroup()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			ssock := NewSocket(conn, 8, 0)
			ssock.Trace = 64
			ssock.Receiver = func(sock *Socket, b []byte) error {
				req := NewWriteReqB(b, func(sock *Socket, err error) bool {
					logger.Info("TEST", "echo done")
					return false
				})
				return sock.Write(req)
			}
			ssock.AddCloseListener(func(sock *Socket) {
				logger.Info("TEST", "server socket close")
			}, "")
			ssock.Start(nil)

			sg.Add(ssock)
		}
	}()

	if true {
		conn, err := net.Dial("tcp", "127.0.0.1:1080")
		if err != nil {
			t.Error("client", err)
			// handle error
		} else {
			timeout := time.Duration(5000) * time.Millisecond
			csock := NewSocket(conn, 8, timeout)
			csock.Trace = 64
			csock.Start(nil)
			csock.SetWriteChanSize(16)
			req := NewWriteReqB([]byte{1, 2, 3, 4}, nil)
			csock.Write(req)
			// csock.Close()
			sg.Add(csock)
		}
	}

	time.Sleep(time.Duration(2) * time.Second)
	ln.Close()
	logger.Info("TEST", "end")
	sg.Close()
	time.Sleep(time.Duration(100) * time.Millisecond)
}
