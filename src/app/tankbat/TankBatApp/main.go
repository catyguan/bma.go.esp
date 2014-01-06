package main

import (
	"app/bombman"
	"bmautil/socket"
	"boot"
	"esp/espnet"
)

func main() {
	root := new(mainObj)
	root.run()
}

type mainObj struct {
	service        *bombman.Service
	channelHandler *espnet.ChannelCoder4Telnet
}

func (this *mainObj) run() {

	cfile := "config/bombman-config.json"

	this.channelHandler = new(espnet.ChannelCoder4Telnet)
	this.channelHandler.Init()
	this.channelHandler.Error2String = this.error2String

	this.service = bombman.NewService("service")
	boot.QuickDefine(this.service, "", true)

	pointSER := espnet.NewListenPoint("servicePoint", nil, this.socketAcceptSer)
	boot.QuickDefine(pointSER, "", true)

	boot.Go(cfile)
}

func (this *mainObj) socketAcceptSer(sock *socket.Socket) error {
	ch := espnet.NewSocketChannelC(sock, this.channelHandler)
	ch.SetProperty(espnet.PROP_SOCKET_TRACE, 64)

	s := this.service
	s.Add(ch)

	return nil
}

func (this *mainObj) error2String(err error) string {
	return "ERROR " + err.Error() + "\n"
}
