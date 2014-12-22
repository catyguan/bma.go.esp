package espservice

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
)

const (
	tag = "espservice"
)

type ServiceHandler func(sock espsocket.Socket, msg *esnp.Message) error
type ServiceEntry func(sock espsocket.Socket)

type ServiceRequestContext struct {
	Sock    espsocket.Socket
	Message *esnp.Message
}

func DoServiceHandle(h ServiceHandler, sock espsocket.Socket, msg *esnp.Message) {
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
		sock.WriteMessage(rmsg)
	}
}
