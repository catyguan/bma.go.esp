package espnetutil

import (
	"errors"
	"esp/espnet"
	"logger"
)

type PublisherFilter func(ch espnet.Channel, msg *espnet.Message) bool

type Publisher struct {
	name  string
	group espnet.ChannelGroup
	c     chan *espnet.Message

	channels espnet.VChannelGroup

	Filter PublisherFilter
}

func NewPublisher(name string, bufsz int) *Publisher {
	this := new(Publisher)
	this.name = name
	this.group.InitGroup()
	this.c = make(chan *espnet.Message, bufsz)
	return this
}

func (this *Publisher) Name() string {
	return this.name
}

func (this *Publisher) String() string {
	return "Publisher[" + this.name + "]"
}

func (this *Publisher) Start() bool {
	go this.run()
	return true
}

func (this *Publisher) run() {
	defer func() {
		close(this.c)
		this.channels.OnClose()
	}()
	var list []espnet.Channel
	var mark int64
	for {
		msg := <-this.c
		if msg == nil {
			return
		}
		if this.group.IsClosing() {
			return
		}
		// send message to all c
		list, mark = this.group.Snapshot(mark, list)
		if list != nil {
			for _, ch := range list {
				if this.Filter != nil && !this.Filter(ch, msg) {
					logger.Debug(tag, "%s ---> (%s) skip", this, ch)
					continue
				}
				logger.Debug(tag, "%s ---> (%s)", this, ch)
				err := ch.SendMessage(msg)
				if err != nil {
					logger.Debug(tag, "%s ---> (%s) fail - %s", this, ch, err)
				}
			}
		}
	}
}

func (this *Publisher) Stop() bool {
	if this.group.AskClose() {
		this.c <- nil
	}
	return true
}

func (this *Publisher) Cleanup() bool {
	this.WaitClose()
	return true
}

func (this *Publisher) WaitClose() {
	this.group.WaitClosed()
}

func (this *Publisher) Add(ch espnet.Channel) {
	this.group.Add(ch)
}

func (this *Publisher) Remove(ch espnet.Channel) {
	this.group.Remove(ch)
}

func (this *Publisher) SendMessage(msg *espnet.Message) (r error) {
	if this.group.IsClosing() {
		return errors.New("closed")
	}
	defer func() {
		if recover() != nil {
			r = errors.New("closed")
		}
	}()
	this.c <- msg
	return nil
}

// Publisher's Channel
func (this *Publisher) NewChannel() (espnet.Channel, error) {
	r := new(espnet.VChannel)
	r.InitVChannel(this.name)
	r.RemoveChannel = this.channels.Remove
	r.Sender = this.SendMessage
	this.channels.Add(r)
	return r, nil
}
