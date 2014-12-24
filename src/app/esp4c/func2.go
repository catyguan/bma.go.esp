package main

import (
	"esp/espnet/esnp"
	"logger"
)

func doReload(address string) {
	sock := createSocket(address)
	if sock == nil {
		return
	}
	defer sock.AskClose()

	msg := esnp.NewMessage()
	msg.GetAddress().SetCall("boot", "reload")
	err := sock.WriteMessage(msg)
	if err != nil {
		logger.Warn(tag, "call 'reload' fail - %s", err)
		return
	}
	logger.Info(tag, "call 'reload' done")
}

func doShutdow(address string) {
	sock := createSocket(address)
	if sock == nil {
		return
	}
	defer sock.AskClose()

	msg := esnp.NewMessage()
	msg.GetAddress().SetCall("boot", "shutdown")
	err := sock.WriteMessage(msg)
	if err != nil {
		logger.Warn(tag, "call 'shutdown' fail - %s", err)
		return
	}
	logger.Info(tag, "call 'shutdown' done")
}
