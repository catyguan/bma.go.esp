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

type ServiceHandler func(msg *esnp.Message, rep ServiceResponser) error

type ServiceRequestContext struct {
	Message   *esnp.Message
	Responser ServiceResponser
}

func DoServiceHandle(h ServiceHandler, msg *esnp.Message, rep ServiceResponser) error {
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
		return h(msg, rep)
	}()
	if err != nil {
		rmsg := msg.ReplyMessage()
		rmsg.BeError(err)
		rep.SendMessage(rmsg)
	}
	return nil
}
