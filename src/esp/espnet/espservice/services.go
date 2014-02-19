package espservice

import (
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"fmt"
)

// Helper
type ServiceResponser4S struct {
	S esnp.MessageSender
}

func (this *ServiceResponser4S) GetChannel() espchannel.Channel {
	return nil
}

func (this *ServiceResponser4S) SendMessage(msg *esnp.Message) error {
	return this.S(msg)
}

type ServiceResponser4C struct {
	C espchannel.Channel
}

func (this *ServiceResponser4C) GetChannel() espchannel.Channel {
	return this.C
}

func (this *ServiceResponser4C) SendMessage(msg *esnp.Message) error {
	return this.C.SendMessage(msg)
}

func ConnectService(ch espchannel.Channel, sh ServiceHandler) error {
	csh := new(ServiceResponser4C)
	csh.C = ch
	ch.SetMessageListner(func(msg *esnp.Message) error {
		return DoServiceHandle(sh, msg, csh)
	})
	return nil
}

func Connect(left espchannel.Channel, right espchannel.Channel, closeOnBreak bool) {
	cid := fmt.Sprintf("%p_%p", left, right)
	left.SetMessageListner(func(msg *esnp.Message) error {
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
	right.SetMessageListner(func(msg *esnp.Message) error {
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
