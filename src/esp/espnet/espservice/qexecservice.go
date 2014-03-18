package espservice

import (
	"bmautil/qexec"
	"bmautil/socket"
	"bmautil/valutil"
	"bytes"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"fmt"
	"logger"
	"sync"
)

type QExecService struct {
	name string

	executor *qexec.QueueExecutor
	rhandler ServiceHandler
	shandler qexec.StopHandler

	lock       sync.Mutex
	properties map[string]interface{}
	channels   espchannel.VChannelGroup
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

	this.properties = make(map[string]interface{})
}

func (this *QExecService) checkInterfaces() {
	espchannel.SupportProp(this).GetProperty("test")
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
		if v.Channel != nil && ctrl.Has(p) {
			info := fmt.Sprintf("%s handle", this)
			rmsg := ctrl.CreateReply(v.Message, info)
			go v.Channel.PostMessage(rmsg)
		}
		if this.rhandler != nil {
			return true, this.rhandler(v.Channel, v.Message)
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
	this.channels.OnClose()
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

func (this *QExecService) GetProperty(name string) (interface{}, bool) {
	if name == PROP_QEXEC_DEBUG {
		return this.executor.EDebug, true
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	r, ok := this.properties[name]
	return r, ok
}

func (this *QExecService) SetProperty(name string, val interface{}) bool {
	if name == PROP_QEXEC_DEBUG {
		this.executor.EDebug = valutil.ToBool(val, false)
		return true
	}
	if name == PROP_QEXEC_QUEUE_SIZE {
		sz := valutil.ToInt(val, 0)
		if sz <= 0 {
			return false
		}
		this.executor.Resize(sz)
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.properties[name] = val
	return true
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

func (this *QExecService) PostRequest(ch espchannel.Channel, msg *esnp.Message) error {
	ctx := &ServiceRequestContext{ch, msg}
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
	ch := espchannel.NewSocketChannel(sock, espchannel.SOCKET_CHANNEL_CODER_ESPNET)
	ConnectService(ch, this.PostRequest)
	return nil
}

// QExecService's Channel
func (this *QExecService) NewChannel() (espchannel.Channel, error) {
	r := new(espchannel.VChannel)
	r.InitVChannel(this.name)
	r.RemoveChannel = this.channels.Remove

	r.Sender = func(msg *esnp.Message) error {
		return DoServiceHandle(this.PostRequest, r, msg)
	}
	this.channels.Add(r)
	return r, nil
}
