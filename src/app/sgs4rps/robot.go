package sgs4rps

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"esp/espnet/mempipeline"
)

type Robot struct {
	psid int
	name string
	sock *espsocket.Socket
}

func newRobot(n string) (*Robot, *espsocket.Socket) {
	this := new(Robot)
	this.name = n
	mp := mempipeline.NewMemPipeline(this.name, 16)
	this.sock = espsocket.NewSocket(mp.ChannelA())
	this.sock.SetMessageListner(this.onMessageIn)
	rsock := espsocket.NewSocket(mp.ChannelB())
	this.run()
	return this, rsock
}

func (this *Robot) run() {
	go func() {

	}()
}

func (this *Robot) close() {
	this.sock.AskClose()
}

func (this *Robot) onMessageIn(msg *esnp.Message) error {
	return nil
}
