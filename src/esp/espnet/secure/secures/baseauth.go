package secures

import (
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"esp/espnet/secure/securep"
	"fmt"
	"logger"
	"time"
)

func NewBaseAuth(key string) espservice.ServiceHandler {
	return func(sock espsocket.Socket, msg1 *esnp.Message) error {
		var req securep.BaseAuthRequest
		var rep securep.BaseAuthResponse

		tk := fmt.Sprintf("%d", time.Now().UnixNano())
		logger.Debug(tag, "BaseAuth request token=%s, key=%s", tk, key)
		rep.Token = tk
		rmsg1 := msg1.ReplyMessage()
		err1 := rep.Encode(rmsg1)
		if err1 != nil {
			logger.Debug(tag, "BaseAuth encode response 1 fail - %s", err1)
			return err1
		}
		err1 = sock.WriteMessage(rmsg1)
		if err1 != nil {
			logger.Debug(tag, "BaseAuth write response 1 fail - %s", err1)
			return err1
		}

		msg2, err2 := sock.ReadMessage()
		if err2 != nil {
			logger.Debug(tag, "BaseAuth read request 2 fail - %s", err2)
			return err2
		}
		err2 = req.Decode(msg2)
		if err2 != nil {
			logger.Debug(tag, "BaseAuth decode request 2 fail - %s", err2)
			return err2
		}
		err2 = req.Valid()
		if err2 != nil {
			logger.Debug(tag, "BaseAuth valid request 2 fail - %s", err2)
			return err2
		}
		atki := req.Token
		atkm := securep.CreateAuthToken(tk, key)
		if atki != atkm {
			return logger.Warn(tag, "BaseAuth %s invalid auth token (in=%s, me=%s)", sock, atki, atkm)
		}
		rep.Reset()
		rmsg2 := msg2.ReplyMessage()
		err2 = rep.Encode(rmsg2)
		if err2 != nil {
			logger.Debug(tag, "BaseAuth encode response 2 fail - %s", err2)
			return err2
		}
		err2 = sock.WriteMessage(rmsg2)
		if err2 != nil {
			logger.Debug(tag, "BaseAuth write response 2 fail - %s", err2)
			return err2
		}
		return nil
	}
}
