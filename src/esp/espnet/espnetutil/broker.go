package espnet

import (
	"errors"
	"fmt"
	"logger"
	"uuid"
)

type BrokerDispatch func(in Channel, list []Channel, msg *Message, left2right bool) (Channel, error)

type brokerREQ struct {
	in  Channel
	msg *Message
	l2r bool
}

type Broker struct {
	id    string
	name  string
	left  ChannelGroup
	right ChannelGroup
	c     chan *brokerREQ

	LeftDispatcher  BrokerDispatch
	RightDispatcher BrokerDispatch

	// runtime
	lcount int
	rcount int
}

func NewBroker(name string, bufsz int) *Broker {
	this := new(Broker)
	uuid, _ := uuid.NewV4()
	this.id = uuid.String()
	this.name = name
	this.left.InitGroup()
	this.right.InitGroup()
	this.c = make(chan *brokerREQ, bufsz)
	return this
}

func (this *Broker) Id() string {
	return this.id
}

func (this *Broker) Name() string {
	return this.name
}

func (this *Broker) String() string {
	return "Broker[" + this.name + "]"
}

func (this *Broker) Start() bool {
	go this.run()
	return true
}

func (this *Broker) DefaultDispatch(in Channel, list []Channel, msg *Message, left2right bool) (Channel, error) {
	mk := msg.GetKind()
	if mk == MK_RESPONSE {
		p := msg.ToPackage()
		if p != nil {
			v, _ := FrameCoders.SessionInfo.Pop(p, this.id, Coders.Uint32)
			if v != nil {
				rv := v.(uint32)
				for _, ch := range list {
					if ch.Id() == rv {
						return ch, nil
					}
				}
				logger.Debug(tag, "%s can't find channel %d", this, rv)
				return nil, nil
			}
		}
		return nil, nil
	} else if mk == MK_REQUEST && in != nil {
		p := msg.ToPackage()
		if p != nil {
			FrameCoders.SessionInfo.Set(p, this.id, in.Id(), Coders.Uint32)
		}
	}

	l := len(list)
	if l == 1 {
		return list[0], nil
	}
	c := 0
	if left2right {
		this.lcount++
		c = this.lcount
	} else {
		this.rcount++
		c = this.rcount
	}
	return list[c%l], nil
}

func (this *Broker) run() {
	defer func() {
		close(this.c)
	}()
	var llist, rlist []Channel
	var lmark, rmark int64
	for {
		req := <-this.c
		if req == nil {
			return
		}
		if this.left.IsClosing() {
			return
		}
		var dis BrokerDispatch
		var list []Channel
		if req.l2r {
			rlist, rmark = this.right.Snapshot(rmark, rlist)
			list = rlist
			dis = this.RightDispatcher
			if dis == nil {
				dis = this.DefaultDispatch
			}
		} else {
			llist, lmark = this.left.Snapshot(lmark, llist)
			list = llist
			dis = this.LeftDispatcher
			if dis == nil {
				dis = this.DefaultDispatch
			}
		}

		var ch Channel
		if len(list) > 0 {
			var err error
			ch, err = dis(req.in, list, req.msg, req.l2r)
			if err != nil {
				logger.Debug(tag, "%s dispatch fail - %s", this, err)
				continue
			}
		}
		if ch == nil {
			logger.Debug(tag, "%s dispatch no channel to send", this)
			err := errors.New("can't deliver")
			req.msg.TryRelyError(req.in, err)
			continue
		}

		if req.in != nil {
			ctrl := FrameCoders.Trace
			p := req.msg.ToPackage()
			if ctrl.Has(p) {
				info := fmt.Sprintf("%s -> %s", this, ch)
				rmsg := ctrl.CreateReply(req.msg, info)
				go req.in.SendMessage(rmsg)
			}
		}

		// send message
		logger.Debug(tag, "%s ===> (%s)", this, ch)
		err := ch.SendMessage(req.msg)
		if err != nil {
			logger.Debug(tag, "%s ===> (%s) fail - %s", this, ch, err)
			req.msg.TryRelyError(ch, err)
		}
	}
}

func (this *Broker) Stop() bool {
	if this.left.AskClose() {
		this.right.AskClose()
		this.c <- nil
	}
	return true
}

func (this *Broker) Cleanup() bool {
	this.WaitClose()
	return true
}

func (this *Broker) WaitClose() {
	this.left.WaitClosed()
	this.right.WaitClosed()
}

func (this *Broker) AddLeft(ch Channel, link bool) {
	if this.left.Add(ch) {
		if link {
			ch.SetMessageListner(func(msg *Message) error {
				return this.LeftSendMessage(ch, msg)
			})
		}
	}
}

func (this *Broker) RemoveLeft(ch Channel) {
	this.left.Remove(ch)
}

func (this *Broker) AddRight(ch Channel, link bool) {
	if this.right.Add(ch) {
		if link {
			ch.SetMessageListner(func(msg *Message) error {
				return this.RightSendMessage(ch, msg)
			})
		}
	}
}

func (this *Broker) RemoveRight(ch Channel) {
	this.right.Remove(ch)
}

func (this *Broker) LeftSendMessage(ch Channel, msg *Message) (r error) {
	if this.left.IsClosing() {
		return errors.New("closed")
	}
	defer func() {
		if recover() != nil {
			r = errors.New("closed")
		}
	}()
	this.c <- &brokerREQ{ch, msg, true}
	return nil
}

func (this *Broker) RightSendMessage(ch Channel, msg *Message) (r error) {
	if this.left.IsClosing() {
		return errors.New("closed")
	}
	defer func() {
		if recover() != nil {
			r = errors.New("closed")
		}
	}()
	this.c <- &brokerREQ{ch, msg, false}
	return nil
}
