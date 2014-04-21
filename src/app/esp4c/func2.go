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
	msg.GetAddress().SetCall("sys", "reload")
	err := sock.SendMessage(msg, nil)
	if err != nil {
		logger.Warn(tag, "call 'reload' fail - %s", err)
		return
	}
	logger.Info(tag, "call 'reload' done")
}
