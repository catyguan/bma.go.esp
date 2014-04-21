package espservice

import (
	"bmautil/qexec"
	"bmautil/socket"
	"bytes"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
)

type QExecService struct {
	name string

	executor *qexec.QueueExecutor
	rhandler ServiceHandler
	shandler qexec.StopHandler
}

func NewQExecService(name string, rhandler ServiceHandler, shandler qexec.StopHandler) *QExecService {
	this := new(QExecService)
	this.Init(name, rhandler, shandler)
	return this
}

func (this *QExecService) Init(name string, rhandler ServiceHandler, shandler qexec.StopHandler) {
	this.name = name

	this.executor = qexec.NewQueueExecutor(tag, 32, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.rhandler = rhandler
	this.shandler = shandler
}

func (this *QExecService) Run() bool {
	return this.executor.Run()
}

func (this *QExecService) requestHandler(req interface{}) (bool, error) {
	switch v := req.(type) {
	case func() error:
		return true, v()
	case *ServiceRequestContext:
		ctrl := esnp.FrameCoders.Trace
		p := v.Message.ToPackage()
		if v.Sock != nil && ctrl.Has(p) {
			info := fmt.Sprintf("%s handle", this)
			rmsg := ctrl.CreateReply(v.Message, info)
			go v.Sock.SendMessage(rmsg, nil)
		}
		if this.rhandler != nil {
			err := DoServiceHandle(this.rhandler, v.Sock, v.Message)
			return true, err
		} else {
			logger.Debug(tag, "%s miss executor", this.name)
		}
	}
	return true, nil
}

func (this *QExecService) stopHandler() {
	if this.shandler != nil {
		func() {
			defer func() {
				recover()
			}()
			this.shandler()
		}()
	}
}

func (this *QExecService) Name() string {
	return this.name
}

func (this *QExecService) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(this.name)
	buf.WriteString("(qexec)")
	return buf.String()
}

func (this *QExecService) SetRequestReceiver(r ServiceHandler) error {
	if this.executor.IsRun() {
		return this.executor.DoNow("SetRequestReceiver", func() error {
			this.rhandler = r
			return nil
		})
	} else {
		this.rhandler = r
		return nil
	}
}

func (this *QExecService) PostRequest(sock *espsocket.Socket, msg *esnp.Message) error {
	ctx := &ServiceRequestContext{sock, msg}
	return this.executor.DoNow("postRequest", ctx)
}

func (this *QExecService) Stop() bool {
	this.AskClose()
	return true
}

func (this *QExecService) AskClose() {
	this.executor.Stop()
}

func (this *QExecService) Cleanup() bool {
	return this.WaitStop()
}

func (this *QExecService) WaitStop() bool {
	return this.executor.WaitStop()
}

func (this *QExecService) AcceptESP(sock *socket.Socket) error {
	ch := espsocket.NewSocketChannel(sock, "")
	s := espsocket.NewSocket(ch)
	return this.AcceptSocket(s)
}

func (this *QExecService) AcceptSocket(sock *espsocket.Socket) error {
	sock.SetMessageListner(func(msg *esnp.Message) error {
		return this.PostRequest(sock, msg)
	})
	return nil
}
