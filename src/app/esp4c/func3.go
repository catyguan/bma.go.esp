package main

import (
	"esp/espnet/esnp"
	"logger"
)

func doPProf(address string, n string) {
	sock := createSocket(address)
	if sock == nil {
		return
	}
	defer sock.AskClose()

	msg := esnp.NewMessage()
	msg.GetAddress().SetCall("pprof", n)
	err := sock.WriteMessage(msg)
	if err != nil {
		logger.Warn(tag, "call 'pprof' fail - %s", err)
		return
	}
	logger.Info(tag, "call 'pprof' done")
}
