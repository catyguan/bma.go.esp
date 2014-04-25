package sgs4rps

import (
	"bmautil/goo"
	"esp/espnet/espsocket"
)

type Service struct {
	name   string
	config *configInfo
	goo    goo.Goo
	m      matrix
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	this.m.s = this
	this.goo.InitGoo(tag, 128, this.clear)
	this.goo.Run()
	return this
}

func (this *Service) clear() {
	this.m.clear()
}

func (this *Service) toPlayerAPI() playerAPI {
	return this
}

func (this *Service) startMatrix() {
	this.m.start()
}

func (this *Service) stopMatrix() {
	this.m.stop()
}

// playerAPI
func (this *Service) NewPlayer(sock *espsocket.Socket) (int, error) {
	return 0, nil
}

func (this *Service) SetNick(psid int, n string) error {
	return nil
}

func (this *Service) Play(psid int, rps int) error {
	return nil
}

func (this *Service) Quit(psid int) error {
	return nil
}
