package tankbat

import (
	"bmautil/qexec"
	"config"
	"esp/espnet"
	"logger"
	"time"
)

type ServiceConfig struct {
	GamePlayerMin int
	TeamPlayerMax int
}

type WUserInfo struct {
	channel  *ServiceChannel
	waitTime time.Time
}

type WUserInfoList []*WUserInfo

func (s WUserInfoList) Len() int      { return len(s) }
func (s WUserInfoList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s WUserInfoList) Less(i, j int) bool {
	return s[i].waitTime.Before(s[j].waitTime)
}

type Service struct {
	name     string
	config   *ServiceConfig
	executor qexec.QueueExecutor
	channels map[uint32]*ServiceChannel
	matrix   *Matrix

	waitingRoom map[uint32]*WUserInfo
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	this.executor.InitQueueExecutor(tag, 32, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.channels = make(map[uint32]*ServiceChannel)
	this.waitingRoom = make(map[uint32]*WUserInfo)
	return this
}

func (this *Service) Name() string {
	return this.name
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
	for k, _ := range this.waitingRoom {
		delete(this.waitingRoom, k)
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
	sch.id = ch.Id()

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
	sch.Close()
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
		this.doMatrixEnd(m)
		return nil
	})

}
