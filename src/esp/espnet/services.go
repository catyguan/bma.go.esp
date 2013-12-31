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
func ConnectService(ch Channel, sh ServiceHandler) error {
	ch.SetMessageListner(func(msg *Message) error {
		return DoServiceHandle(sh, msg, ch.SendMessage)
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
