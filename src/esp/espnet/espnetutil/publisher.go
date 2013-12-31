package espnet

import (
	"errors"
	"logger"
)

type PublisherFilter func(ch Channel, msg *Message) bool

type Publisher struct {
	name  string
	group ChannelGroup
	c     chan *Message

	channels VChannelGroup

	Filter PublisherFilter
}

func NewPublisher(name string, bufsz int) *Publisher {
	this := new(Publisher)
	this.name = name
	this.group.InitGroup()
	this.c = make(chan *Message, bufsz)
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
	var list []Channel
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

func (this *Publisher) Add(ch Channel) {
	this.group.Add(ch)
}

func (this *Publisher) Remove(ch Channel) {
	this.group.Remove(ch)
}

func (this *Publisher) SendMessage(msg *Message) (r error) {
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
func (this *Publisher) NewChannel() (Channel, error) {
	r := new(VChannel)
	r.InitVChannel(this.name)
	r.RemoveChannel = this.channels.Remove
	r.Sender = this.SendMessage
	this.channels.Add(r)
	return r, nil
}
