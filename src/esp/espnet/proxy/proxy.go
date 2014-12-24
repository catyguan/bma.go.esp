package proxy

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"logger"
)

type ForwardSetting struct {
}

func Forward(sin espsocket.Socket, sout espsocket.Socket, msg *esnp.Message, fs *ForwardSetting) (rerr error, failOver bool) {
	err2 := sout.WriteMessage(msg)
	if err2 != nil {
		sout.AskClose()
		return err2, true
	}
	if !msg.IsRequest() {
		logger.Debug(tag, "not request, skip response")
		return nil, false
	}
	rmsg, err3 := sout.ReadMessage(false)
	if err3 != nil {
		sout.AskClose()
		return err3, false
	}
	err4 := sin.WriteMessage(rmsg)
	if err4 != nil {
		sout.AskClose()
		return err4, false
	}
	return nil, false
}
