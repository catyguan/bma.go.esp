package securec

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"esp/espnet/secure/securep"
	"logger"
	"time"
)

const (
	tag = "securec"
)

func DoBaseAuth(sock espsocket.Socket, key string, timeout time.Duration) error {
	var req securep.BaseAuthRequest
	var rep securep.BaseAuthResponse

	espsocket.SetDeadline(sock, time.Now().Add(timeout))
	defer espsocket.ClearDeadline(sock)

	msg1 := esnp.NewRequestMessage()
	err1 := req.Encode(msg1)
	if err1 != nil {
		return err1
	}
	err1 = sock.WriteMessage(msg1)
	if err1 != nil {
		return err1
	}
	rmsg1, err2 := sock.ReadMessage()
	if err2 != nil {
		return err2
	}
	err2 = rmsg1.ToError()
	if err2 != nil {
		return err2
	}
	err2 = rep.Decode(rmsg1)
	if err2 != nil {
		return err2
	}
	err2 = rep.Valid()
	if err2 != nil {
		return err2
	}
	logger.Debug(tag, "BaseAuth request token=%s, key=%s", rep.Token, key)
	autk := securep.CreateAuthToken(rep.Token, key)
	req.Reset()
	req.Token = autk
	msg2 := esnp.NewRequestMessage()
	err2 = req.Encode(msg2)
	if err2 != nil {
		return err2
	}
	err2 = sock.WriteMessage(msg2)
	if err2 != nil {
		return err2
	}
	_, err3 := sock.ReadMessage()
	if err3 != nil {
		return err3
	}
	return nil
}
