package tankbat

import (
	"esp/espnet"
	"fmt"
	"time"
)

type ServiceChannel struct {
	channel espnet.Channel

	name     string
	waiting  bool
	waitTime time.Time
}

func (this *ServiceChannel) Id() uint32 {
	return this.channel.Id()
}

func (this *ServiceChannel) String() string {
	n := this.name
	if n == "" {
		n = fmt.Sprintf("%d", this.Id())
	}
	return fmt.Sprintf("SC_%s", n)
}

func (this *ServiceChannel) BeError(err error) {
	if err != nil {
		rmsg := espnet.NewMessage()
		rmsg.BeError(err)
		this.channel.SendMessage(rmsg)
	}
}

func (this *ServiceChannel) Reply(s string) error {
	rmsg := espnet.NewMessage()
	rmsg.SetPayload([]byte(s))
	return this.channel.SendMessage(rmsg)
}

func (this *ServiceChannel) ReplyOK() error {
	return this.Reply("OK\n")
}

func (this *ServiceChannel) Send(s string) error {
	msg := espnet.NewMessage()
	msg.SetPayload([]byte(s))
	return this.channel.SendMessage(msg)
}

type ServiceChannels []*ServiceChannel

func (s ServiceChannels) Len() int      { return len(s) }
func (s ServiceChannels) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByPlayTime struct{ ServiceChannels }

func (s ByPlayTime) Less(i, j int) bool {
	return s.ServiceChannels[i].waitTime.Before(s.ServiceChannels[j].waitTime)
}
