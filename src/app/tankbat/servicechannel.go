package tankbat

import (
	"esp/espnet"
	"fmt"
	"time"
)

type ServiceChannel struct {
	id      uint32
	closed  bool
	channel espnet.Channel
	cmdTime time.Time

	name       string
	joinTeamId int // 0:unknow, 1:A, 2:B
	playing    bool
}

func (this *ServiceChannel) Id() uint32 {
	return this.id
}

func (this *ServiceChannel) IsClosed() bool {
	return this.closed
}

func (this *ServiceChannel) Close() {
	this.closed = true
}

func (this *ServiceChannel) String() string {
	n := this.name
	if n == "" {
		n = fmt.Sprintf("%d", this.Id())
	}
	return fmt.Sprintf("SC_%s", n)
}

func (this *ServiceChannel) BeError(err error) {
	if err != nil && !this.closed {
		rmsg := espnet.NewMessage()
		rmsg.BeError(err)
		if this.channel.SendMessage(rmsg) != nil {
			this.Close()
		}
	}
}

func (this *ServiceChannel) Reply(s string) error {
	if !this.closed {
		rmsg := espnet.NewMessage()
		rmsg.SetPayload([]byte(s))
		err := this.channel.SendMessage(rmsg)
		if err != nil {
			this.Close()
		}
		return err
	}
	return nil
}

func (this *ServiceChannel) ReplyOK() error {
	return this.Reply("OK\n")
}

func (this *ServiceChannel) Send(s string) error {
	if !this.closed {
		msg := espnet.NewMessage()
		msg.SetPayload([]byte(s))
		err := this.channel.SendMessage(msg)
		if err != nil {
			this.Close()
		}
	}
	return nil
}
