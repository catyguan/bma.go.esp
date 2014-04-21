package espservice

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
)

type ServiceHandler func(sock *espsocket.Socket, msg *esnp.Message) error

type ServiceRequestContext struct {
	Channel espchannel.Channel
	Message *esnp.Message
}

func DoServiceHandle(h ServiceHandler, sock *espsocket.Socket, msg *esnp.Message) error {
	err := func() (r error) {
		defer func() {
			r2 := recover()
			if r2 != nil {
				if r3, ok := r2.(error); ok {
					r = r3
				} else {
					r = fmt.Errorf("%s", r2)
				}
			}
		}()
		return h(sock, msg)
	}()
	if err != nil {
		rmsg := msg.ReplyMessage()
		rmsg.BeError(err)
		sock.SendMessage(rmsg, nil)
	}
	return nil
}
