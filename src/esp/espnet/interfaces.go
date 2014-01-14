package espnet

import (
	"errors"
	"fmt"
)

type SupportProp interface {
	GetProperty(name string) (interface{}, bool)
	SetProperty(name string, val interface{}) bool
}

// Message
type MessageListener func(msg *Message) error
type MessageSender func(msg *Message) error

// Service
type ServiceResponser interface {
	GetChannel() Channel
	SendMessage(replyMsg *Message) error
}

type ServiceHandler func(msg *Message, rep ServiceResponser) error

type ServiceRequestContext struct {
	Message   *Message
	Responser ServiceResponser
}

func DoServiceHandle(h ServiceHandler, msg *Message, rep ServiceResponser) error {
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
