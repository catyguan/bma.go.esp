package espservice

import (
	"errors"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"fmt"
)

// Service
type ServiceResponser interface {
	GetChannel() espchannel.Channel
	SendMessage(replyMsg *esnp.Message) error
}

type ServiceHandler func(ch espchannel.Channel, msg *esnp.Message) error

type ServiceRequestContext struct {
	Channel espchannel.Channel
	Message *esnp.Message
}

func DoServiceHandle(h ServiceHandler, ch espchannel.Channel, msg *esnp.Message) error {
	err := func() (r error) {
		defer func() {
			r2 := recover()
			if r2 != nil {
				if r3, ok := r2.(error); ok {
					r = r3
				} else {
					r = errors.New(fmt.Sprintf("%s", r2))
				}
			}
		}()
		return h(ch, msg)
	}()
	if err != nil {
		rmsg := msg.ReplyMessage()
		rmsg.BeError(err)
		ch.PostMessage(rmsg)
	}
	return nil
}
