package tankbat

import (
	"bmautil/qexec"
	"config"
	"esp/espnet"
	"logger"
)

type ServiceConfig struct {
	GamePlayerMin int
	GamePlayerMax int
}

type Service struct {
	name     string
	config   *ServiceConfig
	executor qexec.QueueExecutor
	channels map[uint32]*ServiceChannel
	matrix   *Matrix
}

func NewService(n string) *Service {
	this := new(Service)
	this.executor.InitQueueExecutor(tag, 32, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.channels = make(map[uint32]*ServiceChannel)
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
	for k, sch := range this.channels {
		delete(this.channels, k)
		ch := sch.channel
		if ch != nil {
			ch.SetCloseListener("app.service", nil)
			ch.AskClose()
		}
	}
}

func (this *Service) Init() bool {
	cfg := ServiceConfig{}
	if config.GetBeanConfig(this.name, &cfg) {
		this.config = &cfg
		return true
	}
	logger.Error(tag, "invalid config")
	return false
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
	sch := new(ServiceChannel)
	sch.channel = ch

	this.channels[ch.Id()] = sch
	ch.SetCloseListener("app.service", func() {
		go this.OnChannelClose(ch.Id())
	})
	ch.SetMessageListner(func(msg *espnet.Message) error {
		return this.OnChannelMessage(sch, msg)
	})
}

func (this *Service) OnChannelClose(cid uint32) {
	this.executor.DoNow("onChannelClose", func() error {
		this.doChannelClose(cid)
		return nil
	})
}

func (this *Service) doChannelClose(cid uint32) {
	sch, ok := this.channels[cid]
	if !ok {
		return
	}
	delete(this.channels, cid)
	ch := sch.channel
	if ch != nil {
		ch.SetMessageListner(nil)
	}
}

func (this *Service) OnChannelMessage(sch *ServiceChannel, msg *espnet.Message) error {
	return this.executor.DoNow("onChannelMessage", func() error {
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
