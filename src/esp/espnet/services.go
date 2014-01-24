package espnet

import (
	"fmt"
	"sync/atomic"
)

var (
	globalChanneIdSeq uint32
	globalMessageId   MessageIdGenerator = MessageIdGenerator{0, true}
)

func NextChanneId() uint32 {
	for {
		v := atomic.AddUint32(&globalChanneIdSeq, 1)
		if v != 0 {
			return v
		}
	}
}

func init() {
	globalMessageId.InitMessageIdGenerator()
}

func NextMessageId() uint64 {
	return globalMessageId.Next()
}

// Helper
type ServiceResponser4S struct {
	S MessageSender
}

func (this *ServiceResponser4S) GetChannel() Channel {
	return nil
}

func (this *ServiceResponser4S) SendMessage(msg *Message) error {
	return this.S(msg)
}

type ServiceResponser4C struct {
	C Channel
}

func (this *ServiceResponser4C) GetChannel() Channel {
	return this.C
}

func (this *ServiceResponser4C) SendMessage(msg *Message) error {
	return this.C.SendMessage(msg)
}

func ConnectService(ch Channel, sh ServiceHandler) error {
	csh := new(ServiceResponser4C)
	csh.C = ch
	ch.SetMessageListner(func(msg *Message) error {
		return DoServiceHandle(sh, msg, csh)
	})
	return nil
}

func Connect(left Channel, right Channel, closeOnBreak bool) {
	cid := fmt.Sprintf("%p_%p", left, right)
	left.SetMessageListner(func(msg *Message) error {
		return right.SendMessage(msg)
	})
	left.SetCloseListener(cid, func() {
		left.SetMessageListner(nil)
		right.SetMessageListner(nil)
		right.SetCloseListener(cid, nil)
		if closeOnBreak {
			right.AskClose()
		}
	})
	right.SetMessageListner(func(msg *Message) error {
		return left.SendMessage(msg)
	})
	right.SetCloseListener(cid, func() {
		right.SetMessageListner(nil)
		left.SetMessageListner(nil)
		left.SetCloseListener(cid, nil)
		if closeOnBreak {
			left.AskClose()
		}
	})
}
