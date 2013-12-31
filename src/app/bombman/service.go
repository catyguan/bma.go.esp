package bombman

import (
	"bmautil/qexec"
	"esp/espnet"
)

type ServiceChanel struct {
	channel espnet.Channel
}

func (this *ServiceChanel) BeError(err error) {
	if err != nil {
		rmsg := espnet.NewMessage()
		rmsg.BeError(err)
		this.channel.SendMessage(rmsg)
	}
}

func (this *ServiceChanel) Replay(s string) {
	rmsg := espnet.NewMessage()
	rmsg.SetPayload([]byte(s))
	this.channel.SendMessage(rmsg)
}

func (this *ServiceChanel) Send(s string) {
	msg := espnet.NewMessage()
	msg.SetPayload([]byte(s))
	this.channel.SendMessage(msg)
}

type Service struct {
	name     string
	executor qexec.QueueExecutor
	channels map[uint32]espnet.Channel
	matrix   *Matrix
}

func NewService(n string) *Service {
	this := new(Service)
	this.executor.InitQueueExecutor(tag, 32, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.channels = make(map[uint32]espnet.Channel)
	return this
}

func (this *Service) requestHandler(ev interface{}) (bool, error) {
	switch rv := ev.(type) {
	case func() error:
		return true, rv()
	}
	return true, nil
}

func (this *Service) stopHandler() {
	if this.matrix != nil {
		this.matrix.AskClose()
	}
	for k, ch := range this.channels {
		delete(this.channels, k)
		ch.SetCloseListener("bombman.service", nil)
		ch.AskClose()
	}
}

func (this *Service) Init() bool {
	return true
}

func (this *Service) Start() bool {
	this.executor.Run()
	return true
}

func (this *Service) Close() bool {
	this.executor.Stop()
	return true
}

func (this *Service) Cleanup() bool {
	this.executor.WaitStop()
	return true
}

func (this *Service) Add(ch espnet.Channel) {
	err := this.executor.DoNow("addChannel", func() error {
		this.doAdd(ch)
		return nil
	})
	if err != nil {
		ch.AskClose()
	}
}

func (this *Service) doAdd(ch espnet.Channel) {
	this.channels[ch.Id()] = ch
	ch.SetCloseListener("bombman.service", func() {
		go this.OnChannelClose(ch.Id())
	})
	ch.SetMessageListner(func(msg *espnet.Message) error {
		return this.OnChannelMessage(ch, msg)
	})
}

func (this *Service) OnChannelClose(cid uint32) {
	this.executor.DoNow("onChannelClose", func() error {
		this.doChannelClose(cid)
		return nil
	})
}

func (this *Service) doChannelClose(cid uint32) {
	ch, ok := this.channels[cid]
	if !ok {
		return
	}
	delete(this.channels, cid)
	ch.SetMessageListner(nil)
	ch.GetProperty(PROP_PLAYER)
}

func (this *Service) OnChannelMessage(ch espnet.Channel, msg *espnet.Message) error {
	return this.executor.DoNow("onChannelMessage", func() error {
		sch := &ServiceChanel{ch}
		this.doChannelMessage(sch, msg)
		return nil
	})
}

func (this *Service) OnMatrixEnd(m *Matrix) {
	this.executor.DoNow("OnMatrixEnd", func() error {
		if this.matrix == m {
			this.matrix = nil
		}
		return nil
	})

}
