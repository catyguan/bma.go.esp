package secures

import (
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"logger"
	"time"
)

func NewAuthEntry(maxAuthTime time.Duration, maxsizeNoAuth, maxsizeAuth int, auth espservice.ServiceHandler, e espservice.ServiceEntry) espservice.ServiceEntry {
	return func(sock espsocket.Socket) {
		defer sock.AskClose()
		espsocket.SetDeadline(sock, time.Now().Add(maxAuthTime))
		sock.SetProperty(espsocket.PROP_MESSAGE_MAXSIZE, maxsizeNoAuth)
		msg, err0 := sock.ReadMessage()
		if err0 != nil {
			logger.Debug(tag, "read auth message fail - %s", err0)
			return
		}
		err2 := auth(sock, msg)
		if err2 != nil {
			logger.Debug(tag, "auth fail - %s", err2)
			return
		}
		espsocket.ClearDeadline(sock)
		sock.SetProperty(espsocket.PROP_MESSAGE_MAXSIZE, maxsizeAuth)
		e(sock)
	}
}
