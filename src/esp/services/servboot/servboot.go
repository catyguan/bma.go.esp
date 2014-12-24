package servboot

import (
	"boot"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"logger"
	"time"
)

const (
	tag = "servboot"

	NAME_SERVICE     = "boot"
	NAME_OP_RELOAD   = "reload"
	NAME_OP_SHUTDOWN = "shutdown"
)

func InitMux(mux *espservice.ServiceMux) {
	mux.AddHandler(NAME_SERVICE, NAME_OP_RELOAD, ServOP_Reload)
	mux.AddHandler(NAME_SERVICE, NAME_OP_SHUTDOWN, ServOP_Shutdown)
}

func ServOP_Reload(sock espsocket.Socket, msg *esnp.Message) error {
	logger.Info(tag, "op reload from %s", sock)
	time.AfterFunc(100*time.Millisecond, func() {
		boot.Restart()
	})
	sock.WriteMessage(msg.ReplyMessage())
	return nil
}

func ServOP_Shutdown(sock espsocket.Socket, msg *esnp.Message) error {
	logger.Info(tag, "op shutdown from %s", sock)
	time.AfterFunc(100*time.Millisecond, func() {
		boot.Shutdown()
	})
	sock.WriteMessage(msg.ReplyMessage())
	return nil
}
