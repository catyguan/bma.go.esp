package espservice

import (
	"bmautil/connutil"
	"bytes"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"io"
	"logger"
	"net"
	"runtime/debug"
)

type GoService struct {
	name    string
	handler ServiceHandler
}

func NewGoService(name string, h ServiceHandler) *GoService {
	this := new(GoService)
	this.name = name
	this.handler = h
	return this
}

func (this *GoService) Name() string {
	return this.name
}

func (this *GoService) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(this.name)
	buf.WriteString("(gos)")
	return buf.String()
}

func (this *GoService) AcceptConn(debug string) func(conn net.Conn) {
	return func(conn net.Conn) {
		ct := connutil.NewConnExt(conn)
		if debug != "" {
			ct.Debuger = connutil.SimpleDebuger(debug)
		}
		sock := espsocket.NewConnSocket(ct, 0)
		this.Serve(sock)
	}
}

func (this *GoService) Serve(sock espsocket.Socket) {
	defer sock.AskFinish()
	this.DoServe(sock)
}

func (this *GoService) DoServe(sock espsocket.Socket) {
	for {
		msg, err := sock.ReadMessage()
		if err != nil {
			if err == io.EOF {
				logger.Debug(tag, "%s closed", sock)
				return
			}
			sock.AskClose()
			return
		}
		DoServiceHandle(this.PostRequest, sock, msg)
	}
}

func (this *GoService) PostRequest(sock espsocket.Socket, msg *esnp.Message) (rerr error) {
	ctrl := esnp.MessageLineCoders.Trace
	p := msg
	if ctrl.Has(p) {
		info := fmt.Sprintf("%s handled", this)
		rmsg := ctrl.CreateReply(msg, info)
		sock.WriteMessage(rmsg)
	}
	defer func() {
		err := recover()
		if err != nil {
			logger.Warn(tag, "execute panic - %s\n%s", err, string(debug.Stack()))
			if terr, ok := err.(error); ok {
				rerr = terr
			} else {
				rerr = fmt.Errorf("%s", err)
			}
		}
	}()
	return this.handler(sock, msg)
}
