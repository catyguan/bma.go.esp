package main

import (
	"app/tankbat"
	"bmautil/socket"
	"boot"
	"esp/espnet"
)

func main() {
	root := new(mainObj)
	root.run()
}

type mainObj struct {
	service        *tankbat.Service
	channelHandler *espnet.ChannelCoder4Telnet
}

func (this *mainObj) run() {

	cfile := "config/sgs4rps-config.json"

	this.channelHandler = new(espnet.ChannelCoder4Telnet)
	this.channelHandler.Init()
	this.channelHandler.Error2String = this.error2String

	this.service = tankbat.NewService("service")
	boot.Add(this.service, "", false)

	pointSER := socket.NewListenPoint("servicePoint", nil, this.socketAcceptSer)
	boot.Add(pointSER, "", false)

	boot.Go(cfile)
}

func (this *mainObj) socketAcceptSer(sock *socket.Socket) error {
	sock.SetWriteChanSize(128)

	ch := espnet.NewSocketChannelC(sock, this.channelHandler)
	ch.SetProperty(espnet.PROP_SOCKET_TRACE, 64)

	s := this.service
	s.Add(ch)

	return nil
}

func (this *mainObj) error2String(err error) string {
	return "ERROR " + err.Error() + "\n"
}
