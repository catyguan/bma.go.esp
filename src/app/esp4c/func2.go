package main

import (
	"esp/espnet/esnp"
	"logger"
)

func doReload(address string) {
	c := createClient(address)
	if c == nil {
		return
	}
	defer c.Close()

	msg := esnp.NewMessage()
	msg.GetAddress().Set(esnp.ADDRESS_OP, "reload")
	err := c.SendMessage(msg)
	if err != nil {
		logger.Warn(tag, "call 'reload' fail - %s", err)
		return
	}
	logger.Info(tag, "call 'reload' done")
}
